package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func CallRandomStrServer(address string, result chan string) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewRandomStrClient(conn)
	response, err := c.GetRandomStr(context.Background(), &RandomStrReqMessage{Message: "Hello from GatherServer."})
	if err != nil {
		log.Printf("Error when calling GetRandomStr: %s", err)
    result <- ""
    return
	}

	result <- response.RandomStr
}

type GatherServer struct {
	RandomStrServerAddress string
}

func (s *GatherServer) GatherRandomStr(ctx context.Context, in *RandomStrReqMessage) (*RandomStrRespMessage, error) {
	log.Printf("Received message %s", in.Message)

	nCalls := 16

	results := make(chan string)
	for i := 0; i < nCalls; i++ {
		go CallRandomStrServer(s.RandomStrServerAddress, results)
	}

	var result string
	for i := 0; i < nCalls; i++ {
		result += <-results
	}

	return &RandomStrRespMessage{
		RandomStr: result,
	}, nil
}
