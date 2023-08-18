package driver

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
	path     string
}

func NewServer(config Config) (*Server, error) {
	if config.Path == "" {
		return nil, errInvalidConfig
	}

	listener, err := net.Listen("unix", config.Path)
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(loggingInterceptor()),
	}
	server := grpc.NewServer(opts...)

	out := &Server{
		server:   server,
		listener: listener,
		path:     config.Path,
	}

	return out, nil
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
