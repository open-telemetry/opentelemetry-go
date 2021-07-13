package asyncint64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric2/sdkapi"
)

type Counter struct {
}

type UpDownCounter struct {
}

type Gauge struct {
}

func (c Counter) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (u UpDownCounter) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (g Gauge) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
}

func (c Counter) Measure(x int64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}

func (u UpDownCounter) Measure(x int64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}

func (g Gauge) Measure(x int64) sdkapi.Measurement {
	return sdkapi.Measurement{}
}
