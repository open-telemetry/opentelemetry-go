package asyncfloat64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type Counter interface {
	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
}

type UpDownCounter interface {
	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
}

type Gauge interface {
	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
}
