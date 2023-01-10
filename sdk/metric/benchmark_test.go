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
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
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
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.String("K", "V"))
	}
}

func BenchmarkCounterAddOneInvalidAttr(b *testing.B) {
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.String("", "V"), attribute.String("K", "V"))
	}
}

func BenchmarkCounterAddSingleUseAttrs(b *testing.B) {
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", i))
	}
}

func BenchmarkCounterAddSingleUseInvalidAttrs(b *testing.B) {
	ctx, _, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("", i), attribute.Int("K", i))
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
		cntr.Add(ctx, 1, attribute.Int("L", i), attribute.Int("K", i))
	}
}

func BenchmarkCounterCollectOneAttr(b *testing.B) {
	ctx, rdr, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", 1))

		_, _ = rdr.Collect(ctx)
	}
}

func BenchmarkCounterCollectTenAttrs(b *testing.B) {
	ctx, rdr, cntr := benchCounter(b)

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			cntr.Add(ctx, 1, attribute.Int("K", j))
		}
		_, _ = rdr.Collect(ctx)
	}
}
