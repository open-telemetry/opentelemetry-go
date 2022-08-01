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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

type testReader struct {
	agg  aggregation.Aggregation
	temp metricdata.Temporality
}

func (t testReader) register(producer)                                       {}
func (t testReader) temporality(view.InstrumentKind) metricdata.Temporality  { return t.temp }
func (t testReader) aggregation(view.InstrumentKind) aggregation.Aggregation { return t.agg } // nolint:revive  // import-shadow for method scoped by type.
func (t testReader) Collect(context.Context) (metricdata.ResourceMetrics, error) {
	return metricdata.ResourceMetrics{}, nil
}
func (t testReader) ForceFlush(context.Context) error { return nil }
func (t testReader) Shutdown(context.Context) error   { return nil }

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

	inst := view.Instrument{
		Name: "foo",
		Kind: view.SyncCounter,
	}

	testcases := []struct {
		name     string
		reader   Reader
		views    []view.View
		wantKind internal.Aggregator[N] //Aggregators should match len and types
		wantLen  int
	}{
		{
			name:   "drop should return 0 aggregators",
			reader: NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.Drop{} })),
			views:  []view.View{{}},
		},
		{
			name:     "default agg should use reader",
			reader:   NewManualReader(WithTemporalitySelector(deltaTemporalitySelector)),
			views:    []view.View{defaultAggView},
			wantKind: internal.NewDeltaSum[N](),
			wantLen:  1,
		},
		{
			name:     "reader should set default agg",
			reader:   NewManualReader(),
			views:    []view.View{{}},
			wantKind: internal.NewCumulativeSum[N](),
			wantLen:  1,
		},
		{
			name:     "view should overwrite reader",
			reader:   NewManualReader(),
			views:    []view.View{changeAggView},
			wantKind: internal.NewCumulativeHistogram[N](aggregation.ExplicitBucketHistogram{}),
			wantLen:  1,
		},
		{
			name: "multiple views should create multiple aggregators",
			reader: testReader{
				agg:  aggregation.Sum{},
				temp: metricdata.DeltaTemporality,
			},
			views:    []view.View{{}, renameView},
			wantKind: internal.NewDeltaSum[N](),
			wantLen:  2,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAggregators[N](tt.reader, tt.views, inst)
			assert.NoError(t, err)
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

func TestPipelineRegistryCreateAggregators(t *testing.T) {
	renameView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithRename("bar"),
	)
	testRdr := NewManualReader()
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	testCases := []struct {
		name      string
		views     map[Reader][]view.View
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
			views: map[Reader][]view.View{
				testRdr: {
					{},
				},
			},
			wantCount: 1,
		},
		{
			name: "1 reader 2 views gets 2 aggregator",
			inst: view.Instrument{Name: "foo"},
			views: map[Reader][]view.View{
				testRdr: {
					{},
					renameView,
				},
			},
			wantCount: 2,
		},
		{
			name: "2 readers 1 view each gets 2 aggregators",
			inst: view.Instrument{Name: "foo"},
			views: map[Reader][]view.View{
				testRdr: {
					{},
				},
				testRdrHistogram: {
					{},
				},
			},
			wantCount: 2,
		},
		{
			name: "2 reader 2 views each gets 4 aggregators",
			inst: view.Instrument{Name: "foo"},
			views: map[Reader][]view.View{
				testRdr: {
					{},
					renameView,
				},
				testRdrHistogram: {
					{},
					renameView,
				},
			},
			wantCount: 4,
		},
		{
			name: "An instrument is duplicated in two views share the same aggregator",
			inst: view.Instrument{Name: "foo"},
			views: map[Reader][]view.View{
				testRdr: {
					{},
					{},
				},
			},
			wantCount: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			intReg, _ := newPipelineRegistries(tt.views)
			testPipelineRegistryCreateAggregators(t, intReg, tt.wantCount)
			_, floatReg := newPipelineRegistries(tt.views)
			testPipelineRegistryCreateAggregators(t, floatReg, tt.wantCount)
		})
	}
}

func testPipelineRegistryCreateAggregators[N int64 | float64](t *testing.T, reg *pipelineRegistry[N], wantCount int) {
	inst := view.Instrument{Name: "foo", Kind: view.SyncCounter}

	aggs, err := reg.createAggregators(inst, unit.Dimensionless)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
}

func TestPipelineRegistryCreateAggregatorsIncompatibleInstrument(t *testing.T) {
	testRdrHistogram := NewManualReader(WithAggregationSelector(func(ik view.InstrumentKind) aggregation.Aggregation { return aggregation.ExplicitBucketHistogram{} }))

	views := map[Reader][]view.View{
		testRdrHistogram: {
			{},
		},
	}
	intReg, floatReg := newPipelineRegistries(views)
	inst := view.Instrument{Name: "foo", Kind: view.AsyncGauge}

	intAggs, err := intReg.createAggregators(inst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, intAggs, 0)

	floatAggs, err := floatReg.createAggregators(inst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 0)
}

func TestPipelineRegistryCreateAggregatorsDuplicateErrors(t *testing.T) {
	renameView, _ := view.New(
		view.MatchInstrumentName("bar"),
		view.WithRename("foo"),
	)
	views := map[Reader][]view.View{
		testReader{agg: aggregation.Sum{}}: {
			{},
			renameView,
		},
	}

	fooInst := view.Instrument{Name: "foo", Kind: view.SyncCounter}
	barInst := view.Instrument{Name: "bar", Kind: view.SyncCounter}

	intReg, floatReg := newPipelineRegistries(views)

	intAggs, err := intReg.createAggregators(fooInst, unit.Dimensionless)
	assert.NoError(t, err)
	assert.Len(t, intAggs, 1)

	// The Rename view should error, because it creates a foo instrument.
	intAggs, err = intReg.createAggregators(barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, intAggs, 2)

	// Creating a float foo instrument should error because there is an int foo instrument.
	floatAggs, err := floatReg.createAggregators(fooInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 1)

	fooInst = view.Instrument{Name: "foo-float", Kind: view.SyncCounter}

	_, err = floatReg.createAggregators(fooInst, unit.Dimensionless)
	assert.NoError(t, err)

	floatAggs, err = floatReg.createAggregators(barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 2)
}

func TestIsAggregatorCompatible(t *testing.T) {
	var undefinedInstrument view.InstrumentKind

	testCases := []struct {
		name    string
		kind    view.InstrumentKind
		agg     aggregation.Aggregation
		wantErr bool
	}{
		{
			name:    "SyncCounter and Drop",
			kind:    view.SyncCounter,
			agg:     aggregation.Drop{},
			wantErr: false,
		},
		{
			name:    "SyncCounter and LastValue",
			kind:    view.SyncCounter,
			agg:     aggregation.LastValue{},
			wantErr: true,
		},
		{
			name:    "SyncCounter and Sum",
			kind:    view.SyncCounter,
			agg:     aggregation.Sum{},
			wantErr: false,
		},
		{
			name:    "SyncCounter and ExplicitBucketHistogram",
			kind:    view.SyncCounter,
			agg:     aggregation.ExplicitBucketHistogram{},
			wantErr: false,
		},
		{
			name:    "SyncUpDownCounter and Drop",
			kind:    view.SyncUpDownCounter,
			agg:     aggregation.Drop{},
			wantErr: false,
		},
		{
			name:    "SyncUpDownCounter and LastValue",
			kind:    view.SyncUpDownCounter,
			agg:     aggregation.LastValue{},
			wantErr: true,
		},
		{
			name:    "SyncUpDownCounter and Sum",
			kind:    view.SyncUpDownCounter,
			agg:     aggregation.Sum{},
			wantErr: false,
		},
		{
			name:    "SyncUpDownCounter and ExplicitBucketHistogram",
			kind:    view.SyncUpDownCounter,
			agg:     aggregation.ExplicitBucketHistogram{},
			wantErr: true,
		},
		{
			name:    "SyncHistogram and Drop",
			kind:    view.SyncHistogram,
			agg:     aggregation.Drop{},
			wantErr: false,
		},
		{
			name:    "SyncHistogram and LastValue",
			kind:    view.SyncHistogram,
			agg:     aggregation.LastValue{},
			wantErr: true,
		},
		{
			name:    "SyncHistogram and Sum",
			kind:    view.SyncHistogram,
			agg:     aggregation.Sum{},
			wantErr: false,
		},
		{
			name:    "SyncHistogram and ExplicitBucketHistogram",
			kind:    view.SyncHistogram,
			agg:     aggregation.ExplicitBucketHistogram{},
			wantErr: false,
		},
		{
			name:    "AsyncCounter and Drop",
			kind:    view.AsyncCounter,
			agg:     aggregation.Drop{},
			wantErr: false,
		},
		{
			name:    "AsyncCounter and LastValue",
			kind:    view.AsyncCounter,
			agg:     aggregation.LastValue{},
			wantErr: true,
		},
		{
			name:    "AsyncCounter and Sum",
			kind:    view.AsyncCounter,
			agg:     aggregation.Sum{},
			wantErr: false,
		},
		{
			name:    "AsyncCounter and ExplicitBucketHistogram",
			kind:    view.AsyncCounter,
			agg:     aggregation.ExplicitBucketHistogram{},
			wantErr: true,
		},
		{
			name:    "AsyncUpDownCounter and Drop",
			kind:    view.AsyncUpDownCounter,
			agg:     aggregation.Drop{},
			wantErr: false,
		},
		{
			name:    "AsyncUpDownCounter and LastValue",
			kind:    view.AsyncUpDownCounter,
			agg:     aggregation.LastValue{},
			wantErr: true,
		},
		{
			name:    "AsyncUpDownCounter and Sum",
			kind:    view.AsyncUpDownCounter,
			agg:     aggregation.Sum{},
			wantErr: false,
		},
		{
			name:    "AsyncUpDownCounter and ExplicitBucketHistogram",
			kind:    view.AsyncUpDownCounter,
			agg:     aggregation.ExplicitBucketHistogram{},
			wantErr: true,
		},
		{
			name:    "AsyncGauge and Drop",
			kind:    view.AsyncGauge,
			agg:     aggregation.Drop{},
			wantErr: false,
		},
		{
			name:    "AsyncGauge and aggregation.LastValue{}",
			kind:    view.AsyncGauge,
			agg:     aggregation.LastValue{},
			wantErr: false,
		},
		{
			name:    "AsyncGauge and Sum",
			kind:    view.AsyncGauge,
			agg:     aggregation.Sum{},
			wantErr: true,
		},
		{
			name:    "AsyncGauge and ExplicitBucketHistogram",
			kind:    view.AsyncGauge,
			agg:     aggregation.ExplicitBucketHistogram{},
			wantErr: true,
		},
		{
			name:    "Default aggregation should error",
			kind:    view.SyncCounter,
			agg:     aggregation.Default{},
			wantErr: true,
		},
		{
			name:    "unknown kind should error",
			kind:    undefinedInstrument,
			agg:     aggregation.Sum{},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := isAggregatorCompatible(tt.kind, tt.agg)

			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}
