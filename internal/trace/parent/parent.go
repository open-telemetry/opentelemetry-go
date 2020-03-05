package parent

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
)

func GetSpanContextAndLinks(ctx context.Context, ignoreContext bool) (core.SpanContext, bool, []trace.Link) {
	lsctx := trace.SpanFromContext(ctx).SpanContext()
	rsctx := trace.RemoteSpanContextFromContext(ctx)

	if ignoreContext {
		links := addLinkIfValid(nil, lsctx, "current")
		links = addLinkIfValid(links, rsctx, "remote")

		return core.EmptySpanContext(), false, links
	}
	if lsctx.IsValid() {
		return lsctx, false, nil
	}
	if rsctx.IsValid() {
		return rsctx, true, nil
	}
	return core.EmptySpanContext(), false, nil
}

func addLinkIfValid(links []trace.Link, sc core.SpanContext, kind string) []trace.Link {
	if !sc.IsValid() {
		return links
	}
	return append(links, trace.Link{
		SpanContext: sc,
		Attributes: []core.KeyValue{
			key.String("ignored-on-demand", kind),
		},
	})
}
