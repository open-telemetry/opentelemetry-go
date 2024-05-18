// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/otlpconf"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/retry"
)

type wrappedOption struct {
	otlpconf.Option
}

func (w wrappedOption) applyOption(cfg config) config {
	cfg.Config = w.Option.ApplyOption(cfg.Config)
	return cfg
}

type fnOpt func(config) config

func (f fnOpt) applyOption(c config) config { return f(c) }

// Option applies an option to the Exporter.
type Option interface {
	applyOption(config) config
}

type config struct {
	otlpconf.Config

	// gRPC configurations
	gRPCCredentials    otlpconf.Setting[credentials.TransportCredentials]
	serviceConfig      otlpconf.Setting[string]
	reconnectionPeriod otlpconf.Setting[time.Duration]
	dialOptions        otlpconf.Setting[[]grpc.DialOption]
	gRPCConn           otlpconf.Setting[*grpc.ClientConn]
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.applyOption(c)
	}

	c.Config = otlpconf.LoadConfig(c.Config)

	if !c.gRPCCredentials.Set && c.Config.TLSCfg.Set {
		c.gRPCCredentials = otlpconf.NewSetting(credentials.NewTLS(c.Config.TLSCfg.Value))
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
	return wrappedOption{Option: otlpconf.WithInsecure()}
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
	return wrappedOption{Option: otlpconf.WithEndpoint(endpoint)}
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
	return wrappedOption{Option: otlpconf.WithEndpointURL(u)}
}

// WithReconnectionPeriod set the minimum amount of time between connection
// attempts to the target endpoint.
//
// This option has no effect if WithGRPCConn is used.
func WithReconnectionPeriod(rp time.Duration) Option {
	return fnOpt(func(c config) config {
		c.reconnectionPeriod = otlpconf.NewSetting(rp)
		return c
	})
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
	return wrappedOption{Option: otlpconf.WithCompression(compressorToCompression(compressor))}
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
	return wrappedOption{Option: otlpconf.WithHeaders(headers)}
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
func WithTLSCredentials(credential credentials.TransportCredentials) Option {
	return fnOpt(func(c config) config {
		c.gRPCCredentials = otlpconf.NewSetting(credential)
		return c
	})
}

// WithServiceConfig defines the default gRPC service config used.
//
// This option has no effect if WithGRPCConn is used.
func WithServiceConfig(serviceConfig string) Option {
	return fnOpt(func(c config) config {
		c.serviceConfig = otlpconf.NewSetting(serviceConfig)
		return c
	})
}

// WithDialOption sets explicit grpc.DialOptions to use when establishing a
// gRPC connection. The options here are appended to the internal grpc.DialOptions
// used so they will take precedence over any other internal grpc.DialOptions
// they might conflict with.
// The [grpc.WithBlock], [grpc.WithTimeout], and [grpc.WithReturnConnectionError]
// grpc.DialOptions are ignored.
//
// This option has no effect if WithGRPCConn is used.
func WithDialOption(opts ...grpc.DialOption) Option {
	return fnOpt(func(c config) config {
		c.dialOptions = otlpconf.NewSetting(opts)
		return c
	})
}

// WithGRPCConn sets conn as the gRPC ClientConn used for all communication.
//
// This option takes precedence over any other option that relates to
// establishing or persisting a gRPC connection to a target endpoint. Any
// other option of those types passed will be ignored.
//
// It is the callers responsibility to close the passed conn. The Exporter
// Shutdown method will not close this connection.
func WithGRPCConn(conn *grpc.ClientConn) Option {
	return fnOpt(func(c config) config {
		c.gRPCConn = otlpconf.NewSetting(conn)
		return c
	})
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
	return wrappedOption{Option: otlpconf.WithTimeout(duration)}
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
	return wrappedOption{Option: otlpconf.WithRetry(retry.Config(settings))}
}

func compressorToCompression(compressor string) otlpconf.Compression {
	if compressor == "gzip" {
		return otlpconf.GzipCompression
	}

	otel.Handle(fmt.Errorf("invalid compression type: '%s', using no compression as default", compressor))
	return otlpconf.NoCompression
}
