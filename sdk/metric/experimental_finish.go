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
	meter           *meter
	registry        *finish.Registry
	int64Counters   cacheWithErr[instID, metric.Int64Counter]
	float64Counters cacheWithErr[instID, metric.Float64Counter]
	int64Gauges     cacheWithErr[instID, metric.Int64Gauge]
	float64Gauges   cacheWithErr[instID, metric.Float64Gauge]
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
		measures, finishers, err := finishAggregators(
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
		measures, finishers, err := finishAggregators(
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

func (f *factory) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	cfg := metric.NewInt64GaugeConfig(options...)
	gauge, err := f.int64Gauges.Lookup(instID{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindGauge,
	}, func() (metric.Int64Gauge, error) {
		inst := Instrument{
			Name:        name,
			Description: cfg.Description(),
			Unit:        cfg.Unit(),
			Kind:        InstrumentKindGauge,
			Scope:       f.meter.scope,
		}
		measures, finishers, err := finishAggregators(
			f.meter.int64Resolver,
			inst,
			defaultAttributes(options),
			f.registry.Register,
		)
		return finish.NewInt64Gauge(
			&int64Inst{measures: measures},
			finishers,
		), err
	})
	if err != nil {
		return gauge, err
	}
	return gauge, validateInstrumentName(name)
}

func (f *factory) Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	cfg := metric.NewFloat64GaugeConfig(options...)
	gauge, err := f.float64Gauges.Lookup(instID{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindGauge,
	}, func() (metric.Float64Gauge, error) {
		inst := Instrument{
			Name:        name,
			Description: cfg.Description(),
			Unit:        cfg.Unit(),
			Kind:        InstrumentKindGauge,
			Scope:       f.meter.scope,
		}
		measures, finishers, err := finishAggregators(
			f.meter.float64Resolver,
			inst,
			defaultAttributes(options),
			f.registry.Register,
		)
		return finish.NewFloat64Gauge(
			&float64Inst{measures: measures},
			finishers,
		), err
	})
	if err != nil {
		return gauge, err
	}
	return gauge, validateInstrumentName(name)
}

// finishAggregators adapts the SDK pipeline resolver to the
// internal Finish implementation.
func finishAggregators[N int64 | float64](
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
			if kind == InstrumentKindCounter {
				if _, ok := agg.(AggregationSum); ok {
					sum := b.FinishSum(true)
					registerShutdown(sum.Stop, sum.Wait)
					return streamAggregation[N]{
						measure: sum.Measure,
						compute: sum.ComputeAggregation,
						finish:  sum.Finish,
					}, nil
				}
			}

			if kind == InstrumentKindGauge {
				if _, ok := agg.(AggregationLastValue); ok {
					gauge := b.FinishLastValue()
					registerShutdown(gauge.Stop, gauge.Wait)
					return streamAggregation[N]{
						measure: gauge.Measure,
						compute: gauge.ComputeAggregation,
						finish:  gauge.Finish,
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
