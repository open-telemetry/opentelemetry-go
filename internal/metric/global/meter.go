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

package global // import "go.opentelemetry.io/otel/internal/metric/global"

import (
	"sync"
	"unsafe"

	"go.opentelemetry.io/otel/internal/metric/registry"
	"go.opentelemetry.io/otel/metric"
)

// This file contains the forwarding implementation of MeterProvider used as
// the default global instance.  Metric events using instruments provided by
// this implementation are no-ops until the first Meter implementation is set
// as the global provider.
//
// The implementation here uses Mutexes to maintain a list of active Meters in
// the MeterProvider and Instruments in each Meter, under the assumption that
// these interfaces are not performance-critical.
//
// We have the invariant that setDelegate() will be called before a new
// MeterProvider implementation is registered as the global provider.  Mutexes
// in the MeterProvider and Meters ensure that each instrument has a delegate
// before the global provider is set.
//
// Metric uniqueness checking is implemented by calling the exported
// methods of the api/metric/registry package.

type meterKey struct {
	InstrumentationName    string
	InstrumentationVersion string
	SchemaURL              string
}

type meterProvider struct {
	delegate metric.MeterProvider

	// lock protects `delegate` and `meters`.
	lock sync.Mutex

	// meters maintains a unique entry for every named Meter
	// that has been registered through the global instance.
	meters map[meterKey]*meterEntry
}

type meterEntry struct {
	unique metric.Meter
	impl   metric.Meter
}

var _ metric.MeterProvider = &meterProvider{}

// MeterProvider interface and delegation

func newMeterProvider() *meterProvider {
	return &meterProvider{
		meters: map[meterKey]*meterEntry{},
	}
}

func (p *meterProvider) setDelegate(provider metric.MeterProvider) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.delegate = provider
	for _, entry := range p.meters {
		entry.impl.(interface {
			setDelegate(provider metric.MeterProvider)
		}).setDelegate(provider)
	}
	p.meters = nil
}

func (p *meterProvider) Meter(instrumentationName string, opts ...metric.MeterOption) metric.Meter {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.delegate != nil {
		return p.delegate.Meter(instrumentationName, opts...)
	}

	cfg := metric.NewMeterConfig(opts...)
	key := meterKey{
		InstrumentationName:    instrumentationName,
		InstrumentationVersion: cfg.InstrumentationVersion(),
		SchemaURL:              cfg.SchemaURL(),
	}
	entry, ok := p.meters[key]
	if !ok {
		entry = &meterEntry{
			impl: newMeterDelegate(instrumentationName, opts...),
		}
		// Note: This code implements its own MeterProvider
		// name-uniqueness logic because there is
		// synchronization required at the moment of
		// delegation.  We use the same instrument-uniqueness
		// checking the real SDK uses here:
		entry.unique = registry.NewUniqueInstrumentMeter(entry.impl)
		p.meters[key] = entry
	}
	return entry.unique
}

func AtomicFieldOffsets() map[string]uintptr {
	return map[string]uintptr{
		"meterProvider.delegate": unsafe.Offsetof(meterProvider{}.delegate),
		"meterImpl.delegate":     unsafe.Offsetof(meterDelegate{}.delegate),
	}
}
