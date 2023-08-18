package driver

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/jnschaeffer/csi-driver-hello/internal/manager"
	"google.golang.org/grpc"
	"k8s.io/utils/mount"
)

type Option func(*Server)

func WithManager(m manager.Interface) Option {
	return func(s *Server) {
		s.manager = m
	}
}

type Server struct {
	server   *grpc.Server
	listener net.Listener
	path     string
	manager  manager.Interface
}

func NewServer(config Config, options ...Option) (*Server, error) {
	if config.Path == "" {
		return nil, errInvalidConfig
	}

	srv := &Server{
		path: config.Path,
	}

	for _, opt := range options {
		opt(srv)
	}

	if err := os.Remove(config.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	listener, err := net.Listen("unix", config.Path)
	if err != nil {
		return nil, err
	}

	srv.listener = listener

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(loggingInterceptor()),
	}

	srv.server = grpc.NewServer(opts...)

	csi.RegisterIdentityServer(srv.server, &identityServer{})

	nodeSrv := &nodeServer{
		nodeName: config.NodeName,
		manager:  srv.manager,
		mounter:  mount.New(""),
	}

	csi.RegisterNodeServer(srv.server, nodeSrv)

	return srv, nil
}

func (s *Server) Run() error {
	log.Printf("listening on %s", s.path)

	return s.server.Serve(s.listener)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func (s *Server) ForceStop() {
	s.server.Stop()
}

func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("received %s: %s", info.FullMethod, req)
		resp, err := handler(ctx, req)
		if err != nil {
			log.Printf("error processing request: %s", err)
		} else {
			log.Printf("request completed, response: %s", resp)
		}
		return resp, err
	}
}
