package metric

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

func BenchmarkCounterAddNoAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1)
	}
}

// Benchmark prints 3 allocs per Add():
//  1. new []attribute.KeyValue for the list of attributes
//  2. interface{} wrapper around attribute.Set
//  3. an attribute array (map key)
func BenchmarkCounterAddOneAttr(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.String("K", "V"))
	}
}

// Benchmark prints 11 allocs per Add(), I see 10 in the profile:
//  1. new []attribute.KeyValue for the list of attributes
//  2. an attribute.Sortable (acquireRecord)
//  3. the attribute.Set underlying array
//  4. interface{} wrapper around attribute.Set value
//  5. internal to sync.Map
//  6. internal sync.Map
//  7. new syncstate.record
//  8. new viewstate.syncAccumulator
//  9. an attribute.Sortable (findOutput)
// 10. an output Aggregator
func BenchmarkCounterAddManyAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", i))
	}
}

func BenchmarkCounterCollectOneAttrNoReuse(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	producer := rdr.Producer()
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", 1))

		_ = producer.Produce(ctx, nil)
	}
}

func BenchmarkCounterCollectOneAttrWithReuse(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	producer := rdr.Producer()
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	var reuse reader.Metrics

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", 1))

		reuse = producer.Produce(ctx, &reuse)
	}
}

func BenchmarkCounterCollectTenAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	producer := rdr.Producer()
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	var reuse reader.Metrics

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			cntr.Add(ctx, 1, attribute.Int("K", j))
		}
		reuse = producer.Produce(ctx, &reuse)
	}
}

func BenchmarkCounterCollectTenAttrsTenTimes(b *testing.B) {
	ctx := context.Background()
	rdr := metrictest.NewReader()
	provider := New(WithReader(rdr))
	producer := rdr.Producer()
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	var reuse reader.Metrics

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10; k++ {
			for j := 0; j < 10; j++ {
				cntr.Add(ctx, 1, attribute.Int("K", j))
			}
			reuse = producer.Produce(ctx, &reuse)
		}
	}
}
