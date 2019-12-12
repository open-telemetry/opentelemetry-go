package parent

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/trace/propagation"
)

func getEffective(ctx context.Context) core.SpanContext {
	if ctx == nil {
		return core.EmptySpanContext()
	}
	rctx := propagation.RemoteContext(ctx)
	sctx := trace.SpanFromContext(ctx).SpanContext()

	if rctx.IsValid() && sctx.IsValid() && rctx.TraceID == sctx.TraceID {
		return sctx
	}
	if rctx.IsValid() {
		return rctx
	}
	return sctx
}

func GetContext(ctx, parent context.Context) (context.Context, core.SpanContext, bool) {
	pctx := getEffective(parent)
	sctx := getEffective(ctx)

	if pctx.IsValid() {
		return parent, pctx, true
	}

	return ctx, sctx, false
}
