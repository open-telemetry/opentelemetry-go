package syncfloat64

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

func (c Counter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (h Histogram) Record(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x float64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}

func (u UpDownCounter) Measure(x float64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}

func (h Histogram) Measure(x float64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}
