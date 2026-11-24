// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"errors"
	"slices"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
)

// bindMeter decorates a meter to construct bound-capable synchronous
// instruments when experimental binding is enabled.
type bindMeter struct {
	metric.Meter
	meter                 *meter
	int64Counters         cacheWithErr[instID, metric.Int64Counter]
	float64Counters       cacheWithErr[instID, metric.Float64Counter]
	int64UpDownCounters   cacheWithErr[instID, metric.Int64UpDownCounter]
	float64UpDownCounters cacheWithErr[instID, metric.Float64UpDownCounter]
}

var _ metric.Meter = (*bindMeter)(nil)

func newBindMeter(base metric.Meter, m *meter) metric.Meter {
	return &bindMeter{Meter: base, meter: m}
}

func (m *bindMeter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	cfg := metric.NewInt64CounterConfig(options...)
	counter, err := m.int64Counters.Lookup(instID{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindCounter,
	}, func() (metric.Int64Counter, error) {
		inst := Instrument{
			Name:        name,
			Description: cfg.Description(),
			Unit:        cfg.Unit(),
			Kind:        InstrumentKindCounter,
			Scope:       m.meter.scope,
		}
		measures, binders, fallbackMeasures, err := bindAggregators(
			m.meter.int64Resolver,
			inst,
			defaultAttributes(options),
			isOptIn(options),
		)
		return &bindInt64Counter{
			bindableAdd: bindableAdd[int64]{
				measures:         measures,
				binders:          binders,
				fallbackMeasures: fallbackMeasures,
			},
		}, err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

func (m *bindMeter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	cfg := metric.NewFloat64CounterConfig(options...)
	counter, err := m.float64Counters.Lookup(instID{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindCounter,
	}, func() (metric.Float64Counter, error) {
		inst := Instrument{
			Name:        name,
			Description: cfg.Description(),
			Unit:        cfg.Unit(),
			Kind:        InstrumentKindCounter,
			Scope:       m.meter.scope,
		}
		measures, binders, fallbackMeasures, err := bindAggregators(
			m.meter.float64Resolver,
			inst,
			defaultAttributes(options),
			isOptIn(options),
		)
		return &bindFloat64Counter{
			bindableAdd: bindableAdd[float64]{
				measures:         measures,
				binders:          binders,
				fallbackMeasures: fallbackMeasures,
			},
		}, err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

func (m *bindMeter) Int64UpDownCounter(
	name string,
	options ...metric.Int64UpDownCounterOption,
) (metric.Int64UpDownCounter, error) {
	cfg := metric.NewInt64UpDownCounterConfig(options...)
	counter, err := m.int64UpDownCounters.Lookup(instID{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindUpDownCounter,
	}, func() (metric.Int64UpDownCounter, error) {
		inst := Instrument{
			Name:        name,
			Description: cfg.Description(),
			Unit:        cfg.Unit(),
			Kind:        InstrumentKindUpDownCounter,
			Scope:       m.meter.scope,
		}
		measures, binders, fallbackMeasures, err := bindAggregators(
			m.meter.int64Resolver,
			inst,
			defaultAttributes(options),
			isOptIn(options),
		)
		return &bindInt64UpDownCounter{
			bindableAdd: bindableAdd[int64]{
				measures:         measures,
				binders:          binders,
				fallbackMeasures: fallbackMeasures,
			},
		}, err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

func (m *bindMeter) Float64UpDownCounter(
	name string,
	options ...metric.Float64UpDownCounterOption,
) (metric.Float64UpDownCounter, error) {
	cfg := metric.NewFloat64UpDownCounterConfig(options...)
	counter, err := m.float64UpDownCounters.Lookup(instID{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindUpDownCounter,
	}, func() (metric.Float64UpDownCounter, error) {
		inst := Instrument{
			Name:        name,
			Description: cfg.Description(),
			Unit:        cfg.Unit(),
			Kind:        InstrumentKindUpDownCounter,
			Scope:       m.meter.scope,
		}
		measures, binders, fallbackMeasures, err := bindAggregators(
			m.meter.float64Resolver,
			inst,
			defaultAttributes(options),
			isOptIn(options),
		)
		return &bindFloat64UpDownCounter{
			bindableAdd: bindableAdd[float64]{
				measures:         measures,
				binders:          binders,
				fallbackMeasures: fallbackMeasures,
			},
		}, err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

// bindAggregators adapts the SDK pipeline resolver to bound-capable
// aggregations. Streams supporting binding populate binders; other compatible
// streams populate fallbackMeasures so bound instruments continue to record
// their pre-bound attributes into every resolved stream.
func bindAggregators[N int64 | float64](
	r resolver[N],
	id Instrument,
	allowedKeys []attribute.Key,
	optIn bool,
) ([]aggregate.Measure[N], []aggregate.Binder[N], []aggregate.Measure[N], error) {
	var (
		measures         []aggregate.Measure[N]
		binders          []aggregate.Binder[N]
		fallbackMeasures []aggregate.Measure[N]
		err              error
	)
	for _, i := range r.inserters {
		newAggregation := func(
			b aggregate.Builder[N],
			agg Aggregation,
			kind InstrumentKind,
		) (streamAggregation[N], error) {
			if kind == InstrumentKindCounter || kind == InstrumentKindUpDownCounter {
				if _, ok := agg.(AggregationSum); ok {
					sum := b.BoundSum(kind == InstrumentKindCounter)
					return streamAggregation[N]{
						measure: sum.Measure,
						compute: sum.ComputeAggregation,
						bind:    sum.Bind,
					}, nil
				}
			}
			return i.newStreamAggregation(b, agg, kind)
		}
		e := i.forEachStream(streamConfig[N]{
			inst:              id,
			allowedKeys:       allowedKeys,
			optIn:             optIn,
			readerAggregation: i.readerDefaultAggregation(id.Kind),
			newAggregation:    newAggregation,
			yield: func(agg aggVal[N]) {
				measures = append(measures, agg.Measure)
				if agg.Bind != nil {
					binders = append(binders, agg.Bind)
				} else {
					fallbackMeasures = append(fallbackMeasures, agg.Measure)
				}
			},
		})
		if e != nil {
			err = errors.Join(err, e)
		}
	}
	return measures, binders, fallbackMeasures, err
}

type bindableAdd[N int64 | float64] struct {
	measures         []aggregate.Measure[N]
	binders          []aggregate.Binder[N]
	fallbackMeasures []aggregate.Measure[N]
}

func (i *bindableAdd[N]) Add(ctx context.Context, val N, opts ...metric.AddOption) {
	c := metric.NewAddConfig(opts)
	s := resolveAttributes(c.Attributes())
	for _, in := range i.measures {
		in(ctx, val, s)
	}
}

func (i *bindableAdd[N]) Enabled(context.Context) bool {
	return len(i.measures) != 0
}

func (i *bindableAdd[N]) bind(attrs []attribute.KeyValue) boundAdd[N] {
	if len(i.measures) == 0 {
		return boundAdd[N]{}
	}
	set := attribute.NewSet()
	if len(attrs) > 0 {
		// Clone before NewSet because attribute.NewSet sorts its input in place.
		set = resolveAttributes(attribute.NewSet(slices.Clone(attrs)...))
	}
	var boundMeasures []aggregate.BoundMeasure[N]
	if len(i.binders) > 0 {
		boundMeasures = make([]aggregate.BoundMeasure[N], len(i.binders))
		for idx, binder := range i.binders {
			boundMeasures[idx] = binder(set)
		}
	}
	return boundAdd[N]{
		boundMeasures:    boundMeasures,
		fallbackMeasures: i.fallbackMeasures,
		measures:         i.measures,
		attrs:            set,
	}
}

type boundAdd[N int64 | float64] struct {
	boundMeasures    []aggregate.BoundMeasure[N]
	fallbackMeasures []aggregate.Measure[N]
	measures         []aggregate.Measure[N]
	attrs            attribute.Set
}

func (b *boundAdd[N]) Add(ctx context.Context, val N, opts ...metric.AddOption) {
	if len(opts) == 0 {
		b.addBound(ctx, val)
		return
	}
	if len(b.measures) == 0 {
		return
	}
	c := metric.NewAddConfig(opts)
	extra := c.Attributes()
	if extra.Len() == 0 {
		b.addBound(ctx, val)
		return
	}
	merged := mergeAttributeSets(b.attrs, resolveAttributes(extra))
	for _, in := range b.measures {
		in(ctx, val, merged)
	}
}

func (b *boundAdd[N]) addBound(ctx context.Context, val N) {
	for _, in := range b.boundMeasures {
		in(ctx, val)
	}
	for _, in := range b.fallbackMeasures {
		in(ctx, val, b.attrs)
	}
}

func (b *boundAdd[N]) Enabled(context.Context) bool {
	return len(b.measures) != 0
}

// mergeAttributeSets returns the union of base and extra.
// Duplicate keys use the value from extra.
func mergeAttributeSets(base, extra attribute.Set) attribute.Set {
	if base.Len() == 0 {
		return extra
	}
	iter := attribute.NewMergeIterator(&extra, &base)
	merged := make([]attribute.KeyValue, 0, base.Len()+extra.Len())
	for iter.Next() {
		merged = append(merged, iter.Attribute())
	}
	return attribute.NewSet(merged...)
}

type bindInt64Counter struct {
	bindableAdd[int64]
	embedded.Int64Counter
}

func (c *bindInt64Counter) Bind(attrs ...attribute.KeyValue) metric.Int64Counter {
	return &boundInt64Counter{boundAdd: c.bind(attrs)}
}

type boundInt64Counter struct {
	boundAdd[int64]
	embedded.Int64Counter
}

type bindFloat64Counter struct {
	bindableAdd[float64]
	embedded.Float64Counter
}

func (c *bindFloat64Counter) Bind(attrs ...attribute.KeyValue) metric.Float64Counter {
	return &boundFloat64Counter{boundAdd: c.bind(attrs)}
}

type boundFloat64Counter struct {
	boundAdd[float64]
	embedded.Float64Counter
}

type bindInt64UpDownCounter struct {
	bindableAdd[int64]
	embedded.Int64UpDownCounter
}

func (c *bindInt64UpDownCounter) Bind(attrs ...attribute.KeyValue) metric.Int64UpDownCounter {
	return &boundInt64UpDownCounter{boundAdd: c.bind(attrs)}
}

type boundInt64UpDownCounter struct {
	boundAdd[int64]
	embedded.Int64UpDownCounter
}

type bindFloat64UpDownCounter struct {
	bindableAdd[float64]
	embedded.Float64UpDownCounter
}

func (c *bindFloat64UpDownCounter) Bind(attrs ...attribute.KeyValue) metric.Float64UpDownCounter {
	return &boundFloat64UpDownCounter{boundAdd: c.bind(attrs)}
}

type boundFloat64UpDownCounter struct {
	boundAdd[float64]
	embedded.Float64UpDownCounter
}

var (
	_ metric.Int64Counter         = (*bindInt64Counter)(nil)
	_ metric.Int64Counter         = (*boundInt64Counter)(nil)
	_ metric.Float64Counter       = (*bindFloat64Counter)(nil)
	_ metric.Float64Counter       = (*boundFloat64Counter)(nil)
	_ metric.Int64UpDownCounter   = (*bindInt64UpDownCounter)(nil)
	_ metric.Int64UpDownCounter   = (*boundInt64UpDownCounter)(nil)
	_ metric.Float64UpDownCounter = (*bindFloat64UpDownCounter)(nil)
	_ metric.Float64UpDownCounter = (*boundFloat64UpDownCounter)(nil)
)
