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
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

type asyncInt64Provider struct {
	scope   instrumentation.Scope
	resolve *resolver[int64]
}

var _ asyncint64.InstrumentProvider = asyncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p asyncInt64Provider) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.AsyncCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}

	return &instrumentImpl[int64]{
		aggregators: aggs,
	}, err
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncInt64Provider) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.AsyncUpDownCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[int64]{
		aggregators: aggs,
	}, err
}

// Gauge creates an instrument for recording the current value.
func (p asyncInt64Provider) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.AsyncGauge,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[int64]{
		aggregators: aggs,
	}, err
}

type asyncFloat64Provider struct {
	scope   instrumentation.Scope
	resolve *resolver[float64]
}

var _ asyncfloat64.InstrumentProvider = asyncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p asyncFloat64Provider) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.AsyncCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[float64]{
		aggregators: aggs,
	}, err
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p asyncFloat64Provider) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.AsyncUpDownCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[float64]{
		aggregators: aggs,
	}, err
}

// Gauge creates an instrument for recording the current value.
func (p asyncFloat64Provider) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.AsyncGauge,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[float64]{
		aggregators: aggs,
	}, err
}

type syncInt64Provider struct {
	scope   instrumentation.Scope
	resolve *resolver[int64]
}

var _ syncint64.InstrumentProvider = syncInt64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncInt64Provider) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.SyncCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[int64]{
		aggregators: aggs,
	}, err
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncInt64Provider) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.SyncUpDownCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[int64]{
		aggregators: aggs,
	}, err
}

// Histogram creates an instrument for recording the current value.
func (p syncInt64Provider) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.SyncHistogram,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[int64]{
		aggregators: aggs,
	}, err
}

type syncFloat64Provider struct {
	scope   instrumentation.Scope
	resolve *resolver[float64]
}

var _ syncfloat64.InstrumentProvider = syncFloat64Provider{}

// Counter creates an instrument for recording increasing values.
func (p syncFloat64Provider) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.SyncCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[float64]{
		aggregators: aggs,
	}, err
}

// UpDownCounter creates an instrument for recording changes of a value.
func (p syncFloat64Provider) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.SyncUpDownCounter,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[float64]{
		aggregators: aggs,
	}, err
}

// Histogram creates an instrument for recording the current value.
func (p syncFloat64Provider) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	cfg := instrument.NewConfig(opts...)

	aggs, err := p.resolve.Aggregators(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        view.SyncHistogram,
	}, cfg.Unit())
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[float64]{
		aggregators: aggs,
	}, err
}
