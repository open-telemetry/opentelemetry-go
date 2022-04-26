package view

import (
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	// ClauseConfig contains each of the configurable aspects of a
	// single Views clause.
	ClauseConfig struct {
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

	// Config describes all the user-configurable aspects of an
	// export pipeline except for the terminal Exporter interface.
	// The configurable aspects are:
	// - Clauses in effect
	// - Defaults by instrument kind for:
	//   - Aggregation Kind
	//   - Aggregation Temporality
	//   - Aggregator configuration for int64, float64
	Config struct {
		Clauses  []ClauseConfig
		Defaults DefaultConfig
	}

	DefaultConfig struct {
		ByInstrumentKind [sdkinstrument.NumKinds]struct {
			Aggregation aggregation.Kind
			Temporality aggregation.Temporality
			Int64       aggregator.Config
			Float64     aggregator.Config
		}
	}

	// Option is a functional option pattern for generating ClauseConfigs.
	ClauseOption func(*ClauseConfig)

	// Option is a functional option pattern for generating Configs.
	Option func(*Config)

	Views struct {
		// Name of these views, used in error reporting
		Name string

		Config
	}
)

const (
	unsetInstrumentKind = sdkinstrument.Kind(-1)
	unsetNumberKind     = number.Kind(-1)
)

// TODO: call views.Validate() to check for:
// - empty (?)
// - duplicate name
// - invalid inst/number/aggregation kind
// - both instrument name and regexp
// - schemaURL or Version without library name
// - empty attribute keys
// - Name w/o SingleInst

// New configures the clauses and default settings of a Views.
func New(name string, opts ...Option) *Views {
	return &Views{
		Name:   name,
		Config: NewConfig(opts...),
	}
}

// NewConfig processes options and returns a reader Config.
func NewConfig(opts ...Option) Config {
	var cfg Config
	WithDefaultAggregationKindSelector(StandardAggregationKind)(&cfg)
	WithDefaultAggregationTemporalitySelector(StandardTemporality)(&cfg)
	WithDefaultAggregationConfigSelector(StandardConfig)(&cfg)
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// WithClause adds a clause to the Views configuration.
func WithClause(opts ...ClauseOption) Option {
	return func(cfg *Config) {
		clause := ClauseConfig{
			instrumentKind: unsetInstrumentKind,
			numberKind:     unsetNumberKind,
		}
		for _, opt := range opts {
			opt(&clause)
		}
		cfg.Clauses = append(cfg.Clauses, clause)
	}
}

// WithDefaultAggregationKindSelector configures the default
// aggregation.Kind to use with each kind of instrument.  This
// overwrites previous settings of the same option.
func WithDefaultAggregationKindSelector(d aggregation.KindSelector) Option {
	return func(cfg *Config) {
		for k := sdkinstrument.Kind(0); k < sdkinstrument.NumKinds; k++ {
			cfg.Defaults.ByInstrumentKind[k].Aggregation = d(k)
		}
	}
}

// WithDefaultAggregationTemporalitySelector configures the default
// aggregation.Temporality to use with each kind of instrument.  This
// overwrites previous settings of the same option.
func WithDefaultAggregationTemporalitySelector(d aggregation.TemporalitySelector) Option {
	return func(cfg *Config) {
		for k := sdkinstrument.Kind(0); k < sdkinstrument.NumKinds; k++ {
			cfg.Defaults.ByInstrumentKind[k].Temporality = d(k)
		}
	}
}

// WithDefaultAggregationConfigSelector configures the default
// aggregator.Config to use with each kind of instrument.  This
// overwrites previous settings of the same option.
func WithDefaultAggregationConfigSelector(d aggregator.ConfigSelector) Option {
	return func(cfg *Config) {
		for k := sdkinstrument.Kind(0); k < sdkinstrument.NumKinds; k++ {
			cfg.Defaults.ByInstrumentKind[k].Int64, cfg.Defaults.ByInstrumentKind[k].Float64 = d(k)
		}
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
