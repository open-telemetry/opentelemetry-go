package asyncint64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type Counter interface {
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}

type UpDownCounter interface {
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}

type Gauge interface {
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}
