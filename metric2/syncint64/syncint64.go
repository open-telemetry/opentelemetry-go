package syncint64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric2/sdkapi"
)

type Counter struct {
}

type UpDownCounter struct {
}

type Histogram struct {
}

func (c Counter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (h Histogram) Record(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x int64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}

func (u UpDownCounter) Measure(x int64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}

func (h Histogram) Measure(x int64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}
