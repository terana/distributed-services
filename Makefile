SERVER_PKG_BUILD := "server"
CLIENT_PKG_BUILD := "client"
GATHER_SERVER_PKG_BUILD := "gather-server"

.PHONY: all api server client

all: server client gather-server

api/random_str_api.pb.go: api/random_str_api.proto
	protoc -I api/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:api \
		api/random_str_api.proto

api: api/random_str_api.pb.go

dep:
	go get -d ./...

server: dep api
	go build -i -o server/server $(SERVER_PKG_BUILD)

client: dep api
	go build -i -o client/client $(CLIENT_PKG_BUILD)

gather-server: dep api
	go build -i -o gather-server/gather-server $(GATHER_SERVER_PKG_BUILD)

build-docker:
	docker build . -t distributed

run-docker: build-docker
	docker run -it --rm distributed bash

start-services: build-docker
	docker swarm init ||:
	docker stack deploy --compose-file docker-compose.yml distributed
	docker service ls
	# docker service logs distributed_client

stop-services:
	docker service rm distributed_client ||:
	docker service rm distributed_server ||:
	docker service rm distributed_gather-server

clean:
	rm server/server client/client gather-server/gather-server api/random_str_api.pb.go
