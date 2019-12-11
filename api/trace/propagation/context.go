package propagation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type spanContextType struct{}

var (
	spanContextKey = &spanContextType{}
)

func WithSpanContext(ctx context.Context, sc core.SpanContext) context.Context {
	return context.WithValue(ctx, spanContextKey, sc)
}

func FromContext(ctx context.Context) core.SpanContext {
	if sc, ok := ctx.Value(spanContextKey).(core.SpanContext); ok {
		return sc
	}
	return core.EmptySpanContext()
}
