package reader

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		defaults DefaultsFunc
	}

	Option func(*Config)

	DefaultsFunc func(sdkapi.InstrumentKind) (aggregation.Kind, aggregation.Temporality)

	Reader struct {
		config   Config
		exporter Exporter
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
		Instrument  sdkapi.Descriptor
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
		Produce() Metrics
	}
)

func WithDefaults(defaults DefaultsFunc) Option {
	return func(cfg *Config) {
		cfg.defaults = defaults
	}
}

func standardDefaults(ik sdkapi.InstrumentKind) (aggregation.Kind, aggregation.Temporality) {
	var ak aggregation.Kind
	switch ik {
	case sdkapi.HistogramInstrumentKind:
		ak = aggregation.HistogramKind
	case sdkapi.GaugeObserverInstrumentKind:
		ak = aggregation.LastValueKind
	case sdkapi.CounterInstrumentKind,
		sdkapi.UpDownCounterInstrumentKind,
		sdkapi.CounterObserverInstrumentKind,
		sdkapi.UpDownCounterObserverInstrumentKind:
		ak = aggregation.SumKind
	}
	return ak, aggregation.CumulativeTemporality
}

func NewConfig(opts ...Option) Config {
	cfg := Config{
		defaults: standardDefaults,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func New(config Config, exporter Exporter) *Reader {
	return &Reader{
		config:   config,
		exporter: exporter,
	}
}

func (r *Reader) Defaults() DefaultsFunc {
	return r.config.defaults
}

func (r *Reader) Exporter() Exporter {
	return r.exporter
}
