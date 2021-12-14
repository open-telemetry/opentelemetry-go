package syncint64

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

func (c Counter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	c.RecordOne(ctx, number.NewInt64Number(x), attrs)
}

func (u UpDownCounter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	u.RecordOne(ctx, number.NewInt64Number(x), attrs)
}

func (h Histogram) Record(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	h.RecordOne(ctx, number.NewInt64Number(x), attrs)
}

func (c Counter) Measure(x int64) sdkapi.Measurement {
	return sdkapi.NewMeasurement(c.Instrument, number.NewInt64Number(x))
}

func (u UpDownCounter) Measure(x int64) sdkapi.Measurement {
	return sdkapi.NewMeasurement(u.Instrument, number.NewInt64Number(x))
}

func (h Histogram) Measure(x int64) sdkapi.Measurement {
	return sdkapi.NewMeasurement(h.Instrument, number.NewInt64Number(x))
}
