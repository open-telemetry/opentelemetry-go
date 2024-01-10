// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log/embedded"
)

// LoggerProvider provides access to named [Logger] instances.
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type LoggerProvider interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.LoggerProvider

	// Logger returns a new [Logger] with the provided name and configuration.
	//
	// This method should:
	//   - be safe to call concurrently,
	//   - use some default name if the passed name is empty.
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

// LoggerOption is an interface for applying Logger options.
type LoggerOption interface {
	// applyLogger is used to set a LoggerOption value of a LoggerConfig.
	applyLogger(LoggerConfig) LoggerConfig
}

// NewLoggerConfig creates a new LoggerConfig and applies
// all the given options.
func NewLoggerConfig(opts ...LoggerOption) LoggerConfig {
	var config LoggerConfig
	for _, o := range opts {
		config = o.applyLogger(config)
	}
	return config
}

type loggerOptionFunc func(LoggerConfig) LoggerConfig

func (fn loggerOptionFunc) applyLogger(cfg LoggerConfig) LoggerConfig {
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
