package mgrpc

import (
	"os"

	"github.com/Leo225/mgrpc/resolver"
	"google.golang.org/grpc"
)

type ServerOptions struct {
	Name                     string
	Address                  string
	Log                      ILogger
	StreamServerInterceptors []grpc.StreamServerInterceptor
	UnaryServerInterceptor   []grpc.UnaryServerInterceptor
	RandomPort               bool
	registry                 resolver.Registry
}

type ServerOptionFunc func(o *ServerOptions)

func NewServerOptions(opts ...ServerOptionFunc) ServerOptions {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	opt := ServerOptions{
		Name:       hostname,
		Address:    ":5000",
		RandomPort: false,
		Log:        new(Logger),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithServerName(name string) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.Name = name
	}
}

func WithServerAddress(addr string) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.Address = addr
	}
}

func WithServerLogger(log ILogger) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.Log = log
	}
}

func WithStreamServerInterceptor(streamServerInterceptors ...grpc.StreamServerInterceptor) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.StreamServerInterceptors = streamServerInterceptors
	}
}

func WithUnaryServerInterceptor(unaryServerInterceptors ...grpc.UnaryServerInterceptor) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.UnaryServerInterceptor = unaryServerInterceptors
	}
}

func WithServerRegister(registry resolver.Registry) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.registry = registry
	}
}

func WithServerRandomPort(randomPort bool) ServerOptionFunc {
	return func(o *ServerOptions) {
		o.RandomPort = randomPort
	}
}

type ClientOptions struct {
	Name                     string
	Address                  string
	Log                      ILogger
	StreamClientInterceptors []grpc.StreamClientInterceptor
	UnaryClientInterceptors  []grpc.UnaryClientInterceptor
	discovery                resolver.Discovery
}

type ClientOptionFunc func(o *ClientOptions)

func NewClientOptions(opts ...ClientOptionFunc) ClientOptions {
	opt := ClientOptions{
		Name:    "unknown",
		Address: ":5000",
		Log:     new(Logger),
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithClientName(name string) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.Name = name
	}
}

func WithClientAddress(addr string) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.Address = addr
	}
}

func WithClientLogger(log ILogger) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.Log = log
	}
}

func WithStreamClientInterceptors(streamClientInterceptors ...grpc.StreamClientInterceptor) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.StreamClientInterceptors = streamClientInterceptors
	}
}

func WithUnaryClientInterceptors(unaryClientInterceptors ...grpc.UnaryClientInterceptor) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.UnaryClientInterceptors = unaryClientInterceptors
	}
}

func WithClientDiscovery(discovery resolver.Discovery) ClientOptionFunc {
	return func(o *ClientOptions) {
		o.discovery = discovery
	}
}
