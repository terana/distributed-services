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
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <port to listen on> <server_host:server_port>", os.Args[0])
	}
	port := os.Args[1]
	serverAddress := os.Args[2]

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.GatherServer{RandomStrServerAddress: serverAddress}
	grpcServer := grpc.NewServer()
	api.RegisterGatherRandomStrServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
