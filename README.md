# distributed-services

A demo visualising long tails of request latency distribution in networked computing.

1. With a help of Docker Swarm, spawn a number of services communicating with each other via gRPC. The services do only one thing: send random strings to each other. 
2. With a probability of 1%, add a 1-second delay to responses. 
3. Send a lot of requests to the services (~15k in my test). My p90 latency distribution is ~87 milliseconds.
4. Add a gathering service that makes requests to a number of other services (16 in my tests) and returns concatenated response strings.
5. Send other 15k requests to the gathering service and look at the latency distribution. My p90 grew up 12 times and became 1.049 seconds.

The demo also gathers latencies with Prometheus and visualises them in Grafana charts. 
