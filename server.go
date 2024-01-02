package mgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Leo225/mgrpc/resolver"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server *grpc.Server
	opts   ServerOptions
	ctx    context.Context
}

func NewGrpcServer(ctx context.Context, opts ...ServerOptionFunc) *GrpcServer {
	srv := new(GrpcServer)
	srv.ctx = ctx
	srv.opts = NewServerOptions(opts...)
	var grpcServerOptions []grpc.ServerOption
	if len(srv.opts.StreamServerInterceptors) > 0 {
		grpcServerOptions = append(grpcServerOptions,
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(srv.opts.StreamServerInterceptors...)))
	}
	if len(srv.opts.UnaryServerInterceptor) > 0 {
		grpcServerOptions = append(grpcServerOptions,
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(srv.opts.UnaryServerInterceptor...)))
	}
	srv.server = grpc.NewServer(grpcServerOptions...)
	return srv
}

func (s *GrpcServer) Run() {
	var addr string
	s.register()

	if s.opts.RandomPort {
		addr = ":0"
	} else {
		addr = s.opts.Address
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.opts.Log.Fatalf(s.ctx, "%s grpc server failed to listen: %v", s.opts.Name, err)
		return
	}
	if s.opts.registry != nil {
		serviceInfo := resolver.NewServiceInfo()
		serviceInfo.Name = s.opts.Name
		if s.opts.RandomPort {
			ipAddr, err := LocalIPv4()
			if err != nil {
				s.opts.Log.Fatalf(s.ctx, "LocalIPv4: %v", err)
				return
			}
			port := lis.Addr().(*net.TCPAddr).Port
			serviceInfo.Address = fmt.Sprintf("%s:%d", ipAddr, port)
		} else {
			serviceInfo.Address = addr
		}

		err = s.opts.registry.Register(serviceInfo)
		if err != nil {
			s.opts.Log.Fatalf(s.ctx, "%s grpc server failed to register: %v", s.opts.Name, err)
		}
	}

	if err := s.server.Serve(lis); err != nil {
		s.opts.Log.Fatalf(s.ctx, "%s grpc server failed to serve: %v", s.opts.Name, err)
	}
}

func (s *GrpcServer) Start() {
	go func() {
		s.Run()
	}()
}

func (s *GrpcServer) Stop() {
	if s.server == nil {
		s.opts.Log.Warnf(s.ctx, "stop %s grpc server is nil", s.opts.Name)
	} else {
		s.server.Stop()
	}

	if s.opts.registry != nil {
		err := s.opts.registry.Close()
		if err != nil {
			s.opts.Log.Errorf(s.ctx, "%s grpc server failed to unregister: %v", s.opts.Name, err)
		}
	}
}

func (s *GrpcServer) Get() *grpc.Server {
	return s.server
}

func (s *GrpcServer) register() {
	reflection.Register(s.server)
}
