package plugin

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

const (
	PluginName    = "csi-templatemap"
	PluginVersion = "0.0.1"
)

type IdentityServer struct{}

func NewIdentityServer() *IdentityServer {
	return &IdentityServer{}
}

func (is *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name:          PluginName,
		VendorVersion: PluginVersion,
	}, nil
}

func (is *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{},
		/*			{
						Type: &csi.PluginCapability_Service_{
							Service: &csi.PluginCapability_Service{
								Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
							},
						},
					},
				}, */
	}, nil
}

func (is *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}
