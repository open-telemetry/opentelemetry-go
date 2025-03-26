// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlptrace/otlptracetest/otlptest.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package otlptracetest provides testing utilties and framework for the
// otlptrace exporters.
package otlptracetest // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/otlptracetest"

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

// RunEndToEndTest can be used by otlptrace.Client tests to validate
// themselves.
func RunEndToEndTest(ctx context.Context, t *testing.T, exp *otlptrace.Exporter, tracesCollector TracesCollector) {
	pOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(10),
		),
	}
	tp1 := sdktrace.NewTracerProvider(append(pOpts,
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("rk1", "rv11)"),
			attribute.Int64("rk2", 5),
		)))...)

	tp2 := sdktrace.NewTracerProvider(append(pOpts,
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("rk1", "rv12)"),
			attribute.Float64("rk3", 6.5),
		)))...)

	tr1 := tp1.Tracer("test-tracer1")
	tr2 := tp2.Tracer("test-tracer2")
	// Now create few spans
	m := 4
	for i := 0; i < m; i++ {
		_, span := tr1.Start(ctx, "AlwaysSample")
		span.SetAttributes(attribute.Int64("i", int64(i)))
		span.End()

		_, span = tr2.Start(ctx, "AlwaysSample")
		span.SetAttributes(attribute.Int64("i", int64(i)))
		span.End()
	}

	func() {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := tp1.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a tracer provider 1: %v", err)
		}
		if err := tp2.Shutdown(ctx); err != nil {
			t.Fatalf("failed to shut down a tracer provider 2: %v", err)
		}
	}()

	// Wait >2 cycles.
	<-time.After(40 * time.Millisecond)

	// Now shutdown the exporter
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := exp.Shutdown(ctx); err != nil {
		t.Fatalf("failed to stop the exporter: %v", err)
	}

	// Shutdown the collector too so that we can begin
	// verification checks of expected data back.
	if err := tracesCollector.Stop(); err != nil {
		t.Fatalf("failed to stop the mock collector: %v", err)
	}

	// Now verify that we only got two resources
	rss := tracesCollector.GetResourceSpans()
	if got, want := len(rss), 2; got != want {
		t.Fatalf("resource span count: got %d, want %d\n", got, want)
	}

	// Now verify spans and attributes for each resource span.
	for _, rs := range rss {
		if len(rs.ScopeSpans) == 0 {
			t.Fatalf("zero ScopeSpans")
		}
		if got, want := len(rs.ScopeSpans[0].Spans), m; got != want {
			t.Fatalf("span counts: got %d, want %d", got, want)
		}
		attrMap := map[int64]bool{}
		for _, s := range rs.ScopeSpans[0].Spans {
			if gotName, want := s.Name, "AlwaysSample"; gotName != want {
				t.Fatalf("span name: got %s, want %s", gotName, want)
			}
			attrMap[s.Attributes[0].Value.Value.(*commonpb.AnyValue_IntValue).IntValue] = true
		}
		if got, want := len(attrMap), m; got != want {
			t.Fatalf("span attribute unique values: got %d  want %d", got, want)
		}
		for i := 0; i < m; i++ {
			_, ok := attrMap[int64(i)]
			if !ok {
				t.Fatalf("span with attribute %d missing", i)
			}
		}
	}
}
