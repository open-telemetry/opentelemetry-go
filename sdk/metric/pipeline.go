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
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/view"
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
		for inst, agg := range instruments {
			data := agg.Aggregation()
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

// pipelineRegistry manages creating pipelines, and aggregators.  Meters retrieve
// new Aggregators from a pipelineRegistry.
type pipelineRegistry[N int64 | float64] struct {
	views     map[Reader][]view.View
	pipelines map[Reader]*pipeline
}

func newPipelineRegistries(views map[Reader][]view.View) (*pipelineRegistry[int64], *pipelineRegistry[float64]) {
	pipelines := map[Reader]*pipeline{}
	for rdr := range views {
		pipe := &pipeline{}
		rdr.register(pipe)
		pipelines[rdr] = pipe
	}
	return &pipelineRegistry[int64]{
			views:     views,
			pipelines: pipelines,
		}, &pipelineRegistry[float64]{
			views:     views,
			pipelines: pipelines,
		}
}

// createAggregators will create all backing aggregators for an instrument.
// It will return an error if an instrument is registered more than once.
// Note: There may be returned aggregators with an error.
func (reg *pipelineRegistry[N]) createAggregators(inst view.Instrument, instUnit unit.Unit) ([]internal.Aggregator[N], error) {
	var aggs []internal.Aggregator[N]

	errs := &multierror{}
	for rdr, views := range reg.views {
		pipe := reg.pipelines[rdr]
		rdrAggs, err := createAggregators[N](rdr, views, inst)
		if err != nil {
			errs.append(err)
		}
		for inst, agg := range rdrAggs {
			err := pipe.addAggregator(inst.scope, inst.name, inst.description, instUnit, agg)
			if err != nil {
				errs.append(err)
			}
			aggs = append(aggs, agg)
		}
	}
	return aggs, errs.errorOrNil()
}

// TODO (#3053) Only register callbacks if any instrument matches in a view.
func (reg *pipelineRegistry[N]) registerCallback(fn func(context.Context)) {
	for _, pipe := range reg.pipelines {
		pipe.addCallback(fn)
	}
}

type multierror struct {
	errors []string
}

func (m *multierror) errorOrNil() error {
	if len(m.errors) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(m.errors, "; "))
}

func (m *multierror) append(err error) {
	m.errors = append(m.errors, err.Error())
}

// instrumentIdentifier is used to identify multiple instruments being mapped to the same aggregator.
// e.g. using an exact match view with a name=* view.
// You can't use a view.Instrument here because not all Aggregators are hashable.
type instrumentIdentifier struct {
	scope       instrumentation.Scope
	name        string
	description string
}

func createAggregators[N int64 | float64](rdr Reader, views []view.View, inst view.Instrument) (map[instrumentIdentifier]internal.Aggregator[N], error) {
	aggs := map[instrumentIdentifier]internal.Aggregator[N]{}
	errs := &multierror{}
	for _, v := range views {
		inst, match := v.TransformInstrument(inst)

		ident := instrumentIdentifier{
			scope:       inst.Scope,
			name:        inst.Name,
			description: inst.Description,
		}

		if _, ok := aggs[ident]; ok || !match {
			continue
		}

		if inst.Aggregation == nil {
			inst.Aggregation = rdr.aggregation(inst.Kind)
		} else if _, ok := inst.Aggregation.(aggregation.Default); ok {
			inst.Aggregation = rdr.aggregation(inst.Kind)
		}

		if err := isAggregatorCompatible(inst.Kind, inst.Aggregation); err != nil {
			global.Error(err, "creating aggregator", "instrumentKind", inst.Kind, "aggregation", inst.Aggregation)
			errs.append(err)
			continue
		}

		agg := createAggregator[N](inst.Aggregation, rdr.temporality(inst.Kind), isMonotonic(inst.Kind))
		if agg != nil {
			// TODO (#3011): If filtering is done at the instrument level add here.
			// This is where the aggregator and the view are both in scope.
			aggs[ident] = agg
		}
	}
	return aggs, errs.errorOrNil()
}

func isMonotonic(kind view.InstrumentKind) bool {
	switch kind {
	case view.AsyncCounter, view.SyncCounter, view.SyncHistogram:
		return true
	}
	return false
}

// createAggregator takes the config (Aggregation and Temporality) and produces a memory backed Aggregator.
// TODO (#3011): If filterting is done by the Aggregator it should be passed here.
func createAggregator[N int64 | float64](agg aggregation.Aggregation, temporality metricdata.Temporality, monotonic bool) internal.Aggregator[N] {
	switch agg := agg.(type) {
	case aggregation.Drop:
		return nil
	case aggregation.LastValue:
		return internal.NewLastValue[N]()
	case aggregation.Sum:
		if temporality == metricdata.CumulativeTemporality {
			return internal.NewCumulativeSum[N](monotonic)
		}
		return internal.NewDeltaSum[N](monotonic)
	case aggregation.ExplicitBucketHistogram:
		if temporality == metricdata.CumulativeTemporality {
			return internal.NewCumulativeHistogram[N](agg)
		}
		return internal.NewDeltaHistogram[N](agg)
	}
	return nil
}

var errIncompatibleAggregation = errors.New("incompatible aggregation")

// is aggregatorCompatible checks if the aggregation can be used by the instrument.
// Current compatibility:
//
// | Instrument Kind      | Drop | LastValue | Sum | Histogram | Exponential Histogram |
// |----------------------|------|-----------|-----|-----------|-----------------------|
// | Sync Counter         | X    |           | X   | X         | X                     |
// | Sync UpDown Counter  | X    |           | X   |           |                       |
// | Sync Histogram       | X    |           | X   | X         | X                     |
// | Async Counter        | X    |           | X   |           |                       |
// | Async UpDown Counter | X    |           | X   |           |                       |
// | Async Gauge          | X    | X         |     |           |                       |.
func isAggregatorCompatible(kind view.InstrumentKind, agg aggregation.Aggregation) error {
	switch agg.(type) {
	case aggregation.ExplicitBucketHistogram:
		if kind == view.SyncCounter || kind == view.SyncHistogram {
			return nil
		}
		return errIncompatibleAggregation
	case aggregation.Sum:
		switch kind {
		case view.AsyncCounter, view.AsyncUpDownCounter, view.SyncCounter, view.SyncHistogram, view.SyncUpDownCounter:
			return nil
		default:
			return errIncompatibleAggregation
		}
	case aggregation.LastValue:
		if kind == view.AsyncGauge {
			return nil
		}
		return errIncompatibleAggregation
	case aggregation.Drop:
		return nil
	default:
		// This is used passed checking for default, it should be an error at this point.
		return errIncompatibleAggregation
	}
}
