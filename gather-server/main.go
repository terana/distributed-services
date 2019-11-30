package main

import (
	"api"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

/* Start a gRPC server */
func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.GatherServer{}
	grpcServer := grpc.NewServer()
	api.RegisterGatherRandomStrServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
