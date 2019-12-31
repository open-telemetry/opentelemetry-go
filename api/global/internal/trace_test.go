package internal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

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
	t.spansStarted = append(t.spansStarted, s.Name)
}

func (t *testSpanProcesor) OnEnd(s *export.SpanData) {
	t.spansEnded = append(t.spansEnded, s.Name)
}

func (t *testSpanProcesor) Shutdown() {}

func TestTraceDefaultSDK(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	gtp := global.TraceProvider()
	tracer1 := gtp.Tracer("pre")
	_, span1 := tracer1.Start(ctx, "span1")

	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	if err != nil {
		t.Fatal(err)
	}
	tsp := &testSpanProcesor{}
	tp.RegisterSpanProcessor(tsp)

	global.SetTraceProvider(tp)

	// This span was started before initialization, it is expected to be dropped.
	span1.End()

	// The existing Tracer should have been configured to now use the configured SDK.
	_, span2 := tracer1.Start(ctx, "span2")
	span2.End()

	// The global trace Provider should now create Tracers that also use the newly configured SDK.
	tracer2 := gtp.Tracer("post")
	_, span3 := tracer2.Start(ctx, "span3")
	span3.End()

	expected := []string{"pre/span2", "post/span3"}
	require.Equal(t, tsp.spansStarted, expected)
	require.Equal(t, tsp.spansEnded, expected)
}
