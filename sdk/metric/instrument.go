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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric/internal"
)

type instrumentImpl[N int64 | float64] struct {
	instrument.Asynchronous
	instrument.Synchronous

	aggregators []internal.Aggregator[N]
}

var _ asyncfloat64.Counter = &instrumentImpl[float64]{}
var _ asyncfloat64.UpDownCounter = &instrumentImpl[float64]{}
var _ asyncfloat64.Gauge = &instrumentImpl[float64]{}
var _ asyncint64.Counter = &instrumentImpl[int64]{}
var _ asyncint64.UpDownCounter = &instrumentImpl[int64]{}
var _ asyncint64.Gauge = &instrumentImpl[int64]{}
var _ syncfloat64.Counter = &instrumentImpl[float64]{}
var _ syncfloat64.UpDownCounter = &instrumentImpl[float64]{}
var _ syncfloat64.Histogram = &instrumentImpl[float64]{}
var _ syncint64.Counter = &instrumentImpl[int64]{}
var _ syncint64.UpDownCounter = &instrumentImpl[int64]{}
var _ syncint64.Histogram = &instrumentImpl[int64]{}

func (i *instrumentImpl[N]) Observe(ctx context.Context, val N, attrs ...attribute.KeyValue) {
	// Only record a value if this is being called from the MetricProvider.
	_, ok := ctx.Value(produceKey).(struct{})
	if !ok {
		return
	}
	i.aggregate(ctx, val, attrs)
}

func (i *instrumentImpl[N]) Add(ctx context.Context, val N, attrs ...attribute.KeyValue) {
	i.aggregate(ctx, val, attrs)
}

func (i *instrumentImpl[N]) Record(ctx context.Context, val N, attrs ...attribute.KeyValue) {
	i.aggregate(ctx, val, attrs)
}

func (i *instrumentImpl[N]) aggregate(ctx context.Context, val N, attrs []attribute.KeyValue) {
	if err := ctx.Err(); err != nil {
		return
	}
	for _, agg := range i.aggregators {
		agg.Aggregate(val, attribute.NewSet(attrs...))
	}
}
