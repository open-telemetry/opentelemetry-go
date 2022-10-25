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
	"go.opentelemetry.io/otel/sdk/metric/view"
)

func BenchmarkCounterAddNoAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1)
	}
}

func BenchmarkCounterAddOneAttr(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.String("K", "V"))
	}
}

func BenchmarkCounterAddOneInvalidAttr(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.String("", "V"), attribute.String("K", "V"))
	}
}

func BenchmarkCounterAddManyAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", i))
	}
}

func BenchmarkCounterAddManyInvalidAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("", i), attribute.Int("K", i))
	}
}

func BenchmarkCounterAddManyFilteredAttrs(b *testing.B) {
	vw, _ := view.New(view.WithFilterAttributes(attribute.Key("K")))

	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr, vw))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("L", i), attribute.Int("K", i))
	}
}

func BenchmarkCounterCollectOneAttr(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		cntr.Add(ctx, 1, attribute.Int("K", 1))

		_, _ = rdr.Collect(ctx)
	}
}

func BenchmarkCounterCollectTenAttrs(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			cntr.Add(ctx, 1, attribute.Int("K", j))
		}
		_, _ = rdr.Collect(ctx)
	}
}

func BenchmarkCounterCollectTenAttrsTenTimes(b *testing.B) {
	ctx := context.Background()
	rdr := NewManualReader()
	provider := NewMeterProvider(WithReader(rdr))
	b.ReportAllocs()

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	for i := 0; i < b.N; i++ {
		for k := 0; k < 10; k++ {
			for j := 0; j < 10; j++ {
				cntr.Add(ctx, 1, attribute.Int("K", j))
			}
			_, _ = rdr.Collect(ctx)
		}
	}
}
