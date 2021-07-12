package syncfloat64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric2"
)

type Builder struct {
}

type Counter struct {
}

type UpDownCounter struct {
}

type Histogram struct {
}

type Instrument interface {
	metric.Instrument

	Measure(x float64) metric.Measurement
}

var (
	_ Instrument = Counter{}
	_ Instrument = UpDownCounter{}
	_ Instrument = Histogram{}
)

func (m Builder) Counter(name string) (Counter, error) {
	return Counter{}, nil
}

func (m Builder) UpDownCounter(name string) (UpDownCounter, error) {
	return UpDownCounter{}, nil
}

func (m Builder) Histogram(name string) (Histogram, error) {
	return Histogram{}, nil
}

func (c Counter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (h Histogram) Record(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x float64) metric.Measurement {
	return metric.Measurement{}
}

func (u UpDownCounter) Measure(x float64) metric.Measurement {
	return metric.Measurement{}
}

func (h Histogram) Measure(x float64) metric.Measurement {
	return metric.Measurement{}
}
