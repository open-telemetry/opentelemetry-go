package asyncfloat64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric2"
)

type Meter struct {
}

type Counter struct {
}

type UpDownCounter struct {
}

type Gauge struct {
}

type Instrument interface {
	metric.Instrument

	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
	Measure(x float64) metric.Measurement
}

var (
	_ Instrument = Counter{}
	_ Instrument = UpDownCounter{}
	_ Instrument = Gauge{}
)

func (m Meter) Counter(name string) (Counter, error) {
	return Counter{}, nil
}

func (m Meter) UpDownCounter(name string) (UpDownCounter, error) {
	return UpDownCounter{}, nil
}

func (m Meter) Gauge(name string) (Gauge, error) {
	return Gauge{}, nil
}

func (c Counter) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (g Gauge) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x float64) metric.Measurement {
	return metric.Measurement{}
}

func (u UpDownCounter) Measure(x float64) metric.Measurement {
	return metric.Measurement{}
}

func (g Gauge) Measure(x float64) metric.Measurement {
	return metric.Measurement{}
}
