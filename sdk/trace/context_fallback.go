// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"slices"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// These values match the experimental modes in sdk/trace/x.
const (
	contextFallbackParent = iota
	contextFallbackLink
)

// spanStartConfig overlays links without changing caller-owned options or
// reconstructing the other, potentially future, fields of trace.SpanConfig.
type spanStartConfig struct {
	trace.SpanConfig
	links []trace.Link
}

func (c *spanStartConfig) Links() []trace.Link {
	if c.links != nil {
		return c.links
	}
	return c.SpanConfig.Links()
}

// contextFallback is immutable after provider construction.
type contextFallback struct {
	parent     trace.SpanContext
	parentSpan trace.Span
	baggage    baggage.Baggage
	members    []baggage.Member
	mode       int
}

func newContextFallback(parent trace.SpanContext, bag baggage.Baggage, mode int) *contextFallback {
	if mode != contextFallbackParent && mode != contextFallbackLink {
		return nil
	}
	if !parent.IsValid() && bag.Len() == 0 {
		return nil
	}
	f := &contextFallback{parent: parent, baggage: bag, members: bag.Members(), mode: mode}
	if parent.IsValid() && mode == contextFallbackParent {
		// Reuse the immutable, non-recording wrapper rather than boxing a new
		// wrapper for each start. This creates no span identity or lifecycle.
		f.parentSpan = trace.SpanFromContext(trace.ContextWithSpanContext(context.Background(), parent))
	}
	return f
}

func (f *contextFallback) apply(ctx context.Context, config *spanStartConfig) context.Context {
	// Even a partial application context is intentional input. WithNewRoot is
	// the explicit opt-out when the context contains no span at all.
	if config.NewRoot() || !trace.SpanContextFromContext(ctx).Equal(trace.SpanContext{}) {
		return ctx
	}
	if f.parent.IsValid() {
		if f.mode == contextFallbackParent {
			ctx = trace.ContextWithSpan(ctx, f.parentSpan)
		} else {
			links := config.Links()
			if !slices.ContainsFunc(links, func(link trace.Link) bool {
				return link.SpanContext.TraceID() == f.parent.TraceID() &&
					link.SpanContext.SpanID() == f.parent.SpanID()
			}) {
				config.links = make([]trace.Link, len(links)+1)
				copy(config.links, links)
				config.links[len(links)] = trace.Link{SpanContext: f.parent}
			}
		}
	}
	if len(f.members) == 0 {
		return ctx
	}
	current := baggage.FromContext(ctx)
	if current.Len() == 0 {
		return baggage.ContextWithBaggage(ctx, f.baggage)
	}

	// Validate additions individually: truncating a combined Baggage could
	// otherwise discard application values when the W3C limits are exceeded.
	members := current.Members()
	changed := false
	for _, member := range f.members {
		if current.Member(member.Key()).Key() != "" {
			continue
		}
		members = append(members, member)
		merged, err := baggage.New(members...)
		if err != nil {
			members = members[:len(members)-1]
			continue
		}
		current = merged
		changed = true
	}
	if changed {
		ctx = baggage.ContextWithBaggage(ctx, current)
	}
	return ctx
}
