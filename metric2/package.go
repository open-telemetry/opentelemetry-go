package metric2

import (
	"go.opentelemetry.io/otel/metric2/asyncmetric"
	"go.opentelemetry.io/otel/metric2/float64metric"
	"go.opentelemetry.io/otel/metric2/int64metric"
)

type MeterProvider interface {
	Meter(instrumentationName string /*, opts ...MeterOption*/) Meter
}

type Meter struct {
}

func (m Meter) Integer() int64metric.Meter {
	return int64metric.Meter{}
}

func (m Meter) FloatingPoint() float64metric.Meter {
	return float64metric.Meter{}
}

func (m Meter) Asynchronous() asyncmetric.Meter {
	return asyncmetric.Meter{}
}
