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

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// instProvider provides all OpenTelemetry instruments.
type instProvider[N int64 | float64] struct {
	scope   instrumentation.Scope
	pipes   pipelines
	resolve resolver[N]
}

func newInstProvider[N int64 | float64](s instrumentation.Scope, p pipelines, c instrumentCache[N]) *instProvider[N] {
	return &instProvider[N]{scope: s, pipes: p, resolve: newResolver(p, c)}
}

func (p *instProvider[N]) lookup(kind InstrumentKind, name, desc string, u unit.Unit) (*instrumentImpl[N], error) {
	inst := Instrument{
		Name:        name,
		Description: desc,
		Unit:        u,
		Kind:        kind,
		Scope:       p.scope,
	}
	aggs, err := p.resolve.Aggregators(inst)
	return &instrumentImpl[N]{aggregators: aggs}, err
}

type asyncInt64Provider struct {
	*instProvider[int64]
}

var _ asyncint64.InstrumentProvider = asyncInt64Provider{}

func (p asyncInt64Provider) registerCallbacks(inst *instrumentImpl[int64], cBacks []instrument.Int64Callback) {
	if inst == nil {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p asyncInt64Provider) callback(i *instrumentImpl[int64], f instrument.Int64Callback) func(context.Context) error {
	return func(ctx context.Context) error { return f(ctx, i) }
}

// Counter creates an instrument for recording increasing values.
func (p asyncInt64Provider) Counter(name string, opts ...instrument.Int64ObserverOption) (asyncint64.Counter, error) {
	cfg := instrument.NewInt64ObserverConfig(opts...)
	inst, err := p.lookup(InstrumentKindAsyncCounter, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncInt64Provider) UpDownCounter(name string, opts ...instrument.Int64ObserverOption) (asyncint64.UpDownCounter, error) {
	cfg := instrument.NewInt64ObserverConfig(opts...)
	inst, err := p.lookup(InstrumentKindAsyncUpDownCounter, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Gauge creates an instrument for recording the current value.
func (p asyncInt64Provider) Gauge(name string, opts ...instrument.Int64ObserverOption) (asyncint64.Gauge, error) {
	cfg := instrument.NewInt64ObserverConfig(opts...)
	inst, err := p.lookup(InstrumentKindAsyncGauge, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

type asyncFloat64Provider struct {
	*instProvider[float64]
}

var _ asyncfloat64.InstrumentProvider = asyncFloat64Provider{}

func (p asyncFloat64Provider) registerCallbacks(inst *instrumentImpl[float64], cBacks []instrument.Float64Callback) {
	if inst == nil {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p asyncFloat64Provider) callback(i *instrumentImpl[float64], f instrument.Float64Callback) func(context.Context) error {
	return func(ctx context.Context) error { return f(ctx, i) }
}

// Counter creates an instrument for recording increasing values.
func (p asyncFloat64Provider) Counter(name string, opts ...instrument.Float64ObserverOption) (asyncfloat64.Counter, error) {
	cfg := instrument.NewFloat64ObserverConfig(opts...)
	inst, err := p.lookup(InstrumentKindAsyncCounter, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncFloat64Provider) UpDownCounter(name string, opts ...instrument.Float64ObserverOption) (asyncfloat64.UpDownCounter, error) {
	cfg := instrument.NewFloat64ObserverConfig(opts...)
	inst, err := p.lookup(InstrumentKindAsyncUpDownCounter, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Gauge creates an instrument for recording the current value.
func (p asyncFloat64Provider) Gauge(name string, opts ...instrument.Float64ObserverOption) (asyncfloat64.Gauge, error) {
	cfg := instrument.NewFloat64ObserverConfig(opts...)
	inst, err := p.lookup(InstrumentKindAsyncGauge, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

type syncInt64Provider struct {
	*instProvider[int64]
}

var _ syncint64.InstrumentProvider = syncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncInt64Provider) Counter(name string, opts ...instrument.Int64Option) (syncint64.Counter, error) {
	cfg := instrument.NewInt64Config(opts...)
	return p.lookup(InstrumentKindSyncCounter, name, cfg.Description(), cfg.Unit())
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncInt64Provider) UpDownCounter(name string, opts ...instrument.Int64Option) (syncint64.UpDownCounter, error) {
	cfg := instrument.NewInt64Config(opts...)
	return p.lookup(InstrumentKindSyncUpDownCounter, name, cfg.Description(), cfg.Unit())
}

// Histogram creates an instrument for recording the current value.
func (p syncInt64Provider) Histogram(name string, opts ...instrument.Int64Option) (syncint64.Histogram, error) {
	cfg := instrument.NewInt64Config(opts...)
	return p.lookup(InstrumentKindSyncHistogram, name, cfg.Description(), cfg.Unit())
}

type syncFloat64Provider struct {
	*instProvider[float64]
}

var _ syncfloat64.InstrumentProvider = syncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncFloat64Provider) Counter(name string, opts ...instrument.Float64Option) (syncfloat64.Counter, error) {
	cfg := instrument.NewFloat64Config(opts...)
	return p.lookup(InstrumentKindSyncCounter, name, cfg.Description(), cfg.Unit())
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncFloat64Provider) UpDownCounter(name string, opts ...instrument.Float64Option) (syncfloat64.UpDownCounter, error) {
	cfg := instrument.NewFloat64Config(opts...)
	return p.lookup(InstrumentKindSyncUpDownCounter, name, cfg.Description(), cfg.Unit())
}

// Histogram creates an instrument for recording the current value.
func (p syncFloat64Provider) Histogram(name string, opts ...instrument.Float64Option) (syncfloat64.Histogram, error) {
	cfg := instrument.NewFloat64Config(opts...)
	return p.lookup(InstrumentKindSyncHistogram, name, cfg.Description(), cfg.Unit())
}
