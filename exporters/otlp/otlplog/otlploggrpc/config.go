// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/retry"
)

// Option applies an option to the Exporter.
type Option interface {
	applyHTTPOption(config) config
}

type config struct {
	// TODO: implement.
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.applyHTTPOption(c)
	}
	return c
}

// RetryConfig defines configuration for retrying the export of log data
// that failed.
//
// This configuration does not define any network retry strategy. That is
// entirely handled by the gRPC ClientConn.
type RetryConfig retry.Config

// WithInsecure disables client transport security for the Exporter's gRPC
// connection, just like grpc.WithInsecure()
// (https://pkg.go.dev/google.golang.org/grpc#WithInsecure) does.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used to determine client security. If the endpoint has a
// scheme of "http" or "unix" client security will be disabled. If both are
// set, OTEL_EXPORTER_OTLP_LOGS_ENDPOINT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, client security will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithInsecure() Option {
	// TODO: implement.
	return nil
}

// WithEndpoint sets the target endpoint the Exporter will connect to.
//
// If the OTEL_EXPORTER_OTLP_ENDPOINT or OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// environment variable is set, and this option is not passed, that variable
// value will be used. If both are set, OTEL_EXPORTER_OTLP_LOGS_ENDPOINT
// will take precedence.
//
// If both this option and WithEndpointURL are used, the last used option will
// take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, "localhost:4317" will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithEndpoint(endpoint string) Option {
	// TODO: implement.
	return nil
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
// passed, "localhost:4317" will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithEndpointURL(u string) Option {
	// TODO: implement.
	return nil
}

// WithReconnectionPeriod set the minimum amount of time between connection
// attempts to the target endpoint.
//
// This option has no effect if WithGRPCConn is used.
func WithReconnectionPeriod(rp time.Duration) Option {
	// TODO: implement.
	return nil
}

// WithCompressor sets the compressor the gRPC client uses.
// Supported compressor values: "gzip".
//
// If the OTEL_EXPORTER_OTLP_COMPRESSION or
// OTEL_EXPORTER_OTLP_LOGS_COMPRESSION environment variable is set, and
// this option is not passed, that variable value will be used. That value can
// be either "none" or "gzip". If both are set,
// OTEL_EXPORTER_OTLP_LOGS_COMPRESSION will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no compressor will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithCompressor(compressor string) Option {
	// TODO: implement.
	return nil
}

// WithHeaders will send the provided headers with each gRPC requests.
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
	// TODO: implement.
	return nil
}

// WithTLSCredentials sets the gRPC connection to use creds.
//
// If the OTEL_EXPORTER_OTLP_CERTIFICATE or
// OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE environment variable is set, and
// this option is not passed, that variable value will be used. The value will
// be parsed the filepath of the TLS certificate chain to use. If both are
// set, OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no TLS credentials will be used.
//
// This option has no effect if WithGRPCConn is used.
func WithTLSCredentials(_ credentials.TransportCredentials) Option {
	// TODO: implement.
	return nil
}

// WithServiceConfig defines the default gRPC service config used.
//
// This option has no effect if WithGRPCConn is used.
func WithServiceConfig(serviceConfig string) Option {
	// TODO: implement.
	return nil
}

// WithDialOption sets explicit grpc.DialOptions to use when establishing a
// gRPC connection. The options here are appended to the internal grpc.DialOptions
// used so they will take precedence over any other internal grpc.DialOptions
// they might conflict with.
// The [grpc.WithBlock], [grpc.WithTimeout], and [grpc.WithReturnConnectionError]
// grpc.DialOptions are ignored.
//
// This option has no effect if WithGRPCConn is used.
func WithDialOption(_ ...grpc.DialOption) Option {
	// TODO: implement.
	return nil
}

// WithGRPCConn sets conn as the gRPC ClientConn used for all communication.
//
// This option takes precedence over any other option that relates to
// establishing or persisting a gRPC connection to a target endpoint. Any
// other option of those types passed will be ignored.
//
// It is the callers responsibility to close the passed conn. The Exporter
// Shutdown method will not close this connection.
func WithGRPCConn(_ *grpc.ClientConn) Option {
	// TODO: implement.
	return nil
}

// WithTimeout sets the max amount of time an Exporter will attempt an export.
//
// This takes precedence over any retry settings defined by WithRetry. Once
// this time limit has been reached the export is abandoned and the log
// data is dropped.
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
	// TODO: implement.
	return nil
}

// WithRetry sets the retry policy for transient retryable errors that are
// returned by the target endpoint.
//
// If the target endpoint responds with not only a retryable error, but
// explicitly returns a backoff time in the response, that time will take
// precedence over these settings.
//
// These settings do not define any network retry strategy. That is entirely
// handled by the gRPC ClientConn.
//
// If unset, the default retry policy will be used. It will retry the export
// 5 seconds after receiving a retryable error and increase exponentially
// after each error for no more than a total time of 1 minute.
func WithRetry(settings RetryConfig) Option {
	// TODO: implement.
	return nil
}
