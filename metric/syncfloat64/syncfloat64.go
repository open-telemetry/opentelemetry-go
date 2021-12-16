package syncfloat64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
)

type Counter struct {
	sdkapi.Instrument
}

type UpDownCounter struct {
	sdkapi.Instrument
}

type Histogram struct {
	sdkapi.Instrument
}

func (c Counter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	c.Capture(ctx, number.NewFloat64(x), attrs)
}

func (u UpDownCounter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	u.Capture(ctx, number.NewFloat64(x), attrs)
}

func (h Histogram) Record(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	h.Capture(ctx, number.NewFloat64(x), attrs)
}
