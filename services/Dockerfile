ARG COMMON
ARG BASE
ARG PROTO 

FROM ${PROTO} as proto
FROM ${COMMON} as builder
COPY --from=proto /go/src/proto /go/src/proto

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=/go/bin:$PATH
ENV GO111MODULE=on
WORKDIR /go/src/services/svc
COPY go.* /go/src/services/svc/
RUN go mod download
COPY . /go/src/services/svc

RUN go generate ./...

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app .

FROM $BASE
COPY --from=builder /app .