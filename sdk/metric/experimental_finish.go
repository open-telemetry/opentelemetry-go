// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
)

// factory adapts the SDK pipeline resolver to Finish instrument construction.
type factory struct {
	meter                 *meter
	registry              *finish.Registry
	int64Counters         cacheWithErr[instID, metric.Int64Counter]
	float64Counters       cacheWithErr[instID, metric.Float64Counter]
	int64UpDownCounters   cacheWithErr[instID, metric.Int64UpDownCounter]
	float64UpDownCounters cacheWithErr[instID, metric.Float64UpDownCounter]
}

var _ finish.Factory = (*factory)(nil)

func (f *factory) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	cfg := metric.NewInt64CounterConfig(options...)
	counter, err := f.int64Counters.Lookup(instID{
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
			Scope:       f.meter.scope,
		}
		measures, finishers, err := finishCounterAggregators(
			f.meter.int64Resolver,
			inst,
			defaultAttributes(options),
			f.registry.Register,
		)
		return finish.NewInt64Counter(
			&int64Inst{measures: measures},
			finishers,
		), err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

func (f *factory) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	cfg := metric.NewFloat64CounterConfig(options...)
	counter, err := f.float64Counters.Lookup(instID{
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
			Scope:       f.meter.scope,
		}
		measures, finishers, err := finishCounterAggregators(
			f.meter.float64Resolver,
			inst,
			defaultAttributes(options),
			f.registry.Register,
		)
		return finish.NewFloat64Counter(
			&float64Inst{measures: measures},
			finishers,
		), err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

func (f *factory) Int64UpDownCounter(
	name string,
	options ...metric.Int64UpDownCounterOption,
) (metric.Int64UpDownCounter, error) {
	cfg := metric.NewInt64UpDownCounterConfig(options...)
	counter, err := f.int64UpDownCounters.Lookup(instID{
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
			Scope:       f.meter.scope,
		}
		measures, finishers, err := finishCounterAggregators(
			f.meter.int64Resolver,
			inst,
			defaultAttributes(options),
			f.registry.Register,
		)
		return finish.NewInt64UpDownCounter(
			&int64Inst{measures: measures},
			finishers,
		), err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

func (f *factory) Float64UpDownCounter(
	name string,
	options ...metric.Float64UpDownCounterOption,
) (metric.Float64UpDownCounter, error) {
	cfg := metric.NewFloat64UpDownCounterConfig(options...)
	counter, err := f.float64UpDownCounters.Lookup(instID{
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
			Scope:       f.meter.scope,
		}
		measures, finishers, err := finishCounterAggregators(
			f.meter.float64Resolver,
			inst,
			defaultAttributes(options),
			f.registry.Register,
		)
		return finish.NewFloat64UpDownCounter(
			&float64Inst{measures: measures},
			finishers,
		), err
	})
	if err != nil {
		return counter, err
	}
	return counter, validateInstrumentName(name)
}

// finishCounterAggregators adapts the SDK pipeline resolver to the
// internal Finish implementation.
func finishCounterAggregators[N int64 | float64](
	r resolver[N],
	id Instrument,
	allowedKeys []attribute.Key,
	registerShutdown func(func(), func(context.Context) error),
) ([]aggregate.Measure[N], []finish.Func, error) {
	var (
		measures  []aggregate.Measure[N]
		finishers []finish.Func
		err       error
	)
	for _, i := range r.inserters {
		newAggregation := func(
			b aggregate.Builder[N],
			agg Aggregation,
			kind InstrumentKind,
		) (streamAggregation[N], error) {
			if kind == InstrumentKindCounter || kind == InstrumentKindUpDownCounter {
				if _, ok := agg.(AggregationSum); ok {
					sum := b.FinishSum(kind == InstrumentKindCounter)
					registerShutdown(sum.Stop, sum.Wait)
					return streamAggregation[N]{
						measure: sum.Measure,
						compute: sum.ComputeAggregation,
						finish:  sum.Finish,
					}, nil
				}
			}
			return i.newStreamAggregation(b, agg, kind)
		}
		e := i.forEachStream(
			id,
			allowedKeys,
			i.readerDefaultAggregation(id.Kind),
			newAggregation,
			func(agg aggVal[N]) {
				measures = append(measures, agg.Measure)
				if agg.Finish != nil {
					finishers = append(finishers, agg.Finish)
				}
			},
		)
		if e != nil {
			err = errors.Join(err, e)
		}
	}
	return measures, finishers, err
}
