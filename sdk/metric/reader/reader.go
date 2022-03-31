package reader

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		defaults DefaultsFunc
	}

	Option func(*Config)

	DefaultsFunc func(sdkinstrument.Kind) (aggregation.Kind, aggregation.Temporality)

	Reader struct {
		config   Config
		exporter Exporter
	}

	Sequence struct {
		Start time.Time
		Last  time.Time
		Now   time.Time
	}

	Metrics struct {
		Resource *resource.Resource
		Scopes   []Scope
	}

	Scope struct {
		Library     instrumentation.Library
		Instruments []Instrument
	}

	Instrument struct {
		Instrument  sdkinstrument.Descriptor
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

func WithDefaults(defaults DefaultsFunc) Option {
	return func(cfg *Config) {
		cfg.defaults = defaults
	}
}

func standardDefaults(ik sdkinstrument.Kind) (aggregation.Kind, aggregation.Temporality) {
	var ak aggregation.Kind
	switch ik {
	case sdkinstrument.HistogramKind:
		ak = aggregation.HistogramKind
	case sdkinstrument.GaugeObserverKind:
		ak = aggregation.GaugeKind
	case sdkinstrument.CounterKind,
		sdkinstrument.UpDownCounterKind,
		sdkinstrument.CounterObserverKind,
		sdkinstrument.UpDownCounterObserverKind:
		ak = aggregation.SumKind
	}
	return ak, aggregation.CumulativeTemporality
}

func New(exporter Exporter, opts ...Option) *Reader {
	cfg := Config{
		defaults: standardDefaults,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Reader{
		config:   cfg,
		exporter: exporter,
	}
}

func (r *Reader) Defaults() DefaultsFunc {
	return r.config.defaults
}

func (r *Reader) Exporter() Exporter {
	return r.exporter
}
