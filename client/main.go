package main

import (
	"api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
			log.Printf("did not connect: %s. Trying again...", err)
		} else {
			break
		}
	}
	defer conn.Close()

	c := api.NewGatherRandomStrClient(conn)

	rpcDurations := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "gather_server_latency_seconds",
			Help:       "Gathered requests latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
	prometheus.MustRegister(rpcDurations)

	go func() {
		i := 0
		for {
			i += 1
			if i%100 == 0 {
				log.Println("request number ", i)
			}
			start := time.Now()
			_, err := c.GatherRandomStr(context.Background(), &api.RandomStrReqMessage{Message: "from Client"})
			latency := time.Since(start)
			rpcDurations.WithLabelValues("normal").Observe(latency.Seconds())
			if err != nil {
				log.Printf("Error when calling GatherRandomStr: %s", err)
			}
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
