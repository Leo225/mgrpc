package mgrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	clientConn *grpc.ClientConn
	opts       ClientOptions
}

func (c *GrpcClient) GetConn() *grpc.ClientConn {
	return c.clientConn
}

func (c *GrpcClient) Close(ctx context.Context) {
	if c.clientConn == nil {
		c.opts.Log.Fatalf(ctx, "close %s grpc client is nil", c.opts.Address)
		return
	}

	err := c.clientConn.Close()
	if err != nil {
		c.opts.Log.Fatalf(ctx, "close %s grpc client: %v", c.opts.Address, err)
		return
	}
	if c.opts.discovery != nil {
		err = c.opts.discovery.Close()
		if err != nil {
			c.opts.Log.Fatalf(ctx, "close %s grpc client: %v", c.opts.Address, err)
			return
		}
	}
}

func NewGrpcClient(ctx context.Context, opts ...ClientOptionFunc) *GrpcClient {
	client := new(GrpcClient)
	client.opts = NewClientOptions(opts...)
	grpcClientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	}

	if len(client.opts.StreamClientInterceptors) > 0 {
		for _, v := range client.opts.StreamClientInterceptors {
			grpcClientOptions = append(grpcClientOptions, grpc.WithStreamInterceptor(v))
		}
	}
	if len(client.opts.UnaryClientInterceptors) > 0 {
		for _, v := range client.opts.UnaryClientInterceptors {
			grpcClientOptions = append(grpcClientOptions, grpc.WithUnaryInterceptor(v))
		}
	}
	if client.opts.discovery != nil {
		builder, err := client.opts.discovery.Discover(client.opts.Name)
		if err != nil {
			client.opts.Log.Fatalf(ctx, "%s grpc client discovery error: %v", client.opts.Name, err)
		}
		grpcClientOptions = append(grpcClientOptions, grpc.WithResolvers(builder))
		cc, err := grpc.Dial(client.opts.discovery.Address(), grpcClientOptions...)
		if err != nil {
			client.opts.Log.Fatalf(ctx, "%s grpc client dial error: %v", client.opts.Name, err)
		}
		client.clientConn = cc
		return client
	}
	cc, err := grpc.Dial(client.opts.Address, grpcClientOptions...)
	if err != nil {
		client.opts.Log.Fatalf(ctx, "%s grpc client dial error: %v", client.opts.Name, err)
	}
	client.clientConn = cc
	return client
}
