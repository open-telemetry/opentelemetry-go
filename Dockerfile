# Copyright The OpenTelemetry Authors
FROM golang:alpine AS base
COPY . /go/src/github.com/open-telemetry/opentelemetry-go/
WORKDIR /go/src/github.com/open-telemetry/opentelemetry-go/

FROM base AS example-http-server
RUN go install ./example/http/server/server.go
CMD ["/go/bin/server"]

FROM base AS example-http-client
RUN go install ./example/http/client/client.go
CMD ["/go/bin/client"]

FROM base AS example-zipkin-client
RUN go install ./example/zipkin/main.go
CMD ["/go/bin/main"]
