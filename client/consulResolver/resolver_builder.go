package consulresolver

import (
	"nameresolver/server/consul"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	CONSUL_RESOLVER_SCHEMA = "consul"
)

func Init() {
	consulConfig, err := consul.NewConsulConfig(0, 0, "")
	if err != nil {
		panic(err.Error())
	}
	consulResolverBuilder := NewConsulResolverBuilder(consulConfig.GetClient())
	resolver.Register(consulResolverBuilder)
}

type consulResolverBuilder struct {
	client *api.Client
}

func NewConsulResolverBuilder(client *api.Client) *consulResolverBuilder {
	return &consulResolverBuilder{
		client: client,
	}
}

func (c *consulResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	service := target.URL.Host
	r := consulResolver{
		cc:          cc,
		client:      c.client,
		serviceName: service,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	go r.Watch()
	return &r, nil
}

func (*consulResolverBuilder) Scheme() string {
	return CONSUL_RESOLVER_SCHEMA
}
