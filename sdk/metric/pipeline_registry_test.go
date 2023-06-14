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
	"sync/atomic"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/resource"
)

var defaultView = NewView(Instrument{Name: "*"}, Stream{})

type invalidAggregation struct {
	aggregation.Aggregation
}

func (invalidAggregation) Copy() aggregation.Aggregation {
	return invalidAggregation{}
}
func (invalidAggregation) Err() error {
	return nil
}

func testCreateAggregators[N int64 | float64](t *testing.T) {
	changeAggView := NewView(
		Instrument{Name: "foo"},
		Stream{Aggregation: aggregation.ExplicitBucketHistogram{}},
	)
	renameView := NewView(
		Instrument{Name: "foo"},
		Stream{Name: "bar"},
	)
	defaultAggView := NewView(
		Instrument{Name: "foo"},
		Stream{Aggregation: aggregation.Default{}},
	)
	invalidAggView := NewView(
		Instrument{Name: "foo"},
		Stream{Aggregation: invalidAggregation{}},
	)

	instruments := []Instrument{
		{Name: "foo", Kind: InstrumentKind(0)}, //Unknown kind
		{Name: "foo", Kind: InstrumentKindCounter},
		{Name: "foo", Kind: InstrumentKindUpDownCounter},
		{Name: "foo", Kind: InstrumentKindHistogram},
		{Name: "foo", Kind: InstrumentKindObservableCounter},
		{Name: "foo", Kind: InstrumentKindObservableUpDownCounter},
		{Name: "foo", Kind: InstrumentKindObservableGauge},
	}

	testcases := []struct {
		name     string
		reader   Reader
		views    []View
		inst     Instrument
		wantKind aggregate.Aggregator[N] //Aggregators should match len and types
		wantLen  int
		wantErr  error
	}{
		{
			name:   "drop should return 0 aggregators",
			reader: NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Drop{} })),
			views:  []View{defaultView},
			inst:   instruments[InstrumentKindCounter],
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindUpDownCounter],
			wantKind: aggregate.NewDeltaSum[N](false),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindHistogram],
			wantKind: aggregate.NewDeltaHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindObservableCounter],
			wantKind: aggregate.NewPrecomputedDeltaSum[N](true),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindObservableUpDownCounter],
			wantKind: aggregate.NewPrecomputedDeltaSum[N](false),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindObservableGauge],
			wantKind: aggregate.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindCounter],
			wantKind: aggregate.NewDeltaSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindUpDownCounter],
			wantKind: aggregate.NewCumulativeSum[N](false),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindHistogram],
			wantKind: aggregate.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindObservableCounter],
			wantKind: aggregate.NewPrecomputedCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindObservableUpDownCounter],
			wantKind: aggregate.NewPrecomputedCumulativeSum[N](false),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindObservableGauge],
			wantKind: aggregate.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindCounter],
			wantKind: aggregate.NewCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "view should overwrite reader",
			reader:   NewManualReader(),
			views:    []View{changeAggView},
			inst:     instruments[InstrumentKindCounter],
			wantKind: aggregate.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "multiple views should create multiple aggregators",
			reader:   NewManualReader(),
			views:    []View{defaultView, renameView},
			inst:     instruments[InstrumentKindCounter],
			wantKind: aggregate.NewCumulativeSum[N](true),
			wantLen:  2,
		},
		{
			name:     "reader with default aggregation should figure out a Counter",
			reader:   NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindCounter],
			wantKind: aggregate.NewCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader with default aggregation should figure out an UpDownCounter",
			reader:   NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindUpDownCounter],
			wantKind: aggregate.NewCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader with default aggregation should figure out an Histogram",
			reader:   NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindHistogram],
			wantKind: aggregate.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "reader with default aggregation should figure out an ObservableCounter",
			reader:   NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindObservableCounter],
			wantKind: aggregate.NewPrecomputedCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader with default aggregation should figure out an ObservableUpDownCounter",
			reader:   NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindObservableUpDownCounter],
			wantKind: aggregate.NewPrecomputedCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader with default aggregation should figure out an ObservableGauge",
			reader:   NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:    []View{defaultView},
			inst:     instruments[InstrumentKindObservableGauge],
			wantKind: aggregate.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:    "view with invalid aggregation should error",
			reader:  NewManualReader(),
			views:   []View{invalidAggView},
			inst:    instruments[InstrumentKindCounter],
			wantErr: errCreatingAggregators,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			var c cache[string, streamID]
			i := newInserter[N](newPipeline(nil, tt.reader, tt.views), &c)
			got, err := i.Instrument(tt.inst)
			assert.ErrorIs(t, err, tt.wantErr)
			require.Len(t, got, tt.wantLen)
			for _, agg := range got {
				assert.IsType(t, tt.wantKind, agg)
			}
		})
	}
}

func TestCreateAggregators(t *testing.T) {
	t.Run("Int64", testCreateAggregators[int64])
	t.Run("Float64", testCreateAggregators[float64])
}

func testInvalidInstrumentShouldPanic[N int64 | float64]() {
	var c cache[string, streamID]
	i := newInserter[N](newPipeline(nil, NewManualReader(), []View{defaultView}), &c)
	inst := Instrument{
		Name: "foo",
		Kind: InstrumentKind(255),
	}
	_, _ = i.Instrument(inst)
}

func TestInvalidInstrumentShouldPanic(t *testing.T) {
	assert.Panics(t, testInvalidInstrumentShouldPanic[int64])
	assert.Panics(t, testInvalidInstrumentShouldPanic[float64])
}

func TestPipelinesAggregatorForEachReader(t *testing.T) {
	r0, r1 := NewManualReader(), NewManualReader()
	pipes := newPipelines(resource.Empty(), []Reader{r0, r1}, nil)
	require.Len(t, pipes, 2, "created pipelines")

	inst := Instrument{Name: "foo", Kind: InstrumentKindCounter}
	var c cache[string, streamID]
	r := newResolver[int64](pipes, &c)
	aggs, err := r.Aggregators(inst)
	require.NoError(t, err, "resolved Aggregators error")
	require.Len(t, aggs, 2, "instrument aggregators")

	for i, p := range pipes {
		var aggN int
		for _, is := range p.aggregations {
			aggN += len(is)
		}
		assert.Equalf(t, 1, aggN, "pipeline %d: number of instrumentSync", i)
	}
}

func TestPipelineRegistryCreateAggregators(t *testing.T) {
	renameView := NewView(Instrument{Name: "foo"}, Stream{Name: "bar"})
	testRdr := NewManualReader()
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	testCases := []struct {
		name      string
		readers   []Reader
		views     []View
		inst      Instrument
		wantCount int
	}{
		{
			name: "No views have no aggregators",
			inst: Instrument{Name: "foo"},
		},
		{
			name:      "1 reader 1 view gets 1 aggregator",
			inst:      Instrument{Name: "foo"},
			readers:   []Reader{testRdr},
			wantCount: 1,
		},
		{
			name:      "1 reader 2 views gets 2 aggregator",
			inst:      Instrument{Name: "foo"},
			readers:   []Reader{testRdr},
			views:     []View{defaultView, renameView},
			wantCount: 2,
		},
		{
			name:      "2 readers 1 view each gets 2 aggregators",
			inst:      Instrument{Name: "foo"},
			readers:   []Reader{testRdr, testRdrHistogram},
			wantCount: 2,
		},
		{
			name:      "2 reader 2 views each gets 4 aggregators",
			inst:      Instrument{Name: "foo"},
			readers:   []Reader{testRdr, testRdrHistogram},
			views:     []View{defaultView, renameView},
			wantCount: 4,
		},
		{
			name:      "An instrument is duplicated in two views share the same aggregator",
			inst:      Instrument{Name: "foo"},
			readers:   []Reader{testRdr},
			views:     []View{defaultView, defaultView},
			wantCount: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p := newPipelines(resource.Empty(), tt.readers, tt.views)
			testPipelineRegistryResolveIntAggregators(t, p, tt.wantCount)
			testPipelineRegistryResolveFloatAggregators(t, p, tt.wantCount)
		})
	}
}

func testPipelineRegistryResolveIntAggregators(t *testing.T, p pipelines, wantCount int) {
	inst := Instrument{Name: "foo", Kind: InstrumentKindCounter}
	var c cache[string, streamID]
	r := newResolver[int64](p, &c)
	aggs, err := r.Aggregators(inst)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func testPipelineRegistryResolveFloatAggregators(t *testing.T, p pipelines, wantCount int) {
	inst := Instrument{Name: "foo", Kind: InstrumentKindCounter}
	var c cache[string, streamID]
	r := newResolver[float64](p, &c)
	aggs, err := r.Aggregators(inst)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func TestPipelineRegistryResource(t *testing.T) {
	v := NewView(Instrument{Name: "bar"}, Stream{Name: "foo"})
	readers := []Reader{NewManualReader()}
	views := []View{defaultView, v}
	res := resource.NewSchemaless(attribute.String("key", "val"))
	pipes := newPipelines(res, readers, views)
	for _, p := range pipes {
		assert.True(t, res.Equal(p.resource), "resource not set")
	}
}

func TestPipelineRegistryCreateAggregatorsIncompatibleInstrument(t *testing.T) {
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	readers := []Reader{testRdrHistogram}
	views := []View{defaultView}
	p := newPipelines(resource.Empty(), readers, views)
	inst := Instrument{Name: "foo", Kind: InstrumentKindObservableGauge}

	var vc cache[string, streamID]
	ri := newResolver[int64](p, &vc)
	intAggs, err := ri.Aggregators(inst)
	assert.Error(t, err)
	assert.Len(t, intAggs, 0)

	rf := newResolver[float64](p, &vc)
	floatAggs, err := rf.Aggregators(inst)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 0)
}

type logCounter struct {
	logr.LogSink

	errN  uint32
	infoN uint32
}

func (l *logCounter) Info(level int, msg string, keysAndValues ...interface{}) {
	atomic.AddUint32(&l.infoN, 1)
	l.LogSink.Info(level, msg, keysAndValues...)
}

func (l *logCounter) InfoN() int {
	return int(atomic.SwapUint32(&l.infoN, 0))
}

func (l *logCounter) Error(err error, msg string, keysAndValues ...interface{}) {
	atomic.AddUint32(&l.errN, 1)
	l.LogSink.Error(err, msg, keysAndValues...)
}

func (l *logCounter) ErrorN() int {
	return int(atomic.SwapUint32(&l.errN, 0))
}

func TestResolveAggregatorsDuplicateErrors(t *testing.T) {
	tLog := testr.NewWithOptions(t, testr.Options{Verbosity: 6})
	l := &logCounter{LogSink: tLog.GetSink()}
	otel.SetLogger(logr.New(l))

	renameView := NewView(Instrument{Name: "bar"}, Stream{Name: "foo"})
	readers := []Reader{NewManualReader()}
	views := []View{defaultView, renameView}

	fooInst := Instrument{Name: "foo", Kind: InstrumentKindCounter}
	barInst := Instrument{Name: "bar", Kind: InstrumentKindCounter}

	p := newPipelines(resource.Empty(), readers, views)

	var vc cache[string, streamID]
	ri := newResolver[int64](p, &vc)
	intAggs, err := ri.Aggregators(fooInst)
	assert.NoError(t, err)
	assert.Equal(t, 0, l.InfoN(), "no info logging should happen")
	assert.Len(t, intAggs, 1)

	// The Rename view should produce the same instrument without an error, the
	// default view should also cause a new aggregator to be returned.
	intAggs, err = ri.Aggregators(barInst)
	assert.NoError(t, err)
	assert.Equal(t, 0, l.InfoN(), "no info logging should happen")
	assert.Len(t, intAggs, 2)

	// Creating a float foo instrument should log a warning because there is an
	// int foo instrument.
	rf := newResolver[float64](p, &vc)
	floatAggs, err := rf.Aggregators(fooInst)
	assert.NoError(t, err)
	assert.Equal(t, 1, l.InfoN(), "instrument conflict not logged")
	assert.Len(t, floatAggs, 1)

	fooInst = Instrument{Name: "foo-float", Kind: InstrumentKindCounter}

	floatAggs, err = rf.Aggregators(fooInst)
	assert.NoError(t, err)
	assert.Equal(t, 0, l.InfoN(), "no info logging should happen")
	assert.Len(t, floatAggs, 1)

	floatAggs, err = rf.Aggregators(barInst)
	assert.NoError(t, err)
	// Both the rename and default view aggregators created above should now
	// conflict. Therefore, 2 warning messages should be logged.
	assert.Equal(t, 2, l.InfoN(), "instrument conflicts not logged")
	assert.Len(t, floatAggs, 2)
}

func TestIsAggregatorCompatible(t *testing.T) {
	var undefinedInstrument InstrumentKind

	testCases := []struct {
		name string
		kind InstrumentKind
		agg  aggregation.Aggregation
		want error
	}{
		{
			name: "SyncCounter and Drop",
			kind: InstrumentKindCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncCounter and LastValue",
			kind: InstrumentKindCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncCounter and Sum",
			kind: InstrumentKindCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncCounter and ExplicitBucketHistogram",
			kind: InstrumentKindCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
		},
		{
			name: "SyncUpDownCounter and Drop",
			kind: InstrumentKindUpDownCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncUpDownCounter and LastValue",
			kind: InstrumentKindUpDownCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncUpDownCounter and Sum",
			kind: InstrumentKindUpDownCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncUpDownCounter and ExplicitBucketHistogram",
			kind: InstrumentKindUpDownCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncHistogram and Drop",
			kind: InstrumentKindHistogram,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncHistogram and LastValue",
			kind: InstrumentKindHistogram,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncHistogram and Sum",
			kind: InstrumentKindHistogram,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncHistogram and ExplicitBucketHistogram",
			kind: InstrumentKindHistogram,
			agg:  aggregation.ExplicitBucketHistogram{},
		},
		{
			name: "ObservableCounter and Drop",
			kind: InstrumentKindObservableCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "ObservableCounter and LastValue",
			kind: InstrumentKindObservableCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "ObservableCounter and Sum",
			kind: InstrumentKindObservableCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "ObservableCounter and ExplicitBucketHistogram",
			kind: InstrumentKindObservableCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "ObservableUpDownCounter and Drop",
			kind: InstrumentKindObservableUpDownCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "ObservableUpDownCounter and LastValue",
			kind: InstrumentKindObservableUpDownCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "ObservableUpDownCounter and Sum",
			kind: InstrumentKindObservableUpDownCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "ObservableUpDownCounter and ExplicitBucketHistogram",
			kind: InstrumentKindObservableUpDownCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "ObservableGauge and Drop",
			kind: InstrumentKindObservableGauge,
			agg:  aggregation.Drop{},
		},
		{
			name: "ObservableGauge and aggregation.LastValue{}",
			kind: InstrumentKindObservableGauge,
			agg:  aggregation.LastValue{},
		},
		{
			name: "ObservableGauge and Sum",
			kind: InstrumentKindObservableGauge,
			agg:  aggregation.Sum{},
			want: errIncompatibleAggregation,
		},
		{
			name: "ObservableGauge and ExplicitBucketHistogram",
			kind: InstrumentKindObservableGauge,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "unknown kind with Sum should error",
			kind: undefinedInstrument,
			agg:  aggregation.Sum{},
			want: errIncompatibleAggregation,
		},
		{
			name: "unknown kind with LastValue should error",
			kind: undefinedInstrument,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "unknown kind with Histogram should error",
			kind: undefinedInstrument,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := isAggregatorCompatible(tt.kind, tt.agg)
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
