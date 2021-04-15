ARG BASE
FROM $BASE

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=/go/bin:$PATH
ENV GO111MODULE=on

WORKDIR /go/src/proto/
COPY go.* /go/src/proto/
RUN go mod download

COPY films /go/src/proto/films
COPY web /go/src/proto/web
RUN go generate .
