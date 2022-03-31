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
	Config struct {
		defaults DefaultsFunc
	}

	Option func(*Config)

	DefaultsFunc func(sdkinstrument.Kind) (aggKind aggregation.Kind, aggTempo aggregation.Temporality, int64Config, float64Config aggregation.Config)

	Reader struct {
		defAggr   [sdkinstrument.NumKinds]aggregation.Kind
		defTempo  [sdkinstrument.NumKinds]aggregation.Temporality
		defI64Cfg [sdkinstrument.NumKinds]aggregation.Config
		defF64Cfg [sdkinstrument.NumKinds]aggregation.Config
		exporter  Exporter
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

func standardDefaults(ik sdkinstrument.Kind) (aggregation.Kind, aggregation.Temporality, aggregation.Config, aggregation.Config) {
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
	return ak, aggregation.CumulativeTemporality, aggregation.Config{}, aggregation.Config{}
}

func New(exporter Exporter, opts ...Option) *Reader {
	cfg := Config{
		defaults: standardDefaults,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	r := &Reader{
		exporter: exporter,
	}
	for i := sdkinstrument.Kind(0); i < sdkinstrument.NumKinds; i++ {
		r.defAggr[i], r.defTempo[i], r.defI64Cfg[i], r.defF64Cfg[i] = cfg.defaults(i)
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
