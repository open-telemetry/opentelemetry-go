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
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
)

func testSumAggregateOutput(dest *metricdata.Aggregation) int {
	*dest = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint[int64]{{Value: 1}},
	}
	return 1
}

func TestNewPipeline(t *testing.T) {
	pipe := newPipeline(nil, nil, nil)

	output := metricdata.ResourceMetrics{}
	err := pipe.produce(context.Background(), &output)
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), output.Resource)
	assert.Len(t, output.ScopeMetrics, 0)

	iSync := instrumentSync{"name", "desc", "1", testSumAggregateOutput}
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

func TestPipelineConcurrentSafe(t *testing.T) {
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
			sync := instrumentSync{name, "desc", "1", testSumAggregateOutput}
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
				var c cache[string, instID]
				i := newInserter[N](test.pipe, &c)
				readerAggregation := i.readerDefaultAggregation(inst.Kind)
				got, err := i.Instrument(inst, readerAggregation)
				require.NoError(t, err)
				assert.Len(t, got, 1, "default view not applied")
				for _, in := range got {
					in(context.Background(), 1, *attribute.EmptySet())
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

func TestLogConflictName(t *testing.T) {
	testcases := []struct {
		existing, name string
		conflict       bool
	}{
		{
			existing: "requestCount",
			name:     "requestCount",
			conflict: false,
		},
		{
			existing: "requestCount",
			name:     "requestDuration",
			conflict: false,
		},
		{
			existing: "requestCount",
			name:     "requestcount",
			conflict: true,
		},
		{
			existing: "requestCount",
			name:     "REQUESTCOUNT",
			conflict: true,
		},
		{
			existing: "requestCount",
			name:     "rEqUeStCoUnT",
			conflict: true,
		},
	}

	var msg string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) {
			msg = args
		}, funcr.Options{Verbosity: 20}))
		return func() { otel.SetLogger(orig) }
	}(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))))

	for _, tc := range testcases {
		var vc cache[string, instID]

		name := strings.ToLower(tc.existing)
		_ = vc.Lookup(name, func() instID {
			return instID{Name: tc.existing}
		})

		i := newInserter[int64](newPipeline(nil, nil, nil), &vc)
		i.logConflict(instID{Name: tc.name})

		if tc.conflict {
			assert.Containsf(
				t, msg, "duplicate metric stream definitions",
				"warning not logged for conflicting names: %s, %s",
				tc.existing, tc.name,
			)
		} else {
			assert.Equalf(
				t, msg, "",
				"warning logged for non-conflicting names: %s, %s",
				tc.existing, tc.name,
			)
		}

		// Reset.
		msg = ""
	}
}

func TestLogConflictSuggestView(t *testing.T) {
	var msg string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) {
			msg = args
		}, funcr.Options{Verbosity: 20}))
		return func() { otel.SetLogger(orig) }
	}(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))))

	orig := instID{
		Name:        "requestCount",
		Description: "number of requests",
		Kind:        InstrumentKindCounter,
		Unit:        "1",
		Number:      "int64",
	}

	var vc cache[string, instID]
	name := strings.ToLower(orig.Name)
	_ = vc.Lookup(name, func() instID { return orig })
	i := newInserter[int64](newPipeline(nil, nil, nil), &vc)

	viewSuggestion := func(inst instID, stream string) string {
		return `"NewView(Instrument{` +
			`Name: \"` + inst.Name +
			`\", Description: \"` + inst.Description +
			`\", Kind: \"InstrumentKind` + inst.Kind.String() +
			`\", Unit: \"` + inst.Unit +
			`\"}, ` +
			stream +
			`)"`
	}

	t.Run("Name", func(t *testing.T) {
		inst := instID{
			Name:        "requestcount",
			Description: orig.Description,
			Kind:        orig.Kind,
			Unit:        orig.Unit,
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.Containsf(t, msg, viewSuggestion(
			inst, `Stream{Name: \"{{NEW_NAME}}\"}`,
		), "no suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Description", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: "alt",
			Kind:        orig.Kind,
			Unit:        orig.Unit,
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.Containsf(t, msg, viewSuggestion(
			inst, `Stream{Description: \"`+orig.Description+`\"}`,
		), "no suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Kind", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: orig.Description,
			Kind:        InstrumentKindHistogram,
			Unit:        orig.Unit,
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.Containsf(t, msg, viewSuggestion(
			inst, `Stream{Name: \"{{NEW_NAME}}\"}`,
		), "no suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Unit", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: orig.Description,
			Kind:        orig.Kind,
			Unit:        "ms",
			Number:      orig.Number,
		}
		i.logConflict(inst)
		assert.NotContains(t, msg, "NewView", "suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})

	t.Run("Number", func(t *testing.T) {
		inst := instID{
			Name:        orig.Name,
			Description: orig.Description,
			Kind:        orig.Kind,
			Unit:        orig.Unit,
			Number:      "float64",
		}
		i.logConflict(inst)
		assert.NotContains(t, msg, "NewView", "suggestion logged: %v", inst)

		// Reset.
		msg = ""
	})
}

func TestInserterCachedAggregatorNameConflict(t *testing.T) {
	const name = "requestCount"
	scope := instrumentation.Scope{Name: "pipeline_test"}
	kind := InstrumentKindCounter
	stream := Stream{
		Name:        name,
		Aggregation: AggregationSum{},
	}

	var vc cache[string, instID]
	pipe := newPipeline(nil, NewManualReader(), nil)
	i := newInserter[int64](pipe, &vc)

	readerAggregation := i.readerDefaultAggregation(kind)
	_, origID, err := i.cachedAggregator(scope, kind, stream, readerAggregation)
	require.NoError(t, err)

	require.Len(t, pipe.aggregations, 1)
	require.Contains(t, pipe.aggregations, scope)
	iSync := pipe.aggregations[scope]
	require.Len(t, iSync, 1)
	require.Equal(t, name, iSync[0].name)

	stream.Name = "RequestCount"
	_, id, err := i.cachedAggregator(scope, kind, stream, readerAggregation)
	require.NoError(t, err)
	assert.Equal(t, origID, id, "multiple aggregators for equivalent name")

	assert.Len(t, pipe.aggregations, 1, "additional scope added")
	require.Contains(t, pipe.aggregations, scope, "original scope removed")
	iSync = pipe.aggregations[scope]
	require.Len(t, iSync, 1, "registered instrumentSync changed")
	assert.Equal(t, name, iSync[0].name, "stream name changed")
}
