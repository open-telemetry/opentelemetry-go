package internal_test

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
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

var _ metric.Provider = &benchFixture{}

func newFixture(b *testing.B) *benchFixture {
	b.ReportAllocs()
	bf := &benchFixture{
		B: b,
	}
	bf.sdk = sdk.New(bf, sdk.NewDefaultLabelEncoder())
	return bf
}

func (*benchFixture) AggregatorFor(descriptor *export.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case export.CounterKind:
		return counter.New()
	case export.GaugeKind:
		return gauge.New()
	case export.MeasureKind:
		if strings.HasSuffix(descriptor.Name(), "minmaxsumcount") {
			return minmaxsumcount.New(descriptor)
		} else if strings.HasSuffix(descriptor.Name(), "ddsketch") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), descriptor)
		} else if strings.HasSuffix(descriptor.Name(), "array") {
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

func (fix *benchFixture) Meter(name string) metric.Meter {
	return fix.sdk
}

func BenchmarkGlobalInt64CounterAddNoSDK(b *testing.B) {
	internal.ResetForTest()
	ctx := context.Background()
	sdk := global.MeterProvider().Meter("test")
	labs := sdk.Labels(key.String("A", "B"))
	cnt := sdk.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs)
	}
}

func BenchmarkGlobalInt64CounterAddWithSDK(b *testing.B) {
	// Comapare with BenchmarkInt64CounterAdd() in ../../sdk/meter/benchmark_test.go
	ctx := context.Background()
	fix := newFixture(b)

	sdk := global.MeterProvider().Meter("test")

	global.SetMeterProvider(fix)

	labs := sdk.Labels(key.String("A", "B"))
	cnt := sdk.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs)
	}
}

func BenchmarkStartEndSpan(b *testing.B) {
	// Comapare with BenchmarkStartEndSpan() in ../../sdk/trace/benchmark_test.go
	traceBenchmark(b, func(b *testing.B) {
		t := global.TraceProvider().Tracer("Benchmark StartEndSpan")
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
		global.SetTraceProvider(traceProvider(b, sdktrace.AlwaysSample()))
		fn(b)
	})
	b.Run("Default SDK (NeverSample)", func(b *testing.B) {
		b.ReportAllocs()
		global.SetTraceProvider(traceProvider(b, sdktrace.NeverSample()))
		fn(b)
	})
}

func traceProvider(b *testing.B, sampler sdktrace.Sampler) trace.Provider {
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}))
	if err != nil {
		b.Fatalf("Failed to create trace provider with sampler: %v", err)
	}
	return tp
}
