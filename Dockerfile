FROM golang:1.13-alpine

RUN apk --update add \
    libpcap-dev \
    git \
    gcc \
    libc-dev \
    make \
    bash

RUN mkdir /protoc && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.0/protoc-3.11.0-linux-x86_64.zip && \
    unzip protoc-3.11.0-linux-x86_64.zip -d /protoc && \
    cp /protoc/bin/protoc /bin/protoc

RUN go get github.com/golang/protobuf/protoc-gen-go

WORKDIR /var/app/src/

ADD . .

RUN export GOPATH="${GOPATH}:/var/app/" && make
