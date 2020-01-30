package correlation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type correlationsType struct{}

var correlationsKey = &correlationsType{}

// WithMap enters a Map into a new Context.
func WithMap(ctx context.Context, m Map) context.Context {
	return context.WithValue(ctx, correlationsKey, m)
}

// WithMap enters a key:value set into a new Context.
func NewContext(ctx context.Context, keyvalues ...core.KeyValue) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(MapUpdate{
		MultiKV: keyvalues,
	}))
}

// FromContext gets the current Map from a Context.
func FromContext(ctx context.Context) Map {
	if m, ok := ctx.Value(correlationsKey).(Map); ok {
		return m
	}
	return NewEmptyMap()
}
