package nacos

import (
	"errors"
	"nameresolver/server/manager"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const (
	nacosIP   = "127.0.0.1"
	nacosPort = 8848

	DEFAULT_WEIGHT = 1
)

type nacosConfig struct {
	clientConfig *constant.ClientConfig
	serverConfig *constant.ServerConfig
	client       naming_client.INamingClient
	servicePort  int
}

func NewNacosConfig(servicePort int) (*nacosConfig, error) {
	clientConfig := constant.NewClientConfig()
	serverConfig := constant.NewServerConfig(nacosIP, nacosPort)
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: clientConfig,
		ServerConfigs: []constant.ServerConfig{
			*serverConfig,
		},
	})
	if err != nil {
		return nil, err
	}
	return &nacosConfig{
		client:       client,
		clientConfig: clientConfig,
		serverConfig: serverConfig,
		servicePort:  servicePort,
	}, nil
}

func (n *nacosConfig) RegisterToCenter(serviceName string) error {
	addresses, err := manager.GetSelfIPv4Addresses()
	if err != nil {
		return err
	}
	success, err := n.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          addresses[0],
		Port:        uint64(n.servicePort),
		Weight:      DEFAULT_WEIGHT,
		Enable:      true,
		Healthy:     true,
		ServiceName: serviceName,
	})
	if err != nil {
		return err
	}
	if !success {
		return errors.New("register failed")
	}
	return nil
}

func (n *nacosConfig) GetClient() naming_client.INamingClient {
	return n.client
}
