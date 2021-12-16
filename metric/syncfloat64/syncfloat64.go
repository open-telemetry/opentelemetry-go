package syncfloat64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
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
	c.RecordOne(ctx, number.NewFloat64Number(x), attribute.Fingerprint(attrs...))
}

func (u UpDownCounter) Add(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	u.RecordOne(ctx, number.NewFloat64Number(x), attribute.Fingerprint(attrs...))
}

func (h Histogram) Record(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	h.RecordOne(ctx, number.NewFloat64Number(x), attribute.Fingerprint(attrs...))
}
