package main

import (
	"api"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

/* Start a gRPC server */
func main() {
	port := os.Args[1]
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.Server{}
	grpcServer := grpc.NewServer()
	api.RegisterRandomStrServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
