package api

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func CallRandomStrServer(port int, result chan string) {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewRandomStrClient(conn)
	response, err := c.GetRandomStr(context.Background(), &RandomStrReqMessage{Message: "Hello from GatherServer."})
	if err != nil {
		log.Fatalf("Error when calling GetRandomStr: %s", err)
	}

	log.Printf("Response from server %d: %s", port, response.RandomStr)
	result <- response.RandomStr
}

type GatherServer struct {
}

func (s *GatherServer) GatherRandomStr(ctx context.Context, in *RandomStrReqMessage) (*RandomStrRespMessage, error) {
	log.Printf("Received message %s", in.Message)

	randomStrServerPorts := [2]int{7778, 7779}
	n_calls := 16

	results := make(chan string)

	for i := 0; i < n_calls; {
		for _, port := range randomStrServerPorts {
			go CallRandomStrServer(port, results)

			i++
			if i >= n_calls {
				break
			}
		}
	}

	var result string
	for i := 0; i < n_calls; i++ {
		result += <-results
	}

	return &RandomStrRespMessage{
		RandomStr: result,
	}, nil
}
