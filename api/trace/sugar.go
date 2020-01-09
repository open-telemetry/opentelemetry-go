package trace

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/internal"
)

type provider interface {
	Tracer() Tracer
}

func getTracer(ctx context.Context) Tracer {
	if p, ok := internal.ScopeImpl(ctx).(provider); ok {
		return p.Tracer()
	}

	if g, ok := (*atomic.Value)(atomic.LoadPointer(&internal.GlobalScope)).Load().(provider); ok {
		return g.Tracer()
	}

	return NoopTracer{}
}

func Start(ctx context.Context, spanName string, opts ...StartOption) (context.Context, Span) {
	return getTracer(ctx).Start(ctx, spanName, opts...)
}

func WithSpan(ctx context.Context, spanName string, fn func(ctx context.Context) error) error {
	return getTracer(ctx).WithSpan(ctx, spanName, fn)
}
