package asyncint64

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

type Gauge struct {
	sdkapi.Instrument
}

func (c Counter) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	c.Capture(ctx, number.NewInt64(x), attrs)
}

func (u UpDownCounter) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	u.Capture(ctx, number.NewInt64(x), attrs)
}

func (g Gauge) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	g.Capture(ctx, number.NewInt64(x), attrs)

}
