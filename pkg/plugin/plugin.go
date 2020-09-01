package plugin

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Plugin struct {
	log      *logrus.Entry
	endpoint string
	nodeId   string
}

func NewPlugin(nodeId string, endpoint string, log *logrus.Logger) (*Plugin, error) {
	entry := log.WithField("name", PluginName).WithField("version", PluginVersion)
	entry.Infof("initialize plugin")

	return &Plugin{
		log:      entry,
		endpoint: endpoint,
		nodeId:   nodeId,
	}, nil
}

func (p *Plugin) Run() error {
	err := os.Remove("/csi/csi.sock")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot remove unix socket for endpoint '%s': %w", p.endpoint, err)
	}

	listener, err := net.Listen("unix", "/csi/csi.sock")
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
	nodeSrv, err := NewNodeServer(p.log, p.nodeId)
	if err != nil {
		return fmt.Errorf("cannot create node server: %w", err)
	}

	csi.RegisterIdentityServer(srv, identSrv)
	csi.RegisterNodeServer(srv, nodeSrv)

	err = srv.Serve(listener)
	if err != nil {
		return fmt.Errorf("cannot listen on specific endpoint '%s': %w", p.endpoint, err)

	}

	return nil
}
