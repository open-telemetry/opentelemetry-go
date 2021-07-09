package syncint64metric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric2/batch"
)

// TODO instrument options

type Meter struct {
}

type Counter struct {
}

type UpDownCounter struct {
}

type Histogram struct {
}

func (m Meter) Counter(name string) (Counter, error) {
	return Counter{}, nil
}

func (m Meter) UpDownCounter(name string) (UpDownCounter, error) {
	return UpDownCounter{}, nil
}

func (m Meter) Histogram(name string) (Histogram, error) {
	return Histogram{}, nil
}

func (c Counter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (h Histogram) Record(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x int64) batch.Measurement {
	return batch.Measurement{}
}

func (u UpDownCounter) Measure(x int64) batch.Measurement {
	return batch.Measurement{}
}

func (h Histogram) Measure(x int64) batch.Measurement {
	return batch.Measurement{}
}
