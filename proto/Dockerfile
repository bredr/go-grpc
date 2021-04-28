ARG BASE
FROM $BASE

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=/go/bin:$PATH
ENV GO111MODULE=on

WORKDIR /go/src/proto/
COPY go.* /go/src/proto/
RUN go mod download

COPY . /go/src/proto/

RUN go generate .
