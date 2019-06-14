package metric

import (
	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type (
	Metric interface {
		Measure() core.Measure

		DefinitionID() core.EventID

		Type() MetricType
		Fields() []core.Key
		Err() error

		base() *baseMetric
	}

	MetricType int
)

const (
	Invalid MetricType = iota
	GaugeInt64
	GaugeFloat64
	DerivedGaugeInt64
	DerivedGaugeFloat64
	CumulativeInt64
	CumulativeFloat64
	DerivedCumulativeInt64
	DerivedCumulativeFloat64
)

type (
	Option func(*baseMetric, *[]tag.Option)
)

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(_ *baseMetric, to *[]tag.Option) {
		*to = append(*to, tag.WithDescription(desc))
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(_ *baseMetric, to *[]tag.Option) {
		*to = append(*to, tag.WithUnit(unit))
	}
}

// WithKeys applies the provided dimension keys.
func WithKeys(keys ...core.Key) Option {
	return func(bm *baseMetric, _ *[]tag.Option) {
		bm.keys = keys
	}
}
