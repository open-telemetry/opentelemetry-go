// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlpgrpc

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/internal/otlpconfig"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Option applies an option to the gRPC driver.
type Option interface {
	otlpconfig.Option
}

// WithInsecure disables client transport security for the exporter's gRPC connection
// just like grpc.WithInsecure() https://pkg.go.dev/google.golang.org/grpc#WithInsecure
// does. Note, by default, client security is required unless WithInsecure is used.
func WithInsecure() Option {
	return otlpconfig.WithInsecure()
}

// WithTracesInsecure disables client transport security for the traces exporter's gRPC connection
// just like grpc.WithInsecure() https://pkg.go.dev/google.golang.org/grpc#WithInsecure
// does. Note, by default, client security is required unless WithInsecure is used.
func WithTracesInsecure() Option {
	return otlpconfig.WithInsecureTraces()
}

// WithInsecureMetrics disables client transport security for the metrics exporter's gRPC connection
// just like grpc.WithInsecure() https://pkg.go.dev/google.golang.org/grpc#WithInsecure
// does. Note, by default, client security is required unless WithInsecure is used.
func WithInsecureMetrics() Option {
	return otlpconfig.WithInsecureMetrics()
}

// WithEndpoint allows one to set the endpoint that the exporter will
// connect to the collector on. If unset, it will instead try to use
// connect to DefaultCollectorHost:DefaultCollectorPort.
func WithEndpoint(endpoint string) Option {
	return otlpconfig.WithEndpoint(endpoint)
}

// WithTracesEndpoint allows one to set the traces endpoint that the exporter will
// connect to the collector on. If unset, it will instead try to use
// connect to DefaultCollectorHost:DefaultCollectorPort.
func WithTracesEndpoint(endpoint string) Option {
	return otlpconfig.WithTracesEndpoint(endpoint)
}

// WithMetricsEndpoint allows one to set the metrics endpoint that the exporter will
// connect to the collector on. If unset, it will instead try to use
// connect to DefaultCollectorHost:DefaultCollectorPort.
func WithMetricsEndpoint(endpoint string) Option {
	return otlpconfig.WithMetricsEndpoint(endpoint)
}

// WithReconnectionPeriod allows one to set the delay between next connection attempt
// after failing to connect with the collector.
func WithReconnectionPeriod(rp time.Duration) Option {
	return otlpconfig.WithReconnectionPeriod(rp)
}

// WithCompressor will set the compressor for the gRPC client to use when sending requests.
// It is the responsibility of the caller to ensure that the compressor set has been registered
// with google.golang.org/grpc/encoding. This can be done by encoding.RegisterCompressor. Some
// compressors auto-register on import, such as gzip, which can be registered by calling
// `import _ "google.golang.org/grpc/encoding/gzip"`
func WithCompressor(compressor string) Option {
	switch compressor {
	case "gzip":
		return otlpconfig.WithCompression(otlp.GzipCompression)
	}

	otel.Handle(fmt.Errorf("invalid compression type: '%s', using no compression.", compressor))

	return otlpconfig.WithCompression(otlp.NoCompression)
}

// WithTracesCompression will set the compressor for the gRPC client to use when sending traces requests.
// It is the responsibility of the caller to ensure that the compressor set has been registered
// with google.golang.org/grpc/encoding. This can be done by encoding.RegisterCompressor. Some
// compressors auto-register on import, such as gzip, which can be registered by calling
// `import _ "google.golang.org/grpc/encoding/gzip"`
func WithTracesCompression(compressor string) Option {
	switch compressor {
	case "gzip":
		return otlpconfig.WithTracesCompression(otlp.GzipCompression)
	}

	otel.Handle(fmt.Errorf("invalid compression type: '%s', using no compression.", compressor))

	return otlpconfig.WithTracesCompression(otlp.NoCompression)
}

// WithMetricsCompression will set the compressor for the gRPC client to use when sending metrics requests.
// It is the responsibility of the caller to ensure that the compressor set has been registered
// with google.golang.org/grpc/encoding. This can be done by encoding.RegisterCompressor. Some
// compressors auto-register on import, such as gzip, which can be registered by calling
// `import _ "google.golang.org/grpc/encoding/gzip"`
func WithMetricsCompression(compressor string) Option {
	switch compressor {
	case "gzip":
		return otlpconfig.WithMetricsCompression(otlp.GzipCompression)
	}

	otel.Handle(fmt.Errorf("invalid compression type: '%s', using no compression.", compressor))

	return otlpconfig.WithMetricsCompression(otlp.NoCompression)
}

// WithHeaders will send the provided headers with gRPC requests
func WithHeaders(headers map[string]string) Option {
	return otlpconfig.WithHeaders(headers)
}

// WithTracesHeaders will send the provided headers with gRPC traces requests
func WithTracesHeaders(headers map[string]string) Option {
	return otlpconfig.WithTracesHeaders(headers)
}

// WithMetricsHeaders will send the provided headers with gRPC metrics requests
func WithMetricsHeaders(headers map[string]string) Option {
	return otlpconfig.WithMetricsHeaders(headers)
}

// WithTLSCredentials allows the connection to use TLS credentials
// when talking to the server. It takes in grpc.TransportCredentials instead
// of say a Certificate file or a tls.Certificate, because the retrieving
// these credentials can be done in many ways e.g. plain file, in code tls.Config
// or by certificate rotation, so it is up to the caller to decide what to use.
func WithTLSCredentials(creds credentials.TransportCredentials) Option {
	return otlpconfig.WithTLSCredentials(creds)
}

// WithTracesTLSCredentials allows the connection to use TLS credentials
// when talking to the traces server. It takes in grpc.TransportCredentials instead
// of say a Certificate file or a tls.Certificate, because the retrieving
// these credentials can be done in many ways e.g. plain file, in code tls.Config
// or by certificate rotation, so it is up to the caller to decide what to use.
func WithTracesTLSCredentials(creds credentials.TransportCredentials) Option {
	return otlpconfig.WithTracesTLSCredentials(creds)
}

// WithMetricsTLSCredentials allows the connection to use TLS credentials
// when talking to the metrics server. It takes in grpc.TransportCredentials instead
// of say a Certificate file or a tls.Certificate, because the retrieving
// these credentials can be done in many ways e.g. plain file, in code tls.Config
// or by certificate rotation, so it is up to the caller to decide what to use.
func WithMetricsTLSCredentials(creds credentials.TransportCredentials) Option {
	return otlpconfig.WithMetricsTLSCredentials(creds)
}

// WithServiceConfig defines the default gRPC service config used.
func WithServiceConfig(serviceConfig string) Option {
	return otlpconfig.WithServiceConfig(serviceConfig)
}

// WithDialOption opens support to any grpc.DialOption to be used. If it conflicts
// with some other configuration the GRPC specified via the collector the ones here will
// take preference since they are set last.
func WithDialOption(opts ...grpc.DialOption) Option {
	return otlpconfig.WithDialOption(opts...)
}
