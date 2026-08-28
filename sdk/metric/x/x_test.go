// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"bytes"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
	"go.opentelemetry.io/otel/trace"
)

func TestBoundCounterCumulativeAndDelta(t *testing.T) {
	cumulative := sdkmetric.NewManualReader()
	delta := sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(sdkmetric.DeltaTemporalitySelector),
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(cumulative),
		sdkmetric.WithReader(delta),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, err := provider.Meter("test").Int64Counter("requests")
	if err != nil {
		t.Fatal(err)
	}
	binder, ok := counter.(metricx.Int64CounterBinder)
	if !ok {
		t.Fatal("counter does not implement Int64CounterBinder")
	}
	attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
	bound := binder.Bind(attrs...)
	if _, ok := bound.(metricx.Int64CounterBinder); ok {
		t.Error("bound counter implements Int64CounterBinder")
	}
	if _, ok := bound.(metricx.Finisher); ok {
		t.Error("bound counter implements Finisher")
	}
	attrs[0] = attribute.String("route", "/mutated")
	bound.Add(t.Context(), 2)
	counter.Add(t.Context(), 3, metric.WithAttributes(attribute.String("route", "/orders")))

	assertSum(t, collect(t, cumulative), metricdata.CumulativeTemporality, "/orders", 5)
	assertSum(t, collect(t, delta), metricdata.DeltaTemporality, "/orders", 5)

	bound.Add(t.Context(), 1)
	assertSum(t, collect(t, cumulative), metricdata.CumulativeTemporality, "/orders", 6)
	assertSum(t, collect(t, delta), metricdata.DeltaTemporality, "/orders", 1)
}

func TestProviderAcceptsStableReader(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(reader),
	)
	counter, err := provider.Meter("test").Int64Counter("requests")
	if err != nil {
		t.Fatal(err)
	}
	counter.Add(t.Context(), 1)
	points := sumPoints(t, collect(t, reader))
	if len(points) != 1 || points[0].Value != 1 {
		t.Fatalf("unexpected points: %#v", points)
	}
}

func TestProviderRetainsStableInstrumentSupport(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(reader),
	)
	counter, err := provider.Meter("test").Int64UpDownCounter("active")
	if err != nil {
		t.Fatal(err)
	}
	counter.Add(t.Context(), -1)

	rm := collect(t, reader)
	points := sumPoints(t, rm)
	if len(points) != 1 || points[0].Value != -1 {
		t.Fatalf("unexpected points: %#v", points)
	}
}

func TestFinishAndLazyReactivation(t *testing.T) {
	cumulative := sdkmetric.NewManualReader()
	delta := sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(sdkmetric.DeltaTemporalitySelector),
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(cumulative),
		sdkmetric.WithReader(delta),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter("requests")
	attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
	bound := counter.(metricx.Int64CounterBinder).Bind(attrs...)
	bound.Add(t.Context(), 4)
	finisher := counter.(metricx.Finisher)
	finisher.Finish(t.Context(), attrs...)
	finisher.Finish(t.Context(), attrs...)
	bound.Add(t.Context(), 2)

	assertSum(t, collect(t, cumulative), metricdata.CumulativeTemporality, "/orders", 4)
	assertSum(t, collect(t, delta), metricdata.DeltaTemporality, "/orders", 4)
	assertSum(t, collect(t, cumulative), metricdata.CumulativeTemporality, "/orders", 2)
	assertSum(t, collect(t, delta), metricdata.DeltaTemporality, "/orders", 2)
}

func TestViewFilteringAndRecordTimeFallback(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	view := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "requests"},
		sdkmetric.Stream{AttributeFilter: attribute.NewAllowKeysFilter("route")},
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(view),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter("requests")
	bound := counter.(metricx.Int64CounterBinder).Bind(
		attribute.String("route", "/old"),
		attribute.String("secret", "drop"),
	)
	bound.Add(t.Context(), 1, metric.WithAttributes(
		attribute.String("route", "/new"),
		attribute.String("other", "drop"),
	))
	assertSum(t, collect(t, reader), metricdata.CumulativeTemporality, "/new", 1)
}

func TestBoundCounterNonSumFallback(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	view := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "requests"},
		sdkmetric.Stream{
			Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0, 10},
			},
		},
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(view),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter("requests")
	bound := counter.(metricx.Int64CounterBinder).Bind(attribute.String("route", "/orders"))
	bound.Add(t.Context(), 5)

	rm := collect(t, reader)
	histogram, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[int64])
	if !ok {
		t.Fatalf("data is %T, want metricdata.Histogram[int64]", rm.ScopeMetrics[0].Metrics[0].Data)
	}
	if len(histogram.DataPoints) != 1 || histogram.DataPoints[0].Count != 1 || histogram.DataPoints[0].Sum != 5 {
		t.Fatalf("unexpected histogram: %#v", histogram.DataPoints)
	}
}

func TestCardinalityFinishAndOverflow(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithCardinalityLimit(2),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter("requests")
	binder := counter.(metricx.Int64CounterBinder)
	finisher := counter.(metricx.Finisher)
	a := binder.Bind(attribute.String("id", "a"))
	b := binder.Bind(attribute.String("id", "b"))
	a.Add(t.Context(), 1)
	b.Add(t.Context(), 2)
	finisher.Finish(t.Context(), attribute.String("id", "b"))
	finisher.Finish(t.Context(), attribute.String("id", "a"))
	c := binder.Bind(attribute.String("id", "c"))
	c.Add(t.Context(), 3)

	rm := collect(t, reader)
	points := sumPoints(t, rm)
	if len(points) != 3 {
		t.Fatalf("got %d points, want 3", len(points))
	}
	if valueFor(points, "otel.metric.overflow", "true") != 2 {
		t.Fatalf("overflow value = %d, want 2", valueFor(points, "otel.metric.overflow", "true"))
	}
	if valueFor(points, "id", "a") != 1 || valueFor(points, "id", "c") != 3 {
		t.Fatalf("unexpected concrete points: %#v", points)
	}
}

func TestDefaultAttributes(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter(
		"requests",
		metricx.WithDefaultAttributes("route"),
	)
	counter.Add(t.Context(), 1, metric.WithAttributes(
		attribute.String("route", "/orders"),
		attribute.String("secret", "drop"),
	))
	points := sumPoints(t, collect(t, reader))
	if len(points) != 1 || points[0].Attributes.Len() != 1 {
		t.Fatalf("unexpected attributes: %v", points)
	}
	if valueFor(points, "route", "/orders") != 1 {
		t.Fatalf("unexpected points: %v", points)
	}
}

func TestBoundCounterExemplarContext(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	view := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "requests"},
		sdkmetric.Stream{AttributeFilter: attribute.NewAllowKeysFilter("route")},
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(view),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOnFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter("requests")
	bound := counter.(metricx.Int64CounterBinder).Bind(
		attribute.String("route", "/orders"),
		attribute.String("secret", "dropped"),
	)
	traceID := trace.TraceID{1}
	spanID := trace.SpanID{2}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})
	bound.Add(trace.ContextWithSpanContext(t.Context(), spanContext), 7)
	points := sumPoints(t, collect(t, reader))
	if len(points[0].Exemplars) == 0 {
		t.Fatal("bound measurement did not retain an exemplar")
	}
	got := points[0].Exemplars[0]
	if !bytes.Equal(got.TraceID, traceID[:]) || !bytes.Equal(got.SpanID, spanID[:]) {
		t.Fatalf("exemplar context = %x/%x, want %x/%x", got.TraceID, got.SpanID, traceID, spanID)
	}
	if len(got.FilteredAttributes) != 1 || got.FilteredAttributes[0].Key != "secret" {
		t.Fatalf("filtered attributes = %v", got.FilteredAttributes)
	}
}

func TestStableCounterDoesNotImplementExperimentalInterfaces(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	counter, _ := provider.Meter("test").Int64Counter("requests")
	if _, ok := counter.(metricx.Int64CounterBinder); ok {
		t.Error("stable counter implements Int64CounterBinder")
	}
	if _, ok := counter.(metricx.Finisher); ok {
		t.Error("stable counter implements Finisher")
	}
}

func TestBoundCounterConcurrentSafe(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithBinding(),
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, _ := provider.Meter("test").Int64Counter("requests")
	attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
	bound := counter.(metricx.Int64CounterBinder).Bind(attrs...)
	finisher := counter.(metricx.Finisher)
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for range 250 {
				bound.Add(t.Context(), 1)
			}
		})
	}
	wg.Go(func() {
		for range 50 {
			finisher.Finish(t.Context(), attrs...)
			_ = reader.Collect(t.Context(), new(metricdata.ResourceMetrics))
		}
	})
	wg.Wait()
}

func collect(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(t.Context(), &rm); err != nil {
		t.Fatal(err)
	}
	return rm
}

func sumPoints(t *testing.T, rm metricdata.ResourceMetrics) []metricdata.DataPoint[int64] {
	t.Helper()
	if len(rm.ScopeMetrics) != 1 || len(rm.ScopeMetrics[0].Metrics) != 1 {
		t.Fatalf("unexpected metric shape: %#v", rm.ScopeMetrics)
	}
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("data is %T, want metricdata.Sum[int64]", rm.ScopeMetrics[0].Metrics[0].Data)
	}
	return sum.DataPoints
}

func assertSum(
	t *testing.T,
	rm metricdata.ResourceMetrics,
	temporality metricdata.Temporality,
	attrValue string,
	want int64,
) {
	t.Helper()
	gotMetric := rm.ScopeMetrics[0].Metrics[0]
	sum := gotMetric.Data.(metricdata.Sum[int64])
	if sum.Temporality != temporality {
		t.Fatalf("temporality = %v, want %v", sum.Temporality, temporality)
	}
	if got := valueFor(sum.DataPoints, "route", attrValue); got != want {
		t.Fatalf("value = %d, want %d; points: %#v", got, want, sum.DataPoints)
	}
}

func valueFor(points []metricdata.DataPoint[int64], key, value string) int64 {
	for _, point := range points {
		if got, ok := point.Attributes.Value(attribute.Key(key)); ok && got.String() == value {
			return point.Value
		}
	}
	return 0
}
