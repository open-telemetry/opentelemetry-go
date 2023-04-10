// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func benchCounter(b *testing.B, views ...View) (context.Context, Reader, instrument.Int64Counter) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr), WithView(views...))
	cntr, _ := provider.Meter("test").Int64Counter("hello")
	b.ResetTimer()
	b.ReportAllocs()
	return ctx, rdr, cntr
}

func BenchmarkCounterAddNoAttrs(b *testing.B) {
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1)
	}
}

func BenchmarkCounterAddOneAttr(b *testing.B) {
	s := attribute.NewSet(attribute.String("K", "V"))
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, instrument.WithAttributeSet(s))
	}
}

func BenchmarkCounterAddOneInvalidAttr(b *testing.B) {
	s := attribute.NewSet(
		attribute.String("", "V"),
		attribute.String("K", "V"),
	)
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, instrument.WithAttributeSet(s))
	}
}

func BenchmarkCounterAddSingleUseAttrs(b *testing.B) {
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		s := attribute.NewSet(attribute.Int("K", i))
		cntr.Add(ctx, 1, instrument.WithAttributeSet(s))
	}
}

func BenchmarkCounterAddSingleUseInvalidAttrs(b *testing.B) {
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		o := instrument.WithAttributes(
			attribute.Int("", i),
			attribute.Int("K", i),
		)
		cntr.Add(ctx, 1, o)
	}
}

func BenchmarkCounterAddSingleUseFilteredAttrs(b *testing.B) {
	ctx, _, cntr := benchCounter(b, NewView(
		Instrument{Name: "*"},
		Stream{AttributeFilter: func(kv attribute.KeyValue) bool {
			return kv.Key == attribute.Key("K")
		}},
	))

	for i := 0; i < b.N; i++ {
		o := instrument.WithAttributes(
			attribute.Int("L", i),
			attribute.Int("K", i),
		)
		cntr.Add(ctx, 1, o)
	}
}

func BenchmarkCounterCollectOneAttr(b *testing.B) {
	s := attribute.NewSet(attribute.Int("K", 1))
	ctx, rdr, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, instrument.WithAttributeSet(s))

		_ = rdr.Collect(ctx, nil)
	}
}

func BenchmarkCounterCollectTenAttrs(b *testing.B) {
	ctx, rdr, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			cntr.Add(ctx, 1, instrument.WithAttributes(attribute.Int("K", j)))
		}
		_ = rdr.Collect(ctx, nil)
	}
}

func BenchmarkCollectHistograms(b *testing.B) {
	b.Run("1", benchCollectHistograms(1))
	b.Run("5", benchCollectHistograms(5))
	b.Run("10", benchCollectHistograms(10))
	b.Run("25", benchCollectHistograms(25))
}

func benchCollectHistograms(count int) func(*testing.B) {
	ctx := context.Background()
	r := NewManualReader()
	mtr := NewMeterProvider(
		WithReader(r),
	).Meter("sdk/metric/bench/histogram")

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("fake data %d", i)
		h, _ := mtr.Int64Histogram(name)

		h.Record(ctx, 1)
	}

	// Store bechmark results in a closure to prevent the compiler from
	// inlining and skipping the function.
	var (
		collectedMetrics metricdata.ResourceMetrics
	)

	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			_ = r.Collect(ctx, &collectedMetrics)
			if len(collectedMetrics.ScopeMetrics[0].Metrics) != count {
				b.Fatalf("got %d metrics, want %d", len(collectedMetrics.ScopeMetrics[0].Metrics), count)
			}
		}
	}
}
