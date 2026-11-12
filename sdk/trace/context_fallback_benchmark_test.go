// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/baggage"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Keep this option independent of sdk/trace/x so the same benchmark can also
// measure SDK revisions that ignore experimental context-fallback options.
type benchmarkContextFallbackOption struct {
	sdktrace.TracerProviderOption
	spanContext trace.SpanContext
	baggage     baggage.Baggage
	mode        int
}

func (benchmarkContextFallbackOption) Experimental() {}

func (o benchmarkContextFallbackOption) ContextFallback() (trace.SpanContext, baggage.Baggage, int) {
	return o.spanContext, o.baggage, o.mode
}

func BenchmarkContextFallback(b *testing.B) {
	state, err := trace.ParseTraceState("vendor=value")
	if err != nil {
		b.Fatal(err)
	}
	incoming := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1},
		SpanID:     trace.SpanID{2},
		TraceFlags: trace.FlagsSampled,
		TraceState: state,
		Remote:     true,
	})
	otherLink := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{3},
		SpanID:  trace.SpanID{4},
		Remote:  true,
	})
	incomingBaggage := benchmarkContextFallbackBaggage(b, "incoming", "environment")
	localBaggage := benchmarkContextFallbackBaggage(b, "local", "application")
	parentOption := benchmarkContextFallbackOption{spanContext: incoming}
	linkOption := benchmarkContextFallbackOption{spanContext: incoming, mode: 1}
	withLinks := []trace.SpanStartOption{trace.WithLinks(trace.Link{SpanContext: otherLink})}
	withNewRoot := []trace.SpanStartOption{trace.WithNewRoot()}

	for _, sampling := range []struct {
		name    string
		sampler sdktrace.Sampler
	}{
		{name: "sampled", sampler: sdktrace.AlwaysSample()},
		{name: "dropped", sampler: sdktrace.NeverSample()},
	} {
		b.Run(sampling.name, func(b *testing.B) {
			parentProvider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sampling.sampler))
			b.Cleanup(func() {
				if err := parentProvider.Shutdown(b.Context()); err != nil {
					b.Error(err)
				}
			})
			activeCtx, parent := parentProvider.Tracer("benchmark-parent").Start(b.Context(), "parent")
			b.Cleanup(func() { parent.End() })
			baggageCtx := baggage.ContextWithBaggage(b.Context(), localBaggage)

			for _, test := range []struct {
				name    string
				option  sdktrace.TracerProviderOption
				ctx     context.Context
				options []trace.SpanStartOption
			}{
				{name: "disabled"},
				{name: "none", option: benchmarkContextFallbackOption{spanContext: incoming, mode: 2}},
				{name: "empty", option: benchmarkContextFallbackOption{}},
				{name: "parent", option: parentOption},
				{name: "disabled-active", ctx: activeCtx},
				{name: "parent-active", option: parentOption, ctx: activeCtx},
				{name: "disabled-explicit-root", options: withNewRoot},
				{name: "explicit-root", option: parentOption, options: withNewRoot},
				{name: "link", option: linkOption},
				{name: "disabled-existing-links", options: withLinks},
				{name: "link-existing", option: linkOption, options: withLinks},
				{name: "disabled-baggage", ctx: baggageCtx},
				{
					name:   "baggage",
					option: benchmarkContextFallbackOption{spanContext: incoming, baggage: incomingBaggage},
					ctx:    baggageCtx,
				},
			} {
				b.Run(test.name, func(b *testing.B) {
					options := []sdktrace.TracerProviderOption{sdktrace.WithSampler(sampling.sampler)}
					if test.option != nil {
						options = append(options, test.option)
					}
					provider := sdktrace.NewTracerProvider(options...)
					b.Cleanup(func() {
						if err := provider.Shutdown(b.Context()); err != nil {
							b.Error(err)
						}
					})
					tracer := provider.Tracer("context-fallback-benchmark")
					ctx := test.ctx
					if ctx == nil {
						ctx = b.Context()
					}
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, span := tracer.Start(ctx, "work", test.options...)
						span.End()
					}
				})
			}
		})
	}
}

func benchmarkContextFallbackBaggage(b *testing.B, key, value string) baggage.Baggage {
	b.Helper()
	member, err := baggage.NewMember(key, value)
	if err != nil {
		b.Fatal(err)
	}
	bag, err := baggage.New(member)
	if err != nil {
		b.Fatal(err)
	}
	return bag
}
