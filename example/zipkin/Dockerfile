# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0
FROM golang:1.23-alpine
COPY . /go/src/github.com/open-telemetry/opentelemetry-go/
WORKDIR /go/src/github.com/open-telemetry/opentelemetry-go/example/zipkin/
RUN go install ./main.go
CMD ["/go/bin/main"]
