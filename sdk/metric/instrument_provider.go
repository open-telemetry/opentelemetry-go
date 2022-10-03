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
	"fmt"

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

// instProviderKey uniquely describes an instrument creation request received
// by an instrument provider.
type instProviderKey struct {
	// Name is the name of the instrument.
	Name string
	// Description is the description of the instrument.
	Description string
	// Unit is the unit of the instrument.
	Unit unit.Unit
	// Kind is the instrument Kind provided.
	Kind view.InstrumentKind
}

// viewInst returns the instProviderKey as a view Instrument using scope s.
func (k instProviderKey) viewInst(s instrumentation.Scope) view.Instrument {
	return view.Instrument{
		Scope:       s,
		Name:        k.Name,
		Description: k.Description,
		Kind:        k.Kind,
	}
}

// instProvider provides all OpenTelemetry instruments.
type instProvider[N int64 | float64] struct {
	resolve resolver[N]
}

func newInstProvider[N int64 | float64](r resolver[N]) *instProvider[N] {
	return &instProvider[N]{resolve: r}
}

// lookup returns the resolved instrumentImpl.
func (p *instProvider[N]) lookup(kind view.InstrumentKind, name string, opts []instrument.Option) (*instrumentImpl[N], error) {
	cfg := instrument.NewConfig(opts...)
	key := instProviderKey{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        kind,
	}

	aggs, err := p.resolve.Aggregators(key)
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[N]{aggregators: aggs}, err
}

type asyncInt64Provider struct {
	*instProvider[int64]
}

var _ asyncint64.InstrumentProvider = asyncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p asyncInt64Provider) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	return p.lookup(view.AsyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncInt64Provider) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	return p.lookup(view.AsyncUpDownCounter, name, opts)
}

// Gauge creates an instrument for recording the current value.
func (p asyncInt64Provider) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	return p.lookup(view.AsyncGauge, name, opts)
}

type asyncFloat64Provider struct {
	*instProvider[float64]
}

var _ asyncfloat64.InstrumentProvider = asyncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p asyncFloat64Provider) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	return p.lookup(view.AsyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncFloat64Provider) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	return p.lookup(view.AsyncUpDownCounter, name, opts)
}

// Gauge creates an instrument for recording the current value.
func (p asyncFloat64Provider) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	return p.lookup(view.AsyncGauge, name, opts)
}

type syncInt64Provider struct {
	*instProvider[int64]
}

var _ syncint64.InstrumentProvider = syncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncInt64Provider) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	return p.lookup(view.SyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncInt64Provider) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	return p.lookup(view.SyncUpDownCounter, name, opts)
}

// Histogram creates an instrument for recording the current value.
func (p syncInt64Provider) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	return p.lookup(view.SyncHistogram, name, opts)
}

type syncFloat64Provider struct {
	*instProvider[float64]
}

var _ syncfloat64.InstrumentProvider = syncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncFloat64Provider) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	return p.lookup(view.SyncCounter, name, opts)
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncFloat64Provider) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	return p.lookup(view.SyncUpDownCounter, name, opts)
}

// Histogram creates an instrument for recording the current value.
func (p syncFloat64Provider) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	return p.lookup(view.SyncHistogram, name, opts)
}
