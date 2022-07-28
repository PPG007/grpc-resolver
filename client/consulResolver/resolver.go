package consulresolver

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

type consulResolver struct {
	cc          resolver.ClientConn
	serviceName string
	client      *api.Client
	lastIndex   uint64
}

func (c *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	services, _, err := c.client.Health().Service(c.serviceName, "", true, &api.QueryOptions{})
	if err != nil {
		return
	}
	c.UpdateAddresses(services)
}

func (*consulResolver) Close() {
	log.Println("ConsulResolver will be closed")
}

func (c *consulResolver) Watch() {
	c.lastIndex = 0
	for {
		services, meta, err := c.client.Health().Service(c.serviceName, "", true, &api.QueryOptions{
			WaitIndex: c.lastIndex,
		})
		if err != nil {
			return
		}
		c.lastIndex = meta.LastIndex
		c.UpdateAddresses(services)
	}
}

func (c *consulResolver) UpdateAddresses(services []*api.ServiceEntry) {

	addresses := []resolver.Address{}
	for _, service := range services {
		addresses = append(addresses, resolver.Address{
			Addr: fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port),
		})
	}
	c.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
}
