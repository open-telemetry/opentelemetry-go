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
		name           string
		views          map[Reader][]view.View
		inst           view.Instrument
		wantInt64Agg   []internal.Aggregator[int64]   // Should match len and type
		wantFloat64Agg []internal.Aggregator[float64] // Should match len and type
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
			wantInt64Agg: []internal.Aggregator[int64]{
				internal.NewLastValue[int64](),
			},
			wantFloat64Agg: []internal.Aggregator[float64]{
				internal.NewLastValue[float64](),
			},
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
			wantInt64Agg: []internal.Aggregator[int64]{
				internal.NewLastValue[int64](),
				internal.NewLastValue[int64](),
			},
			wantFloat64Agg: []internal.Aggregator[float64]{
				internal.NewLastValue[float64](),
				internal.NewLastValue[float64](),
			},
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
			wantInt64Agg: []internal.Aggregator[int64]{
				internal.NewLastValue[int64](),
				internal.NewLastValue[int64](),
			},
			wantFloat64Agg: []internal.Aggregator[float64]{
				internal.NewLastValue[float64](),
				internal.NewLastValue[float64](),
			},
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
			wantInt64Agg: []internal.Aggregator[int64]{
				internal.NewLastValue[int64](),
				internal.NewLastValue[int64](),
				internal.NewLastValue[int64](),
				internal.NewLastValue[int64](),
			},
			wantFloat64Agg: []internal.Aggregator[float64]{
				internal.NewLastValue[float64](),
				internal.NewLastValue[float64](),
				internal.NewLastValue[float64](),
				internal.NewLastValue[float64](),
			},
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
			wantInt64Agg: []internal.Aggregator[int64]{
				internal.NewLastValue[int64](),
			},
			wantFloat64Agg: []internal.Aggregator[float64]{
				internal.NewLastValue[float64](),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			reg := newPipelineRegistry(tt.views)

			intAggs, err := reg.createInt64Aggregators(tt.inst, unit.Dimensionless)
			assert.NoError(t, err)

			require.Len(t, intAggs, len(tt.wantInt64Agg))
			for i, agg := range intAggs {
				assert.IsType(t, tt.wantInt64Agg[i], agg)
			}

			reg = newPipelineRegistry(tt.views)

			floatAggs, err := reg.createFloat64Aggregators(tt.inst, unit.Dimensionless)
			assert.NoError(t, err)

			require.Len(t, floatAggs, len(tt.wantFloat64Agg))
			for i, agg := range floatAggs {
				assert.IsType(t, tt.wantFloat64Agg[i], agg)
			}
		})
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

	reg := newPipelineRegistry(views)

	_, err := reg.createInt64Aggregators(fooInst, unit.Dimensionless)
	assert.NoError(t, err)

	intAggs, err := reg.createInt64Aggregators(barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, intAggs, 2)

	reg = newPipelineRegistry(views)
	_, err = reg.createFloat64Aggregators(fooInst, unit.Dimensionless)
	assert.NoError(t, err)

	floatAggs, err := reg.createFloat64Aggregators(barInst, unit.Dimensionless)
	assert.Error(t, err)
	assert.Len(t, floatAggs, 2)
}
