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

package otlp // import "go.opentelemetry.io/otel/exporters/otlp"

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
)

const (
	// DefaultCollectorPort is the port the Exporter will attempt connect to
	// if no collector port is provided.
	DefaultCollectorPort uint16 = 55680
	// DefaultCollectorHost is the host address the Exporter will attempt
	// connect to if no collector address is provided.
	DefaultCollectorHost string = "localhost"
	// DefaultGRPCServiceConfig is the gRPC service config used if none is
	// provided by the user.
	//
	// For more info on gRPC service configs:
	// https://github.com/grpc/proposal/blob/master/A6-client-retries.md
	//
	// For more info on the RetryableStatusCodes we allow here:
	// https://github.com/open-telemetry/oteps/blob/be2a3fcbaa417ebbf5845cd485d34fdf0ab4a2a4/text/0035-opentelemetry-protocol.md#export-response
	//
	// Note: MaxAttempts > 5 are treated as 5. See
	// https://github.com/grpc/proposal/blob/master/A6-client-retries.md#validation-of-retrypolicy
	// for more details.
	DefaultGRPCServiceConfig = `{
	"methodConfig":[{
		"name":[
			{ "service":"opentelemetry.proto.collector.metrics.v1.MetricsService" },
			{ "service":"opentelemetry.proto.collector.trace.v1.TraceService" }
		],
		"retryPolicy":{
			"MaxAttempts":5,
			"InitialBackoff":"0.3s",
			"MaxBackoff":"5s",
			"BackoffMultiplier":2,
			"RetryableStatusCodes":[
				"UNAVAILABLE",
				"CANCELLED",
				"DEADLINE_EXCEEDED",
				"RESOURCE_EXHAUSTED",
				"ABORTED",
				"OUT_OF_RANGE",
				"UNAVAILABLE",
				"DATA_LOSS"
			]
		}
	}]
}`
)

// ExporterOption are setting options passed to an Exporter on creation.
type ExporterOption func(*config)

type config struct {
	canDialInsecure    bool
	collectorAddr      string
	compressor         string
	reconnectionPeriod time.Duration
	grpcServiceConfig  string
	grpcDialOptions    []grpc.DialOption
	headers            map[string]string
	clientCredentials  credentials.TransportCredentials
	exportKindSelector metricsdk.ExportKindSelector
}

// WithInsecure disables client transport security for the exporter's gRPC connection
// just like grpc.WithInsecure() https://pkg.go.dev/google.golang.org/grpc#WithInsecure
// does. Note, by default, client security is required unless WithInsecure is used.
func WithInsecure() ExporterOption {
	return func(cfg *config) {
		cfg.canDialInsecure = true
	}
}

// WithAddress allows one to set the address that the exporter will
// connect to the collector on. If unset, it will instead try to use
// connect to DefaultCollectorHost:DefaultCollectorPort.
func WithAddress(addr string) ExporterOption {
	return func(cfg *config) {
		cfg.collectorAddr = addr
	}
}

// WithReconnectionPeriod allows one to set the delay between next connection attempt
// after failing to connect with the collector.
func WithReconnectionPeriod(rp time.Duration) ExporterOption {
	return func(cfg *config) {
		cfg.reconnectionPeriod = rp
	}
}

// WithCompressor will set the compressor for the gRPC client to use when sending requests.
// It is the responsibility of the caller to ensure that the compressor set has been registered
// with google.golang.org/grpc/encoding. This can be done by encoding.RegisterCompressor. Some
// compressors auto-register on import, such as gzip, which can be registered by calling
// `import _ "google.golang.org/grpc/encoding/gzip"`
func WithCompressor(compressor string) ExporterOption {
	return func(cfg *config) {
		cfg.compressor = compressor
	}
}

// WithHeaders will send the provided headers with gRPC requests
func WithHeaders(headers map[string]string) ExporterOption {
	return func(cfg *config) {
		cfg.headers = headers
	}
}

// WithTLSCredentials allows the connection to use TLS credentials
// when talking to the server. It takes in grpc.TransportCredentials instead
// of say a Certificate file or a tls.Certificate, because the retrieving
// these credentials can be done in many ways e.g. plain file, in code tls.Config
// or by certificate rotation, so it is up to the caller to decide what to use.
func WithTLSCredentials(creds credentials.TransportCredentials) ExporterOption {
	return func(cfg *config) {
		cfg.clientCredentials = creds
	}
}

// WithGRPCServiceConfig defines the default gRPC service config used.
func WithGRPCServiceConfig(serviceConfig string) ExporterOption {
	return func(cfg *config) {
		cfg.grpcServiceConfig = serviceConfig
	}
}

// WithGRPCDialOption opens support to any grpc.DialOption to be used. If it conflicts
// with some other configuration the GRPC specified via the collector the ones here will
// take preference since they are set last.
func WithGRPCDialOption(opts ...grpc.DialOption) ExporterOption {
	return func(cfg *config) {
		cfg.grpcDialOptions = opts
	}
}

// WithMetricExportKindSelector defines the ExportKindSelector used for selecting
// AggregationTemporality (i.e., Cumulative vs. Delta aggregation).
func WithMetricExportKindSelector(selector metricsdk.ExportKindSelector) ExporterOption {
	return func(cfg *config) {
		cfg.exportKindSelector = selector
	}
}
