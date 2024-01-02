package mgrpc

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type GrpcContextHandler func(ctx context.Context) context.Context

func UnaryServerInterceptor(f GrpcContextHandler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = f(ctx)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(f GrpcContextHandler) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wss := grpc_middleware.WrapServerStream(ss)
		wss.WrappedContext = f(ss.Context())
		return handler(srv, wss)
	}
}
