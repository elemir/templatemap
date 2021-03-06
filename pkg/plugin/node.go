package plugin

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/elemir/templatemap/pkg/k8su"
	"github.com/elemir/templatemap/pkg/template"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeServer struct {
	nodeId    string
	log       *logrus.Entry
	clientset *kubernetes.Clientset
	config    *rest.Config
	csi.UnimplementedNodeServer
}

func NewNodeServer(log *logrus.Entry, kubeconfig string, nodeId string) (*NodeServer, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &NodeServer{
		nodeId:    nodeId,
		log:       log,
		clientset: clientset,
		config:    config,
	}, nil
}

func (ns *NodeServer) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{}, nil
}

func (ns *NodeServer) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: ns.nodeId,
	}, nil
}

func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability missing in request")
	}
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}

	configMapName := req.GetVolumeContext()["configMapName"]
	namespace := req.GetVolumeContext()["csi.storage.k8s.io/pod.namespace"]
	podName := req.GetVolumeContext()["csi.storage.k8s.io/pod.name"]
	serviceAccountName := req.GetVolumeContext()["serviceAccountName"]

	if serviceAccountName == "" {
		serviceAccountName = req.GetVolumeContext()["csi.storage.k8s.io/serviceAccount.name"]
	}

	targetPath := req.GetTargetPath()

	clientset, err := k8su.RunAs(ns.clientset, ns.config, namespace, serviceAccountName)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot run as specific SA: %s", err))
	}

	tmpl, err := template.NewTemplate(clientset, namespace, podName)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot create templater from config: %s", err))
	}

	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(configMapName, v1.GetOptions{})
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Cannot found specific ConfigMap '%s': %s", configMapName, err))
	}

	err = os.Mkdir(targetPath, 0750)
	if err != nil && !os.IsExist(err) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for key, data := range cm.Data {
		err = tmpl.GenerateFile(targetPath, key, data)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}

	targetPath := req.GetTargetPath()

	err := os.RemoveAll(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ns.log.WithField("targetPath", targetPath).Debug("removed target path")

	return &csi.NodeUnpublishVolumeResponse{}, nil
}
