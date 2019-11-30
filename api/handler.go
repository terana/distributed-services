package api

import (
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func Delay() {
	/* 1 request from n will be slow */
	n := 100

	if seededRand.Intn(n) == 0 {
		time.Sleep(time.Second)
	}
}

type Server struct {
}

func (s *Server) GetRandomStr(ctx context.Context, in *RandomStrReqMessage) (*RandomStrRespMessage, error) {
	log.Printf("Received message %s", in.Message)

	Delay()

	return &RandomStrRespMessage{
		RandomStr: StringWithCharset(128, charset),
	}, nil
}

func CallRandomStrServer() string {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(":7778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := NewRandomStrClient(conn)
	response, err := c.GetRandomStr(context.Background(), &RandomStrReqMessage{Message: "Hello from GatherServer."})
	if err != nil {
		log.Fatalf("Error when calling GetRandomStr: %s", err)
	}

	log.Printf("Response from server: %s", response.RandomStr)
	return response.RandomStr
}

type GatherServer struct {
}

func (s *GatherServer) GatherRandomStr(ctx context.Context, in *RandomStrReqMessage) (*RandomStrRespMessage, error) {
	log.Printf("Received message %s", in.Message)

	var result string
	for i := 0; i < 16; i++ {
		result += CallRandomStrServer()
	}

	return &RandomStrRespMessage{
		RandomStr: result,
	}, nil
}
