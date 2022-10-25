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
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/resource"
)

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
		InstrumentProperties{Name: "foo"},
		DataStream{Aggregation: aggregation.ExplicitBucketHistogram{}},
	)
	renameView := NewView(
		InstrumentProperties{Name: "foo"},
		DataStream{InstrumentProperties: InstrumentProperties{Name: "bar"}},
	)
	defaultAggView := NewView(
		InstrumentProperties{Name: "foo"},
		DataStream{Aggregation: aggregation.Default{}},
	)
	invalidAggView := NewView(
		InstrumentProperties{Name: "foo"},
		DataStream{Aggregation: invalidAggregation{}},
	)

	instruments := []InstrumentProperties{
		{Name: "foo", Kind: instrumentKindUndefined}, //Unknown kind
		{Name: "foo", Kind: InstrumentKindSyncCounter},
		{Name: "foo", Kind: InstrumentKindSyncUpDownCounter},
		{Name: "foo", Kind: InstrumentKindSyncHistogram},
		{Name: "foo", Kind: InstrumentKindAsyncCounter},
		{Name: "foo", Kind: InstrumentKindAsyncUpDownCounter},
		{Name: "foo", Kind: InstrumentKindAsyncGauge},
	}

	testcases := []struct {
		name     string
		reader   Reader
		views    []View
		inst     InstrumentProperties
		wantKind internal.Aggregator[N] //Aggregators should match len and types
		wantLen  int
		wantErr  error
	}{
		{
			name:   "drop should return 0 aggregators",
			reader: NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Drop{} })),
			inst:   instruments[InstrumentKindSyncCounter],
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindSyncUpDownCounter],
			wantKind: internal.NewDeltaSum[N](false),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindSyncHistogram],
			wantKind: internal.NewDeltaHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindAsyncCounter],
			wantKind: internal.NewPrecomputedDeltaSum[N](true),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindAsyncUpDownCounter],
			wantKind: internal.NewPrecomputedDeltaSum[N](false),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindAsyncGauge],
			wantKind: internal.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []View{defaultAggView},
			inst:     instruments[InstrumentKindSyncCounter],
			wantKind: internal.NewDeltaSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			inst:     instruments[InstrumentKindSyncUpDownCounter],
			wantKind: internal.NewCumulativeSum[N](false),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			inst:     instruments[InstrumentKindSyncHistogram],
			wantKind: internal.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			inst:     instruments[InstrumentKindAsyncCounter],
			wantKind: internal.NewPrecomputedCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			inst:     instruments[InstrumentKindAsyncUpDownCounter],
			wantKind: internal.NewPrecomputedCumulativeSum[N](false),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			inst:     instruments[InstrumentKindAsyncGauge],
			wantKind: internal.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			inst:     instruments[InstrumentKindSyncCounter],
			wantKind: internal.NewCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "view should overwrite reader",
			reader:   NewManualReader(),
			views:    []View{changeAggView},
			inst:     instruments[InstrumentKindSyncCounter],
			wantKind: internal.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "multiple views should create multiple aggregators",
			reader:   NewManualReader(),
			views:    []View{defaultAggView, renameView},
			inst:     instruments[InstrumentKindSyncCounter],
			wantKind: internal.NewCumulativeSum[N](true),
			wantLen:  2,
		},
		{
			name:    "reader with invalid aggregation should error",
			reader:  NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			inst:    instruments[InstrumentKindSyncCounter],
			wantErr: errCreatingAggregators,
		},
		{
			name:    "view with invalid aggregation should error",
			reader:  NewManualReader(),
			views:   []View{invalidAggView},
			inst:    instruments[InstrumentKindSyncCounter],
			wantErr: errCreatingAggregators,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			c := newInstrumentCache[N](nil, nil)
			i := newInserter(newPipeline(nil, tt.reader, tt.views), c)
			got, err := i.Instrument(tt.inst)
			assert.ErrorIs(t, err, tt.wantErr)
			require.Len(t, got, tt.wantLen)
			for _, agg := range got {
				assert.IsType(t, tt.wantKind, agg)
			}
		})
	}
}

func testInvalidInstrumentShouldPanic[N int64 | float64]() {
	c := newInstrumentCache[N](nil, nil)
	i := newInserter(newPipeline(nil, NewManualReader(), nil), c)
	inst := InstrumentProperties{
		Name: "foo",
		Kind: InstrumentKind(255),
	}
	_, _ = i.Instrument(inst)
}

func TestInvalidInstrumentShouldPanic(t *testing.T) {
	assert.Panics(t, testInvalidInstrumentShouldPanic[int64])
	assert.Panics(t, testInvalidInstrumentShouldPanic[float64])
}

func TestCreateAggregators(t *testing.T) {
	t.Run("Int64", testCreateAggregators[int64])
	t.Run("Float64", testCreateAggregators[float64])
}

func TestPipelineRegistryCreateAggregators(t *testing.T) {
	defaultView := NewView(
		InstrumentProperties{Name: "*"},
		DataStream{},
	)
	renameView := NewView(
		InstrumentProperties{Name: "foo"},
		DataStream{InstrumentProperties: InstrumentProperties{Name: "bar"}},
	)
	testRdr := NewManualReader()
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	testCases := []struct {
		name      string
		views     map[Reader][]View
		inst      InstrumentProperties
		wantCount int
	}{
		{
			name: "No views have no aggregators",
			inst: InstrumentProperties{Name: "foo"},
		},
		{
			name: "1 reader 1 view gets 1 aggregator",
			inst: InstrumentProperties{Name: "foo"},
			views: map[Reader][]View{
				testRdr: {defaultView},
			},
			wantCount: 1,
		},
		{
			name: "1 reader 2 views gets 2 aggregator",
			inst: InstrumentProperties{Name: "foo"},
			views: map[Reader][]View{
				testRdr: {defaultView, renameView},
			},
			wantCount: 2,
		},
		{
			name: "2 readers 1 view each gets 2 aggregators",
			inst: InstrumentProperties{Name: "foo"},
			views: map[Reader][]View{
				testRdr:          {defaultView},
				testRdrHistogram: {defaultView},
			},
			wantCount: 2,
		},
		{
			name: "2 reader 2 views each gets 4 aggregators",
			inst: InstrumentProperties{Name: "foo"},
			views: map[Reader][]View{
				testRdr:          {defaultView, renameView},
				testRdrHistogram: {defaultView, renameView},
			},
			wantCount: 4,
		},
		{
			name: "An instrument is duplicated in two views share the same aggregator",
			inst: InstrumentProperties{Name: "foo"},
			views: map[Reader][]View{
				testRdr: {defaultView, defaultView},
			},
			wantCount: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p := newPipelines(resource.Empty(), tt.views)
			testPipelineRegistryResolveIntAggregators(t, p, tt.wantCount)
			p = newPipelines(resource.Empty(), tt.views)
			testPipelineRegistryResolveFloatAggregators(t, p, tt.wantCount)
		})
	}
}

func testPipelineRegistryResolveIntAggregators(t *testing.T, p pipelines, wantCount int) {
	inst := InstrumentProperties{Name: "foo", Kind: InstrumentKindSyncCounter}

	c := newInstrumentCache[int64](nil, nil)
	r := newResolver(p, c)
	aggs, err := r.Aggregators(inst)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func testPipelineRegistryResolveFloatAggregators(t *testing.T, p pipelines, wantCount int) {
	inst := InstrumentProperties{Name: "foo", Kind: InstrumentKindSyncCounter}

	c := newInstrumentCache[float64](nil, nil)
	r := newResolver(p, c)
	aggs, err := r.Aggregators(inst)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func TestPipelineRegistryResource(t *testing.T) {
	views := map[Reader][]View{
		NewManualReader(): {
			NewView(InstrumentProperties{Name: "*"}, DataStream{}),
			NewView(
				InstrumentProperties{Name: "foo"},
				DataStream{InstrumentProperties: InstrumentProperties{Name: "bar"}},
			),
		},
	}
	res := resource.NewSchemaless(attribute.String("key", "val"))
	pipes := newPipelines(res, views)
	for _, p := range pipes {
		assert.True(t, res.Equal(p.resource), "resource not set")
	}
}

func TestPipelineRegistryCreateAggregatorsIncompatibleInstrument(t *testing.T) {
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	views := map[Reader][]View{
		testRdrHistogram: {
			NewView(InstrumentProperties{Name: "*"}, DataStream{}),
		},
	}
	p := newPipelines(resource.Empty(), views)
	inst := InstrumentProperties{Name: "foo", Kind: InstrumentKindAsyncGauge}

	vc := cache[string, instrumentID]{}
	ri := newResolver(p, newInstrumentCache[int64](nil, &vc))
	intAggs, err := ri.Aggregators(inst)
	assert.Error(t, err)
	assert.Len(t, intAggs, 0)

	p = newPipelines(resource.Empty(), views)

	rf := newResolver(p, newInstrumentCache[float64](nil, &vc))
	floatAggs, err := rf.Aggregators(inst)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 0)
}

type logCounter struct {
	logr.LogSink

	infoN uint32
}

func (l *logCounter) Info(level int, msg string, keysAndValues ...interface{}) {
	atomic.AddUint32(&l.infoN, 1)
	l.LogSink.Info(level, msg, keysAndValues...)
}

func (l *logCounter) InfoN() int {
	return int(atomic.SwapUint32(&l.infoN, 0))
}

func TestResolveAggregatorsDuplicateErrors(t *testing.T) {
	tLog := testr.NewWithOptions(t, testr.Options{Verbosity: 6})
	l := &logCounter{LogSink: tLog.GetSink()}
	otel.SetLogger(logr.New(l))

	views := map[Reader][]View{
		NewManualReader(): {
			NewView(InstrumentProperties{Name: "*"}, DataStream{}),
			NewView(
				InstrumentProperties{Name: "bar"},
				DataStream{InstrumentProperties: InstrumentProperties{Name: "foo"}},
			),
		},
	}

	fooInst := InstrumentProperties{Name: "foo", Kind: InstrumentKindSyncCounter}
	barInst := InstrumentProperties{Name: "bar", Kind: InstrumentKindSyncCounter}

	p := newPipelines(resource.Empty(), views)

	vc := cache[string, instrumentID]{}
	ri := newResolver(p, newInstrumentCache[int64](nil, &vc))
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
	rf := newResolver(p, newInstrumentCache[float64](nil, &vc))
	floatAggs, err := rf.Aggregators(fooInst)
	assert.NoError(t, err)
	assert.Equal(t, 1, l.InfoN(), "instrument conflict not logged")
	assert.Len(t, floatAggs, 1)

	fooInst = InstrumentProperties{Name: "foo-float", Kind: InstrumentKindSyncCounter}

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
	testCases := []struct {
		name string
		kind InstrumentKind
		agg  aggregation.Aggregation
		want error
	}{
		{
			name: "SyncCounter and Drop",
			kind: InstrumentKindSyncCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncCounter and LastValue",
			kind: InstrumentKindSyncCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncCounter and Sum",
			kind: InstrumentKindSyncCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncCounter and ExplicitBucketHistogram",
			kind: InstrumentKindSyncCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
		},
		{
			name: "SyncUpDownCounter and Drop",
			kind: InstrumentKindSyncUpDownCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncUpDownCounter and LastValue",
			kind: InstrumentKindSyncUpDownCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncUpDownCounter and Sum",
			kind: InstrumentKindSyncUpDownCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncUpDownCounter and ExplicitBucketHistogram",
			kind: InstrumentKindSyncUpDownCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncHistogram and Drop",
			kind: InstrumentKindSyncHistogram,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncHistogram and LastValue",
			kind: InstrumentKindSyncHistogram,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncHistogram and Sum",
			kind: InstrumentKindSyncHistogram,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncHistogram and ExplicitBucketHistogram",
			kind: InstrumentKindSyncHistogram,
			agg:  aggregation.ExplicitBucketHistogram{},
		},
		{
			name: "AsyncCounter and Drop",
			kind: InstrumentKindAsyncCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "AsyncCounter and LastValue",
			kind: InstrumentKindAsyncCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncCounter and Sum",
			kind: InstrumentKindAsyncCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "AsyncCounter and ExplicitBucketHistogram",
			kind: InstrumentKindAsyncCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncUpDownCounter and Drop",
			kind: InstrumentKindAsyncUpDownCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "AsyncUpDownCounter and LastValue",
			kind: InstrumentKindAsyncUpDownCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncUpDownCounter and Sum",
			kind: InstrumentKindAsyncUpDownCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "AsyncUpDownCounter and ExplicitBucketHistogram",
			kind: InstrumentKindAsyncUpDownCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncGauge and Drop",
			kind: InstrumentKindAsyncGauge,
			agg:  aggregation.Drop{},
		},
		{
			name: "AsyncGauge and aggregation.LastValue{}",
			kind: InstrumentKindAsyncGauge,
			agg:  aggregation.LastValue{},
		},
		{
			name: "AsyncGauge and Sum",
			kind: InstrumentKindAsyncGauge,
			agg:  aggregation.Sum{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncGauge and ExplicitBucketHistogram",
			kind: InstrumentKindAsyncGauge,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "Default aggregation should error",
			kind: InstrumentKindSyncCounter,
			agg:  aggregation.Default{},
			want: errUnknownAggregation,
		},
		{
			name: "unknown kind with Sum should error",
			kind: instrumentKindUndefined,
			agg:  aggregation.Sum{},
			want: errIncompatibleAggregation,
		},
		{
			name: "unknown kind with LastValue should error",
			kind: instrumentKindUndefined,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "unknown kind with Histogram should error",
			kind: instrumentKindUndefined,
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
