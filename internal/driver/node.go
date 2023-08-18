package driver

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

type nodeServer struct {
	csi.UnimplementedNodeServer
}

func (nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Printf("publishing volume at %s", req.TargetPath)

	if err := os.Mkdir(req.TargetPath, 0755); err != nil {
		return nil, err
	}

	path := filepath.Join(req.TargetPath, "hello")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if _, err := f.Write([]byte("hello, world!\n")); err != nil {
		return nil, err
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if err := os.RemoveAll(req.TargetPath); err != nil {
		return nil, err
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}
