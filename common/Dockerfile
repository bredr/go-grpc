FROM golang:1.16.5-buster@sha256:be0e3a0f3ffa448b0bcbb9019edca692b8278407a44dc138c60e6f12f0218f87 as builder

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=/go/bin:$PATH
ENV GO111MODULE=on
COPY . /go/src/common