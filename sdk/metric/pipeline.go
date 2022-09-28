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
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/view"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	errCreatingAggregators     = errors.New("could not create all aggregators")
	errIncompatibleAggregation = errors.New("incompatible aggregation")
	errUnknownAggregation      = errors.New("unrecognized aggregation")
)

type aggregator interface {
	Aggregation() metricdata.Aggregation
}

// instrumentID is used to identify multiple instruments being mapped to the
// same aggregator. e.g. using an exact match view with a name=* view. You
// can't use a view.Instrument here because not all Aggregators are comparable.
type instrumentID struct {
	scope       instrumentation.Scope
	name        string
	description string
}

type instrumentKey struct {
	name string
	unit unit.Unit
}

type instrumentValue struct {
	description string
	aggregator  aggregator
}

func newPipeline(res *resource.Resource, reader Reader, views []view.View) *pipeline {
	if res == nil {
		res = resource.Empty()
	}
	return &pipeline{
		resource:     res,
		reader:       reader,
		views:        views,
		aggregations: make(map[instrumentation.Scope]map[instrumentKey]instrumentValue),
	}
}

// pipeline connects all of the instruments created by a meter provider to a Reader.
// This is the object that will be `Reader.register()` when a meter provider is created.
//
// As instruments are created the instrument should be checked if it exists in the
// views of a the Reader, and if so each aggregator should be added to the pipeline.
type pipeline struct {
	resource *resource.Resource

	reader Reader
	views  []view.View

	sync.Mutex
	aggregations map[instrumentation.Scope]map[instrumentKey]instrumentValue
	callbacks    []func(context.Context)
}

var errAlreadyRegistered = errors.New("instrument already registered")

// addAggregator will stores an aggregator with an instrument description.  The aggregator
// is used when `produce()` is called.
func (p *pipeline) addAggregator(scope instrumentation.Scope, name, description string, instUnit unit.Unit, agg aggregator) error {
	p.Lock()
	defer p.Unlock()
	if p.aggregations == nil {
		p.aggregations = map[instrumentation.Scope]map[instrumentKey]instrumentValue{}
	}
	if p.aggregations[scope] == nil {
		p.aggregations[scope] = map[instrumentKey]instrumentValue{}
	}
	inst := instrumentKey{
		name: name,
		unit: instUnit,
	}
	if _, ok := p.aggregations[scope][inst]; ok {
		return fmt.Errorf("%w: name %s, scope: %s", errAlreadyRegistered, name, scope)
	}

	p.aggregations[scope][inst] = instrumentValue{
		description: description,
		aggregator:  agg,
	}
	return nil
}

// addCallback registers a callback to be run when `produce()` is called.
func (p *pipeline) addCallback(callback func(context.Context)) {
	p.Lock()
	defer p.Unlock()
	p.callbacks = append(p.callbacks, callback)
}

// callbackKey is a context key type used to identify context that came from the SDK.
type callbackKey int

// produceKey is the context key to tell if a Observe is called within a callback.
// Its value of zero is arbitrary. If this package defined other context keys,
// they would have different integer values.
const produceKey callbackKey = 0

// produce returns aggregated metrics from a single collection.
//
// This method is safe to call concurrently.
func (p *pipeline) produce(ctx context.Context) (metricdata.ResourceMetrics, error) {
	p.Lock()
	defer p.Unlock()

	ctx = context.WithValue(ctx, produceKey, struct{}{})

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
		for inst, instValue := range instruments {
			data := instValue.aggregator.Aggregation()
			if data != nil {
				metrics = append(metrics, metricdata.Metrics{
					Name:        inst.name,
					Description: instValue.description,
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

// inserter facilitates inserting of new instruments into a pipeline.
type inserter[N int64 | float64] struct {
	pipeline *pipeline
}

func newInserter[N int64 | float64](p *pipeline) *inserter[N] {
	return &inserter[N]{p}
}

// Instrument inserts instrument inst with instUnit returning the Aggregators
// that need to be updated with measurments for that instrument.
func (i *inserter[N]) Instrument(inst view.Instrument, instUnit unit.Unit) ([]internal.Aggregator[N], error) {
	seen := map[instrumentID]struct{}{}
	var aggs []internal.Aggregator[N]
	errs := &multierror{wrapped: errCreatingAggregators}
	for _, v := range i.pipeline.views {
		inst, match := v.TransformInstrument(inst)

		id := instrumentID{
			scope:       inst.Scope,
			name:        inst.Name,
			description: inst.Description,
		}

		if _, ok := seen[id]; ok || !match {
			continue
		}

		if inst.Aggregation == nil {
			inst.Aggregation = i.pipeline.reader.aggregation(inst.Kind)
		} else if _, ok := inst.Aggregation.(aggregation.Default); ok {
			inst.Aggregation = i.pipeline.reader.aggregation(inst.Kind)
		}

		if err := isAggregatorCompatible(inst.Kind, inst.Aggregation); err != nil {
			err = fmt.Errorf("creating aggregator with instrumentKind: %d, aggregation %v: %w", inst.Kind, inst.Aggregation, err)
			errs.append(err)
			continue
		}

		agg, err := i.aggregator(inst)
		if err != nil {
			errs.append(err)
			continue
		}
		if agg == nil { // Drop aggregator.
			continue
		}
		// TODO (#3011): If filtering is done at the instrument level add here.
		// This is where the aggregator and the view are both in scope.
		aggs = append(aggs, agg)
		seen[id] = struct{}{}
		err = i.pipeline.addAggregator(inst.Scope, inst.Name, inst.Description, instUnit, agg)
		if err != nil {
			errs.append(err)
		}
	}
	// TODO(#3224): handle when no views match. Default should be reader
	// aggregation returned.
	return aggs, errs.errorOrNil()
}

// aggregator returns the Aggregator for an instrument configuration. If the
// instrument defines an unknown aggregation, an error is returned.
func (i *inserter[N]) aggregator(inst view.Instrument) (internal.Aggregator[N], error) {
	// TODO (#3011): If filtering is done by the Aggregator it should be passed
	// here.
	var (
		temporality = i.pipeline.reader.temporality(inst.Kind)
		monotonic   bool
	)

	switch inst.Kind {
	case view.AsyncCounter, view.SyncCounter, view.SyncHistogram:
		monotonic = true
	}

	switch agg := inst.Aggregation.(type) {
	case aggregation.Drop:
		return nil, nil
	case aggregation.LastValue:
		return internal.NewLastValue[N](), nil
	case aggregation.Sum:
		if temporality == metricdata.CumulativeTemporality {
			return internal.NewCumulativeSum[N](monotonic), nil
		}
		return internal.NewDeltaSum[N](monotonic), nil
	case aggregation.ExplicitBucketHistogram:
		if temporality == metricdata.CumulativeTemporality {
			return internal.NewCumulativeHistogram[N](agg), nil
		}
		return internal.NewDeltaHistogram[N](agg), nil
	}
	return nil, errUnknownAggregation
}

// isAggregatorCompatible checks if the aggregation can be used by the instrument.
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
		// TODO: review need for aggregation check after
		// https://github.com/open-telemetry/opentelemetry-specification/issues/2710
		return errIncompatibleAggregation
	case aggregation.Sum:
		switch kind {
		case view.AsyncCounter, view.AsyncUpDownCounter, view.SyncCounter, view.SyncHistogram, view.SyncUpDownCounter:
			return nil
		default:
			// TODO: review need for aggregation check after
			// https://github.com/open-telemetry/opentelemetry-specification/issues/2710
			return errIncompatibleAggregation
		}
	case aggregation.LastValue:
		if kind == view.AsyncGauge {
			return nil
		}
		// TODO: review need for aggregation check after
		// https://github.com/open-telemetry/opentelemetry-specification/issues/2710
		return errIncompatibleAggregation
	case aggregation.Drop:
		return nil
	default:
		// This is used passed checking for default, it should be an error at this point.
		return fmt.Errorf("%w: %v", errUnknownAggregation, agg)
	}
}

// pipelines is the group of pipelines connecting Readers with instrument
// measurement.
type pipelines []*pipeline

func newPipelines(res *resource.Resource, readers map[Reader][]view.View) pipelines {
	pipes := make([]*pipeline, 0, len(readers))
	for r, v := range readers {
		p := &pipeline{
			resource: res,
			reader:   r,
			views:    v,
		}
		r.register(p)
		pipes = append(pipes, p)
	}
	return pipes
}

// TODO (#3053) Only register callbacks if any instrument matches in a view.
func (p pipelines) registerCallback(fn func(context.Context)) {
	for _, pipe := range p {
		pipe.addCallback(fn)
	}
}

// resolver facilitates resolving Aggregators an instrument needs to aggregate
// measurements with while updating all pipelines that need to pull from those
// aggregations.
type resolver[N int64 | float64] struct {
	cache     instrumentCache[N]
	inserters []*inserter[N]
}

func newResolver[N int64 | float64](p pipelines, c instrumentCache[N]) *resolver[N] {
	in := make([]*inserter[N], len(p))
	for i := range in {
		in[i] = newInserter[N](p[i])
	}
	return &resolver[N]{cache: c, inserters: in}
}

// Aggregators returns the Aggregators instrument inst needs to update when it
// makes a measurement.
func (r *resolver[N]) Aggregators(inst view.Instrument, instUnit unit.Unit) ([]internal.Aggregator[N], error) {
	id := instrumentID{
		scope:       inst.Scope,
		name:        inst.Name,
		description: inst.Description,
	}

	return r.cache.Lookup(id, func() ([]internal.Aggregator[N], error) {
		var aggs []internal.Aggregator[N]
		errs := &multierror{}
		for _, i := range r.inserters {
			a, err := i.Instrument(inst, instUnit)
			if err != nil {
				errs.append(err)
			}
			aggs = append(aggs, a...)
		}
		return aggs, errs.errorOrNil()
	})
}

// resolvedAggregators is the result of resolving aggregators for an instrument.
type resolvedAggregators[N int64 | float64] struct {
	aggregators []internal.Aggregator[N]
	err         error
}

type instrumentCache[N int64 | float64] struct {
	cache *cache[instrumentID, any]
}

func newInstrumentCache[N int64 | float64](c *cache[instrumentID, any]) instrumentCache[N] {
	if c == nil {
		c = &cache[instrumentID, any]{}
	}
	return instrumentCache[N]{cache: c}
}

var errExists = errors.New("instrument already exists for different number type")

// Lookup returns the Aggregators and error for a cached instrumentID if they
// exist in the cache. Otherwise, f is called and its returned values are set
// in the cache and returned.
//
// If an instrumentID has been stored in the cache for a different N, an error
// is returned describing the conflict.
//
// Lookup is safe to call concurrently.
func (c instrumentCache[N]) Lookup(key instrumentID, f func() ([]internal.Aggregator[N], error)) (aggs []internal.Aggregator[N], err error) {
	vAny := c.cache.Lookup(key, func() any {
		a, err := f()
		return &resolvedAggregators[N]{
			aggregators: a,
			err:         err,
		}
	})

	switch v := vAny.(type) {
	case *resolvedAggregators[N]:
		aggs = v.aggregators
	default:
		err = errExists
	}
	return aggs, err
}

type multierror struct {
	wrapped error
	errors  []string
}

func (m *multierror) errorOrNil() error {
	if len(m.errors) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %s", m.wrapped, strings.Join(m.errors, "; "))
}

func (m *multierror) append(err error) {
	m.errors = append(m.errors, err.Error())
}
