// Copyright 2019, OpenTelemetry Authors
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

package metric

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/api/core"
)

type synchronousInstrument struct {
	instrument SynchronousImpl
}

type synchronousBoundInstrument struct {
	boundInstrument BoundSynchronousImpl
}

type asynchronousInstrument struct {
	instrument AsynchronousImpl
}

var ErrSDKReturnedNilImpl = errors.New("SDK returned a nil implementation")

func (s synchronousInstrument) bind(labels LabelSet) synchronousBoundInstrument {
	return newSynchronousBoundInstrument(s.instrument.Bind(labels))
}

func (s synchronousInstrument) float64Measurement(value float64) Measurement {
	return newMeasurement(s.instrument, core.NewFloat64Number(value))
}

func (s synchronousInstrument) int64Measurement(value int64) Measurement {
	return newMeasurement(s.instrument, core.NewInt64Number(value))
}

func (s synchronousInstrument) directRecord(ctx context.Context, number core.Number, labels LabelSet) {
	s.instrument.RecordOne(ctx, number, labels)
}

func (s synchronousInstrument) SynchronousImpl() SynchronousImpl {
	return s.instrument
}

func (h synchronousBoundInstrument) directRecord(ctx context.Context, number core.Number) {
	h.boundInstrument.RecordOne(ctx, number)
}

func (h synchronousBoundInstrument) Unbind() {
	h.boundInstrument.Unbind()
}

func (a asynchronousInstrument) AsynchronousImpl() AsynchronousImpl {
	return a.instrument
}

func checkSynchronous(instrument SynchronousImpl, err error) (synchronousInstrument, error) {
	if instrument == nil {
		if err == nil {
			err = ErrSDKReturnedNilImpl
		}
		// Note: an alternate behavior would be to synthesize a new name
		// or group all duplicately-named instruments of a certain type
		// together and use a tag for the original name, e.g.,
		//   name = 'invalid.counter.int64'
		//   label = 'original-name=duplicate-counter-name'
		instrument = noopSynchronous{}
	}
	return synchronousInstrument{
		instrument: instrument,
	}, err
}

func newSynchronousBoundInstrument(boundInstrument BoundSynchronousImpl) synchronousBoundInstrument {
	return synchronousBoundInstrument{
		boundInstrument: boundInstrument,
	}
}

func newMeasurement(instrument SynchronousImpl, number core.Number) Measurement {
	return Measurement{
		instrument: instrument,
		number:     number,
	}
}

func checkAsynchronous(instrument AsynchronousImpl, err error) (asynchronousInstrument, error) {
	if instrument == nil {
		if err == nil {
			err = ErrSDKReturnedNilImpl
		}
		instrument = noopAsynchronous{}
	}
	return asynchronousInstrument{
		instrument: instrument,
	}, err
}
