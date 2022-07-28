package consul

import (
	"fmt"
	"nameresolver/server/manager"

	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	consulHost = "127.0.0.1"
	consulPort = 8500
)

type consulConfig struct {
	serverHost          string
	serverPort          int
	servicePort         int
	client              *api.Client
	healthCheckEndpoint string
}

func NewConsulConfig(healthCheckEndpoint string, servicePort int) (*consulConfig, error) {
	config := api.DefaultConfig()
	result := &consulConfig{
		serverHost:          consulHost,
		serverPort:          consulPort,
		servicePort:         servicePort,
		healthCheckEndpoint: healthCheckEndpoint,
	}
	config.Address = fmt.Sprintf("%s:%d", result.serverHost, result.serverPort)
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	result.client = c
	return result, nil
}

func (c *consulConfig) RegisterToCenter(serviceName string) error {
	addresses, err := manager.GetSelfIPv4Addresses()
	if err != nil {
		return err
	}
	registration := api.AgentServiceRegistration{
		Name:    serviceName,
		ID:      fmt.Sprintf("%s_%s", serviceName, primitive.NewObjectID().Hex()),
		Port:    c.servicePort,
		Address: addresses[0],
		Check: &api.AgentServiceCheck{
			HTTP:                           c.healthCheckEndpoint,
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "20s",
		},
	}
	return c.client.Agent().ServiceRegister(&registration)
}

func (c *consulConfig) GetClient() *api.Client {
	return c.client
}
