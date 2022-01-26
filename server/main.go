package main

import (
	"api"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"encoding/json"
)

type Config struct {
	every_nth_request_slow int
	seconds_delay int
}

/* Start a gRPC server */
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <port to listen on>", os.Args[0])
	}
	port := os.Args[1]

	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatalf("failed to get configuration: %v", err)
	}


	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.Server{
		n: config.every_nth_request_slow,
		delay: config.seconds_delay,
	}
	grpcServer := grpc.NewServer()
	api.RegisterRandomStrServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
