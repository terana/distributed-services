package api

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func CallRandomStrServer(address string, summaryVec *prometheus.SummaryVec, result chan string) {
	var conn *grpc.ClientConn
	var err error

	for ; ; {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("did not connect: %s. Trying again...", err)
		} else {
			break
		}
	}
	defer conn.Close()

	c := NewRandomStrClient(conn)
	start := time.Now()
	response, err := c.GetRandomStr(context.Background(), &RandomStrReqMessage{Message: "from GatherServer."})
	latency := time.Since(start)
	summaryVec.WithLabelValues("normal").Observe(latency.Seconds())
	if err != nil {
		log.Printf("Error when calling GetRandomStr: %s", err)
		result <- ""
		return
	}

	result <- response.RandomStr
}

type GatherServer struct {
	RandomStrServerAddress string
	PrometheusSummaryVec   *prometheus.SummaryVec
}

func (s *GatherServer) GatherRandomStr(ctx context.Context, in *RandomStrReqMessage) (*RandomStrRespMessage, error) {
	log.Printf("Received message %s", in.Message)

	nCalls := 16

	results := make(chan string)
	for i := 0; i < nCalls; i++ {
		go CallRandomStrServer(s.RandomStrServerAddress, s.PrometheusSummaryVec, results)
	}

	var result string
	for i := 0; i < nCalls; i++ {
		result += <-results
	}

	return &RandomStrRespMessage{
		RandomStr: result,
	}, nil
}
