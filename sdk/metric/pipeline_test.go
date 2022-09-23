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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testSumAggregator struct{}

func (testSumAggregator) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint[int64]{}}
}

func TestEmptyPipeline(t *testing.T) {
	pipe := &pipeline{}

	output, err := pipe.produce(context.Background())
	require.NoError(t, err)
	assert.Nil(t, output.Resource)
	assert.Len(t, output.ScopeMetrics, 0)

	err = pipe.addAggregator(instrumentation.Scope{}, "name", "desc", unit.Dimensionless, testSumAggregator{})
	assert.NoError(t, err)

	require.NotPanics(t, func() {
		pipe.addCallback(func(ctx context.Context) {})
	})

	output, err = pipe.produce(context.Background())
	require.NoError(t, err)
	assert.Nil(t, output.Resource)
	require.Len(t, output.ScopeMetrics, 1)
	require.Len(t, output.ScopeMetrics[0].Metrics, 1)
}

func TestNewPipeline(t *testing.T) {
	pipe := newPipeline(nil, nil, nil)

	output, err := pipe.produce(context.Background())
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	assert.Len(t, output.ScopeMetrics, 0)

	err = pipe.addAggregator(instrumentation.Scope{}, "name", "desc", unit.Dimensionless, testSumAggregator{})
	assert.NoError(t, err)

	require.NotPanics(t, func() {
		pipe.addCallback(func(ctx context.Context) {})
	})

	output, err = pipe.produce(context.Background())
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	require.Len(t, output.ScopeMetrics, 1)
	require.Len(t, output.ScopeMetrics[0].Metrics, 1)
}

func TestPipelineDuplicateRegistration(t *testing.T) {
	type instrumentID struct {
		scope       instrumentation.Scope
		name        string
		description string
		unit        unit.Unit
	}
	testCases := []struct {
		name           string
		secondInst     instrumentID
		want           error
		wantScopeLen   int
		wantMetricsLen int
	}{
		{
			name: "exact should error",
			secondInst: instrumentID{
				scope:       instrumentation.Scope{},
				name:        "name",
				description: "desc",
				unit:        unit.Dimensionless,
			},
			want:           errAlreadyRegistered,
			wantScopeLen:   1,
			wantMetricsLen: 1,
		},
		{
			name: "description should not be identifying",
			secondInst: instrumentID{
				scope:       instrumentation.Scope{},
				name:        "name",
				description: "other desc",
				unit:        unit.Dimensionless,
			},
			want:           errAlreadyRegistered,
			wantScopeLen:   1,
			wantMetricsLen: 1,
		},
		{
			name: "scope should be identifying",
			secondInst: instrumentID{
				scope: instrumentation.Scope{
					Name: "newScope",
				},
				name:        "name",
				description: "desc",
				unit:        unit.Dimensionless,
			},
			wantScopeLen:   2,
			wantMetricsLen: 1,
		},
		{
			name: "name should be identifying",
			secondInst: instrumentID{
				scope:       instrumentation.Scope{},
				name:        "newName",
				description: "desc",
				unit:        unit.Dimensionless,
			},
			wantScopeLen:   1,
			wantMetricsLen: 2,
		},
		{
			name: "unit should be identifying",
			secondInst: instrumentID{
				scope:       instrumentation.Scope{},
				name:        "name",
				description: "desc",
				unit:        unit.Bytes,
			},
			wantScopeLen:   1,
			wantMetricsLen: 2,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			pipe := newPipeline(nil, nil, nil)
			err := pipe.addAggregator(instrumentation.Scope{}, "name", "desc", unit.Dimensionless, testSumAggregator{})
			require.NoError(t, err)

			err = pipe.addAggregator(tt.secondInst.scope, tt.secondInst.name, tt.secondInst.description, tt.secondInst.unit, testSumAggregator{})
			assert.ErrorIs(t, err, tt.want)

			if tt.wantScopeLen > 0 {
				output, err := pipe.produce(context.Background())
				assert.NoError(t, err)
				require.Len(t, output.ScopeMetrics, tt.wantScopeLen)
				require.Len(t, output.ScopeMetrics[0].Metrics, tt.wantMetricsLen)
			}
		})
	}
}

func TestPipelineUsesResource(t *testing.T) {
	res := resource.NewWithAttributes("noSchema", attribute.String("test", "resource"))
	pipe := newPipeline(res, nil, nil)

	output, err := pipe.produce(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, res, output.Resource)
}

func TestPipelineConcurrency(t *testing.T) {
	pipe := newPipeline(nil, nil, nil)
	ctx := context.Background()

	var wg sync.WaitGroup
	const threads = 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = pipe.produce(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = pipe.addAggregator(instrumentation.Scope{}, "name", "desc", unit.Dimensionless, testSumAggregator{})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			pipe.addCallback(func(ctx context.Context) {})
		}()
	}
	wg.Wait()
}
