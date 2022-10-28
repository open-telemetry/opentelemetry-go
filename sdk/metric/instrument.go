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
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// instrumentID are the identifying properties of an instrument.
type instrumentID struct {
	// Name is the name of the instrument.
	Name string
	// Description is the description of the instrument.
	Description string
	// Unit is the unit of the instrument.
	Unit unit.Unit
	// Aggregation is the aggregation data type of the instrument.
	Aggregation string
	// Monotonic is the monotonicity of an instruments data type. This field is
	// not used for all data types, so a zero value needs to be understood in the
	// context of Aggregation.
	Monotonic bool
	// Temporality is the temporality of an instrument's data type. This field
	// is not used by some data types.
	Temporality metricdata.Temporality
	// Number is the number type of the instrument.
	Number string
}

type instrumentImpl[N int64 | float64] struct {
	instrument.Asynchronous
	instrument.Synchronous

	aggregators []internal.Aggregator[N]
}

//instrument.Sync[int64]
//instrument.Sync[float64]

var _ instrument.AsyncCounter[float64] = &instrumentImpl[float64]{}
var _ instrument.AsyncUpDownCounter[float64] = &instrumentImpl[float64]{}
var _ instrument.AsyncGauge[float64] = &instrumentImpl[float64]{}
var _ instrument.AsyncCounter[int64] = &instrumentImpl[int64]{}
var _ instrument.AsyncUpDownCounter[int64] = &instrumentImpl[int64]{}
var _ instrument.AsyncGauge[int64] = &instrumentImpl[int64]{}
var _ instrument.SyncCounter[float64] = &instrumentImpl[float64]{}
var _ instrument.SyncUpDownCounter[float64] = &instrumentImpl[float64]{}
var _ instrument.SyncHistogram[float64] = &instrumentImpl[float64]{}
var _ instrument.SyncCounter[int64] = &instrumentImpl[int64]{}
var _ instrument.SyncUpDownCounter[int64] = &instrumentImpl[int64]{}
var _ instrument.SyncHistogram[int64] = &instrumentImpl[int64]{}

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
