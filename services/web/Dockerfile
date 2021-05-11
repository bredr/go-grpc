ARG COMMON
ARG BASE
ARG PROTO 


FROM node:16 as www 
WORKDIR /usr/src/app
COPY ./www/package.json ./www/package-lock.json ./www/codegen.yml ./www/tsconfig.json ./
RUN npm ci
COPY ./graph/schema /usr/src/graph/schema
COPY ./www .
RUN npm run generate
RUN npm run build

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

COPY ./server.go ./gqlgen.yml /go/src/services/svc/
COPY ./graph /go/src/services/svc/graph
COPY ./services /go/src/services/svc/services

RUN go generate ./...

COPY --from=www /usr/src/app/build /go/src/services/svc/www/build

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app .

FROM $BASE
COPY --from=builder /app .