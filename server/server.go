package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"nameresolver/proto"
	"nameresolver/server/consul"
	"nameresolver/server/manager"
	"nameresolver/server/nacos"
	"nameresolver/server/service"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	port        int
	center      string
	gatewayPort int
)

const (
	host        = "127.0.0.1"
	serviceName = "HelloService"

	CENTER_CONSUL = "consul"
	CENTER_NACOS  = "nacos"
)

func init() {
	flag.IntVar(&port, "port", 0, "server port")
	flag.IntVar(&gatewayPort, "gatewayPort", 0, "grpc gateway port")
	flag.StringVar(&center, "center", "", "nacos or consul")
	flag.Parse()
}

func register(center string) {
	var serviceManager manager.ServerManager
	switch center {
	case CENTER_CONSUL:
		cc, err := consul.NewConsulConfig(port, gatewayPort, "/ping")
		if err != nil {
			panic(err.Error())
		}
		serviceManager = cc
	case CENTER_NACOS:
		nc, err := nacos.NewNacosConfig(port)
		if err != nil {
			panic(err.Error())
		}
		serviceManager = nc
	default:
		return
	}
	serviceManager.RegisterToCenter(serviceName)
}

func startGrpcGateway(ctx context.Context, serverPort, startPort int) {
	mux := runtime.NewServeMux()
	err := proto.RegisterNameResolverServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("127.0.0.1:%d", serverPort), []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		panic(err.Error())
	}
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", startPort), mux)
}

func main() {
	log.Printf("Server will start at port: %d\n", port)
	log.Printf("Grpc gateway will start at port: %d\n", gatewayPort)
	grpcServer := grpc.NewServer()
	proto.RegisterNameResolverServiceServer(grpcServer, service.NameResolverService{})
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	reflection.Register(grpcServer)
	go startGrpcGateway(context.Background(), port, gatewayPort)
	register(center)
	err = grpcServer.Serve(l)
	if err != nil {
		panic(err)
	}
}
