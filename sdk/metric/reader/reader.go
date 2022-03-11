package reader

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/views"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		views    []views.View
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
		ByScope  []ScopeMetrics
	}

	ScopeMetrics struct {
		Library      instrumentation.Library
		ByInstrument []InstrumentMetrics
	}

	InstrumentMetrics struct {
		Instrument   *sdkapi.Descriptor
		Temporality  aggregation.Temporality
		ByAttributes []AttributeMetrics
	}

	AttributeMetrics struct {
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
		Produce(*Metrics) Metrics
	}
)

func WithView(view views.View) Option {
	return func(cfg *Config) {
		cfg.views = append(cfg.views, view)
	}
}

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

func (r *Reader) Views() []views.View {
	return r.config.views
}

func (r *Reader) Defaults() DefaultsFunc {
	return r.config.defaults
}

func (r *Reader) Exporter() Exporter {
	return r.exporter
}
