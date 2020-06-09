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

package metric

import (
	"go.opentelemetry.io/otel/api/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// Config contains some options for metrics of any kind.
type Config struct {
	// Description is an optional field describing the metric
	// instrument.
	Description string
	// Unit is an optional field describing the metric instrument.
	Unit unit.Unit
	// InstrumentationLibrary describes the library that provided
	// instrumentation.
	InstrumentationLibrary instrumentation.Library
}

// Option is an interface for applying metric options.
type Option interface {
	// Apply is used to set the Option value of a Config.
	Apply(*Config)
}

// Configure is a helper that applies all the options to a Config.
func Configure(opts []Option) Config {
	var config Config
	for _, o := range opts {
		o.Apply(&config)
	}
	return config
}

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return descriptionOption(desc)
}

type descriptionOption string

func (d descriptionOption) Apply(config *Config) {
	config.Description = string(d)
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return unitOption(unit)
}

type unitOption unit.Unit

func (u unitOption) Apply(config *Config) {
	config.Unit = unit.Unit(u)
}

// WithInstrumentationLibrary sets the library used to provided
// instrumentation. This is meant for use in `Provider` implementations that
// have not used `WrapMeterImpl`.  Implementations built using
// `WrapMeterImpl` have instrument descriptors taken care of through this
// package.
//
// This option will have no effect when supplied by the user.
// Provider implementations are expected to append this option after
// the user-supplied options when building instrument descriptors.
func WithInstrumentationLibrary(il instrumentation.Library) Option {
	return instrumentationLibraryOption(il)
}

type instrumentationLibraryOption instrumentation.Library

func (i instrumentationLibraryOption) Apply(config *Config) {
	config.InstrumentationLibrary = instrumentation.Library(i)
}

// Config contains options for a Meter.
type MeterConfig struct {
	// InstrumentationVersion is the version of the library providing
	// instrumentation.
	InstrumentationVersion string
}

// MeterOption is an interface for applying meter options.
type MeterOption interface {
	// Apply is used to set the MeterOption value of a MeterConfig.
	Apply(*MeterConfig)
}

// MeterConfigure is a helper that applies all the MeterOptions to a
// MeterConfig.
func MeterConfigure(opts []MeterOption) MeterConfig {
	var config MeterConfig
	for _, o := range opts {
		o.Apply(&config)
	}
	return config
}

// WithInstrumentationVersion sets the instrumentation version.
func WithInstrumentationVersion(version string) MeterOption {
	return instrumentationVersionOption(version)
}

type instrumentationVersionOption string

func (i instrumentationVersionOption) Apply(config *MeterConfig) {
	config.InstrumentationVersion = string(i)
}
