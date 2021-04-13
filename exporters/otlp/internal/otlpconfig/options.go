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

package otlpconfig // import "go.opentelemetry.io/otel/exporters/otlp/internal/otlpconfig"

import (
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/exporters/otlp"
)

const (
	// DefaultMaxAttempts describes how many times the driver
	// should retry the sending of the payload in case of a
	// retryable error.
	DefaultMaxAttempts int = 5
	// DefaultTracesPath is a default URL path for endpoint that
	// receives spans.
	DefaultTracesPath string = "/v1/traces"
	// DefaultMetricsPath is a default URL path for endpoint that
	// receives metrics.
	DefaultMetricsPath string = "/v1/metrics"
	// DefaultBackoff is a default base backoff time used in the
	// exponential backoff strategy.
	DefaultBackoff time.Duration = 300 * time.Millisecond
	// DefaultTimeout is a default max waiting time for the backend to process
	// each span or metrics batch.
	DefaultTimeout time.Duration = 10 * time.Second
	// DefaultServiceConfig is the gRPC service config used if none is
	// provided by the user.
	DefaultServiceConfig = `{
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

type (
	SignalConfig struct {
		Endpoint    string
		Insecure    bool
		TLSCfg      *tls.Config
		Headers     map[string]string
		Compression otlp.Compression
		Timeout     time.Duration
		URLPath     string

		// gRPC configurations
		GrpcCredentials credentials.TransportCredentials
	}

	Config struct {
		// Signal specific configurations
		Metrics SignalConfig
		Traces  SignalConfig

		// General configurations
		MaxAttempts int
		Backoff     time.Duration

		// HTTP configuration
		Marshaler otlp.Marshaler

		// gRPC configurations
		ReconnectionPeriod time.Duration
		ServiceConfig      string
		DialOptions        []grpc.DialOption
	}
)

func NewDefaultConfig() Config {
	c := Config{
		Traces: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", otlp.DefaultCollectorHost, otlp.DefaultCollectorPort),
			URLPath:     DefaultTracesPath,
			Compression: otlp.NoCompression,
			Timeout:     DefaultTimeout,
		},
		Metrics: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", otlp.DefaultCollectorHost, otlp.DefaultCollectorPort),
			URLPath:     DefaultMetricsPath,
			Compression: otlp.NoCompression,
			Timeout:     DefaultTimeout,
		},
		MaxAttempts:   DefaultMaxAttempts,
		Backoff:       DefaultBackoff,
		ServiceConfig: DefaultServiceConfig,
	}

	return c
}

type (
	// GenericOption applies an option to the HTTP or gRPC driver.
	GenericOption interface {
		ApplyHTTPOption(*Config)
		ApplyGRPCOption(*Config)

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}

	// HTTPOption applies an option to the HTTP driver.
	HTTPOption interface {
		ApplyHTTPOption(*Config)

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}

	// GRPCOption applies an option to the gRPC driver.
	GRPCOption interface {
		ApplyGRPCOption(*Config)

		// A private method to prevent users implementing the
		// interface and so future additions to it will not
		// violate compatibility.
		private()
	}
)

type genericOption struct {
	fn func(*Config)
}

func (g *genericOption) ApplyGRPCOption(cfg *Config) {
	g.fn(cfg)
}

func (g *genericOption) ApplyHTTPOption(cfg *Config) {
	g.fn(cfg)
}

func (genericOption) private() {}

func newGenericOption(fn func(cfg *Config)) GenericOption {
	return &genericOption{fn: fn}
}

type splitOption struct {
	httpFn func(*Config)
	grpcFn func(*Config)
}

func (g *splitOption) ApplyGRPCOption(cfg *Config) {
	g.grpcFn(cfg)
}

func (g *splitOption) ApplyHTTPOption(cfg *Config) {
	g.httpFn(cfg)
}

func (splitOption) private() {}

func newSplitOption(httpFn func(cfg *Config), grpcFn func(cfg *Config)) GenericOption {
	return &splitOption{httpFn: httpFn, grpcFn: grpcFn}
}

type httpOption struct {
	fn func(*Config)
}

func (h *httpOption) ApplyHTTPOption(cfg *Config) {
	h.fn(cfg)
}

func (httpOption) private() {}

func NewHTTPOption(fn func(cfg *Config)) HTTPOption {
	return &httpOption{fn: fn}
}

type grpcOption struct {
	fn func(*Config)
}

func (h *grpcOption) ApplyGRPCOption(cfg *Config) {
	h.fn(cfg)
}

func (grpcOption) private() {}

func NewGRPCOption(fn func(cfg *Config)) GRPCOption {
	return &grpcOption{fn: fn}
}

// WithEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send metrics and spans. If
// unset, it will instead try to use
// DefaultCollectorHost:DefaultCollectorPort. Note that the endpoint
// must not contain any URL path.
func WithEndpoint(endpoint string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Endpoint = endpoint
		cfg.Metrics.Endpoint = endpoint
	})
}

// WithTracesEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send spans. If
// unset, it will instead try to use the Endpoint configuration.
// Note that the endpoint must not contain any URL path.
func WithTracesEndpoint(endpoint string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Endpoint = endpoint
	})
}

// WithMetricsEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send metrics. If
// unset, it will instead try to use the Endpoint configuration.
// Note that the endpoint must not contain any URL path.
func WithMetricsEndpoint(endpoint string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Metrics.Endpoint = endpoint
	})
}

// WithCompression tells the driver to compress the sent data.
func WithCompression(compression otlp.Compression) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Compression = compression
		cfg.Metrics.Compression = compression
	})
}

// WithTracesCompression tells the driver to compress the sent traces data.
func WithTracesCompression(compression otlp.Compression) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Compression = compression
	})
}

// WithMetricsCompression tells the driver to compress the sent metrics data.
func WithMetricsCompression(compression otlp.Compression) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Metrics.Compression = compression
	})
}

// WithTracesURLPath allows one to override the default URL path used
// for sending traces. If unset, DefaultTracesPath will be used.
func WithTracesURLPath(urlPath string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.URLPath = urlPath
	})
}

// WithMetricsURLPath allows one to override the default URL path used
// for sending metrics. If unset, DefaultMetricsPath will be used.
func WithMetricsURLPath(urlPath string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Metrics.URLPath = urlPath
	})
}

// WithMaxAttempts allows one to override how many times the driver
// will try to send the payload in case of retryable errors. If unset,
// DefaultMaxAttempts will be used.
func WithMaxAttempts(maxAttempts int) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.MaxAttempts = maxAttempts
	})
}

// WithBackoff tells the driver to use the duration as a base of the
// exponential backoff strategy. If unset, DefaultBackoff will be
// used.
func WithBackoff(duration time.Duration) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Backoff = duration
	})
}

// WithTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send payloads to the
// collector. Use it if you want to use a custom certificate.
func WithTLSClientConfig(tlsCfg *tls.Config) GenericOption {
	return newSplitOption(func(cfg *Config) {
		cfg.Traces.TLSCfg = tlsCfg
		cfg.Metrics.TLSCfg = tlsCfg
	}, func(cfg *Config) {
		cfg.Traces.GrpcCredentials = credentials.NewTLS(tlsCfg)
		cfg.Metrics.GrpcCredentials = credentials.NewTLS(tlsCfg)
	})
}

// WithTracesTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send traces.
// Use it if you want to use a custom certificate.
func WithTracesTLSClientConfig(tlsCfg *tls.Config) GenericOption {
	return newSplitOption(func(cfg *Config) {
		cfg.Traces.TLSCfg = tlsCfg
	}, func(cfg *Config) {
		cfg.Traces.GrpcCredentials = credentials.NewTLS(tlsCfg)
	})
}

// WithMetricsTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send metrics.
// Use it if you want to use a custom certificate.
func WithMetricsTLSClientConfig(tlsCfg *tls.Config) GenericOption {
	return newSplitOption(func(cfg *Config) {
		cfg.Metrics.TLSCfg = tlsCfg
	}, func(cfg *Config) {
		cfg.Metrics.GrpcCredentials = credentials.NewTLS(tlsCfg)
	})
}

// WithInsecure tells the driver to connect to the collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecure() GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Insecure = true
		cfg.Metrics.Insecure = true
	})
}

// WithInsecureTraces tells the driver to connect to the traces collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecureTraces() GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Insecure = true
	})
}

// WithInsecure tells the driver to connect to the metrics collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecureMetrics() GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Metrics.Insecure = true
	})
}

// WithHeaders allows one to tell the driver to send additional HTTP
// headers with the payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithHeaders(headers map[string]string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Headers = headers
		cfg.Metrics.Headers = headers
	})
}

// WithTracesHeaders allows one to tell the driver to send additional HTTP
// headers with the trace payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithTracesHeaders(headers map[string]string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Headers = headers
	})
}

// WithMetricsHeaders allows one to tell the driver to send additional HTTP
// headers with the metrics payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithMetricsHeaders(headers map[string]string) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Metrics.Headers = headers
	})
}

// WithMarshal tells the driver which wire format to use when sending to the
// collector.  If unset, MarshalProto will be used
func WithMarshal(m otlp.Marshaler) HTTPOption {
	return NewHTTPOption(func(cfg *Config) {
		cfg.Marshaler = m
	})
}

// WithTimeout tells the driver the max waiting time for the backend to process
// each spans or metrics batch.  If unset, the default will be 10 seconds.
func WithTimeout(duration time.Duration) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Timeout = duration
		cfg.Metrics.Timeout = duration
	})
}

// WithTracesTimeout tells the driver the max waiting time for the backend to process
// each spans batch.  If unset, the default will be 10 seconds.
func WithTracesTimeout(duration time.Duration) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Traces.Timeout = duration
	})
}

// WithMetricsTimeout tells the driver the max waiting time for the backend to process
// each metrics batch.  If unset, the default will be 10 seconds.
func WithMetricsTimeout(duration time.Duration) GenericOption {
	return newGenericOption(func(cfg *Config) {
		cfg.Metrics.Timeout = duration
	})
}
