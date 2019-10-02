package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"path"
	"sync/atomic"
	"time"
)

const FlashServiceTimeDuration = 10 * time.Second
const ServiceTTl = 10

type EtcdRegistry struct {
	client  *clientv3.Client
	options *OptionsRegistry
	cache   atomic.Value
}

type EtcdServiceCache struct {
	srvs map[string]*Service
}

func (*EtcdRegistry) Name() string {
	return "etcd"
}

func (registry *EtcdRegistry) Init(ctx context.Context, options ...OptionRegistry) error {
	var err error

	registry.options = &OptionsRegistry{}
	for _, option := range options {
		option(registry.options)
	}

	if len(registry.options.Endpoints) == 0 {
		return fmt.Errorf("there is no endpoints given for initing etcd client")
	}

	registry.client, err = clientv3.New(clientv3.Config{
		Endpoints:   registry.options.Endpoints,
		DialTimeout: registry.options.Timeout,
	})
	if err != nil {
		return fmt.Errorf("fail to create the etcd client, error: %v\n", err)
	}

	registry.cache.Store(&EtcdServiceCache{srvs:make(map[string]*Service)})
	go registry.flushServiceCache()

	return nil
}

func (registry *EtcdRegistry) Register(ctx context.Context, srv *Service) {
	resp, err := registry.client.Grant(ctx, ServiceTTl)
	if err != nil {
		log.Printf("an error occurs during creating etcd lease, error: %v\n", err)
		return
	}
	leaseId := resp.ID

	for _, node := range srv.Nodes {
		key := path.Join("/",registry.options.RootPath, srv.Name, fmt.Sprintf("%s:%d", node.Ip, node.Port))
		tempService := Service{Name: srv.Name, Version: srv.Version, Nodes: []*Node{node}}
		valueBytes, err := json.Marshal(tempService)
		if err != nil {
			log.Printf("an error occurs during marshal the service meta data to json, error: %v\n", err)
			return
		}

		_, err = registry.client.Put(ctx, key, string(valueBytes), clientv3.WithLease(leaseId))
		if err != nil {
			log.Printf("an error occurs during during putting value into etcd, error: %v\n", err)
			return
		}

	}

	leaseChan, err := registry.client.KeepAlive(ctx, leaseId)
	if err != nil {
		log.Printf("an error occurs during keeping alive a lease, error: %v\n", err)
		return
	}

	go registry.keepAlive(ctx, leaseChan, srv)
}

func (registry *EtcdRegistry) Discover(ctx context.Context, name string) (*Service, error) {
	srv, ok := registry.discoverFromCache(name)
	if ok {
		return srv, nil
	}

	keyPrefix := path.Join("/", registry.options.RootPath, name)
	response, err := registry.client.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	srv = &Service{Name: name}
	for _, kv := range response.Kvs {
		tempSrv := new(Service)
		err = json.Unmarshal(kv.Value, tempSrv)
		if err != nil {
			fmt.Printf("an error occurs during parse the service meta data from json, error: %v\n", err)
			return nil, err
		}
		node := &Node{Ip: tempSrv.Nodes[0].Ip, Port: tempSrv.Nodes[0].Port, Weight: tempSrv.Nodes[0].Weight}
		srv.Nodes = append(srv.Nodes, node)
	}

	srvsCache := registry.cache.Load().(*EtcdServiceCache)
	srvsCache.srvs[name] = srv
	registry.cache.Store(srvsCache)

	return srv, nil
}

func (registry *EtcdRegistry) Withdraw(ctx context.Context, service *Service) {

}

func (registry *EtcdRegistry) keepAlive(ctx context.Context, leaseChan <-chan *clientv3.LeaseKeepAliveResponse, srv *Service) {
	for {
		select {
		case response := <-leaseChan:
			if response == nil {
				go registry.Register(ctx, srv)
				return
			}
			fmt.Printf("Success to keep alive a lease,resp:%+v,service:%+v,ttl:%v\n", response, srv, response.TTL)
		}
	}
}

func (registry *EtcdRegistry) discoverFromCache(name string) (service *Service, ok bool) {
	serviceCache := registry.cache.Load().(*EtcdServiceCache)
	service, ok = serviceCache.srvs[name]
	return
}

func (registry *EtcdRegistry) flushServiceCache() {
	for {
		serviceCache := registry.cache.Load().(*EtcdServiceCache)
		for _, oldService := range serviceCache.srvs {
			response, err := registry.client.Get(context.TODO(), path.Join(registry.options.RootPath, oldService.Name), clientv3.WithPrefix())
			newService := &Service{Name: oldService.Name}
			for _, kv := range response.Kvs {
				tmpService := &Service{Name: oldService.Name}
				err = json.Unmarshal(kv.Value, tmpService)
				if err != nil {
					continue
				}
				newService.Nodes = append(newService.Nodes, tmpService.Nodes...)
			}
			serviceCache.srvs[newService.Name] = newService
		}

		registry.cache.Store(serviceCache)

		time.Sleep(FlashServiceTimeDuration)
	}
}
