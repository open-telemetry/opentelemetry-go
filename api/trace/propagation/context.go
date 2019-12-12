package propagation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type spanContextType struct{}

var spanContextKey = &spanContextType{}

// WithUpstreamContext enters a core.SpanContext into a new Context.
func WithUpstreamContext(ctx context.Context, sc core.SpanContext) context.Context {
	return context.WithValue(ctx, spanContextKey, sc)
}

// FromContext gets the current core.SpanContext from a Context.
func UpstreamContext(ctx context.Context) core.SpanContext {
	if sc, ok := ctx.Value(spanContextKey).(core.SpanContext); ok {
		return sc
	}
	return core.EmptySpanContext()
}
