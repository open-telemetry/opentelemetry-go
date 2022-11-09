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

package global // import "go.opentelemetry.io/otel/metric/internal/global"

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type afCounter struct {
	name string
	opts []metric.ObservableOption

	delegate atomic.Value // metric.Float64ObservableCounter

	metric.Observable
}

func (i *afCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Float64ObservableCounter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *afCounter) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Float64ObservableCounter).Observe(ctx, x, attrs...)
	}
}

func (i *afCounter) unwrap() metric.Observable {
	if ctr := i.delegate.Load(); ctr != nil {
		return ctr.(metric.Float64ObservableCounter)
	}
	return nil
}

type afUpDownCounter struct {
	name string
	opts []metric.ObservableOption

	delegate atomic.Value // metric.Float64ObservableUpDownCounter

	metric.Observable
}

func (i *afUpDownCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Float64ObservableUpDownCounter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *afUpDownCounter) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Float64ObservableUpDownCounter).Observe(ctx, x, attrs...)
	}
}

func (i *afUpDownCounter) unwrap() metric.Observable {
	if ctr := i.delegate.Load(); ctr != nil {
		return ctr.(metric.Float64ObservableUpDownCounter)
	}
	return nil
}

type afGauge struct {
	name string
	opts []metric.ObservableOption

	delegate atomic.Value // metric.Float64ObservableGauge

	metric.Observable
}

func (i *afGauge) setDelegate(m metric.Meter) {
	ctr, err := m.Float64ObservableGauge(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *afGauge) Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Float64ObservableGauge).Observe(ctx, x, attrs...)
	}
}

func (i *afGauge) unwrap() metric.Observable {
	if ctr := i.delegate.Load(); ctr != nil {
		return ctr.(metric.Float64ObservableGauge)
	}
	return nil
}

type aiCounter struct {
	name string
	opts []metric.ObservableOption

	delegate atomic.Value // metric.Int64ObservableCounter

	metric.Observable
}

func (i *aiCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Int64ObservableCounter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *aiCounter) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Int64ObservableCounter).Observe(ctx, x, attrs...)
	}
}

func (i *aiCounter) unwrap() metric.Observable {
	if ctr := i.delegate.Load(); ctr != nil {
		return ctr.(metric.Int64ObservableCounter)
	}
	return nil
}

type aiUpDownCounter struct {
	name string
	opts []metric.ObservableOption

	delegate atomic.Value // metric.Int64ObservableUpDownCounter

	metric.Observable
}

func (i *aiUpDownCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Int64ObservableUpDownCounter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *aiUpDownCounter) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Int64ObservableUpDownCounter).Observe(ctx, x, attrs...)
	}
}

func (i *aiUpDownCounter) unwrap() metric.Observable {
	if ctr := i.delegate.Load(); ctr != nil {
		return ctr.(metric.Int64ObservableUpDownCounter)
	}
	return nil
}

type aiGauge struct {
	name string
	opts []metric.ObservableOption

	delegate atomic.Value // metric.Int64ObservableGauge

	metric.Observable
}

func (i *aiGauge) setDelegate(m metric.Meter) {
	ctr, err := m.Int64ObservableGauge(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *aiGauge) Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Int64ObservableGauge).Observe(ctx, x, attrs...)
	}
}

func (i *aiGauge) unwrap() metric.Observable {
	if ctr := i.delegate.Load(); ctr != nil {
		return ctr.(metric.Int64ObservableGauge)
	}
	return nil
}

// Sync Instruments.
type sfCounter struct {
	name string
	opts []metric.InstrumentOption

	delegate atomic.Value // metric.Float64Counter
}

func (i *sfCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Float64Counter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *sfCounter) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Float64Counter).Add(ctx, incr, attrs...)
	}
}

type sfUpDownCounter struct {
	name string
	opts []metric.InstrumentOption

	delegate atomic.Value // metric.Float64UpDownCounter
}

func (i *sfUpDownCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Float64UpDownCounter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *sfUpDownCounter) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Float64UpDownCounter).Add(ctx, incr, attrs...)
	}
}

type sfHistogram struct {
	name string
	opts []metric.InstrumentOption

	delegate atomic.Value // metric.Float64Histogram
}

func (i *sfHistogram) setDelegate(m metric.Meter) {
	ctr, err := m.Float64Histogram(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *sfHistogram) Record(ctx context.Context, x float64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Float64Histogram).Record(ctx, x, attrs...)
	}
}

type siCounter struct {
	name string
	opts []metric.InstrumentOption

	delegate atomic.Value // metric.Int64Counter
}

func (i *siCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Int64Counter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *siCounter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Int64Counter).Add(ctx, x, attrs...)
	}
}

type siUpDownCounter struct {
	name string
	opts []metric.InstrumentOption

	delegate atomic.Value // metric.Int64UpDownCounter
}

func (i *siUpDownCounter) setDelegate(m metric.Meter) {
	ctr, err := m.Int64UpDownCounter(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *siUpDownCounter) Add(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Int64UpDownCounter).Add(ctx, x, attrs...)
	}
}

type siHistogram struct {
	name string
	opts []metric.InstrumentOption

	delegate atomic.Value // metric.Int64Histogram
}

func (i *siHistogram) setDelegate(m metric.Meter) {
	ctr, err := m.Int64Histogram(i.name, i.opts...)
	if err != nil {
		otel.Handle(err)
		return
	}
	i.delegate.Store(ctr)
}

func (i *siHistogram) Record(ctx context.Context, x int64, attrs ...attribute.KeyValue) {
	if ctr := i.delegate.Load(); ctr != nil {
		ctr.(metric.Int64Histogram).Record(ctx, x, attrs...)
	}
}
