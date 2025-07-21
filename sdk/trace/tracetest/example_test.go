// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func doSomething(ctx context.Context) {
	tr := otel.Tracer("example.com/test")
	_, span := tr.Start(ctx, "doSomething")
	defer span.End()
}

func TestDoSomething(t *testing.T) {
	sr := tracetest.NewSpanRecorder()

	tp := trace.NewTracerProvider()
	tp.RegisterSpanProcessor(sr)

	otel.SetTracerProvider(tp)

	ctx := context.Background()
	doSomething(ctx)

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("Expected 1 span, got %d", len(spans))
	}
	if spans[0].Name() != "doSomething" {
		t.Errorf("Expected span name 'doSomething', got %s", spans[0].Name())
	}
}
