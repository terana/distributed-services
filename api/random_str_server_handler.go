package api

import (
	"log"
	"math/rand"
	"time"
	"encoding/json"
	"os"

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

/* 1 request from n will be slower by seconds */
func Delay(n int, seconds int) {
	if seededRand.Intn(n) == 0 {
		log.Println("Going to sleep...")
		time.Sleep(time.Second * seconds)
	}
}

type Server struct {
	n int # Every nth request slow
	delay int # seconds
}

func (s *Server) GetRandomStr(ctx context.Context, in *RandomStrReqMessage) (*RandomStrRespMessage, error) {
	log.Printf("Received message %s", in.Message)

	Delay(s.n, s.delay)

	return &RandomStrRespMessage{
		RandomStr: StringWithCharset(128, charset),
	}, nil
}
