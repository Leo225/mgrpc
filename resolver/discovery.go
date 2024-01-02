package resolver

import (
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	gresolver "google.golang.org/grpc/resolver"
)

type Discovery interface {
	Discover(service string) (builder gresolver.Builder, err error)
	Close() (err error)
	Address() (addr string)
}

type EtcdDiscovery struct {
	etcdClient *clientv3.Client
	builder    gresolver.Builder
	domain     string
	service    string
}

func NewEtcdDiscovery(client *clientv3.Client, domain string) *EtcdDiscovery {
	builder, err := eresolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}

	return &EtcdDiscovery{
		etcdClient: client,
		builder:    builder,
		domain:     domain,
	}
}

func (d *EtcdDiscovery) Discover(service string) (builder gresolver.Builder, err error) {
	d.service = service
	builder = d.builder
	return
}

func (d *EtcdDiscovery) Close() (err error) {
	return
}
func (d *EtcdDiscovery) Address() (addr string) {
	addr = fmt.Sprintf("%s:///%s/%s", d.builder.Scheme(), d.domain, d.service)
	return
}
