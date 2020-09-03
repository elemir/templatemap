package plugin

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type ControllerServer struct {
	csi.UnimplementedControllerServer
}

func NewControllerServer() *ControllerServer {
	return &ControllerServer{}
}

func (c *ControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: []*csi.ControllerServiceCapability{
			{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_UNKNOWN,
					},
				},
			},
		},
	}, nil
}
