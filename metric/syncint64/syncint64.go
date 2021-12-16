package syncint64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type Counter interface {
	Add(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}

type UpDownCounter interface {
	Add(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}

type Histogram interface {
	Record(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}
