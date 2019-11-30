package main

import (
	"api"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	port := os.Args[1]

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := api.NewGatherRandomStrClient(conn)
	response, err := c.GatherRandomStr(context.Background(), &api.RandomStrReqMessage{Message: "ping"})
	if err != nil {
		log.Fatalf("Error when calling GetRandomStr: %s", err)
	}

	log.Printf("Response from server: %s", response.RandomStr)
}
