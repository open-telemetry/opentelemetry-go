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

package otlphttp

import (
	"crypto/tls"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp"
)

// Compression describes the compression used for payloads sent to the
// collector.
type Compression int

const (
	// NoCompression tells the driver to send payloads without
	// compression.
	NoCompression Compression = iota
	// GzipCompression tells the driver to send payloads after
	// compressing them with gzip.
	GzipCompression
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
)

// Marshaler describes the kind of message format sent to the collector
type Marshaler int

const (
	// MarshalProto tells the driver to send using the protobuf binary format.
	MarshalProto Marshaler = iota
	// MarshalJSON tells the driver to send using json format.
	MarshalJSON
)

type signalConfig struct {
	endpoint    string
	insecure    bool
	tlsCfg      *tls.Config
	headers     map[string]string
	compression Compression
	timeout     time.Duration
	urlPath     string
}

type config struct {
	metrics signalConfig
	traces  signalConfig

	maxAttempts int
	backoff     time.Duration
	marshaler   Marshaler
}

func newDefaultConfig() config {
	c := config{
		traces: signalConfig{
			endpoint:    fmt.Sprintf("%s:%d", otlp.DefaultCollectorHost, otlp.DefaultCollectorPort),
			urlPath:     DefaultTracesPath,
			compression: NoCompression,
			timeout:     DefaultTimeout,
		},
		metrics: signalConfig{
			endpoint:    fmt.Sprintf("%s:%d", otlp.DefaultCollectorHost, otlp.DefaultCollectorPort),
			urlPath:     DefaultMetricsPath,
			compression: NoCompression,
			timeout:     DefaultTimeout,
		},
		maxAttempts: DefaultMaxAttempts,
		backoff:     DefaultBackoff,
	}

	return c
}

// Option applies an option to the HTTP driver.
type Option interface {
	Apply(*config)

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

type genericOption struct {
	fn func(*config)
}

func (g *genericOption) Apply(cfg *config) {
	g.fn(cfg)
}

func (genericOption) private() {}

func newGenericOption(fn func(cfg *config)) Option {
	return &genericOption{fn: fn}
}

// WithEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send metrics and spans. If
// unset, it will instead try to use
// DefaultCollectorHost:DefaultCollectorPort. Note that the endpoint
// must not contain any URL path.
func WithEndpoint(endpoint string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.endpoint = endpoint
		cfg.metrics.endpoint = endpoint
	})
}

// WithTracesEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send spans. If
// unset, it will instead try to use the Endpoint configuration.
// Note that the endpoint must not contain any URL path.
func WithTracesEndpoint(endpoint string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.endpoint = endpoint
	})
}

// WithMetricsEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send metrics. If
// unset, it will instead try to use the Endpoint configuration.
// Note that the endpoint must not contain any URL path.
func WithMetricsEndpoint(endpoint string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.endpoint = endpoint
	})
}

// WithCompression tells the driver to compress the sent data.
func WithCompression(compression Compression) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.compression = compression
		cfg.metrics.compression = compression
	})
}

// WithTracesCompression tells the driver to compress the sent traces data.
func WithTracesCompression(compression Compression) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.compression = compression
	})
}

// WithMetricsCompression tells the driver to compress the sent metrics data.
func WithMetricsCompression(compression Compression) Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.compression = compression
	})
}

// WithTracesURLPath allows one to override the default URL path used
// for sending traces. If unset, DefaultTracesPath will be used.
func WithTracesURLPath(urlPath string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.urlPath = urlPath
	})
}

// WithMetricsURLPath allows one to override the default URL path used
// for sending metrics. If unset, DefaultMetricsPath will be used.
func WithMetricsURLPath(urlPath string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.urlPath = urlPath
	})
}

// WithMaxAttempts allows one to override how many times the driver
// will try to send the payload in case of retryable errors. If unset,
// DefaultMaxAttempts will be used.
func WithMaxAttempts(maxAttempts int) Option {
	return newGenericOption(func(cfg *config) {
		cfg.maxAttempts = maxAttempts
	})
}

// WithBackoff tells the driver to use the duration as a base of the
// exponential backoff strategy. If unset, DefaultBackoff will be
// used.
func WithBackoff(duration time.Duration) Option {
	return newGenericOption(func(cfg *config) {
		cfg.backoff = duration
	})
}

// WithTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send payloads to the
// collector. Use it if you want to use a custom certificate.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.tlsCfg = tlsCfg
		cfg.metrics.tlsCfg = tlsCfg
	})
}

// WithTracesTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send traces.
// Use it if you want to use a custom certificate.
func WithTracesTLSClientConfig(tlsCfg *tls.Config) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.tlsCfg = tlsCfg
	})
}

// WithMetricsTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send metrics.
// Use it if you want to use a custom certificate.
func WithMetricsTLSClientConfig(tlsCfg *tls.Config) Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.tlsCfg = tlsCfg
	})
}

// WithInsecure tells the driver to connect to the collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecure() Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.insecure = true
		cfg.metrics.insecure = true
	})
}

// WithInsecureTraces tells the driver to connect to the traces collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecureTraces() Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.insecure = true
	})
}

// WithInsecure tells the driver to connect to the metrics collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecureMetrics() Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.insecure = true
	})
}

// WithHeaders allows one to tell the driver to send additional HTTP
// headers with the payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithHeaders(headers map[string]string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.headers = headers
		cfg.metrics.headers = headers
	})
}

// WithTracesHeaders allows one to tell the driver to send additional HTTP
// headers with the trace payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithTracesHeaders(headers map[string]string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.headers = headers
	})
}

// WithMetricsHeaders allows one to tell the driver to send additional HTTP
// headers with the metrics payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithMetricsHeaders(headers map[string]string) Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.headers = headers
	})
}

// WithMarshal tells the driver which wire format to use when sending to the
// collector.  If unset, MarshalProto will be used
func WithMarshal(m Marshaler) Option {
	return newGenericOption(func(cfg *config) {
		cfg.marshaler = m
	})
}

// WithTimeout tells the driver the max waiting time for the backend to process
// each spans or metrics batch.  If unset, the default will be 10 seconds.
func WithTimeout(duration time.Duration) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.timeout = duration
		cfg.metrics.timeout = duration
	})
}

// WithTracesTimeout tells the driver the max waiting time for the backend to process
// each spans batch.  If unset, the default will be 10 seconds.
func WithTracesTimeout(duration time.Duration) Option {
	return newGenericOption(func(cfg *config) {
		cfg.traces.timeout = duration
	})
}

// WithMetricsTimeout tells the driver the max waiting time for the backend to process
// each metrics batch.  If unset, the default will be 10 seconds.
func WithMetricsTimeout(duration time.Duration) Option {
	return newGenericOption(func(cfg *config) {
		cfg.metrics.timeout = duration
	})
}
