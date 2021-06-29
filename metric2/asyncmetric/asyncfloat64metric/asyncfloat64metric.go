package asyncfloat64metric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric2/batch"
)

type Meter struct {
}

type Counter struct {
}

type UpDownCounter struct {
}

type Gauge struct {
}

func (m Meter) Counter(name string) (Counter, error) {
	return Counter{}, nil
}

func (m Meter) UpDownCounter(name string) (UpDownCounter, error) {
	return UpDownCounter{}, nil
}

func (m Meter) Gauge(name string) (Gauge, error) {
	return Gauge{}, nil
}

func (c Counter) Set(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Set(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (g Gauge) Set(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x int64) batch.Measurement {
	return batch.Measurement{}
}

func (u UpDownCounter) Measure(x int64) batch.Measurement {
	return batch.Measurement{}
}

func (g Gauge) Measure(x int64) batch.Measurement {
	return batch.Measurement{}
}
