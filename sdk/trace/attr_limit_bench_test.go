// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// BenchmarkAttributeLimitNormalization compares value-length truncation with
// attribute count limits.
//
// retained keeps every attribute, so a finite length limit pays to truncate
// it. count=0 and discard-suffix drop attributes before that work. The
// discarded value is a 1,024-element string slice, which truncation would
// otherwise copy. SetAttributes covers deferred deduplication, while the span
// is under its attribute cap, and eager deduplication once new keys would
// exceed that cap.
func BenchmarkAttributeLimitNormalization(b *testing.B) {
	long := strings.Repeat("abcdefghij", 5) // 50 bytes, over a limit of 3
	elements := make([]string, 1024)
	for i := range elements {
		elements[i] = long
	}
	sliceAttr := attribute.StringSlice("arr", elements)

	retained := []attribute.KeyValue{
		attribute.String("s0", long),
		attribute.String("s1", "ab"),
		attribute.Bool("ok", true),
		sliceAttr,
	}
	// The suffix is the oversized string slice. count=1 keeps s0 only.
	discardSuffix := []attribute.KeyValue{
		attribute.String("keep", long),
		sliceAttr,
	}

	smallElements := elements[:32]
	deferred := []attribute.KeyValue{
		attribute.String("s0", long),
		attribute.String("s1", "ab"),
		attribute.Bool("ok", true),
		attribute.Int("n", 7),
		attribute.StringSlice("arr", smallElements),
	}
	eager := make([]attribute.KeyValue, 0, 34)
	eager = append(eager, attribute.StringSlice("kept", smallElements))
	for i := range 32 {
		eager = append(eager, attribute.String(fmt.Sprintf("k%d", i), long))
	}
	// Duplicate of the retained key, so the over-cap path normalizes an update
	// after it has started dropping new keys.
	eager = append(eager, attribute.StringSlice("kept", []string{"ab", "cd"}))

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})

	lengths := []lengthLimitCase{
		{name: "length=finite", limit: 3},
		{name: "length=unlimited", limit: -1},
	}

	b.Run("AddEvent", func(b *testing.B) {
		runLengthModes(b, lengths, retained, discardSuffix, benchAddEvent)
	})
	b.Run("AddLink", func(b *testing.B) {
		runLengthModes(b, lengths, retained, discardSuffix, func(
			b *testing.B,
			lengthLimit, countLimit int,
			attrs []attribute.KeyValue,
		) {
			benchAddLink(b, lengthLimit, countLimit, attrs, sc)
		})
	})

	b.Run("SetAttributes", func(b *testing.B) {
		for _, length := range lengths {
			b.Run(length.name, func(b *testing.B) {
				b.Run("deferred", func(b *testing.B) {
					benchSetAttributes(b, length.limit, 128, deferred)
				})
				b.Run("eager", func(b *testing.B) {
					benchSetAttributes(b, length.limit, 1, eager)
				})
			})
		}
	})
}

type lengthLimitCase struct {
	name  string
	limit int
}

func runLengthModes(
	b *testing.B,
	lengths []lengthLimitCase,
	retained, discard []attribute.KeyValue,
	run func(*testing.B, int, int, []attribute.KeyValue),
) {
	b.Helper()
	for _, length := range lengths {
		b.Run(length.name, func(b *testing.B) {
			b.Run("retained", func(b *testing.B) {
				run(b, length.limit, -1, retained)
			})
			b.Run("count=0", func(b *testing.B) {
				run(b, length.limit, 0, retained)
			})
			b.Run("discard-suffix", func(b *testing.B) {
				run(b, length.limit, 1, discard)
			})
		})
	}
}

func startAttrLimitSpan(b *testing.B, lengthLimit, countLimit int) trace.Span {
	b.Helper()

	limits := sdktrace.NewSpanLimits()
	limits.AttributeValueLengthLimit = lengthLimit
	limits.AttributePerEventCountLimit = countLimit
	limits.AttributePerLinkCountLimit = countLimit
	// Hold the event and link queues at one element so steady state does not
	// allocate. Those limits are independent of the attribute count limit.
	limits.EventCountLimit = 1
	limits.LinkCountLimit = 1

	tp := sdktrace.NewTracerProvider(sdktrace.WithRawSpanLimits(limits))
	b.Cleanup(func() {
		//nolint:usetesting // The benchmark context is canceled before Cleanup.
		_ = tp.Shutdown(context.Background())
	})
	_, span := tp.Tracer(b.Name()).Start(b.Context(), "span")
	b.Cleanup(func() { span.End() })
	return span
}

func benchAddEvent(b *testing.B, lengthLimit, countLimit int, attrs []attribute.KeyValue) {
	b.Helper()
	span := startAttrLimitSpan(b, lengthLimit, countLimit)
	span.AddEvent("prime", trace.WithAttributes(attrs...))
	span.AddEvent("prime", trace.WithAttributes(attrs...))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span.AddEvent("event", trace.WithAttributes(attrs...))
	}
}

func benchAddLink(
	b *testing.B,
	lengthLimit, countLimit int,
	attrs []attribute.KeyValue,
	sc trace.SpanContext,
) {
	b.Helper()
	span := startAttrLimitSpan(b, lengthLimit, countLimit)
	l := trace.Link{SpanContext: sc, Attributes: attrs}
	span.AddLink(l)
	span.AddLink(l)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span.AddLink(l)
	}
}

func benchSetAttributes(b *testing.B, lengthLimit, countLimit int, attrs []attribute.KeyValue) {
	b.Helper()

	limits := sdktrace.NewSpanLimits()
	limits.AttributeValueLengthLimit = lengthLimit
	limits.AttributeCountLimit = countLimit

	tp := sdktrace.NewTracerProvider(sdktrace.WithRawSpanLimits(limits))
	b.Cleanup(func() {
		//nolint:usetesting // The benchmark context is canceled before Cleanup.
		_ = tp.Shutdown(context.Background())
	})
	tracer := tp.Tracer(b.Name())
	ctx := b.Context()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := tracer.Start(ctx, "span")
		span.SetAttributes(attrs...)
		span.End()
	}
}
