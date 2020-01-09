package internal_test

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// benchFixture is copied from sdk/metric/benchmark_test.go.
// TODO refactor to share this code.
type benchFixture struct {
	sdk *sdk.SDK
	B   *testing.B
}

func newFixture(b *testing.B) *benchFixture {
	b.ReportAllocs()
	bf := &benchFixture{
		B: b,
	}
	bf.sdk = sdk.New(bf, label.NewDefaultEncoder())
	return bf
}

func (*benchFixture) AggregatorFor(descriptor *export.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case export.CounterKind:
		return counter.New()
	case export.GaugeKind:
		return gauge.New()
	case export.MeasureKind:
		if strings.HasSuffix(descriptor.Name().String(), "minmaxsumcount") {
			return minmaxsumcount.New(descriptor)
		} else if strings.HasSuffix(descriptor.Name().String(), "ddsketch") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), descriptor)
		} else if strings.HasSuffix(descriptor.Name().String(), "array") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), descriptor)
		}
	}
	return nil
}

func (*benchFixture) Process(context.Context, export.Record) error {
	return nil
}

func (*benchFixture) CheckpointSet() export.CheckpointSet {
	return nil
}

func (*benchFixture) FinishedCollection() {
}

func BenchmarkGlobalInt64CounterAddNoSDK(b *testing.B) {
	internal.ResetForTest()

	sdk := global.Scope().WithNamespace("test").Meter()

	ctx := scope.ContextWithScope(
		context.Background(),
		scope.Empty().AddResources(key.String("A", "B")))

	cnt := sdk.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1)
	}
}

func BenchmarkGlobalInt64CounterAddWithSDK(b *testing.B) {
	// Comapare with BenchmarkInt64CounterAdd() in ../../sdk/meter/benchmark_test.go
	fix := newFixture(b)

	sdk := global.Scope().WithNamespace("test").Meter()

	ctx := scope.ContextWithScope(
		context.Background(),
		scope.Empty().AddResources(key.String("A", "B")))

	global.SetScope(scope.Empty().WithMeter(fix.sdk))

	cnt := sdk.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1)
	}
}

func BenchmarkStartEndSpan(b *testing.B) {
	// Comapare with BenchmarkStartEndSpan() in ../../sdk/trace/benchmark_test.go
	traceBenchmark(b, func(b *testing.B) {
		t := global.Scope().WithNamespace("Benchmark StartEndSpan").Tracer()
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.End()
		}
	})
}

func traceBenchmark(b *testing.B, fn func(*testing.B)) {
	internal.ResetForTest()
	b.Run("No SDK", func(b *testing.B) {
		b.ReportAllocs()
		fn(b)
	})
	b.Run("Default SDK (AlwaysSample)", func(b *testing.B) {
		b.ReportAllocs()
		global.SetScope(scope.Empty().WithTracer(newTracer(b, sdktrace.AlwaysSample())))
		fn(b)
	})
	b.Run("Default SDK (NeverSample)", func(b *testing.B) {
		b.ReportAllocs()
		global.SetScope(scope.Empty().WithTracer(newTracer(b, sdktrace.NeverSample())))
		fn(b)
	})
}

func newTracer(b *testing.B, sampler sdktrace.Sampler) trace.TracerSDK {
	tpi, err := sdktrace.NewTracer(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}))
	if err != nil {
		b.Fatalf("Failed to create trace provider with sampler: %v", err)
	}
	return tpi
}
