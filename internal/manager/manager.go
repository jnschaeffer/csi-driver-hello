package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Interface interface {
	PathFromVolumeID(id string) string
	ManageVolume(ctx context.Context, id string) error
	UnmanageVolume(ctx context.Context, id string) error
}

var _ Interface = &manager{}

func NewManager(config Config) (Interface, error) {
	if config.Path == "" {
		return nil, ErrInvalidConfig
	}

	out := &manager{
		path: config.Path,
	}

	return out, nil
}

type manager struct {
	path string
}

func (m *manager) PathFromVolumeID(id string) string {
	volumePath := filepath.Join(m.path, id)

	return volumePath
}

func (m *manager) ManageVolume(ctx context.Context, id string) error {
	volumePath := m.PathFromVolumeID(id)

	if err := os.MkdirAll(volumePath, 0755); err != nil {
		return err
	}

	path := filepath.Join(volumePath, "hello")

	log.Printf("writing file at %s", path)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := fmt.Fprintf(f, "hello, world!\n"); err != nil {
		return err
	}

	log.Printf("successfully wrote file")

	return nil
}

func (m *manager) UnmanageVolume(ctx context.Context, id string) error {
	if err := os.RemoveAll(m.PathFromVolumeID(id)); err != nil {
		return err
	}

	return nil
}
