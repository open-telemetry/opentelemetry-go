// Copyright 2020, OpenTelemetry Authors
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

package registry // import "go.opentelemetry.io/otel/api/metric/registry"

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
)

type uniqueInstrumentMeterImpl struct {
	impl  metric.MeterImpl
	state map[mapKey]metric.InstrumentImpl
	lock  sync.Mutex
}

type mapKey struct {
	name        string
	libraryName string
}

var ErrMetricTypeMismatch = fmt.Errorf("A metric was already registered by this name with another kind or number type")

var _ metric.MeterImpl = (*uniqueInstrumentMeterImpl)(nil)

func NewUniqueInstrumentMeter(impl metric.MeterImpl) metric.MeterImpl {
	return &uniqueInstrumentMeterImpl{
		impl:  impl,
		state: map[mapKey]metric.InstrumentImpl{},
	}
}

func (u *uniqueInstrumentMeterImpl) Labels(kvs ...core.KeyValue) metric.LabelSet {
	return u.impl.Labels(kvs...)
}

func (u *uniqueInstrumentMeterImpl) RecordBatch(ctx context.Context, labels metric.LabelSet, ms ...metric.Measurement) {
	u.impl.RecordBatch(ctx, labels, ms...)
}

func keyOf(descriptor metric.Descriptor) mapKey {
	return mapKey{
		descriptor.Name(),
		descriptor.LibraryName(),
	}
}

func (u *uniqueInstrumentMeterImpl) checkUniqueness(descriptor metric.Descriptor) (metric.InstrumentImpl, error) {
	impl, ok := u.state[keyOf(descriptor)]
	if !ok {
		return nil, nil
	}

	if impl.Descriptor().MetricKind() != descriptor.MetricKind() ||
		impl.Descriptor().NumberKind() != descriptor.NumberKind() {
		return nil, fmt.Errorf("Metric was %s registered as a %s %s: %w",
			impl.Descriptor().Name(),
			impl.Descriptor().NumberKind(),
			impl.Descriptor().MetricKind(),
			ErrMetricTypeMismatch,
		)
	}

	return impl, nil
}

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
	u.state[keyOf(descriptor)] = syncInst
	return syncInst, nil
}

func (u *uniqueInstrumentMeterImpl) NewAsyncInstrument(
	descriptor metric.Descriptor,
	callback func(func(core.Number, metric.LabelSet)),
) (metric.AsyncImpl, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	impl, err := u.checkUniqueness(descriptor)

	if err != nil {
		return nil, err
	} else if impl != nil {
		return impl.(metric.AsyncImpl), nil
	}

	asyncInst, err := u.impl.NewAsyncInstrument(descriptor, callback)
	if err != nil {
		return nil, err
	}
	u.state[keyOf(descriptor)] = asyncInst
	return asyncInst, nil
}
