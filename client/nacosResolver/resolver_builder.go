package nacosresolver

import (
	"nameresolver/server/nacos"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"google.golang.org/grpc/resolver"
)

const (
	NACOS_RESOLVER_SCHEMA = "nacos"
)

func Init() {
	nc, err := nacos.NewNacosConfig(0)
	if err != nil {
		panic(err)
	}
	builder := &nacosResolverBuilder{
		client: nc.GetClient(),
	}
	resolver.Register(builder)
}

type nacosResolverBuilder struct {
	client naming_client.INamingClient
}

func NewNacosResolverBuilder() *nacosResolverBuilder {
	return nil
}

func (n *nacosResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	r := &nacosResolver{
		cc:          cc,
		target:      target,
		client:      n.client,
		serviceName: target.URL.Host,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	go r.Watch()
	return nil, nil
}

func (*nacosResolverBuilder) Scheme() string {
	return NACOS_RESOLVER_SCHEMA
}
