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

import "go.opentelemetry.io/otel/api/unit"

// Config contains some options for metrics of any kind.
type Config struct {
	// Description is an optional field describing the metric
	// instrument.
	Description string
	// Unit is an optional field describing the metric instrument.
	Unit unit.Unit
	// LibraryName is the name given to the Meter that created
	// this instrument.  See `Provider`.
	LibraryName string
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

// WithLibraryName applies provided library name.  This is meant for
// use in `Provider` implementations that have not used
// `WrapMeterImpl`.  Implementations built using `WrapMeterImpl` have
// instrument descriptors taken care of through this package.
//
// This option will have no effect when supplied by the user.
// Provider implementations are expected to append this option after
// the user-supplied options when building instrument descriptors.
func WithLibraryName(name string) Option {
	return libraryNameOption(name)
}

type libraryNameOption string

func (r libraryNameOption) Apply(config *Config) {
	config.LibraryName = string(r)
}
