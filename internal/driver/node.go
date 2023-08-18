package driver

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/jnschaeffer/csi-driver-hello/internal/manager"
	"k8s.io/utils/mount"
)

type nodeServer struct {
	nodeName string
	mounter  mount.Interface
	manager  manager.Interface

	csi.UnimplementedNodeServer
}

func (nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	resp := &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_UNKNOWN,
					},
				},
			},
		},
	}

	return resp, nil
}

func (s *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	out := &csi.NodeGetInfoResponse{
		NodeId: s.nodeName,
	}

	return out, nil
}

func (s *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Printf("publishing volume at %s", req.TargetPath)

	var success bool

	defer func() {
		if !success {
			s.mounter.Unmount(req.TargetPath)
		}
	}()

	volumeID := req.GetVolumeId()

	if err := s.manager.ManageVolume(ctx, volumeID); err != nil {
		return nil, err
	}

	targetPath := req.GetTargetPath()

	notMnt, err := mount.IsNotMountPoint(s.mounter, targetPath)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		log.Print("mount point does not exist, creating")
		if err := os.MkdirAll(targetPath, 0440); err != nil {
			return nil, err
		}
		notMnt = true
	case err != nil:
		return nil, err
	}

	if !notMnt {
		log.Print("already mounted, nothing to do")
		success = true
		return &csi.NodePublishVolumeResponse{}, nil
	}

	log.Printf("mounting volume at %s", targetPath)

	if err := s.mounter.Mount(s.manager.PathFromVolumeID(volumeID), targetPath, "", []string{"bind", "ro"}); err != nil {
		return nil, err
	}

	success = true

	return &csi.NodePublishVolumeResponse{}, nil
}

func (s *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if err := s.mounter.Unmount(req.GetTargetPath()); err != nil {
		log.Printf("error unmounting volume: %s", err)
	}

	if err := s.manager.UnmanageVolume(ctx, req.GetVolumeId()); err != nil {
		log.Printf("error unmanaging volume: %s", err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}
