package main

import (
	"api"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

/* Start a gRPC server */
func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <port to listen on> <server_host:server_port>", os.Args[0])
	}
	port := os.Args[1]
	serverAddress := os.Args[2]

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpcDurations := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "service_latency_seconds",
			Help:       "Service latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"service"},
	)
	prometheus.MustRegister(rpcDurations)

	s := api.GatherServer{
		RandomStrServerAddress: serverAddress,
		PrometheusSummaryVec:   rpcDurations,
	}
	grpcServer := grpc.NewServer()
	api.RegisterGatherRandomStrServer(grpcServer, &s)

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				// Opt into OpenMetrics to support exemplars.
				EnableOpenMetrics: true,
			},
		))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
