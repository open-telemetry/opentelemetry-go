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
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type aggregator interface {
	Aggregation() metricdata.Aggregation
}
type instrumentKey struct {
	name        string
	description string
	unit        unit.Unit
}

func newPipeline(res *resource.Resource) *pipeline {
	if res == nil {
		res = resource.Empty()
	}
	return &pipeline{
		resource:     res,
		aggregations: make(map[instrumentation.Scope]map[instrumentKey]aggregator),
	}
}

// pipeline connects all of the instruments created by a meter provider to a Reader.
// This is the object that will be `Reader.register()` when a meter provider is created.
//
// As instruments are created the instrument should be checked if it exists in the
// views of a the Reader, and if so each aggregator should be added to the pipeline.
type pipeline struct {
	resource *resource.Resource

	sync.Mutex
	aggregations map[instrumentation.Scope]map[instrumentKey]aggregator
	callbacks    []func(context.Context)
}

// addAggregator will stores an aggregator with an instrument description.  The aggregator
// is used when `produce()` is called.
func (p *pipeline) addAggregator(scope instrumentation.Scope, name, description string, instUnit unit.Unit, agg aggregator) error {
	p.Lock()
	defer p.Unlock()
	if p.aggregations == nil {
		p.aggregations = map[instrumentation.Scope]map[instrumentKey]aggregator{}
	}
	if p.aggregations[scope] == nil {
		p.aggregations[scope] = map[instrumentKey]aggregator{}
	}
	inst := instrumentKey{
		name:        name,
		description: description,
		unit:        instUnit,
	}
	if _, ok := p.aggregations[scope][inst]; ok {
		return fmt.Errorf("instrument already registered: name %s, scope: %s", name, scope)
	}

	p.aggregations[scope][inst] = agg
	return nil
}

// addCallback registers a callback to be run when `produce()` is called.
func (p *pipeline) addCallback(callback func(context.Context)) {
	p.Lock()
	defer p.Unlock()
	p.callbacks = append(p.callbacks, callback)
}

// produce returns aggregated metrics from a single collection.
//
// This method is safe to call concurrently.
func (p *pipeline) produce(ctx context.Context) (metricdata.ResourceMetrics, error) {
	p.Lock()
	defer p.Unlock()

	for _, callback := range p.callbacks {
		// TODO make the callbacks parallel. ( #3034 )
		callback(ctx)
		if err := ctx.Err(); err != nil {
			// This means the context expired before we finished running callbacks.
			return metricdata.ResourceMetrics{}, err
		}
	}

	sm := make([]metricdata.ScopeMetrics, 0, len(p.aggregations))
	for scope, instruments := range p.aggregations {
		metrics := make([]metricdata.Metrics, 0, len(instruments))
		for inst, aggregation := range instruments {
			data := aggregation.Aggregation()
			if data != nil {
				metrics = append(metrics, metricdata.Metrics{
					Name:        inst.name,
					Description: inst.description,
					Unit:        inst.unit,
					Data:        data,
				})
			}
		}
		if len(metrics) > 0 {
			sm = append(sm, metricdata.ScopeMetrics{
				Scope:   scope,
				Metrics: metrics,
			})
		}
	}

	return metricdata.ResourceMetrics{
		Resource:     p.resource,
		ScopeMetrics: sm,
	}, nil
}
