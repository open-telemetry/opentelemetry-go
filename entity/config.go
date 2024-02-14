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

package entity // import "go.opentelemetry.io/otel/entity"

import (
	"go.opentelemetry.io/otel/attribute"
)

// EntityEmitterConfig is a group of options for a EntityEmitter.
type EntityEmitterConfig struct {
	instrumentationVersion string
	// Schema URL of the telemetry emitted by the EntityEmitter.
	schemaURL string
	attrs     attribute.Set
}

// InstrumentationVersion returns the version of the library providing instrumentation.
func (t *EntityEmitterConfig) InstrumentationVersion() string {
	return t.instrumentationVersion
}

// InstrumentationAttributes returns the attributes associated with the library
// providing instrumentation.
func (t *EntityEmitterConfig) InstrumentationAttributes() attribute.Set {
	return t.attrs
}

// SchemaURL returns the Schema URL of the telemetry emitted by the EntityEmitter.
func (t *EntityEmitterConfig) SchemaURL() string {
	return t.schemaURL
}

// NewEntityEmitterConfig applies all the options to a returned EntityEmitterConfig.
func NewEntityEmitterConfig(options ...EntityEmitterOption) EntityEmitterConfig {
	var config EntityEmitterConfig
	for _, option := range options {
		config = option.apply(config)
	}
	return config
}

// EntityEmitterOption applies an option to a EntityEmitterConfig.
type EntityEmitterOption interface {
	apply(EntityEmitterConfig) EntityEmitterConfig
}

type entityEmitterOptionFunc func(EntityEmitterConfig) EntityEmitterConfig

func (fn entityEmitterOptionFunc) apply(cfg EntityEmitterConfig) EntityEmitterConfig {
	return fn(cfg)
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) EntityEmitterOption {
	return entityEmitterOptionFunc(
		func(cfg EntityEmitterConfig) EntityEmitterConfig {
			cfg.instrumentationVersion = version
			return cfg
		},
	)
}

// WithInstrumentationAttributes sets the instrumentation attributes.
//
// The passed attributes will be de-duplicated.
func WithInstrumentationAttributes(attr ...attribute.KeyValue) EntityEmitterOption {
	return entityEmitterOptionFunc(
		func(config EntityEmitterConfig) EntityEmitterConfig {
			config.attrs = attribute.NewSet(attr...)
			return config
		},
	)
}

// WithSchemaURL sets the schema URL for the EntityEmitter.
func WithSchemaURL(schemaURL string) EntityEmitterOption {
	return entityEmitterOptionFunc(
		func(cfg EntityEmitterConfig) EntityEmitterConfig {
			cfg.schemaURL = schemaURL
			return cfg
		},
	)
}
