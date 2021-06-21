FROM alpine:3.14.0@sha256:8d99168167baa6a6a0d7851b9684625df9c1455116a9601835c2127df2aaa2f5 as builder

ENV GRPC_HEALTH_PROBE_VERSION v0.3.1

RUN set -ex \
  && apk add --no-cache curl \
  && curl -fSL https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 -o grpc_health_probe \
  && chmod +x grpc_health_probe

FROM gcr.io/distroless/base@sha256:6ec6da1888b18dd971802c2a58a76a7702902b4c9c1be28f38e75e871cedc2df

COPY --from=builder /grpc_health_probe /bin/grpc_health_probe

# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single
CMD ["./app"]