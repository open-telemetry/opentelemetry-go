package trace

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type remoteContextType struct{}

var remoteContextKey = &remoteContextType{}

// WithRemoteContext enters a core.SpanContext into a new Context.
func WithRemoteContext(ctx context.Context, sc core.SpanContext) context.Context {
	return context.WithValue(ctx, remoteContextKey, sc)
}

// RemoteContext gets the current core.SpanContext from a Context.
func RemoteContext(ctx context.Context) core.SpanContext {
	if sc, ok := ctx.Value(remoteContextKey).(core.SpanContext); ok {
		return sc
	}
	return core.EmptySpanContext()
}
