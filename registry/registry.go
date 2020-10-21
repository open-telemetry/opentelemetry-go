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

package registry // import "go.opentelemetry.io/otel/registry"

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

// MeterProvider is a standard MeterProvider for wrapping `MeterImpl`
type MeterProvider struct {
	impl otel.MeterImpl
}

var _ otel.MeterProvider = (*MeterProvider)(nil)

// uniqueInstrumentMeterImpl implements the otel.MeterImpl interface, adding
// uniqueness checking for instrument descriptors.  Use NewUniqueInstrumentMeter
// to wrap an implementation with uniqueness checking.
type uniqueInstrumentMeterImpl struct {
	lock  sync.Mutex
	impl  otel.MeterImpl
	state map[key]otel.InstrumentImpl
}

var _ otel.MeterImpl = (*uniqueInstrumentMeterImpl)(nil)

type key struct {
	instrumentName         string
	instrumentationName    string
	InstrumentationVersion string
}

// NewMeterProvider returns a new provider that implements instrument
// name-uniqueness checking.
func NewMeterProvider(impl otel.MeterImpl) *MeterProvider {
	return &MeterProvider{
		impl: NewUniqueInstrumentMeterImpl(impl),
	}
}

// Meter implements MeterProvider.
func (p *MeterProvider) Meter(instrumentationName string, opts ...otel.MeterOption) otel.Meter {
	return otel.WrapMeterImpl(p.impl, instrumentationName, opts...)
}

// ErrMetricKindMismatch is the standard error for mismatched metric
// instrument definitions.
var ErrMetricKindMismatch = fmt.Errorf(
	"A metric was already registered by this name with another kind or number type")

// NewUniqueInstrumentMeterImpl returns a wrapped otel.MeterImpl with
// the addition of uniqueness checking.
func NewUniqueInstrumentMeterImpl(impl otel.MeterImpl) otel.MeterImpl {
	return &uniqueInstrumentMeterImpl{
		impl:  impl,
		state: map[key]otel.InstrumentImpl{},
	}
}

// RecordBatch implements otel.MeterImpl.
func (u *uniqueInstrumentMeterImpl) RecordBatch(ctx context.Context, labels []label.KeyValue, ms ...otel.Measurement) {
	u.impl.RecordBatch(ctx, labels, ms...)
}

func keyOf(descriptor otel.Descriptor) key {
	return key{
		descriptor.Name(),
		descriptor.InstrumentationName(),
		descriptor.InstrumentationVersion(),
	}
}

// NewMetricKindMismatchError formats an error that describes a
// mismatched metric instrument definition.
func NewMetricKindMismatchError(desc otel.Descriptor) error {
	return fmt.Errorf("Metric was %s (%s %s)registered as a %s %s: %w",
		desc.Name(),
		desc.InstrumentationName(),
		desc.InstrumentationVersion(),
		desc.NumberKind(),
		desc.InstrumentKind(),
		ErrMetricKindMismatch)
}

// Compatible determines whether two otel.Descriptors are considered
// the same for the purpose of uniqueness checking.
func Compatible(candidate, existing otel.Descriptor) bool {
	return candidate.InstrumentKind() == existing.InstrumentKind() &&
		candidate.NumberKind() == existing.NumberKind()
}

// checkUniqueness returns an ErrMetricKindMismatch error if there is
// a conflict between a descriptor that was already registered and the
// `descriptor` argument.  If there is an existing compatible
// registration, this returns the already-registered instrument.  If
// there is no conflict and no prior registration, returns (nil, nil).
func (u *uniqueInstrumentMeterImpl) checkUniqueness(descriptor otel.Descriptor) (otel.InstrumentImpl, error) {
	impl, ok := u.state[keyOf(descriptor)]
	if !ok {
		return nil, nil
	}

	if !Compatible(descriptor, impl.Descriptor()) {
		return nil, NewMetricKindMismatchError(impl.Descriptor())
	}

	return impl, nil
}

// NewSyncInstrument implements otel.MeterImpl.
func (u *uniqueInstrumentMeterImpl) NewSyncInstrument(descriptor otel.Descriptor) (otel.SyncImpl, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	impl, err := u.checkUniqueness(descriptor)

	if err != nil {
		return nil, err
	} else if impl != nil {
		return impl.(otel.SyncImpl), nil
	}

	syncInst, err := u.impl.NewSyncInstrument(descriptor)
	if err != nil {
		return nil, err
	}
	u.state[keyOf(descriptor)] = syncInst
	return syncInst, nil
}

// NewAsyncInstrument implements otel.MeterImpl.
func (u *uniqueInstrumentMeterImpl) NewAsyncInstrument(
	descriptor otel.Descriptor,
	runner otel.AsyncRunner,
) (otel.AsyncImpl, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	impl, err := u.checkUniqueness(descriptor)

	if err != nil {
		return nil, err
	} else if impl != nil {
		return impl.(otel.AsyncImpl), nil
	}

	asyncInst, err := u.impl.NewAsyncInstrument(descriptor, runner)
	if err != nil {
		return nil, err
	}
	u.state[keyOf(descriptor)] = asyncInst
	return asyncInst, nil
}
