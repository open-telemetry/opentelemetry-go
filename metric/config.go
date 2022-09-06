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

package metric // import "go.opentelemetry.io/otel/metric"

import "go.opentelemetry.io/otel/attribute"

// MeterConfig contains options for Meters.
type MeterConfig struct {
	instrumentationVersion string
	schemaURL              string
	attributes             attribute.Set
}

// InstrumentationVersion is the version of the library providing instrumentation.
func (cfg MeterConfig) InstrumentationVersion() string {
	return cfg.instrumentationVersion
}

// SchemaURL is the schema_url of the library providing instrumentation.
func (cfg MeterConfig) SchemaURL() string {
	return cfg.schemaURL
}

// Attributes returns the scope attribute set of the Meter.
func (t MeterConfig) Attributes() attribute.Set {
	return t.attributes
}

// MeterOption is an interface for applying Meter options.
type MeterOption interface {
	// applyMeter is used to set a MeterOption value of a MeterConfig.
	applyMeter(MeterConfig) MeterConfig
}

// NewMeterConfig creates a new MeterConfig and applies
// all the given options.
func NewMeterConfig(opts ...MeterOption) MeterConfig {
	var config MeterConfig
	for _, o := range opts {
		config = o.applyMeter(config)
	}
	return config
}

type meterOptionFunc func(MeterConfig) MeterConfig

func (fn meterOptionFunc) applyMeter(cfg MeterConfig) MeterConfig {
	return fn(cfg)
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) MeterOption {
	return meterOptionFunc(func(config MeterConfig) MeterConfig {
		config.instrumentationVersion = version
		return config
	})
}

// WithSchemaURL sets the schema URL.
func WithSchemaURL(schemaURL string) MeterOption {
	return meterOptionFunc(func(config MeterConfig) MeterConfig {
		config.schemaURL = schemaURL
		return config
	})
}

// WithScopeAttributes sets the attributes for the scope of a Meter. The
// attributes are stored as an attribute set. Duplicate values are removed, the
// last value is used.
func WithScopeAttributes(attr ...attribute.KeyValue) MeterOption {
	return meterOptionFunc(func(cfg MeterConfig) MeterConfig {
		cfg.attributes = attribute.NewSet(attr...)
		return cfg
	})
}
