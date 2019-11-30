API_OUT := "api/api.pb.go"

SERVER_PKG_BUILD := "server"
CLIENT_PKG_BUILD := "client"

.PHONY: all api server client

all: server client

api/random_str_api.pb.go: api/random_str_api.proto
	@protoc -I api/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:api \
		api/random_str_api.proto

api: api/random_str_api.pb.go

dep:
	@go get -v -d ./...

server: dep api
	@go build -i -v -o server $(SERVER_PKG_BUILD)

client: dep api
	@go build -i -v -o client $(CLIENT_PKG_BUILD)

clean:
	@rm server client api/random_str_api.pb.go
