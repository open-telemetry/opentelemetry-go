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

package registry // import "go.opentelemetry.io/otel/sdk/metric/registry"

import (
	"context"
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

func (u *uniqueInstrumentMeterImpl) uniqCheck(desc metric.Descriptor) (error, mapKey, metric.InstrumentImpl) {
	key := mapKey{
		desc.Name(),
		desc.LibraryName(),
	}
	// TODO: Finish this. @@@
	return nil, key, nil
}

func (u *uniqueInstrumentMeterImpl) NewSynchronousInstrument(descriptor metric.Descriptor) (metric.SynchronousImpl, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	// TODO: Finish this.  This is not really implemented.
	err, key, value := u.uniqCheck(descriptor)

	if err != nil {
		return nil, err
	}

	u.state[key] = value

	return u.impl.NewSynchronousInstrument(descriptor)
}

func (u *uniqueInstrumentMeterImpl) NewAsynchronousInstrument(
	descriptor metric.Descriptor,
	callback func(func(core.Number, metric.LabelSet)),
) (metric.AsynchronousImpl, error) {
	return u.impl.NewAsynchronousInstrument(descriptor, callback)
}
