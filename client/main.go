package main

import (
	"api"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
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
