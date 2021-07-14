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

package registry // import "go.opentelemetry.io/otel/metric/registry"

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// UniqueMeterProvider is a standard MeterProvider for wrapping `MeterImpl`
type UniqueMeterProvider struct {
	lock   sync.Mutex
	impl   metric.MeterProvider
	meters map[uniqueMeterKey]metric.Meter
}

type uniqueMeterKey struct {
	name      string
	version   string
	schemaURL string
}

func keyOf(name string, opts ...metric.MeterOption) uniqueMeterKey {
	cfg := metric.NewMeterConfig(opts...)
	return uniqueMeterKey{
		name:      name,
		version:   cfg.InstrumentationVersion(),
		schemaURL: cfg.SchemaURL(),
	}
}

var _ metric.MeterProvider = (*UniqueMeterProvider)(nil)

// NewUniqueMeterProvider returns a new provider that implements meter
// name-uniqueness checking.
func NewUniqueMeterProvider(impl metric.MeterProvider) *UniqueMeterProvider {
	return &UniqueMeterProvider{
		impl:   impl,
		meters: map[uniqueMeterKey]metric.Meter{},
	}
}

// Meter implements MeterProvider.
func (p *UniqueMeterProvider) Meter(instrumentationName string, opts ...metric.MeterOption) metric.Meter {
	k := keyOf(instrumentationName, opts...)
	p.lock.Lock()
	defer p.lock.Unlock()
	m, ok := p.meters[k]
	if !ok {
		m = metric.WrapMeterImpl(
			NewUniqueInstrumentMeterImpl(
				p.impl.Meter(instrumentationName, opts...).MeterImpl()))
		p.meters[k] = m
	}
	return m

}

// List provides a list of MeterImpl objects created through this
// provider.
func (p *UniqueMeterProvider) List() []metric.MeterImpl {
	p.lock.Lock()
	defer p.lock.Unlock()

	var r []metric.MeterImpl
	for _, meter := range p.meters {
		r = append(r, meter.MeterImpl().(*uniqueInstrumentMeterImpl).impl)
	}
	return r
}

// uniqueInstrumentMeterImpl implements the metric.MeterImpl interface, adding
// uniqueness checking for instrument descriptors.  Use NewUniqueInstrumentMeter
// to wrap an implementation with uniqueness checking.
type uniqueInstrumentMeterImpl struct {
	lock  sync.Mutex
	impl  metric.MeterImpl
	state map[string]metric.InstrumentImpl
}

var _ metric.MeterImpl = (*uniqueInstrumentMeterImpl)(nil)

// ErrMetricKindMismatch is the standard error for mismatched metric
// instrument definitions.
var ErrMetricKindMismatch = fmt.Errorf(
	"a metric was already registered by this name with another kind or number type")

// NewUniqueInstrumentMeterImpl returns a wrapped metric.MeterImpl with
// the addition of uniqueness checking.
func NewUniqueInstrumentMeterImpl(impl metric.MeterImpl) metric.MeterImpl {
	return &uniqueInstrumentMeterImpl{
		impl:  impl,
		state: map[string]metric.InstrumentImpl{},
	}
}

// RecordBatch implements metric.MeterImpl.
func (u *uniqueInstrumentMeterImpl) RecordBatch(ctx context.Context, labels []attribute.KeyValue, ms ...metric.Measurement) {
	u.impl.RecordBatch(ctx, labels, ms...)
}

// NewMetricKindMismatchError formats an error that describes a
// mismatched metric instrument definition.
func NewMetricKindMismatchError(desc metric.Descriptor) error {
	return fmt.Errorf("metric %s registered as %s %s: %w",
		desc.Name(),
		desc.NumberKind(),
		desc.InstrumentKind(),
		ErrMetricKindMismatch)
}

// Compatible determines whether two metric.Descriptors are considered
// the same for the purpose of uniqueness checking.
func Compatible(candidate, existing metric.Descriptor) bool {
	return candidate.InstrumentKind() == existing.InstrumentKind() &&
		candidate.NumberKind() == existing.NumberKind()
}

// checkUniqueness returns an ErrMetricKindMismatch error if there is
// a conflict between a descriptor that was already registered and the
// `descriptor` argument.  If there is an existing compatible
// registration, this returns the already-registered instrument.  If
// there is no conflict and no prior registration, returns (nil, nil).
func (u *uniqueInstrumentMeterImpl) checkUniqueness(descriptor metric.Descriptor) (metric.InstrumentImpl, error) {
	impl, ok := u.state[descriptor.Name()]
	if !ok {
		return nil, nil
	}

	if !Compatible(descriptor, impl.Descriptor()) {
		return nil, NewMetricKindMismatchError(impl.Descriptor())
	}

	return impl, nil
}

// NewSyncInstrument implements metric.MeterImpl.
func (u *uniqueInstrumentMeterImpl) NewSyncInstrument(descriptor metric.Descriptor) (metric.SyncImpl, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	impl, err := u.checkUniqueness(descriptor)

	if err != nil {
		return nil, err
	} else if impl != nil {
		return impl.(metric.SyncImpl), nil
	}

	syncInst, err := u.impl.NewSyncInstrument(descriptor)
	if err != nil {
		return nil, err
	}
	u.state[descriptor.Name()] = syncInst
	return syncInst, nil
}

// NewAsyncInstrument implements metric.MeterImpl.
func (u *uniqueInstrumentMeterImpl) NewAsyncInstrument(
	descriptor metric.Descriptor,
	runner metric.AsyncRunner,
) (metric.AsyncImpl, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	impl, err := u.checkUniqueness(descriptor)

	if err != nil {
		return nil, err
	} else if impl != nil {
		return impl.(metric.AsyncImpl), nil
	}

	asyncInst, err := u.impl.NewAsyncInstrument(descriptor, runner)
	if err != nil {
		return nil, err
	}
	u.state[descriptor.Name()] = asyncInst
	return asyncInst, nil
}
