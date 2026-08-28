// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"slices"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/internal/attrdedup"
)

type int64CounterBindingWrapper func(
	api.Int64Counter,
	func(...attribute.KeyValue) api.Int64Counter,
) api.Int64Counter

type int64CounterFinishingWrapper func(
	api.Int64Counter,
	func(context.Context, ...attribute.KeyValue),
) api.Int64Counter

type experimentalInt64 struct {
	fallbackMeasures []aggregate.Measure[int64]
	binders          []func(attribute.Set) func(context.Context, int64)
	finishers        []func(attribute.Set)
}

func newExperimentalInt64Inst(aggs []aggregate.Aggregator[int64]) *int64Inst {
	state := &experimentalInt64{}
	i := &int64Inst{
		measures:     make([]aggregate.Measure[int64], 0, len(aggs)),
		experimental: state,
	}
	for _, agg := range aggs {
		measure := agg.Measure()
		i.measures = append(i.measures, measure)
		if binder, ok := agg.(interface {
			Bind(attribute.Set) func(context.Context, int64)
		}); ok {
			state.binders = append(state.binders, binder.Bind)
		} else {
			state.fallbackMeasures = append(state.fallbackMeasures, measure)
		}
		if finisher, ok := agg.(interface {
			Finish(attribute.Set)
		}); ok {
			state.finishers = append(state.finishers, finisher.Finish)
		}
	}
	return i
}

func lifecycleAttributes(attrs []attribute.KeyValue) attribute.Set {
	set, _ := attrdedup.Set(attribute.NewSet(slices.Clone(attrs)...))
	return set
}

func (i *int64Inst) bind(attrs ...attribute.KeyValue) api.Int64Counter {
	set := lifecycleAttributes(attrs)
	bound := &boundInt64Inst{instrument: i, attrs: set}
	for _, bind := range i.experimental.binders {
		bound.measures = append(bound.measures, bind(set))
	}
	return bound
}

func (i *int64Inst) finish(_ context.Context, attrs ...attribute.KeyValue) {
	set := lifecycleAttributes(attrs)
	for _, finisher := range i.experimental.finishers {
		finisher(set)
	}
}

type boundInt64Inst struct {
	embedded.Int64Counter
	instrument *int64Inst
	attrs      attribute.Set
	measures   []func(context.Context, int64)
}

var _ api.Int64Counter = (*boundInt64Inst)(nil)

func (i *boundInt64Inst) Enabled(ctx context.Context) bool {
	return i.instrument.Enabled(ctx)
}

func (i *boundInt64Inst) Add(ctx context.Context, value int64, opts ...api.AddOption) {
	if len(opts) == 0 {
		for _, measure := range i.measures {
			measure(ctx, value)
		}
		for _, measure := range i.instrument.experimental.fallbackMeasures {
			measure(ctx, value, i.attrs)
		}
		return
	}
	cfg := api.NewAddConfig(opts)
	extra := resolveAttributes(cfg.Attributes(), extractRawKVs(opts))
	if extra.Len() == 0 {
		i.add(ctx, value)
		return
	}
	merged := make([]attribute.KeyValue, 0, i.attrs.Len()+extra.Len())
	merged = append(merged, i.attrs.ToSlice()...)
	merged = append(merged, extra.ToSlice()...)
	set, _ := attrdedup.Set(attribute.NewSet(merged...))
	i.instrument.aggregate(ctx, value, set)
}

func (i *boundInt64Inst) add(ctx context.Context, value int64) {
	for _, measure := range i.measures {
		measure(ctx, value)
	}
	for _, measure := range i.instrument.experimental.fallbackMeasures {
		measure(ctx, value, i.attrs)
	}
}
