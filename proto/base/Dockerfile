FROM golang:1.16 as builder
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=/go/bin:$PATH
ENV GO111MODULE=on

RUN apt-get update && \
  apt-get install unzip

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip && \
  unzip -o protoc-3.9.1-linux-x86_64.zip -d /usr/local bin/protoc && \
  unzip -o protoc-3.9.1-linux-x86_64.zip -d /usr/local include/* && \
  rm -rf protoc-3.9.1-linux-x86_64.zip

RUN go get google.golang.org/protobuf/cmd/protoc-gen-go \
  google.golang.org/grpc/cmd/protoc-gen-go-grpc