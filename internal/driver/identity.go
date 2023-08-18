package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
)

const (
	pluginName    = "csi.experimental.systems"
	pluginVersion = "v0.0.0"
)

var _ csi.IdentityServer = identityServer{}

type identityServer struct {
}

func (identityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	out := &csi.GetPluginInfoResponse{
		Name:          pluginName,
		VendorVersion: pluginVersion,
	}

	return out, nil
}

func (identityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	out := &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_UNKNOWN,
					},
				},
			},
		},
	}

	return out, nil
}

func (identityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	out := &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{
			Value: true,
		},
	}

	return out, nil
}
