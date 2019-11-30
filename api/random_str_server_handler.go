package api

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

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
