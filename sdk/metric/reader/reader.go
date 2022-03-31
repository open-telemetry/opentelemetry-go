package reader

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// Config includes Reader-specific options that the user
	// configures, which are generally per-instrument-kind defaults.
	Config struct {
		DefaultAggregationKindFunc
		DefaultAggregationTemporalityFunc
		DefaultAggregationConfigFunc
	}

	// Option is a functional option pattern for generating Configs.
	Option func(*Config)

	// DefaultAggregationKindFunc is a per-instrument-kind aggregation.Kind choice.
	DefaultAggregationKindFunc func(sdkinstrument.Kind) aggregation.Kind

	// DefaultAggregationKindFunc is a per-instrument-kind aggregation.Temporality choice.
	DefaultAggregationTemporalityFunc func(sdkinstrument.Kind) aggregation.Temporality

	// DefaultAggregationKindFunc is a per-instrument-kind, per-number-kind aggregation.Config choice.
	DefaultAggregationConfigFunc func(sdkinstrument.Kind) (int64Config, float64Config aggregation.Config)

	// Reader represents the connection between an Exporter and
	// the MeterProvider.  Readers give the internal View compiler
	// all the necessary information to construct a metrics pipeline.
	Reader struct {
		// exporter is the output terminal of a metrics pipeline.
		exporter Exporter

		// These four arrays are the evaluted
		// per-instrument-kind choices described in the
		// Config.

		defAggr   [sdkinstrument.NumKinds]aggregation.Kind
		defTempo  [sdkinstrument.NumKinds]aggregation.Temporality
		defI64Cfg [sdkinstrument.NumKinds]aggregation.Config
		defF64Cfg [sdkinstrument.NumKinds]aggregation.Config
	}

	// Sequence provides the three relevant timestamps to the Reader
	Sequence struct {
		// Start is the time when the MeterProvider was initialized.
		Start time.Time
		// Last is the time when the previous collection
		// happened.  If there was no previous collection,
		// this will match Start.
		Last time.Time
		// Now is the moment the current collection began.  This value
		// will be used as the subsequent value for Last.
		Now time.Time
	}

	// Metrics is the result of a single collection.
	//
	// This struct supports re-use of the nested memory structure
	// underlying Scopes such that repeated calls Produce will not
	// reallocate the same quantity of memory over and over.
	//
	// To re-use the memory from a previous Metrics value, pass a
	// pointer pointer to Produce(). This is safe for push-based
	// exporters that perform sequential collection.
	Metrics struct {
		// Resource is the MeterProvider's configured Resource.
		Resource *resource.Resource
		Scopes   []Scope
	}

	Scope struct {
		Library     instrumentation.Library
		Instruments []Instrument
	}

	Instrument struct {
		Descriptor  sdkinstrument.Descriptor
		Temporality aggregation.Temporality
		Series      []Series
	}

	Series struct {
		Attributes  attribute.Set
		Aggregation aggregation.Aggregation
		Start       time.Time
		End         time.Time
	}

	Exporter interface {
		Registeree

		Flush(context.Context) error
		Shutdown(context.Context) error
	}

	Registeree interface {
		Register(Producer)
	}

	Producer interface {
		Produce(in *Metrics) Metrics
	}
)

func WithDefaultAggregationKindFunc(d DefaultAggregationKindFunc) Option {
	return func(cfg *Config) {
		cfg.DefaultAggregationKindFunc = d
	}
}

func WithDefaultAggregationTemporalityFunc(d DefaultAggregationTemporalityFunc) Option {
	return func(cfg *Config) {
		cfg.DefaultAggregationTemporalityFunc = d
	}
}

func WithDefaultAggregationConfigFunc(d DefaultAggregationConfigFunc) Option {
	return func(cfg *Config) {
		cfg.DefaultAggregationConfigFunc = d
	}
}

func standardAggregation(ik sdkinstrument.Kind) aggregation.Kind {
	switch ik {
	case sdkinstrument.HistogramKind:
		return aggregation.HistogramKind
	case sdkinstrument.GaugeObserverKind:
		return aggregation.GaugeKind
	default:
		return aggregation.SumKind
	}
}

func standardTemporality(ik sdkinstrument.Kind) aggregation.Temporality {
	return aggregation.CumulativeTemporality
}

func standardConfig(ik sdkinstrument.Kind) (ints, floats aggregation.Config) {
	return aggregation.Config{}, aggregation.Config{}
}

func New(exporter Exporter, opts ...Option) *Reader {
	cfg := Config{
		DefaultAggregationKindFunc:        standardAggregation,
		DefaultAggregationTemporalityFunc: standardTemporality,
		DefaultAggregationConfigFunc:      standardConfig,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	r := &Reader{
		exporter: exporter,
	}
	for i := sdkinstrument.Kind(0); i < sdkinstrument.NumKinds; i++ {
		r.defAggr[i] = cfg.DefaultAggregationKindFunc(i)
		r.defTempo[i] = cfg.DefaultAggregationTemporalityFunc(i)
		r.defI64Cfg[i], r.defF64Cfg[i] = cfg.DefaultAggregationConfigFunc(i)
	}
	return r
}

func (r *Reader) DefaultAggregation(k sdkinstrument.Kind) aggregation.Kind {
	return r.defAggr[k]
}

func (r *Reader) DefaultTemporality(k sdkinstrument.Kind) aggregation.Temporality {
	return r.defTempo[k]
}

func (r *Reader) DefaultAggregationConfig(k sdkinstrument.Kind, nk number.Kind) aggregation.Config {
	if nk == number.Int64Kind {
		return r.defI64Cfg[k]
	}
	return r.defF64Cfg[k]
}

func (r *Reader) Exporter() Exporter {
	return r.exporter
}
