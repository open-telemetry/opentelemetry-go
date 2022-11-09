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

import (
	"go.opentelemetry.io/otel/metric/unit"
)

// MeterConfig contains options for Meters.
type MeterConfig struct {
	instrumentationVersion string
	schemaURL              string
}

// InstrumentationVersion is the version of the library providing instrumentation.
func (cfg MeterConfig) InstrumentationVersion() string {
	return cfg.instrumentationVersion
}

// SchemaURL is the schema_url of the library providing instrumentation.
func (cfg MeterConfig) SchemaURL() string {
	return cfg.schemaURL
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

// InstrumentConfig contains options for all instruments.
type InstrumentConfig struct {
	description string
	unit        unit.Unit
}

// NewInstrumentConfig returns a new InstrumentConfig with all opts applied.
func NewInstrumentConfig(opts ...InstrumentOption) InstrumentConfig {
	var config InstrumentConfig
	for _, o := range opts {
		config = o.applyInstrument(config)
	}
	return config
}

// Description returns the InstrumentConfig description.
func (c InstrumentConfig) Description() string {
	return c.description
}

// Unit returns the InstrumentConfig unit.
func (c InstrumentConfig) Unit() unit.Unit {
	return c.unit
}

// InstrumentOption applies options to all instrument configuration.
type InstrumentOption interface {
	ObservableOption
	applyInstrument(InstrumentConfig) InstrumentConfig
}

type descriptionOption string

func (o descriptionOption) applyInstrument(cfg InstrumentConfig) InstrumentConfig {
	cfg.description = string(o)
	return cfg
}

func (o descriptionOption) applyObservable(cfg ObservableConfig) ObservableConfig {
	cfg.InstrumentConfig.description = string(o)
	return cfg
}

// WithDescription sets the instrument description.
func WithDescription(desc string) InstrumentOption {
	return descriptionOption(desc)
}

type unitOption unit.Unit

func (o unitOption) applyInstrument(cfg InstrumentConfig) InstrumentConfig {
	cfg.unit = unit.Unit(o)
	return cfg
}

func (o unitOption) applyObservable(cfg ObservableConfig) ObservableConfig {
	cfg.InstrumentConfig.unit = unit.Unit(o)
	return cfg
}

// WithUnit sets the instrument unit.
func WithUnit(u unit.Unit) InstrumentOption {
	return WithUnit(u)
}

// ObservableConfig contains options for Observable instruments.
type ObservableConfig struct {
	InstrumentConfig

	callbacks []Callback
}

// NewObservableConfig returns a new ObservableConfig with all opts applied.
func NewObservableConfig(opts ...ObservableOption) ObservableConfig {
	var config ObservableConfig
	for _, o := range opts {
		config = o.applyObservable(config)
	}
	return config
}

// Callbacks returns the ObservableConfig callbacks.
func (c ObservableConfig) Callbacks() []Callback {
	return c.callbacks
}

// ObservableOption applies options to Observable instruments.
type ObservableOption interface {
	applyObservable(ObservableConfig) ObservableConfig
}

type callbackOption Callback

func (o callbackOption) applyInstrument(cfg InstrumentConfig) InstrumentConfig {
	return cfg
}

func (o callbackOption) applyObservable(cfg ObservableConfig) ObservableConfig {
	cfg.callbacks = append(cfg.callbacks, (Callback)(o))
	return cfg
}

// WithCallback adds callback to be called for an Observable instrument.
func WithCallback(callback Callback) ObservableOption {
	return WithCallback(callback)
}
