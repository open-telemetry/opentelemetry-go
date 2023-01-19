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

package log // import "go.opentelemetry.io/otel/log"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// LoggerConfig is a group of options for a Logger.
type LoggerConfig struct {
	instrumentationVersion string
	// Schema URL of the telemetry emitted by the Logger.
	schemaURL string
}

// InstrumentationVersion returns the version of the library providing instrumentation.
func (t *LoggerConfig) InstrumentationVersion() string {
	return t.instrumentationVersion
}

// SchemaURL returns the Schema URL of the telemetry emitted by the Logger.
func (t *LoggerConfig) SchemaURL() string {
	return t.schemaURL
}

// NewLoggerConfig applies all the options to a returned LoggerConfig.
func NewLoggerConfig(options ...LoggerOption) LoggerConfig {
	var config LoggerConfig
	for _, option := range options {
		config = option.apply(config)
	}
	return config
}

// LoggerOption applies an option to a LoggerConfig.
type LoggerOption interface {
	apply(LoggerConfig) LoggerConfig
}

type loggerOptionFunc func(LoggerConfig) LoggerConfig

func (fn loggerOptionFunc) apply(cfg LoggerConfig) LoggerConfig {
	return fn(cfg)
}

// LogRecordConfig is a group of options for a LogRecord.
type LogRecordConfig struct {
	attributes []attribute.KeyValue
	timestamp  time.Time
}

// Attributes describe the associated qualities of a LogRecord.
func (cfg *LogRecordConfig) Attributes() []attribute.KeyValue {
	return cfg.attributes
}

// Timestamp is a time in a LogRecord life-cycle.
func (cfg *LogRecordConfig) Timestamp() time.Time {
	return cfg.timestamp
}

// NewLogRecordConfig applies all the options to a returned LogRecordConfig.
// No validation is performed on the returned LogRecordConfig (e.g. no uniqueness
// checking or bounding of data), it is left to the SDK to perform this
// action.
func NewLogRecordConfig(options ...LogRecordOption) LogRecordConfig {
	var c LogRecordConfig
	for _, option := range options {
		c = option.applyLogRecord(c)
	}
	return c
}

// LogRecordOption applies an option to a LogRecordConfig. These options are applicable
// only when the span is created.
type LogRecordOption interface {
	applyLogRecord(LogRecordConfig) LogRecordConfig
}

type logRecordOptionFunc func(LogRecordConfig) LogRecordConfig

func (fn logRecordOptionFunc) applyLogRecord(cfg LogRecordConfig) LogRecordConfig {
	return fn(cfg)
}

type attributeOption []attribute.KeyValue

func (o attributeOption) applyLogRecord(c LogRecordConfig) LogRecordConfig {
	c.attributes = append(c.attributes, []attribute.KeyValue(o)...)
	return c
}

// WithAttributes adds the attributes related to a span life-cycle event.
// These attributes are used to describe the work a LogRecord represents when this
// option is provided to a LogRecord's start or end events. Otherwise, these
// attributes provide additional information about the event being recorded
// (e.g. error, state change, processing progress, system event).
//
// If multiple of these options are passed the attributes of each successive
// option will extend the attributes instead of overwriting. There is no
// guarantee of uniqueness in the resulting attributes.
func WithAttributes(attributes ...attribute.KeyValue) LogRecordOption {
	return attributeOption(attributes)
}

type timestampOption time.Time

func (o timestampOption) applyLogRecord(c LogRecordConfig) LogRecordConfig {
	c.timestamp = time.Time(o)
	return c
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) LoggerOption {
	return loggerOptionFunc(
		func(cfg LoggerConfig) LoggerConfig {
			cfg.instrumentationVersion = version
			return cfg
		},
	)
}

// WithSchemaURL sets the schema URL for the Logger.
func WithSchemaURL(schemaURL string) LoggerOption {
	return loggerOptionFunc(
		func(cfg LoggerConfig) LoggerConfig {
			cfg.schemaURL = schemaURL
			return cfg
		},
	)
}
