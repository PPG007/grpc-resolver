package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"nameresolver/proto"
	"nameresolver/server/consul"
	"nameresolver/server/nacos"
	"nameresolver/server/service"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	port                int
	center              string
	healthCheckEndpoint string
)

const (
	gatewayPort = 80
	host        = "127.0.0.1"
	serviceName = "HelloService"

	CENTER_CONSUL = "consul"
	CENTER_NACOS  = "nacos"
)

func init() {
	flag.IntVar(&port, "port", 0, "server port")
	flag.StringVar(&center, "center", "", "nacos or consul")
	flag.Parse()
	healthCheckEndpoint = fmt.Sprintf("http://%s:%d/ping", host, gatewayPort)
}

func registerToNacos() {
	nc, err := nacos.NewNacosConfig(port)
	if err != nil {
		panic(err.Error())
	}
	err = nc.RegisterToCenter(serviceName)
	if err != nil {
		panic(err.Error())
	}
}

func registerToConsul() {
	cc, err := consul.NewConsulConfig(healthCheckEndpoint, port)
	if err != nil {
		panic(err.Error())
	}
	err = cc.RegisterToCenter(serviceName)
	if err != nil {
		panic(err.Error())
	}
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
	crt, err := credentials.NewServerTLSFromFile("/home/user/playground/grpc-resolver/cert/localhost.crt", "/home/user/playground/grpc-resolver/cert/prikey.pem")
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(crt),
	)
	proto.RegisterNameResolverServiceServer(grpcServer, service.NameResolverService{})
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
	reflection.Register(grpcServer)
	go startGrpcGateway(context.Background(), port, gatewayPort)
	if center == CENTER_CONSUL {
		registerToConsul()
	}
	if center == CENTER_NACOS {
		registerToNacos()
	}
	err = grpcServer.Serve(l)
	if err != nil {
		panic(err)
	}
}
