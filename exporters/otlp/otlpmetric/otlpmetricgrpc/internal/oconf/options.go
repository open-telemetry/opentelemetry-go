// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oconf // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/oconf"

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/retry"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	// DefaultMaxAttempts describes how many times the driver
	// should retry the sending of the payload in case of a
	// retryable error.
	DefaultMaxAttempts int = 5
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

type (
	// HTTPTransportProxyFunc is a function that resolves which URL to use as proxy for a given request.
	// This type is compatible with `http.Transport.Proxy` and can be used to set a custom proxy function to the OTLP HTTP client.
	HTTPTransportProxyFunc func(*http.Request) (*url.URL, error)

	// SignalConfig represents signal specific configuration.
	SignalConfig struct {
		Endpoint    string
		Insecure    bool
		TLSCfg      *tls.Config
		Headers     map[string]string
		Compression Compression
		Timeout     time.Duration
		URLPath     string

		// gRPC configurations
		GRPCCredentials credentials.TransportCredentials

		TemporalitySelector metric.TemporalitySelector
		AggregationSelector metric.AggregationSelector

		Proxy HTTPTransportProxyFunc
	}

	// Config represents exporter configuration.
	Config struct {
		Metrics SignalConfig

		RetryConfig retry.Config

		// gRPC configurations
		ReconnectionPeriod time.Duration
		ServiceConfig      string
		DialOptions        []grpc.DialOption
		GRPCConn           *grpc.ClientConn
	}
)

// cleanPath returns a path with all spaces trimmed and all redundancies
// removed. If urlPath is empty or cleaning it results in an empty string,
// defaultPath is returned instead.
func cleanPath(urlPath string, defaultPath string) string {
	tmp := path.Clean(strings.TrimSpace(urlPath))
	if tmp == "." {
		return defaultPath
	}
	if !path.IsAbs(tmp) {
		tmp = fmt.Sprintf("/%s", tmp)
	}
	return tmp
}

// NewGRPCConfig returns a new Config with all settings applied from opts and
// any unset setting using the default gRPC config values.
func NewGRPCConfig(opts ...GRPCOption) Config {
	cfg := Config{
		Metrics: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorGRPCPort),
			URLPath:     DefaultMetricsPath,
			Compression: NoCompression,
			Timeout:     DefaultTimeout,

			TemporalitySelector: metric.DefaultTemporalitySelector,
			AggregationSelector: metric.DefaultAggregationSelector,
		},
		RetryConfig: retry.DefaultConfig,
	}
	cfg = ApplyGRPCEnvConfigs(cfg)
	for _, opt := range opts {
		cfg = opt.ApplyGRPCOption(cfg)
	}

	if cfg.ServiceConfig != "" {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithDefaultServiceConfig(cfg.ServiceConfig))
	}
	// Priroritize GRPCCredentials over Insecure (passing both is an error).
	if cfg.Metrics.GRPCCredentials != nil {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithTransportCredentials(cfg.Metrics.GRPCCredentials))
	} else if cfg.Metrics.Insecure {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Default to using the host's root CA.
		creds := credentials.NewTLS(nil)
		cfg.Metrics.GRPCCredentials = creds
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithTransportCredentials(creds))
	}
	if cfg.Metrics.Compression == GzipCompression {
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}
	if cfg.ReconnectionPeriod != 0 {
		p := grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: cfg.ReconnectionPeriod,
		}
		cfg.DialOptions = append(cfg.DialOptions, grpc.WithConnectParams(p))
	}

	return cfg
}

// GRPCOption applies an option to the gRPC driver.
type GRPCOption interface {
	ApplyGRPCOption(Config) Config

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

// grpcOption is an option that is only applied to the gRPC driver.
type grpcOption struct {
	fn func(Config) Config
}

func (h *grpcOption) ApplyGRPCOption(cfg Config) Config {
	return h.fn(cfg)
}

func (grpcOption) private() {}

func NewGRPCOption(fn func(cfg Config) Config) GRPCOption {
	return &grpcOption{fn: fn}
}

// Generic Options

func WithEndpoint(endpoint string) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Endpoint = endpoint
		return cfg
	})
}

func WithEndpointURL(v string) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		u, err := url.Parse(v)
		if err != nil {
			global.Error(err, "otlpmetric: parse endpoint url", "url", v)
			return cfg
		}

		cfg.Metrics.Endpoint = u.Host
		cfg.Metrics.URLPath = u.Path
		if u.Scheme != "https" {
			cfg.Metrics.Insecure = true
		}

		return cfg
	})
}

func WithCompression(compression Compression) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Compression = compression
		return cfg
	})
}

func WithURLPath(urlPath string) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.URLPath = urlPath
		return cfg
	})
}

func WithRetry(rc retry.Config) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.RetryConfig = rc
		return cfg
	})
}

func WithTLSClientConfig(tlsCfg *tls.Config) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.GRPCCredentials = credentials.NewTLS(tlsCfg)
		return cfg
	})
}

func WithInsecure() GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Insecure = true
		return cfg
	})
}

func WithSecure() GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Insecure = false
		return cfg
	})
}

func WithHeaders(headers map[string]string) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Headers = headers
		return cfg
	})
}

func WithTimeout(duration time.Duration) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Timeout = duration
		return cfg
	})
}

func WithTemporalitySelector(selector metric.TemporalitySelector) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.TemporalitySelector = selector
		return cfg
	})
}

func WithAggregationSelector(selector metric.AggregationSelector) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.AggregationSelector = selector
		return cfg
	})
}

func WithProxy(pf HTTPTransportProxyFunc) GRPCOption {
	return NewGRPCOption(func(cfg Config) Config {
		cfg.Metrics.Proxy = pf
		return cfg
	})
}
