// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log/embedded"
)

// LoggerProvider TODO: comment.
type LoggerProvider interface {
	embedded.LoggerProvider

	// Logger TODO: comment.
	Logger(name string, options ...LoggerOption) Logger
}

// LoggerConfig contains options for Logger.
type LoggerConfig struct {
	instrumentationVersion string
	schemaURL              string
	attrs                  attribute.Set

	// Ensure forward compatibility by explicitly making this not comparable.
	noCmp [0]func() //nolint: unused  // This is indeed used.
}

// InstrumentationVersion returns the version of the library providing
// instrumentation.
func (cfg LoggerConfig) InstrumentationVersion() string {
	return cfg.instrumentationVersion
}

// InstrumentationAttributes returns the attributes associated with the library
// providing instrumentation.
func (cfg LoggerConfig) InstrumentationAttributes() attribute.Set {
	return cfg.attrs
}

// SchemaURL is the schema_url of the library providing instrumentation.
func (cfg LoggerConfig) SchemaURL() string {
	return cfg.schemaURL
}

// LoggerOption is an interface for applying Meter options.
type LoggerOption interface {
	// applyMeter is used to set a LoggerOption value of a LoggerConfig.
	applyMeter(LoggerConfig) LoggerConfig
}

// NewLoggerConfig creates a new LoggerConfig and applies
// all the given options.
// TODO: Add unit tests.
func NewLoggerConfig(opts ...LoggerOption) LoggerConfig {
	var config LoggerConfig
	for _, o := range opts {
		config = o.applyMeter(config)
	}
	return config
}

type loggerOptionFunc func(LoggerConfig) LoggerConfig

func (fn loggerOptionFunc) applyMeter(cfg LoggerConfig) LoggerConfig {
	return fn(cfg)
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) LoggerOption {
	return loggerOptionFunc(func(config LoggerConfig) LoggerConfig {
		config.instrumentationVersion = version
		return config
	})
}

// WithInstrumentationAttributes sets the instrumentation attributes.
//
// The passed attributes will be de-duplicated.
func WithInstrumentationAttributes(attr ...attribute.KeyValue) LoggerOption {
	return loggerOptionFunc(func(config LoggerConfig) LoggerConfig {
		config.attrs = attribute.NewSet(attr...)
		return config
	})
}

// WithSchemaURL sets the schema URL.
func WithSchemaURL(schemaURL string) LoggerOption {
	return loggerOptionFunc(func(config LoggerConfig) LoggerConfig {
		config.schemaURL = schemaURL
		return config
	})
}
