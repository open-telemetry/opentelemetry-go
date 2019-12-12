package parent

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/trace/propagation"
)

func getEffective(ctx context.Context) (core.SpanContext, bool) {
	if ctx == nil {
		return core.EmptySpanContext(), false
	}
	rctx := propagation.RemoteContext(ctx)
	sctx := trace.SpanFromContext(ctx).SpanContext()

	if rctx.IsValid() && sctx.IsValid() && rctx.TraceID == sctx.TraceID {
		return sctx, false
	}
	if rctx.IsValid() {
		return rctx, true
	}
	return sctx, false
}

func GetContext(ctx, parent context.Context) (context.Context, core.SpanContext, bool) {
	if pctx, remote := getEffective(parent); pctx.IsValid() {
		return parent, pctx, remote
	}

	sctx, remote := getEffective(ctx)
	return ctx, sctx, remote
}
