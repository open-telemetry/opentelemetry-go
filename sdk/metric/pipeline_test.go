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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testSumAggregator struct{}

func (testSumAggregator) Aggregation() metricdata.Aggregation {
	return metricdata.Sum{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint{}}
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
	pipe := newPipeline(nil)

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
	pipe := newPipeline(nil)

	err := pipe.addAggregator(instrumentation.Scope{}, "name", "desc", unit.Dimensionless, testSumAggregator{})
	require.NoError(t, err)

	err = pipe.addAggregator(instrumentation.Scope{}, "name", "desc", unit.Dimensionless, testSumAggregator{})
	assert.Error(t, err)

	output, err := pipe.produce(context.Background())
	assert.NoError(t, err)
	require.Len(t, output.ScopeMetrics, 1)
	require.Len(t, output.ScopeMetrics[0].Metrics, 1)
}

func TestPipelineConcurrency(t *testing.T) {
	pipe := newPipeline(nil)
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
