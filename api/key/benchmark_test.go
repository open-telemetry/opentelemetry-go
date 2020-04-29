package key_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Note: The tests below load a real SDK to ensure the compiler isn't optimizing
// the test based on global analysis limited to the NoopSpan implementation.
func getSpan() trace.Span {
	_, sp := global.Tracer("Test").Start(context.Background(), "Span")
	tr, _ := sdktrace.NewProvider()
	global.SetTraceProvider(tr)
	return sp
}

func BenchmarkMultiNoKeyInference(b *testing.B) {
	sp := getSpan()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sp.SetAttributes(key.Int("Attr", 1))
	}
}

func BenchmarkMultiWithKeyInference(b *testing.B) {
	sp := getSpan()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sp.SetAttributes(key.Infer("Attr", 1))
	}
}

func BenchmarkSingleWithKeyInferenceValue(b *testing.B) {
	sp := getSpan()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sp.SetAttribute("Attr", 1)
	}

	b.StopTimer()
}
