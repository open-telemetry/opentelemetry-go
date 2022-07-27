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
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	renameView, _ := view.New(
		view.MatchInstrumentName("foo"),
		view.WithRename("bar"),
	)
	testcases := []struct {
		name     string
		reader   Reader
		views    []view.View
		inst     view.Instrument
		wantKind internal.Aggregator[N] //Aggregators should match len and types
		wantLen  int
	}{
		{
			name: "drop should return 0 aggregators",
			reader: testReader{
				agg: aggregation.Drop{},
			},
			views: []view.View{{}},
			inst:  view.Instrument{Name: "foo"},
		},
		{
			name: "reader should set default agg",
			reader: testReader{
				agg:  aggregation.Sum{},
				temp: metricdata.DeltaTemporality,
			},
			views:    []view.View{{}},
			inst:     view.Instrument{Name: "foo"},
			wantKind: internal.NewDeltaSum[N](),
			wantLen:  1,
		},
		{
			name: "view should overwrite reader",
			reader: testReader{
				agg:  aggregation.Sum{},
				temp: metricdata.DeltaTemporality,
			},
			views:    []view.View{changeAggView},
			inst:     view.Instrument{Name: "foo"},
			wantKind: internal.NewLastValue[N](),
			wantLen:  1,
		},
		{
			name: "multiple views should create multiple aggregators",
			reader: testReader{
				agg:  aggregation.Sum{},
				temp: metricdata.DeltaTemporality,
			},
			views:    []view.View{{}, renameView},
			inst:     view.Instrument{Name: "foo"},
			wantKind: internal.NewDeltaSum[N](),
			wantLen:  2,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			got := createAggregators[N](tt.reader, tt.views, tt.inst)
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
				testReader{agg: aggregation.LastValue{}}: {
					{},
				},
			},
			wantCount: 1,
		},
		{
			name: "1 reader 2 views gets 2 aggregator",
			inst: view.Instrument{Name: "foo"},
			views: map[Reader][]view.View{
				testReader{agg: aggregation.LastValue{}}: {
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
				testReader{agg: aggregation.LastValue{}}: {
					{},
				},
				testReader{agg: aggregation.LastValue{}, temp: metricdata.CumulativeTemporality}: {
					{},
				},
			},
			wantCount: 2,
		},
		{
			name: "2 reader 2 views each gets 4 aggregators",
			inst: view.Instrument{Name: "foo"},
			views: map[Reader][]view.View{
				testReader{agg: aggregation.LastValue{}}: {
					{},
					renameView,
				},
				testReader{agg: aggregation.LastValue{}, temp: metricdata.CumulativeTemporality}: {
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
				testReader{agg: aggregation.LastValue{}}: {
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
	inst := view.Instrument{Name: "foo"}
	want := make([]internal.Aggregator[N], wantCount)
	for i := range want {
		want[i] = internal.NewLastValue[N]()
	}

	aggs, err := reg.createAggregators(inst, unit.Dimensionless)
	assert.NoError(t, err)

	require.Len(t, aggs, wantCount)
	for i, agg := range aggs {
		assert.IsType(t, want[i], agg)
	}
}

func TestPipelineRegistryCreateAggregatorsDuplicateErrors(t *testing.T) {
	renameView, _ := view.New(
		view.MatchInstrumentName("bar"),
		view.WithRename("foo"),
	)
	views := map[Reader][]view.View{
		testReader{agg: aggregation.LastValue{}}: {
			{},
			renameView,
		},
	}

	fooInst := view.Instrument{Name: "foo"}
	barInst := view.Instrument{Name: "bar"}

	intReg, floatReg := newPipelineRegistries(views)

	_, err := intReg.createAggregators(fooInst, unit.Dimensionless)
	assert.NoError(t, err)

	// The Rename view should error, because it creates a foo instrument.
	intAggs, err := intReg.createAggregators(barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, intAggs, 2)

	// Creating a float foo instrument should error because there is an int foo instrument.
	floatAggs, err := floatReg.createAggregators(fooInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 1)

	fooInst = view.Instrument{Name: "foo-float"}

	_, err = floatReg.createAggregators(fooInst, unit.Dimensionless)
	assert.NoError(t, err)

	floatAggs, err = floatReg.createAggregators(barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 2)
}
