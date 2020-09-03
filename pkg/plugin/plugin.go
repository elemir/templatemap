package plugin

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Plugin struct {
	log        *logrus.Entry
	endpoint   string
	scheme     string
	kubeconfig string
	nodeId     string
}

func NewPlugin(nodeId string, endpoint string, log *logrus.Logger) (*Plugin, error) {
	entry := log.WithField("name", PluginName).WithField("version", PluginVersion)
	entry.Infof("initialize plugin")

	es := strings.SplitN(endpoint, "://", 2)

	if len(es) != 2 {
		return nil, fmt.Errorf("cannot parse endpoint '%s'", endpoint)
	}

	return &Plugin{
		log:      entry,
		endpoint: es[1],
		scheme:   es[0],
		nodeId:   nodeId,
	}, nil
}

func (p *Plugin) SetKubeConfig(kubeconfig string) {
	p.kubeconfig = kubeconfig
}

func (p *Plugin) Run() error {
	if p.scheme == "unix" {
		err := os.Remove(p.endpoint)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("cannot remove unix socket for endpoint '%s': %w", p.endpoint, err)
		}
	}

	listener, err := net.Listen(p.scheme, p.endpoint)
	if err != nil {
		return fmt.Errorf("cannot listen on specific endpoint '%s': %w", p.endpoint, err)
	}

	errInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			p.log.WithError(err).WithField("method", info.FullMethod).Error("method failed")
		} else {
			p.log.WithField("method", info.FullMethod).Debug("running method")
		}
		return resp, err
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(errInterceptor))

	identSrv := NewIdentityServer()
	ctrlSrv := NewControllerServer()
	nodeSrv, err := NewNodeServer(p.log, p.kubeconfig, p.nodeId)
	if err != nil {
		return fmt.Errorf("cannot create node server: %w", err)
	}

	csi.RegisterIdentityServer(srv, identSrv)
	csi.RegisterControllerServer(srv, ctrlSrv)
	csi.RegisterNodeServer(srv, nodeSrv)

	err = srv.Serve(listener)
	if err != nil {
		return fmt.Errorf("cannot listen on specific endpoint '%s': %w", p.endpoint, err)

	}

	return nil
}
