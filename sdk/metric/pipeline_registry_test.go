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

//go:build go1.18
// +build go1.18

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/metric/view"
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
	changeAggView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithSetAggregation(aggregation.ExplicitBucketHistogram{}),
	)
	renameView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithRename("bar"),
	)
	defaultAggView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithSetAggregation(aggregation.Default{}),
	)
	invalidAggView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithSetAggregation(invalidAggregation{}),
	)

	instruments := []view.Instrument{
		{Name: "foo", Kind: view.InstrumentKind(0)}, //Unknown kind
		{Name: "foo", Kind: view.SyncCounter},
		{Name: "foo", Kind: view.SyncUpDownCounter},
		{Name: "foo", Kind: view.SyncHistogram},
		{Name: "foo", Kind: view.AsyncCounter},
		{Name: "foo", Kind: view.AsyncUpDownCounter},
		{Name: "foo", Kind: view.AsyncGauge},
	}

	testcases := []struct {
		name     string
		reader   Reader
		views    []view.View
		inst     view.Instrument
		wantKind internal.Aggregator[N] //Aggregators should match len and types
		wantLen  int
		wantErr  error
	}{
		{
			name:   "drop should return 0 aggregators",
			reader: NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.Drop{} })),
			views:  []view.View{{}},
			inst:   instruments[view.SyncCounter],
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			inst:     instruments[view.SyncUpDownCounter],
			wantKind: internal.NewDeltaSum[N](false),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			inst:     instruments[view.SyncHistogram],
			wantKind: internal.NewDeltaHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			inst:     instruments[view.AsyncCounter],
			wantKind: internal.NewDeltaSum[N](true),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			inst:     instruments[view.AsyncUpDownCounter],
			wantKind: internal.NewDeltaSum[N](false),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			inst:     instruments[view.AsyncGauge],
			wantKind: internal.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			inst:     instruments[view.SyncCounter],
			wantKind: internal.NewDeltaSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			inst:     instruments[view.SyncUpDownCounter],
			wantKind: internal.NewCumulativeSum[N](false),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			inst:     instruments[view.SyncHistogram],
			wantKind: internal.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			inst:     instruments[view.AsyncCounter],
			wantKind: internal.NewCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			inst:     instruments[view.AsyncUpDownCounter],
			wantKind: internal.NewCumulativeSum[N](false),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			inst:     instruments[view.AsyncGauge],
			wantKind: internal.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			inst:     instruments[view.SyncCounter],
			wantKind: internal.NewCumulativeSum[N](true),
			wantLen:  1,
		},
		{
			name:     "view should overwrite reader",
			reader:   NewManualReader(),
			views:    []view.View{changeAggView},
			inst:     instruments[view.SyncCounter],
			wantKind: internal.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name:     "multiple views should create multiple aggregators",
			reader:   NewManualReader(),
			views:    []view.View{{}, renameView},
			inst:     instruments[view.SyncCounter],
			wantKind: internal.NewCumulativeSum[N](true),
			wantLen:  2,
		},
		{
			name:    "reader with invalid aggregation should error",
			reader:  NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.Default{} })),
			views:   []view.View{{}},
			inst:    instruments[view.SyncCounter],
			wantErr: errCreatingAggregators,
		},
		{
			name:    "view with invalid aggregation should error",
			reader:  NewManualReader(),
			views:   []view.View{invalidAggView},
			inst:    instruments[view.SyncCounter],
			wantErr: errCreatingAggregators,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAggregatorsForReader[N](tt.reader, tt.views, tt.inst)
			assert.ErrorIs(t, err, tt.wantErr)
			require.Len(t, got, tt.wantLen)
			for _, agg := range got {
				assert.IsType(t, tt.wantKind, agg)
			}
		})
	}
}

func testInvalidInstrumentShouldPanic[N int64 | float64]() {
	reader := NewManualReader()
	views := []view.View{{}}
	inst := view.Instrument{
		Name: "foo",
		Kind: view.InstrumentKind(255),
	}
	_, _ = createAggregatorsForReader[N](reader, views, inst)
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
	renameView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithRename("bar"),
	)
	testRdr := NewManualReader()
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	testCases := []struct {
		name      string
		entries   []registryEntry
		inst      view.Instrument
		wantCount int
	}{
		{
			name: "No views have no aggregators",
			inst: view.Instrument{Name: "foo"},
		},
		{
			name: "1 reader 1 view gets 1 aggregator",
			inst: view.Instrument{Name: "foo"},
			entries: []registryEntry{
				{reader: testRdr, views: []view.View{{}}},
			},
			wantCount: 1,
		},
		{
			name: "1 reader 2 views gets 2 aggregator",
			inst: view.Instrument{Name: "foo"},
			entries: []registryEntry{
				{
					reader: testRdr,
					views: []view.View{
						{},
						renameView,
					},
				},
			},
			wantCount: 2,
		},
		{
			name: "2 readers 1 view each gets 2 aggregators",
			inst: view.Instrument{Name: "foo"},
			entries: []registryEntry{
				{
					reader: testRdr,
					views: []view.View{
						{},
					},
				},
				{
					reader: testRdrHistogram,
					views: []view.View{
						{},
					},
				},
			},
			wantCount: 2,
		},
		{
			name: "2 reader 2 views each gets 4 aggregators",
			inst: view.Instrument{Name: "foo"},
			entries: []registryEntry{
				{
					reader: testRdr,
					views: []view.View{
						{},
						renameView,
					},
				},
				{
					reader: testRdrHistogram,
					views: []view.View{
						{},
						renameView,
					},
				},
			},
			wantCount: 4,
		},
		{
			name: "An instrument is duplicated in two views share the same aggregator",
			inst: view.Instrument{Name: "foo"},
			entries: []registryEntry{
				{
					reader: testRdr,
					views: []view.View{
						{},
						{},
					},
				},
			},
			wantCount: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			reg := newPipelineRegistries(tt.entries)
			testPipelineRegistryCreateIntAggregators(t, reg, tt.wantCount)
			reg = newPipelineRegistries(tt.entries)
			testPipelineRegistryCreateFloatAggregators(t, reg, tt.wantCount)
		})
	}
}

func testPipelineRegistryCreateIntAggregators(t *testing.T, reg *pipelineRegistry, wantCount int) {
	inst := view.Instrument{Name: "foo", Kind: view.SyncCounter}

	aggs, err := createAggregators[int64](reg, inst, unit.Dimensionless)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func testPipelineRegistryCreateFloatAggregators(t *testing.T, reg *pipelineRegistry, wantCount int) {
	inst := view.Instrument{Name: "foo", Kind: view.SyncCounter}

	aggs, err := createAggregators[float64](reg, inst, unit.Dimensionless)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func TestPipelineRegistryCreateAggregatorsIncompatibleInstrument(t *testing.T) {
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	entries := []registryEntry{
		registryEntry{reader: testRdrHistogram, views: []view.View{{}}},
	}
	reg := newPipelineRegistries(entries)
	inst := view.Instrument{Name: "foo", Kind: view.AsyncGauge}

	intAggs, err := createAggregators[int64](reg, inst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, intAggs, 0)

	reg = newPipelineRegistries(entries)

	floatAggs, err := createAggregators[float64](reg, inst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 0)
}

func TestPipelineRegistryCreateAggregatorsDuplicateErrors(t *testing.T) {
	renameView, _ := view.New(
		view.MatchInstrumentName("bar"),
		view.WithRename("foo"),
	)

	entries := []registryEntry{
		{
			reader: NewManualReader(),
			views: []view.View{
				{},
				renameView,
			},
		},
	}

	fooInst := view.Instrument{Name: "foo", Kind: view.SyncCounter}
	barInst := view.Instrument{Name: "bar", Kind: view.SyncCounter}

	reg := newPipelineRegistries(entries)

	intAggs, err := createAggregators[int64](reg, fooInst, unit.Dimensionless)
	assert.NoError(t, err)
	assert.Len(t, intAggs, 1)

	// The Rename view should error, because it creates a foo instrument.
	intAggs, err = createAggregators[int64](reg, barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, intAggs, 2)

	// Creating a float foo instrument should error because there is an int foo instrument.
	floatAggs, err := createAggregators[float64](reg, fooInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 1)

	fooInst = view.Instrument{Name: "foo-float", Kind: view.SyncCounter}

	_, err = createAggregators[float64](reg, fooInst, unit.Dimensionless)
	assert.NoError(t, err)

	floatAggs, err = createAggregators[float64](reg, barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 2)
}

func TestIsAggregatorCompatible(t *testing.T) {
	var undefinedInstrument view.InstrumentKind

	testCases := []struct {
		name string
		kind view.InstrumentKind
		agg  aggregation.Aggregation
		want error
	}{
		{
			name: "SyncCounter and Drop",
			kind: view.SyncCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncCounter and LastValue",
			kind: view.SyncCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncCounter and Sum",
			kind: view.SyncCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncCounter and ExplicitBucketHistogram",
			kind: view.SyncCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
		},
		{
			name: "SyncUpDownCounter and Drop",
			kind: view.SyncUpDownCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncUpDownCounter and LastValue",
			kind: view.SyncUpDownCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncUpDownCounter and Sum",
			kind: view.SyncUpDownCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncUpDownCounter and ExplicitBucketHistogram",
			kind: view.SyncUpDownCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncHistogram and Drop",
			kind: view.SyncHistogram,
			agg:  aggregation.Drop{},
		},
		{
			name: "SyncHistogram and LastValue",
			kind: view.SyncHistogram,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "SyncHistogram and Sum",
			kind: view.SyncHistogram,
			agg:  aggregation.Sum{},
		},
		{
			name: "SyncHistogram and ExplicitBucketHistogram",
			kind: view.SyncHistogram,
			agg:  aggregation.ExplicitBucketHistogram{},
		},
		{
			name: "AsyncCounter and Drop",
			kind: view.AsyncCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "AsyncCounter and LastValue",
			kind: view.AsyncCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncCounter and Sum",
			kind: view.AsyncCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "AsyncCounter and ExplicitBucketHistogram",
			kind: view.AsyncCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncUpDownCounter and Drop",
			kind: view.AsyncUpDownCounter,
			agg:  aggregation.Drop{},
		},
		{
			name: "AsyncUpDownCounter and LastValue",
			kind: view.AsyncUpDownCounter,
			agg:  aggregation.LastValue{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncUpDownCounter and Sum",
			kind: view.AsyncUpDownCounter,
			agg:  aggregation.Sum{},
		},
		{
			name: "AsyncUpDownCounter and ExplicitBucketHistogram",
			kind: view.AsyncUpDownCounter,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncGauge and Drop",
			kind: view.AsyncGauge,
			agg:  aggregation.Drop{},
		},
		{
			name: "AsyncGauge and aggregation.LastValue{}",
			kind: view.AsyncGauge,
			agg:  aggregation.LastValue{},
		},
		{
			name: "AsyncGauge and Sum",
			kind: view.AsyncGauge,
			agg:  aggregation.Sum{},
			want: errIncompatibleAggregation,
		},
		{
			name: "AsyncGauge and ExplicitBucketHistogram",
			kind: view.AsyncGauge,
			agg:  aggregation.ExplicitBucketHistogram{},
			want: errIncompatibleAggregation,
		},
		{
			name: "Default aggregation should error",
			kind: view.SyncCounter,
			agg:  aggregation.Default{},
			want: errUnknownAggregation,
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
