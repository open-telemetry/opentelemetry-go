// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpconf // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/otlpconf"

import (
	"crypto/tls"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/retry"
	"go.opentelemetry.io/otel/internal/global"
)

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
	return fnOpt(func(c Config) Config {
		c.Insecure = NewSetting(true)
		return c
	})
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
	return fnOpt(func(c Config) Config {
		c.Endpoint = NewSetting(endpoint)
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
		return fnOpt(func(c Config) Config { return c })
	}
	return fnOpt(func(c Config) Config {
		c.Endpoint = NewSetting(u.Host)
		c.Path = NewSetting(u.Path)
		if u.Scheme != "https" {
			c.Insecure = NewSetting(true)
		} else {
			c.Insecure = NewSetting(false)
		}
		return c
	})
}

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
	return fnOpt(func(c Config) Config {
		c.Compression = NewSetting(compression)
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
	return fnOpt(func(c Config) Config {
		c.Path = NewSetting(urlPath)
		return c
	})
}

// WithTLSClientConfig sets the TLS Configuration the Exporter will use for
// HTTP requests.
//
// If the OTEL_EXPORTER_OTLP_CERTIFICATE or
// OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE environment variable is set, and
// this option is not passed, that variable value will be used. The value will
// be parsed the filepath of the TLS certificate chain to use. If both are
// set, OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, the system default Configuration is used.
func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return fnOpt(func(c Config) Config {
		c.TLSCfg = NewSetting(tlsCfg.Clone())
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
	return fnOpt(func(c Config) Config {
		c.Headers = NewSetting(headers)
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
	return fnOpt(func(c Config) Config {
		c.Timeout = NewSetting(duration)
		return c
	})
}

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
func WithRetry(rc retry.Config) Option {
	return fnOpt(func(c Config) Config {
		c.RetryCfg = NewSetting(rc)
		return c
	})
}
