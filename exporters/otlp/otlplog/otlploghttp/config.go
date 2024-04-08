// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal/retry"
	"go.opentelemetry.io/otel/internal/global"
)

// Default values.
var (
	defaultEndpoint                        = "localhost:4318"
	defaultPath                            = "/v1/logs"
	defaultTimeout                         = 10 * time.Second
	defaultProxy    HTTPTransportProxyFunc = http.ProxyFromEnvironment
	defaultRetryCfg                        = RetryConfig{} // TODO: define.
)

// Option applies an option to the Exporter.
type Option interface {
	applyHTTPOption(config) config
}

type fnOpt func(config) config

func (f fnOpt) applyHTTPOption(c config) config { return f(c) }

type config struct {
	endpoint    setting[string]
	path        setting[string]
	insecure    setting[bool]
	tlsCfg      setting[*tls.Config]
	headers     setting[map[string]string]
	compression setting[Compression]
	timeout     setting[time.Duration]
	proxy       setting[HTTPTransportProxyFunc]
	retryCfg    setting[RetryConfig]
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.applyHTTPOption(c)
	}

	c.endpoint = c.endpoint.Resolve(
		fallback[string](defaultEndpoint),
	)
	c.path = c.path.Resolve(
		fallback[string](defaultPath),
	)
	c.timeout = c.timeout.Resolve(
		fallback[time.Duration](defaultTimeout),
	)
	c.proxy = c.proxy.Resolve(
		fallback[HTTPTransportProxyFunc](defaultProxy),
	)
	c.retryCfg = c.retryCfg.Resolve(
		fallback[RetryConfig](defaultRetryCfg),
	)

	return c
}

// WithEndpoint sets the target endpoint the Exporter will connect to. This
// endpoint is specified as a host and optional port, no path or scheme should
// be included (see WithInsecure and WithURLPath).
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used. If both are set, OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, "localhost:4318" will be used.
func WithEndpoint(endpoint string) Option {
	return fnOpt(func(c config) config {
		c.endpoint = newSetting(endpoint)
		return c
	})
}

// WithEndpointURL sets the target endpoint URL the Exporter will connect to.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used. If both are set, OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// will take precedence.
//
// If both this option and WithEndpoint are used, the last used option will
// take precedence.
//
// If an invalid URL is provided, the default value will be kept.
//
// By default, if an environment variable is not set, and this option is not
// passed, "localhost:4318" will be used.
func WithEndpointURL(rawURL string) Option {
	u, err := url.Parse(rawURL)
	if err != nil {
		global.Error(err, "otlpmetric: parse endpoint url", "url", rawURL)
		return fnOpt(func(c config) config { return c })
	}
	return fnOpt(func(c config) config {
		c.endpoint = newSetting(u.Host)
		c.path = newSetting(u.Path)
		if u.Scheme != "https" {
			c.insecure = newSetting(true)
		} else {
			c.insecure = newSetting(false)
		}
		return c
	})
}

// Compression describes the compression used for exported payloads.
type Compression int

const (
	// NoCompression represents that no compression should be used.
	NoCompression Compression = iota
	// GzipCompression represents that gzip compression should be used.
	GzipCompression
)

// WithCompression sets the compression strategy the Exporter will use to
// compress the HTTP body.
//
// If the OTEL_EXPORTER_OTLP_COMPRESSION or
// OTEL_EXPORTER_OTLP_LOGS_COMPRESSION environment variable is set, and
// this option is not passed, that variable value will be used. That value can
// be either "none" or "gzip". If both are set,
// OTEL_EXPORTER_OTLP_LOGS_COMPRESSION will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no compression strategy will be used.
func WithCompression(compression Compression) Option {
	return fnOpt(func(c config) config {
		c.compression = newSetting(compression)
		return c
	})
}

// WithURLPath sets the URL path the Exporter will send requests to.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// environment variable is set, and this option is not passed, the path
// contained in that variable value will be used. If both are set,
// OTEL_EXPORTER_OTLP_LOGS_ENDPOINT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, "/v1/logs" will be used.
func WithURLPath(urlPath string) Option {
	return fnOpt(func(c config) config {
		c.path = newSetting(urlPath)
		return c
	})
}

// WithTLSClientConfig sets the TLS configuration the Exporter will use for
// HTTP requests.
//
// If the OTEL_EXPORTER_OTLP_CERTIFICATE or
// OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE environment variable is set, and
// this option is not passed, that variable value will be used. The value will
// be parsed the filepath of the TLS certificate chain to use. If both are
// set, OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, the system default configuration is used.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return fnOpt(func(c config) config {
		c.tlsCfg = newSetting(tlsCfg.Clone())
		return c
	})
}

// WithInsecure disables client transport security for the Exporter's HTTP
// connection.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used to determine client security. If the endpoint has a
// scheme of "http" or "unix" client security will be disabled. If both are
// set, OTEL_EXPORTER_OTLP_LOGS_ENDPOINT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, client security will be used.
func WithInsecure() Option {
	return fnOpt(func(c config) config {
		c.insecure = newSetting(true)
		return c
	})
}

// WithHeaders will send the provided headers with each HTTP requests.
//
// If the OTEL_EXPORTER_OTLP_HEADERS or OTEL_EXPORTER_OTLP_LOGS_HEADERS
// environment variable is set, and this option is not passed, that variable
// value will be used. The value will be parsed as a list of key value pairs.
// These pairs are expected to be in the W3C Correlation-Context format
// without additional semi-colon delimited metadata (i.e. "k1=v1,k2=v2"). If
// both are set, OTEL_EXPORTER_OTLP_LOGS_HEADERS will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no user headers will be set.
func WithHeaders(headers map[string]string) Option {
	return fnOpt(func(c config) config {
		c.headers = newSetting(headers)
		return c
	})
}

// WithTimeout sets the max amount of time an Exporter will attempt an export.
//
// This takes precedence over any retry settings defined by WithRetry. Once
// this time limit has been reached the export is abandoned and the log data is
// dropped.
//
// If the OTEL_EXPORTER_OTLP_TIMEOUT or OTEL_EXPORTER_OTLP_LOGS_TIMEOUT
// environment variable is set, and this option is not passed, that variable
// value will be used. The value will be parsed as an integer representing the
// timeout in milliseconds. If both are set,
// OTEL_EXPORTER_OTLP_LOGS_TIMEOUT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, a timeout of 10 seconds will be used.
func WithTimeout(duration time.Duration) Option {
	return fnOpt(func(c config) config {
		c.timeout = newSetting(duration)
		return c
	})
}

// RetryConfig defines configuration for retrying the export of log data that
// failed.
type RetryConfig retry.Config

// WithRetry sets the retry policy for transient retryable errors that are
// returned by the target endpoint.
//
// If the target endpoint responds with not only a retryable error, but
// explicitly returns a backoff time in the response, that time will take
// precedence over these settings.
//
// If unset, the default retry policy will be used. It will retry the export
// 5 seconds after receiving a retryable error and increase exponentially
// after each error for no more than a total time of 1 minute.
func WithRetry(rc RetryConfig) Option {
	return fnOpt(func(c config) config {
		c.retryCfg = newSetting(rc)
		return c
	})
}

// HTTPTransportProxyFunc is a function that resolves which URL to use as proxy
// for a given request. This type is compatible with http.Transport.Proxy and
// can be used to set a custom proxy function to the OTLP HTTP client.
type HTTPTransportProxyFunc func(*http.Request) (*url.URL, error)

// WithProxy sets the Proxy function the client will use to determine the
// proxy to use for an HTTP request. If this option is not used, the client
// will use [http.ProxyFromEnvironment].
func WithProxy(pf HTTPTransportProxyFunc) Option {
	return fnOpt(func(c config) config {
		c.proxy = newSetting(pf)
		return c
	})
}

// setting is a configuration setting value.
type setting[T any] struct {
	Value T
	Set   bool
}

// newSetting returns a new [setting] with the value set.
func newSetting[T any](value T) setting[T] {
	return setting[T]{Value: value, Set: true}
}

// resolver returns an updated setting after applying an resolution operation.
type resolver[T any] func(setting[T]) setting[T]

// Resolve returns a resolved version of s.
//
// It will apply all the passed fn in the order provided, chaining together the
// return setting to the next input. The setting s is used as the initial
// argument to the first fn.
//
// Each fn needs to validate if it should apply given the Set state of the
// setting. This will not perform any checks on the set state when chaining
// function.
func (s setting[T]) Resolve(fn ...resolver[T]) setting[T] {
	for _, f := range fn {
		s = f(s)
	}
	return s
}

// fallback returns a resolve that will set a setting value to val if it is not
// already set.
//
// This is usually passed at the end of a resolver chain to ensure a default is
// applied if the setting has not already been set.
func fallback[T any](val T) resolver[T] {
	return func(s setting[T]) setting[T] {
		if !s.Set {
			s.Value = val
			s.Set = true
		}
		return s
	}
}
