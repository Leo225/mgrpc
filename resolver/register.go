package resolver

import (
	"context"
	"path"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type ServiceInfo struct {
	Name     string
	Address  string
	Tags     []string
	Interval time.Duration
}

func NewServiceInfo() *ServiceInfo {
	return &ServiceInfo{
		Tags:     []string{"grpc"},
		Interval: 10,
	}
}

type Registry interface {
	Register(info *ServiceInfo) (err error)
	Close() (err error)
}

type EtcdRegistry struct {
	etcdClient *clientv3.Client
	domain     string
	leaseID    clientv3.LeaseID
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewEtcdRegistry(ctx context.Context, client *clientv3.Client, domain string) *EtcdRegistry {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &EtcdRegistry{
		etcdClient: client,
		domain:     domain,
		ctx:        ctx,
		cancel:     cancelFunc,
	}
}

func (r *EtcdRegistry) Register(info *ServiceInfo) (err error) {
	var resp *clientv3.LeaseGrantResponse
	resp, err = r.etcdClient.Grant(r.ctx, int64(info.Interval))
	if err != nil {
		return
	}

	r.leaseID = resp.ID
	target := path.Join(r.domain, info.Name)
	var em endpoints.Manager
	em, err = endpoints.NewManager(r.etcdClient, target)
	if err != nil {
		return
	}

	key := path.Join(target, info.Address)
	endpoint := endpoints.Endpoint{
		Addr: info.Address,
		Metadata: map[string]string{
			"name": info.Name,
			"tags": strings.Join(info.Tags, ","),
		},
	}

	err = em.AddEndpoint(r.ctx, key, endpoint, clientv3.WithLease(r.leaseID))
	if err != nil {
		return
	}

	var keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	keepAliveChan, err = r.etcdClient.KeepAlive(r.ctx, r.leaseID)
	if err != nil {
		return
	}
	go r.watcher(keepAliveChan)
	return
}

func (r *EtcdRegistry) Close() (err error) {
	r.cancel()
	r.etcdClient.Revoke(r.ctx, r.leaseID)
	return
}

func (r *EtcdRegistry) watcher(resp <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case c := <-resp:
			if c == nil {
				return
			}
		case <-r.ctx.Done():
			return
		}
	}
}
