// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package selfobservability // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/selfobservability"

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"

	"github.com/stretchr/testify/assert"

	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
)

func TestServerAddrAttrs(t *testing.T) {
	testcases := []struct {
		name   string
		target string
		want   []attribute.KeyValue
	}{
		{
			name:   "Unix socket",
			target: "unix:///tmp/grpc.sock",
			want:   []attribute.KeyValue{semconv.ServerAddress("/tmp/grpc.sock")},
		},
		{
			name:   "DNS with port",
			target: "dns:///localhost:8080",
			want:   []attribute.KeyValue{semconv.ServerAddress("localhost"), semconv.ServerPort(8080)},
		},
		{
			name:   "Dns with endpoint host:port",
			target: "dns://8.8.8.8/example.com:4",
			want:   []attribute.KeyValue{semconv.ServerAddress("example.com"), semconv.ServerPort(4)},
		},
		{
			name:   "Simple host port",
			target: "localhost:10001",
			want:   []attribute.KeyValue{semconv.ServerAddress("localhost"), semconv.ServerPort(10001)},
		},
		{
			name:   "Host without port",
			target: "example.com",
			want:   []attribute.KeyValue{semconv.ServerAddress("example.com")},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			attrs := ServerAddrAttrs(tc.target)
			assert.Equal(t, tc.want, attrs)
		})
	}
}
