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

package metric

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/api/kv"
)

// Measurement is used for reporting a batch of metric
// values. Instances of this type should be created by instruments
// (e.g., Int64Counter.Measurement()).
type Measurement struct {
	// number needs to be aligned for 64-bit atomic operations.
	number     Number
	instrument SyncImpl
}

// SyncImpl returns the instrument that created this measurement.
// This returns an implementation-level object for use by the SDK,
// users should not refer to this.
func (m Measurement) SyncImpl() SyncImpl {
	return m.instrument
}

// Number returns a number recorded in this measurement.
func (m Measurement) Number() Number {
	return m.number
}

// Observation is used for reporting a batch of metric
// values. Instances of this type should be created by Observer
// instruments (e.g., Int64Observer.Observation()).
type Observation struct {
	// number needs to be aligned for 64-bit atomic operations.
	number     Number
	instrument AsyncImpl
}

// AsyncImpl returns the instrument that created this observation.
// This returns an implementation-level object for use by the SDK,
// users should not refer to this.
func (m Observation) AsyncImpl() AsyncImpl {
	return m.instrument
}

// Number returns a number recorded in this observation.
func (m Observation) Number() Number {
	return m.number
}

type syncInstrument struct {
	instrument SyncImpl
}

type syncBoundInstrument struct {
	boundInstrument BoundSyncImpl
}

type asyncInstrument struct {
	instrument AsyncImpl
}

var ErrSDKReturnedNilImpl = errors.New("SDK returned a nil implementation")

func (s syncInstrument) bind(labels []kv.KeyValue) syncBoundInstrument {
	return newSyncBoundInstrument(s.instrument.Bind(labels))
}

func (s syncInstrument) float64Measurement(value float64) Measurement {
	return newMeasurement(s.instrument, NewFloat64Number(value))
}

func (s syncInstrument) int64Measurement(value int64) Measurement {
	return newMeasurement(s.instrument, NewInt64Number(value))
}

func (s syncInstrument) directRecord(ctx context.Context, number Number, labels []kv.KeyValue) {
	s.instrument.RecordOne(ctx, number, labels)
}

func (s syncInstrument) SyncImpl() SyncImpl {
	return s.instrument
}

func (h syncBoundInstrument) directRecord(ctx context.Context, number Number) {
	h.boundInstrument.RecordOne(ctx, number)
}

func (h syncBoundInstrument) Unbind() {
	h.boundInstrument.Unbind()
}

func (a asyncInstrument) AsyncImpl() AsyncImpl {
	return a.instrument
}

// checkNewSync receives an SyncImpl and potential
// error, and returns the same types, checking for and ensuring that
// the returned interface is not nil.
func checkNewSync(instrument SyncImpl, err error) (syncInstrument, error) {
	if instrument == nil {
		if err == nil {
			err = ErrSDKReturnedNilImpl
		}
		// Note: an alternate behavior would be to synthesize a new name
		// or group all duplicately-named instruments of a certain type
		// together and use a tag for the original name, e.g.,
		//   name = 'invalid.counter.int64'
		//   label = 'original-name=duplicate-counter-name'
		instrument = NoopSync{}
	}
	return syncInstrument{
		instrument: instrument,
	}, err
}

func newSyncBoundInstrument(boundInstrument BoundSyncImpl) syncBoundInstrument {
	return syncBoundInstrument{
		boundInstrument: boundInstrument,
	}
}

func newMeasurement(instrument SyncImpl, number Number) Measurement {
	return Measurement{
		instrument: instrument,
		number:     number,
	}
}

// checkNewAsync receives an AsyncImpl and potential
// error, and returns the same types, checking for and ensuring that
// the returned interface is not nil.
func checkNewAsync(instrument AsyncImpl, err error) (asyncInstrument, error) {
	if instrument == nil {
		if err == nil {
			err = ErrSDKReturnedNilImpl
		}
		instrument = NoopAsync{}
	}
	return asyncInstrument{
		instrument: instrument,
	}, err
}
