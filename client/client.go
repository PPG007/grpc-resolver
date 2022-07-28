package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	consulresolver "nameresolver/client/consulResolver"
	nacosresolver "nameresolver/client/nacosResolver"
	"nameresolver/proto"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	cc                        *grpc.ClientConn
	err                       error
	nameResolverServiceClient proto.NameResolverServiceClient
	center                    string
)

const (
	CENTER_CONSUL = "consul"
	CENTER_NACOS  = "nacos"

	serviceName = "HelloService"
)

func init() {
	flag.StringVar(&center, "center", "", "nacos or consul")
	flag.Parse()
	target := ""
	switch center {
	case CENTER_NACOS:
		target = fmt.Sprintf("%s://%s", nacosresolver.NACOS_RESOLVER_SCHEMA, serviceName)
		nacosresolver.Init()
	case CENTER_CONSUL:
		target = fmt.Sprintf("%s://%s", consulresolver.CONSUL_RESOLVER_SCHEMA, serviceName)
		consulresolver.Init()
	default:
		panic("Unsupported center")
	}
	cc, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	nameResolverServiceClient = proto.NewNameResolverServiceClient(cc)
	log.Println("Client init over")
}

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		resp, err := nameResolverServiceClient.Hello(r.Context(), &proto.EmptyRequest{})
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(resp.String()))
	})
	err = http.ListenAndServe("0.0.0.0:9090", nil)
	if err != nil {
		panic(err.Error())
	}
}
