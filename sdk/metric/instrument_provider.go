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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
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

// lookupSync returns the resolved synchronous instrumentImpl.
func (p *instProvider[N]) lookupSync(kind InstrumentKind, name string, opts []instrument.SynchronousOption) (*instrumentImpl[N], error) {
	cfg := instrument.NewSynchronousConfig(opts...)
	return p.lookup(Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        kind,
		Scope:       p.scope,
	})
}

// lookupAsync returns the resolved asyncchronous instrumentImpl.
func (p *instProvider[N]) lookupAsync(kind InstrumentKind, name string, opts []instrument.AsynchronousOption) (*instrumentImpl[N], error) {
	cfg := instrument.NewAsynchronousConfig(opts...)
	inst, err := p.lookup(Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        kind,
		Scope:       p.scope,
	})
	if err != nil {
		return nil, err
	}

	if inst == nil {
		// Drop aggregator.
		return nil, nil
	}

	for _, cBack := range cfg.Callbacks() {
		p.pipes.registerCallback(callback{inst: inst, f: cBack})
	}

	return inst, nil
}

func (p *instProvider[N]) lookup(inst Instrument) (*instrumentImpl[N], error) {
	aggs, err := p.resolve.Aggregators(inst)
	return &instrumentImpl[N]{aggregators: aggs}, err
}

type asyncInt64Provider struct {
	*instProvider[int64]
}

var _ asyncint64.InstrumentProvider = asyncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p asyncInt64Provider) Counter(name string, opts ...instrument.AsynchronousOption) (asyncint64.Counter, error) {
	return p.lookupAsync(InstrumentKindAsyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncInt64Provider) UpDownCounter(name string, opts ...instrument.AsynchronousOption) (asyncint64.UpDownCounter, error) {
	return p.lookupAsync(InstrumentKindAsyncUpDownCounter, name, opts)
}

// Gauge creates an instrument for recording the current value.
func (p asyncInt64Provider) Gauge(name string, opts ...instrument.AsynchronousOption) (asyncint64.Gauge, error) {
	return p.lookupAsync(InstrumentKindAsyncGauge, name, opts)
}

type asyncFloat64Provider struct {
	*instProvider[float64]
}

var _ asyncfloat64.InstrumentProvider = asyncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p asyncFloat64Provider) Counter(name string, opts ...instrument.AsynchronousOption) (asyncfloat64.Counter, error) {
	return p.lookupAsync(InstrumentKindAsyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncFloat64Provider) UpDownCounter(name string, opts ...instrument.AsynchronousOption) (asyncfloat64.UpDownCounter, error) {
	return p.lookupAsync(InstrumentKindAsyncUpDownCounter, name, opts)
}

// Gauge creates an instrument for recording the current value.
func (p asyncFloat64Provider) Gauge(name string, opts ...instrument.AsynchronousOption) (asyncfloat64.Gauge, error) {
	return p.lookupAsync(InstrumentKindAsyncGauge, name, opts)
}

type syncInt64Provider struct {
	*instProvider[int64]
}

var _ syncint64.InstrumentProvider = syncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncInt64Provider) Counter(name string, opts ...instrument.SynchronousOption) (syncint64.Counter, error) {
	return p.lookupSync(InstrumentKindSyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncInt64Provider) UpDownCounter(name string, opts ...instrument.SynchronousOption) (syncint64.UpDownCounter, error) {
	return p.lookupSync(InstrumentKindSyncUpDownCounter, name, opts)
}

// Histogram creates an instrument for recording the current value.
func (p syncInt64Provider) Histogram(name string, opts ...instrument.SynchronousOption) (syncint64.Histogram, error) {
	return p.lookupSync(InstrumentKindSyncHistogram, name, opts)
}

type syncFloat64Provider struct {
	*instProvider[float64]
}

var _ syncfloat64.InstrumentProvider = syncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncFloat64Provider) Counter(name string, opts ...instrument.SynchronousOption) (syncfloat64.Counter, error) {
	return p.lookupSync(InstrumentKindSyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncFloat64Provider) UpDownCounter(name string, opts ...instrument.SynchronousOption) (syncfloat64.UpDownCounter, error) {
	return p.lookupSync(InstrumentKindSyncUpDownCounter, name, opts)
}

// Histogram creates an instrument for recording the current value.
func (p syncFloat64Provider) Histogram(name string, opts ...instrument.SynchronousOption) (syncfloat64.Histogram, error) {
	return p.lookupSync(InstrumentKindSyncHistogram, name, opts)
}
