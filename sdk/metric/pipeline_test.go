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
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testSumAggregator struct{}

func (testSumAggregator) Aggregation() metricdata.Aggregation {
	return metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint[int64]{}}
}

func TestNewPipeline(t *testing.T) {
	pipe := newPipeline(nil, nil, nil)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(context.Background(), &output)
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	assert.Len(t, output.ScopeMetrics, 0)

	iSync := instrumentSync{"name", "desc", "1", testSumAggregator{}}
	assert.NotPanics(t, func() {
		pipe.addSync(instrumentation.Scope{}, iSync)
	})

	require.NotPanics(t, func() {
		pipe.addMultiCallback(func(context.Context) error { return nil })
	})

	err = pipe.produce(context.Background(), &output)
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	require.Len(t, output.ScopeMetrics, 1)
	require.Len(t, output.ScopeMetrics[0].Metrics, 1)
}

func TestPipelineUsesResource(t *testing.T) {
	res := resource.NewWithAttributes("noSchema", attribute.String("test", "resource"))
	pipe := newPipeline(res, nil, nil)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(context.Background(), &output)
	assert.NoError(t, err)
	assert.Equal(t, res, output.Resource)
}

func TestPipelineConcurrency(t *testing.T) {
	pipe := newPipeline(nil, nil, nil)
	ctx := context.Background()
	var output metricdata.ResourceMetrics

	var wg sync.WaitGroup
	const threads = 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = pipe.produce(ctx, &output)
		}()

		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			name := fmt.Sprintf("name %d", n)
			sync := instrumentSync{name, "desc", "1", testSumAggregator{}}
			pipe.addSync(instrumentation.Scope{}, sync)
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			pipe.addMultiCallback(func(context.Context) error { return nil })
		}()
	}
	wg.Wait()
}

func TestDefaultViewImplicit(t *testing.T) {
	t.Run("Int64", testDefaultViewImplicit[int64]())
	t.Run("Float64", testDefaultViewImplicit[float64]())
}

func testDefaultViewImplicit[N int64 | float64]() func(t *testing.T) {
	inst := Instrument{
		Name:        "requests",
		Description: "count of requests received",
		Kind:        InstrumentKindCounter,
		Unit:        "1",
	}
	return func(t *testing.T) {
		reader := NewManualReader()
		tests := []struct {
			name string
			pipe *pipeline
		}{
			{
				name: "NoView",
				pipe: newPipeline(nil, reader, nil),
			},
			{
				name: "NoMatchingView",
				pipe: newPipeline(nil, reader, []View{
					NewView(Instrument{Name: "foo"}, Stream{Name: "bar"}),
				}),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var c cache[string, streamID]
				i := newInserter[N](test.pipe, &c)
				got, err := i.Instrument(inst)
				require.NoError(t, err)
				assert.Len(t, got, 1, "default view not applied")
				for _, a := range got {
					a.Aggregate(1, *attribute.EmptySet())
				}

				out := metricdata.ResourceMetrics{}
				err = test.pipe.produce(context.Background(), &out)
				require.NoError(t, err)
				require.Len(t, out.ScopeMetrics, 1, "Aggregator not registered with pipeline")
				sm := out.ScopeMetrics[0]
				require.Len(t, sm.Metrics, 1, "metrics not produced from default view")
				metricdatatest.AssertEqual(t, metricdata.Metrics{
					Name:        inst.Name,
					Description: inst.Description,
					Unit:        "1",
					Data: metricdata.Sum[N]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints:  []metricdata.DataPoint[N]{{Value: N(1)}},
					},
				}, sm.Metrics[0], metricdatatest.IgnoreTimestamp())
			})
		}
	}
}
