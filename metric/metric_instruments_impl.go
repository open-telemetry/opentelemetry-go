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

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

type asyncInstrumentImpl struct {
	instrument sdkapi.AsyncImpl
}

// AsyncImpl returns the implementation object for asynchronous instruments.
func (a asyncInstrumentImpl) AsyncImpl() sdkapi.AsyncImpl {
	return a.instrument
}

type syncInstrumentImpl struct {
	instrument sdkapi.SyncImpl
}

// SyncImpl returns the implementation object for synchronous instruments.
func (s syncInstrumentImpl) SyncImpl() sdkapi.SyncImpl {
	return s.instrument
}

func (s syncInstrumentImpl) float64Measurement(value float64) Measurement {
	return sdkapi.NewMeasurement(s.instrument, number.NewFloat64Number(value))
}

func (s syncInstrumentImpl) int64Measurement(value int64) Measurement {
	return sdkapi.NewMeasurement(s.instrument, number.NewInt64Number(value))
}

func (s syncInstrumentImpl) directRecord(ctx context.Context, number number.Number, labels []attribute.KeyValue) {
	s.instrument.RecordOne(ctx, number, labels)
}

type float64CounterImpl struct {
	syncInstrumentImpl
}

// wrapFloat64CounterInstrument converts a SyncImpl into Float64Counter.
func wrapFloat64CounterInstrument(syncInst sdkapi.SyncImpl, err error) (Float64Counter, error) {
	common, err := checkNewSync(syncInst, err)
	return float64CounterImpl{syncInstrumentImpl: common}, err
}

func (c float64CounterImpl) Measurement(value float64) Measurement {
	return c.float64Measurement(value)
}

func (c float64CounterImpl) Add(ctx context.Context, value float64, labels ...attribute.KeyValue) {
	c.directRecord(ctx, number.NewFloat64Number(value), labels)
}

type int64CounterImpl struct {
	syncInstrumentImpl
}

// wrapInt64CounterInstrument converts a SyncImpl into Int64Counter.
func wrapInt64CounterInstrument(syncInst sdkapi.SyncImpl, err error) (Int64Counter, error) {
	common, err := checkNewSync(syncInst, err)
	return int64CounterImpl{syncInstrumentImpl: common}, err
}

func (c int64CounterImpl) Measurement(value int64) Measurement {
	return c.int64Measurement(value)
}

func (c int64CounterImpl) Add(ctx context.Context, value int64, labels ...attribute.KeyValue) {
	c.directRecord(ctx, number.NewInt64Number(value), labels)
}

type float64UpDownCounterImpl struct {
	syncInstrumentImpl
}

// wrapFloat64UpDownCounterInstrument converts a SyncImpl into Float64UpDownCounter.
func wrapFloat64UpDownCounterInstrument(syncInst sdkapi.SyncImpl, err error) (Float64UpDownCounter, error) {
	common, err := checkNewSync(syncInst, err)
	return float64UpDownCounterImpl{syncInstrumentImpl: common}, err
}

func (c float64UpDownCounterImpl) Measurement(value float64) Measurement {
	return c.float64Measurement(value)
}

func (c float64UpDownCounterImpl) Add(ctx context.Context, value float64, labels ...attribute.KeyValue) {
	c.directRecord(ctx, number.NewFloat64Number(value), labels)
}

type int64UpDownCounterImpl struct {
	syncInstrumentImpl
}

// wrapInt64UpDownCounterInstrument converts a SyncImpl into Int64UpDownCounter.
func wrapInt64UpDownCounterInstrument(syncInst sdkapi.SyncImpl, err error) (Int64UpDownCounter, error) {
	common, err := checkNewSync(syncInst, err)
	return int64UpDownCounterImpl{syncInstrumentImpl: common}, err
}

func (c int64UpDownCounterImpl) Measurement(value int64) Measurement {
	return c.int64Measurement(value)
}

func (c int64UpDownCounterImpl) Add(ctx context.Context, value int64, labels ...attribute.KeyValue) {
	c.directRecord(ctx, number.NewInt64Number(value), labels)
}

type float64HistogramImpl struct {
	syncInstrumentImpl
}

// wrapFloat64HistogramInstrument converts a SyncImpl into Float64Histogram.
func wrapFloat64HistogramInstrument(syncInst sdkapi.SyncImpl, err error) (Float64Histogram, error) {
	common, err := checkNewSync(syncInst, err)
	return float64HistogramImpl{syncInstrumentImpl: common}, err
}

func (c float64HistogramImpl) Measurement(value float64) Measurement {
	return c.float64Measurement(value)
}

func (c float64HistogramImpl) Record(ctx context.Context, value float64, labels ...attribute.KeyValue) {
	c.directRecord(ctx, number.NewFloat64Number(value), labels)
}

type int64HistogramImpl struct {
	syncInstrumentImpl
}

// wrapInt64HistogramInstrument converts a SyncImpl into Int64Histogram.
func wrapInt64HistogramInstrument(syncInst sdkapi.SyncImpl, err error) (Int64Histogram, error) {
	common, err := checkNewSync(syncInst, err)
	return int64HistogramImpl{syncInstrumentImpl: common}, err
}

func (c int64HistogramImpl) Measurement(value int64) Measurement {
	return c.int64Measurement(value)
}

func (c int64HistogramImpl) Record(ctx context.Context, value int64, labels ...attribute.KeyValue) {
	c.directRecord(ctx, number.NewInt64Number(value), labels)
}

// checkNewSync receives an SyncImpl and potential
// error, and returns the same types, checking for and ensuring that
// the returned interface is not nil.
func checkNewSync(instrument sdkapi.SyncImpl, err error) (syncInstrumentImpl, error) {
	if instrument == nil {
		if err == nil {
			err = ErrSDKReturnedNilImpl
		}
		// Note: an alternate behavior would be to synthesize a new name
		// or group all duplicately-named instruments of a certain type
		// together and use a tag for the original name, e.g.,
		//   name = 'invalid.counter.int64'
		//   label = 'original-name=duplicate-counter-name'
		instrument = sdkapi.NewNoopSyncInstrument()
	}
	return syncInstrumentImpl{
		instrument: instrument,
	}, err
}

type float64GaugeObserverImpl struct {
	asyncInstrumentImpl
}

// wrapFloat64GaugeObserverInstrument converts an AsyncImpl into Float64GaugeObserver.
func wrapFloat64GaugeObserverInstrument(asyncInst sdkapi.AsyncImpl, err error) (Float64GaugeObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return float64GaugeObserverImpl{asyncInstrumentImpl: common}, err
}

func (f float64GaugeObserverImpl) Observation(v float64) Observation {
	return sdkapi.NewObservation(f.instrument, number.NewFloat64Number(v))
}

type int64GaugeObserverImpl struct {
	asyncInstrumentImpl
}

// wrapInt64GaugeObserverInstrument converts an AsyncImpl into Int64GaugeObserver.
func wrapInt64GaugeObserverInstrument(asyncInst sdkapi.AsyncImpl, err error) (Int64GaugeObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return int64GaugeObserverImpl{asyncInstrumentImpl: common}, err
}

func (i int64GaugeObserverImpl) Observation(v int64) Observation {
	return sdkapi.NewObservation(i.instrument, number.NewInt64Number(v))
}

type float64CounterObserverImpl struct {
	asyncInstrumentImpl
}

// wrapFloat64CounterObserverInstrument converts an AsyncImpl into Float64CounterObserver.
func wrapFloat64CounterObserverInstrument(asyncInst sdkapi.AsyncImpl, err error) (Float64CounterObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return float64CounterObserverImpl{asyncInstrumentImpl: common}, err
}

func (f float64CounterObserverImpl) Observation(v float64) Observation {
	return sdkapi.NewObservation(f.instrument, number.NewFloat64Number(v))
}

type int64CounterObserverImpl struct {
	asyncInstrumentImpl
}

// wrapInt64CounterObserverInstrument converts an AsyncImpl into Int64CounterObserver.
func wrapInt64CounterObserverInstrument(asyncInst sdkapi.AsyncImpl, err error) (Int64CounterObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return int64CounterObserverImpl{asyncInstrumentImpl: common}, err
}

func (i int64CounterObserverImpl) Observation(v int64) Observation {
	return sdkapi.NewObservation(i.instrument, number.NewInt64Number(v))
}

type float64UpDownCounterObserverImpl struct {
	asyncInstrumentImpl
}

// wrapFloat64UpDownCounterObserverInstrument converts an AsyncImpl into Float64UpDownCounterObserver.
func wrapFloat64UpDownCounterObserverInstrument(asyncInst sdkapi.AsyncImpl, err error) (Float64UpDownCounterObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return float64UpDownCounterObserverImpl{asyncInstrumentImpl: common}, err
}

func (f float64UpDownCounterObserverImpl) Observation(v float64) Observation {
	return sdkapi.NewObservation(f.instrument, number.NewFloat64Number(v))
}

type int64UpDownCounterObserverImpl struct {
	asyncInstrumentImpl
}

// wrapInt64UpDownCounterObserverInstrument converts an AsyncImpl into Int64UpDownCounterObserver.
func wrapInt64UpDownCounterObserverInstrument(asyncInst sdkapi.AsyncImpl, err error) (Int64UpDownCounterObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return int64UpDownCounterObserverImpl{asyncInstrumentImpl: common}, err
}

func (i int64UpDownCounterObserverImpl) Observation(v int64) Observation {
	return sdkapi.NewObservation(i.instrument, number.NewInt64Number(v))
}

// checkNewAsync receives an AsyncImpl and potential
// error, and returns the same types, checking for and ensuring that
// the returned interface is not nil.
func checkNewAsync(instrument sdkapi.AsyncImpl, err error) (asyncInstrumentImpl, error) {
	if instrument == nil {
		if err == nil {
			err = ErrSDKReturnedNilImpl
		}
		instrument = sdkapi.NewNoopAsyncInstrument()
	}
	return asyncInstrumentImpl{
		instrument: instrument,
	}, err
}
