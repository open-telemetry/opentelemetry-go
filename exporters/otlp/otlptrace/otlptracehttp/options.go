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

package otlptracehttp // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

import (
	"crypto/tls"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/internal/envconfig"
	"go.opentelemetry.io/otel/exporters/otlp/internal/retry"
)

// Compression describes the compression used for payloads sent to the
// collector.
type Compression envconfig.Compression

const (
	// NoCompression tells the driver to send payloads without
	// compression.
	NoCompression = Compression(envconfig.NoCompression)
	// GzipCompression tells the driver to send payloads after
	// compressing them with gzip.
	GzipCompression = Compression(envconfig.GzipCompression)
)

// Option applies an option to the HTTP client.
type Option interface {
	applyHTTPOption(envconfig.Config) envconfig.Config
}

func asHTTPOptions(opts []Option) []envconfig.HTTPOption {
	converted := make([]envconfig.HTTPOption, len(opts))
	for i, o := range opts {
		converted[i] = envconfig.NewHTTPOption(o.applyHTTPOption)
	}
	return converted
}

// RetryConfig defines configuration for retrying batches in case of export
// failure using an exponential backoff.
type RetryConfig retry.Config

type wrappedOption struct {
	envconfig.HTTPOption
}

func (w wrappedOption) applyHTTPOption(cfg envconfig.Config) envconfig.Config {
	return w.ApplyHTTPOption(cfg)
}

// WithEndpoint allows one to set the address of the collector
// endpoint that the driver will use to send spans. If
// unset, it will instead try to use
// the default endpoint (localhost:4318). Note that the endpoint
// must not contain any URL path.
func WithEndpoint(endpoint string) Option {
	return wrappedOption{envconfig.WithEndpoint(endpoint)}
}

// WithCompression tells the driver to compress the sent data.
func WithCompression(compression Compression) Option {
	return wrappedOption{envconfig.WithCompression(envconfig.Compression(compression))}
}

// WithURLPath allows one to override the default URL path used
// for sending traces. If unset, default ("/v1/traces") will be used.
func WithURLPath(urlPath string) Option {
	return wrappedOption{envconfig.WithURLPath(urlPath)}
}

// WithTLSClientConfig can be used to set up a custom TLS
// configuration for the client used to send payloads to the
// collector. Use it if you want to use a custom certificate.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return wrappedOption{envconfig.WithTLSClientConfig(tlsCfg)}
}

// WithInsecure tells the driver to connect to the collector using the
// HTTP scheme, instead of HTTPS.
func WithInsecure() Option {
	return wrappedOption{envconfig.WithInsecure()}
}

// WithHeaders allows one to tell the driver to send additional HTTP
// headers with the payloads. Specifying headers like Content-Length,
// Content-Encoding and Content-Type may result in a broken driver.
func WithHeaders(headers map[string]string) Option {
	return wrappedOption{envconfig.WithHeader(headers)}
}

// WithTimeout tells the driver the max waiting time for the backend to process
// each spans batch.  If unset, the default will be 10 seconds.
func WithTimeout(duration time.Duration) Option {
	return wrappedOption{envconfig.WithTimeout(duration)}
}

// WithRetry configures the retry policy for transient errors that may occurs
// when exporting traces. An exponential back-off algorithm is used to ensure
// endpoints are not overwhelmed with retries. If unset, the default retry
// policy will retry after 5 seconds and increase exponentially after each
// error for a total of 1 minute.
func WithRetry(rc RetryConfig) Option {
	return wrappedOption{envconfig.WithRetry(retry.Config(rc))}
}
