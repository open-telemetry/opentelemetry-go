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
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var (
	zeroUnit           unit.Unit
	zeroInstrumentKind InstrumentKind
	zeroScope          instrumentation.Scope
)

// InstrumentKind is the identifier of a group of instruments that all
// perfroming the same function.
type InstrumentKind uint8

const (
	// instrumentKindUndefined is an undefined instrument kind, it should not
	// be used by any initialized type.
	instrumentKindUndefined InstrumentKind = iota // nolint:deadcode,varcheck
	// InstrumentKindSyncCounter identifies a group of instruments that record
	// increasing values synchronously with the code path they are measuring.
	InstrumentKindSyncCounter
	// InstrumentKindSyncUpDownCounter identifies a group of instruments that
	// record increasing and decreasing values synchronously with the code path
	// they are measuring.
	InstrumentKindSyncUpDownCounter
	// InstrumentKindSyncHistogram identifies a group of instruments that
	// record a distribution of values synchronously with the code path they
	// are measuring.
	InstrumentKindSyncHistogram
	// InstrumentKindAsyncCounter identifies a group of instruments that record
	// increasing values in an asynchronous callback.
	InstrumentKindAsyncCounter
	// InstrumentKindAsyncUpDownCounter identifies a group of instruments that
	// record increasing and decreasing values in an asynchronous callback.
	InstrumentKindAsyncUpDownCounter
	// InstrumentKindAsyncGauge identifies a group of instruments that record
	// current values in an asynchronous callback.
	InstrumentKindAsyncGauge
)

type nonComparable [0]func()

// Instrument describes properies an instrument is created with.
type Instrument struct {
	// Name is the human-readable identifier of the instrument.
	Name string
	// Description describes the purpose of the instrument.
	Description string
	// Kind defines the functional group of the instrument.
	Kind InstrumentKind
	// Unit is the unit of measurment recorded by the instrument.
	Unit unit.Unit
	// Scope identifies the instrumentation that created the instrument.
	Scope instrumentation.Scope

	nonComparable
}

func (p Instrument) mask(m Instrument) Instrument {
	if m.Name != "" {
		p.Name = m.Name
	}
	if m.Description != "" {
		p.Description = m.Description
	}
	if m.Kind != zeroInstrumentKind {
		p.Kind = m.Kind
	}
	if m.Unit != zeroUnit {
		p.Unit = m.Unit
	}
	if m.Scope.Name != "" {
		p.Scope.Name = m.Scope.Name
	}
	if m.Scope.Version != "" {
		p.Scope.Version = m.Scope.Version
	}
	if m.Scope.SchemaURL != "" {
		p.Scope.SchemaURL = m.Scope.SchemaURL
	}
	return p
}

func (p Instrument) empty() bool {
	return p.Name == "" &&
		p.Description == "" &&
		p.Kind == zeroInstrumentKind &&
		p.Unit == zeroUnit &&
		p.Scope == zeroScope
}

func (p Instrument) matches(o Instrument) bool {
	return p.matchesName(o) &&
		p.matchesDescription(o) &&
		p.matchesKind(o) &&
		p.matchesUnit(o) &&
		p.matchesScope(o)
}

func (p Instrument) matchesName(o Instrument) bool {
	return p.Name == "" || p.Name == o.Name
}

func (p Instrument) matchesDescription(o Instrument) bool {
	return p.Description == "" || p.Description == o.Description
}

func (p Instrument) matchesKind(o Instrument) bool {
	return p.Kind == zeroInstrumentKind || p.Kind == o.Kind
}

func (p Instrument) matchesUnit(o Instrument) bool {
	return p.Unit == zeroUnit || p.Unit == o.Unit
}

func (p Instrument) matchesScope(o Instrument) bool {
	return (p.Scope.Name == "" || p.Scope.Name == o.Scope.Name) &&
		(p.Scope.Version == "" || p.Scope.Version == o.Scope.Version) &&
		(p.Scope.SchemaURL == "" || p.Scope.SchemaURL == o.Scope.SchemaURL)
}

// DataStream defines the stream of data an instrument produces.
type DataStream struct {
	Instrument

	// Temporality defines the calculation window for the data aggregation.
	Temporality metricdata.Temporality
	// Aggregation the stream uses for an instrument.
	Aggregation aggregation.Aggregation
	// AttributeFilter applied to all attributes recorded for an instrument.
	AttributeFilter attribute.Filter

	nonComparable
}

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
