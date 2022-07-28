package nacosresolver

import (
	"fmt"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc/resolver"
)

type nacosResolver struct {
	cc          resolver.ClientConn
	target      resolver.Target
	client      naming_client.INamingClient
	serviceName string
}

func (n *nacosResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	instances, err := n.client.SelectAllInstances(vo.SelectAllInstancesParam{})
	if err != nil {
		return
	}
	n.UpdateAddress(instances, nil)
}

func (*nacosResolver) Close() {
	log.Println("Nacos nameresolver will be closed")
}

func (n *nacosResolver) Watch() {
	err := n.client.Subscribe(&vo.SubscribeParam{
		ServiceName: n.serviceName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			if err != nil {
				return
			}
			n.UpdateAddress(nil, services)
		},
	})
	if err != nil {
		panic(err.Error())
	}
}

func (n *nacosResolver) UpdateAddress(instances []model.Instance, services []model.SubscribeService) {
	addresses := []resolver.Address{}
	for _, instance := range instances {
		addresses = append(addresses, resolver.Address{
			Addr: fmt.Sprintf("%s:%d", instance.Ip, instance.Port),
		})
	}
	for _, service := range services {
		addresses = append(addresses, resolver.Address{
			Addr: fmt.Sprintf("%s:%d", service.Ip, service.Port),
		})
	}
	n.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
}
