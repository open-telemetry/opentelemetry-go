package internal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testSpanProcesor struct {
	// Names of Spans started.
	spansStarted []string
	// Names of Spans ended.
	spansEnded []string
}

func (t *testSpanProcesor) OnStart(s *export.SpanData) {
	t.spansStarted = append(t.spansStarted, s.Namespace.Name(s.Name).String())
}

func (t *testSpanProcesor) OnEnd(s *export.SpanData) {
	t.spansEnded = append(t.spansEnded, s.Namespace.Name(s.Name).String())
}

func (t *testSpanProcesor) Shutdown() {}

func TestTraceDefaultSDK(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	tracer1 := global.Scope().WithNamespace("pre").Tracer()
	_, span1 := tracer1.Start(ctx, "span1")

	// This should be dropped.
	if err := tracer1.WithSpan(ctx, "withSpan1", func(context.Context) error { return nil }); err != nil {
		t.Errorf("failed to wrap function with span prior to initialization: %v", err)
	}

	tpi, err := sdktrace.NewTracer(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	if err != nil {
		t.Fatal(err)
	}
	tsp := &testSpanProcesor{}
	tpi.RegisterSpanProcessor(tsp)

	global.SetScope(scope.Empty().WithTracer(tpi))

	// This span was started before initialization, it is expected to be dropped.
	span1.End()

	// The existing Tracer should have been configured to now use the configured SDK.
	_, span2 := tracer1.Start(ctx, "span2")
	span2.End()
	if err := tracer1.WithSpan(ctx, "withSpan2", func(context.Context) error { return nil }); err != nil {
		t.Errorf("failed to wrap function with span post initialization: %v", err)
	}

	// The global trace Provider should now create Tracers that also use the newly configured SDK.
	tracer2 := global.Scope().WithNamespace("post").Tracer()
	_, span3 := tracer2.Start(ctx, "span3")
	span3.End()
	if err := tracer2.WithSpan(ctx, "withSpan3", func(context.Context) error { return nil }); err != nil {
		t.Errorf("failed to wrap function with span post initialization with new tracer: %v", err)
	}

	expected := []string{"pre/span2", "pre/withSpan2", "post/span3", "post/withSpan3"}
	require.Equal(t, expected, tsp.spansStarted)
	require.Equal(t, expected, tsp.spansEnded)
}
