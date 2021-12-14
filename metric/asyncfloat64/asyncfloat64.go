package asyncfloat64

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

type Gauge struct {
	sdkapi.Instrument
}

func (c Counter) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	c.Instrument.RecordOne(ctx, number.NewFloat64Number(x), attrs)
}

func (u UpDownCounter) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	u.Instrument.RecordOne(ctx, number.NewFloat64Number(x), attrs)
}

func (g Gauge) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	g.Instrument.RecordOne(ctx, number.NewFloat64Number(x), attrs)
}

func (c Counter) Measure(x float64) sdkapi.Measurement {
	return sdkapi.NewMeasurement(c.Instrument, number.NewFloat64Number(x))
}

func (u UpDownCounter) Measure(x float64) sdkapi.Measurement {
	return sdkapi.NewMeasurement(u.Instrument, number.NewFloat64Number(x))
}

func (g Gauge) Measure(x float64) sdkapi.Measurement {
	return sdkapi.NewMeasurement(g.Instrument, number.NewFloat64Number(x))
}
