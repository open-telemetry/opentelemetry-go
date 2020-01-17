package parent

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
)

func getEffective(ctx context.Context) (core.SpanContext, bool) {
	if sctx := trace.SpanFromContext(ctx).SpanContext(); sctx.IsValid() {
		return sctx, false
	}
	return trace.RemoteContext(ctx), true
}

func GetContext(ctx, parent context.Context) (context.Context, core.SpanContext, bool) {
	if parent != nil {
		pctx, remote := getEffective(parent)
		return parent, pctx, remote
	}

	sctx, remote := getEffective(ctx)
	return ctx, sctx, remote
}
