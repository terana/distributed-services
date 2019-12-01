package main

import (
	"api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <gather_server_host:gather_server_port>", os.Args[0])
	}
	gatherServerAddress := os.Args[1]

	var conn *grpc.ClientConn
	var err error

	for ; ; {
		conn, err = grpc.Dial(gatherServerAddress, grpc.WithInsecure())
		if err != nil {
			log.Println("did not connect: %s. Trying again...", err)
		} else {
			break
		}
	}
	defer conn.Close()

	c := api.NewGatherRandomStrClient(conn)
	response, err := c.GatherRandomStr(context.Background(), &api.RandomStrReqMessage{Message: "ping"})
	if err != nil {
		log.Fatalf("Error when calling GetRandomStr: %s", err)
	}

	log.Printf("Response from gather server: %s", response.RandomStr)
}
