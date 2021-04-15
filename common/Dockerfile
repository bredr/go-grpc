FROM golang:1.16 as builder

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=/go/bin:$PATH
ENV GO111MODULE=on
COPY . /go/src/common