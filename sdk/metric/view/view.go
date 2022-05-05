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

package view // import "go.opentelemetry.io/otel/sdk/metric/view"

import (
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// Config contains configuration options for a view.
//
// The configurable aspects are:
// - Clauses in effect
// - Defaults by instrument kind for:
//   - Aggregation Kind
//   - Aggregation Temporality
//   - Aggregator configuration for int64, float64
type Config struct {
	Clauses  []ClauseConfig
	Defaults DefaultConfig
}

// ClauseConfig contains each of the configurable aspects of a
// single Views clause.
type ClauseConfig struct {
	// Matchers for the instrument
	instrumentName       string
	instrumentNameRegexp *regexp.Regexp
	instrumentKind       sdkinstrument.Kind
	numberKind           number.Kind
	library              instrumentation.Library

	// Properties of the view
	keys        []attribute.Key // nil implies all keys, []attribute.Key{} implies none
	name        string
	description string
	aggregation aggregation.Kind
	temporality aggregation.Temporality
	acfg        aggregator.Config
}

// DefaultConfig contains configurable aspects that apply to all
// instruments in a View.
type DefaultConfig struct {
	ByInstrumentKind [sdkinstrument.NumKinds]struct {
		Aggregation aggregation.Kind
		Temporality aggregation.Temporality
		Int64       aggregator.Config
		Float64     aggregator.Config
	}
}

// Aggregation returns the default aggregation.Kind for each instrument kind.
func (d *DefaultConfig) Aggregation(k sdkinstrument.Kind) aggregation.Kind {
	return d.ByInstrumentKind[k].Aggregation
}

// DefaultTemporality returns the default aggregation.Temporality for each instrument kind.
func (c *DefaultConfig) Temporality(k sdkinstrument.Kind) aggregation.Temporality {
	return c.ByInstrumentKind[k].Temporality
}

// AggregationConfig returns the default aggregation.Temporality for each instrument kind.
func (c *DefaultConfig) AggregationConfig(k sdkinstrument.Kind, nk number.Kind) aggregator.Config {
	if nk == number.Int64Kind {
		return c.ByInstrumentKind[k].Int64
	}
	return c.ByInstrumentKind[k].Float64
}

// WithClause adds a clause to the Views configuration.
func WithClause(options ...ClauseOption) Option {
	return optionFunction(func(cfg Config) Config {
		clause := ClauseConfig{
			instrumentKind: unsetInstrumentKind,
			numberKind:     unsetNumberKind,
		}
		for _, option := range options {
			clause = option.apply(clause)
		}
		cfg.Clauses = append(cfg.Clauses, clause)
		return cfg
	})
}

// WithDefaultAggregationKindSelector configures the default
// aggregation.Kind to use with each kind of instrument.  This
// overwrites previous settings of the same option.
func WithDefaultAggregationKindSelector(d aggregation.KindSelector) Option {
	return optionFunction(func(cfg Config) Config {
		for k := sdkinstrument.Kind(0); k < sdkinstrument.NumKinds; k++ {
			cfg.Defaults.ByInstrumentKind[k].Aggregation = d(k)
		}
		return cfg
	})
}

// WithDefaultAggregationTemporalitySelector configures the default
// aggregation.Temporality to use with each kind of instrument.  This
// overwrites previous settings of the same option.
func WithDefaultAggregationTemporalitySelector(d aggregation.TemporalitySelector) Option {
	return optionFunction(func(cfg Config) Config {
		for k := sdkinstrument.Kind(0); k < sdkinstrument.NumKinds; k++ {
			cfg.Defaults.ByInstrumentKind[k].Temporality = d(k)
		}
		return cfg
	})
}

// WithDefaultAggregationConfigSelector configures the default
// aggregator.Config to use with each kind of instrument.  This
// overwrites previous settings of the same option.
func WithDefaultAggregationConfigSelector(d aggregator.ConfigSelector) Option {
	return optionFunction(func(cfg Config) Config {
		for k := sdkinstrument.Kind(0); k < sdkinstrument.NumKinds; k++ {
			cfg.Defaults.ByInstrumentKind[k].Int64, cfg.Defaults.ByInstrumentKind[k].Float64 = d(k)
		}
		return cfg
	})
}

// Option applies a configuration option value to a view Config.
type Option interface {
	apply(Config) Config
}

// optionFunction makes a functional Option out of a function object.
type optionFunction func(cfg Config) Config

// apply implements Option.
func (of optionFunction) apply(in Config) Config {
	return of(in)
}

// NewConfig returns a new and configured view Config.
func NewConfig(options ...Option) Config {
	standard := []Option{
		WithDefaultAggregationKindSelector(StandardAggregationKind),
		WithDefaultAggregationTemporalitySelector(StandardTemporality),
		WithDefaultAggregationConfigSelector(StandardConfig),
	}
	var cfg Config
	for _, option := range append(standard, options...) {
		cfg = option.apply(cfg)
	}
	return cfg
}
