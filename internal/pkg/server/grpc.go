package server

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port) // Listen on any IP address on the specified port
	RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(
			func(fullMethodName string, err error) bool {
				if fullMethodName == "grpc.health.v1.Health.Check" {
					return false
				}
				return true
			},
		),
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
		), grpc_middleware.WithStreamServerChain(
			grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
		),
	)
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting: gRPC Listener")
	log.Fatal(grpcServer.Serve(listen))
}
