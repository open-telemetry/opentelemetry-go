// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/aggregate"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
)

func newFinishMeterFactory() (meterFactoryFn meterFactory, shutdown func()) {
	registry := &finish.Registry{}
	return func(s instrumentation.Scope, p pipelines) metric.Meter {
		m := newMeter(s, p)
		return finish.NewMeter(m, &factory{meter: m, registry: registry})
	}, registry.Shutdown
}

// factory adapts the SDK pipeline resolver to Finish instrument construction.
type factory struct {
	meter         *meter
	registry      *finish.Registry
	int64Counters cacheWithErr[instID, metric.Int64Counter]
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
		measures, finishers, err := finishInt64CounterAggregators(
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

// finishInt64CounterAggregators adapts the SDK pipeline resolver to the
// internal Finish implementation.
func finishInt64CounterAggregators(
	r resolver[int64],
	id Instrument,
	allowedKeys []attribute.Key,
	registerShutdown func(func()),
) ([]aggregate.Measure[int64], []finish.Func, error) {
	var (
		measures  []aggregate.Measure[int64]
		finishers []finish.Func
		err       error
	)
	for _, i := range r.inserters {
		newAggregation := func(
			b aggregate.Builder[int64],
			agg Aggregation,
			kind InstrumentKind,
		) (streamAggregation[int64], error) {
			if kind == InstrumentKindCounter {
				if _, ok := agg.(AggregationSum); ok {
					sum := b.FinishSum(true)
					registerShutdown(sum.Shutdown)
					return streamAggregation[int64]{
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
			func(agg aggVal[int64]) {
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
