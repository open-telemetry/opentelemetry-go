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
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// meterProvider is a placeholder for a configured SDK MeterProvider.
//
// All MeterProvider functionality is forwarded to a delegate once
// configured.
type meterProvider struct {
	mtx    sync.Mutex
	meters map[il]*meter

	delegate metric.MeterProvider
}

type il struct {
	name    string
	version string
}

// setDelegate configures p to delegate all MeterProvider functionality to
// provider.
//
// All Meters provided prior to this function call are switched out to be
// Meters provided by provider. All instruments and callbacks are recreated and
// delegated.
//
// It is guaranteed by the caller that this happens only once.
func (p *meterProvider) setDelegate(provider metric.MeterProvider) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.delegate = provider

	if len(p.meters) == 0 {
		return
	}

	for _, meter := range p.meters {
		meter.setDelegate(provider)
	}

	p.meters = nil
}

// Meter implements MeterProvider.
func (p *meterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.delegate != nil {
		return p.delegate.Meter(name, opts...)
	}

	// At this moment it is guaranteed that no sdk is installed, save the meter in the meters map.

	c := metric.NewMeterConfig(opts...)
	key := il{
		name:    name,
		version: c.InstrumentationVersion(),
	}

	if p.meters == nil {
		p.meters = make(map[il]*meter)
	}

	if val, ok := p.meters[key]; ok {
		return val
	}

	t := &meter{name: name, opts: opts}
	p.meters[key] = t
	return t
}

// meter is a placeholder for a metric.Meter.
//
// All Meter functionality is forwarded to a delegate once configured.
// Otherwise, all functionality is forwarded to a NoopMeter.
type meter struct {
	name string
	opts []metric.MeterOption

	mtx         sync.Mutex
	instruments []delegatedInstrument
	callbacks   []delegatedCallback

	delegate atomic.Value // metric.Meter
}

type delegatedInstrument interface {
	setDelegate(metric.Meter)
}

// setDelegate configures m to delegate all Meter functionality to Meters
// created by provider.
//
// All subsequent calls to the Meter methods will be passed to the delegate.
//
// It is guaranteed by the caller that this happens only once.
func (m *meter) setDelegate(provider metric.MeterProvider) {
	meter := provider.Meter(m.name, m.opts...)
	m.delegate.Store(meter)

	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, inst := range m.instruments {
		inst.setDelegate(meter)
	}

	for _, callback := range m.callbacks {
		callback.setDelegate(meter)
	}

	m.instruments = nil
	m.callbacks = nil
}

func (m *meter) Float64Counter(name string, opts ...metric.InstrumentOption) (metric.Float64Counter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Float64Counter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &sfCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Float64UpDownCounter(name string, opts ...metric.InstrumentOption) (metric.Float64UpDownCounter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Float64UpDownCounter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &sfUpDownCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Float64Histogram(name string, opts ...metric.InstrumentOption) (metric.Float64Histogram, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Float64Histogram(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	hist := &sfHistogram{name: name, opts: opts}
	m.instruments = append(m.instruments, hist)
	return hist, nil
}

func (m *meter) Float64ObservableCounter(name string, opts ...metric.ObservableOption) (metric.Float64ObservableCounter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Float64ObservableCounter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &afCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Float64ObservableUpDownCounter(name string, opts ...metric.ObservableOption) (metric.Float64ObservableUpDownCounter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Float64ObservableUpDownCounter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &afUpDownCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Float64ObservableGauge(name string, opts ...metric.ObservableOption) (metric.Float64ObservableGauge, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Float64ObservableGauge(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	g := &afGauge{name: name, opts: opts}
	m.instruments = append(m.instruments, g)
	return g, nil
}

func (m *meter) Int64Counter(name string, opts ...metric.InstrumentOption) (metric.Int64Counter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Int64Counter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &siCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Int64UpDownCounter(name string, opts ...metric.InstrumentOption) (metric.Int64UpDownCounter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Int64UpDownCounter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &siUpDownCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Int64Histogram(name string, opts ...metric.InstrumentOption) (metric.Int64Histogram, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Int64Histogram(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	hist := &siHistogram{name: name, opts: opts}
	m.instruments = append(m.instruments, hist)
	return hist, nil
}

func (m *meter) Int64ObservableCounter(name string, opts ...metric.ObservableOption) (metric.Int64ObservableCounter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Int64ObservableCounter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &aiCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Int64ObservableUpDownCounter(name string, opts ...metric.ObservableOption) (metric.Int64ObservableUpDownCounter, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Int64ObservableUpDownCounter(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	ctr := &aiUpDownCounter{name: name, opts: opts}
	m.instruments = append(m.instruments, ctr)
	return ctr, nil
}

func (m *meter) Int64ObservableGauge(name string, opts ...metric.ObservableOption) (metric.Int64ObservableGauge, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.Int64ObservableGauge(name, opts...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	g := &aiGauge{name: name, opts: opts}
	m.instruments = append(m.instruments, g)
	return g, nil
}

func (m *meter) RegisterCallback(f metric.Callback, instrument metric.Observable, additional ...metric.Observable) (metric.Unregisterer, error) {
	if del, ok := m.delegate.Load().(metric.Meter); ok {
		return del.RegisterCallback(f, unwrapInst(instrument), unwrapInsts(additional)...)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	cback := newDelegatedCallback(f, instrument, additional...)
	m.callbacks = append(m.callbacks, cback)

	return &cback, nil
}

type unregisterer []func() error

func (u unregisterer) Unregister() error {
	for _, f := range u {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

type wrapped interface {
	unwrap() metric.Observable
}

func unwrapInsts(instruments []metric.Observable) []metric.Observable {
	for i := range instruments {
		instruments[i] = unwrapInst(instruments[i])
	}
	return instruments
}

func unwrapInst(inst metric.Observable) metric.Observable {
	if in, ok := inst.(wrapped); ok {
		return in.unwrap()
	}
	return inst
}

type delegatedCallback struct {
	instrument metric.Observable
	additional []metric.Observable
	f          metric.Callback

	unregistered uint32
	unregister   atomic.Value // func() error
}

func newDelegatedCallback(f metric.Callback, instrument metric.Observable, additional ...metric.Observable) delegatedCallback {
	cback := delegatedCallback{f: f, instrument: instrument, additional: additional}
	cback.unregister.Store(func() error {
		atomic.StoreUint32(&cback.unregistered, 1)
		return nil
	})
	return cback
}

func (c *delegatedCallback) setDelegate(m metric.Meter) {
	u, err := m.RegisterCallback(c.f, unwrapInst(c.instrument), unwrapInsts(c.additional)...)
	if err != nil {
		otel.Handle(err)
		return
	}
	if u != nil {
		c.unregister.Store(u.Unregister)
	}

	if unregistered := atomic.LoadUint32(&c.unregistered); unregistered > 0 {
		err := u.Unregister()
		if err != nil {
			otel.Handle(err)
		}
		return
	}
}

func (c *delegatedCallback) Unregister() error {
	f := c.unregister.Load().(func() error)
	return f()
}
