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
	"go.opentelemetry.io/otel/metric/instrument"
)

// NewNoopMeterProvider creates a MeterProvider that does not record any metrics.
func NewNoopMeterProvider() MeterProvider {
	return noopMeterProvider{}
}

type noopMeterProvider struct{}

func (noopMeterProvider) Meter(string, ...MeterOption) Meter {
	return noopMeter{}
}

// NewNoopMeter creates a Meter that does not record any metrics.
func NewNoopMeter() Meter {
	return noopMeter{}
}

type noopMeter struct{}

// AsyncInt64 creates an instrument that does not record any metrics.
func (noopMeter) AsyncInt64() instrument.AsyncInstrumentProvider[int64] {
	return nonrecordingAsyncInt64Instrument{}
}

// AsyncFloat64 creates an instrument that does not record any metrics.
func (noopMeter) AsyncFloat64() instrument.AsyncInstrumentProvider[float64] {
	return nonrecordingAsyncFloat64Instrument{}
}

// SyncInt64 creates an instrument that does not record any metrics.
func (noopMeter) SyncInt64() instrument.SyncInstrumentProvider[int64] {
	return nonrecordingSyncInt64Instrument{}
}

// SyncFloat64 creates an instrument that does not record any metrics.
func (noopMeter) SyncFloat64() instrument.SyncInstrumentProvider[float64] {
	return nonrecordingSyncFloat64Instrument{}
}

// RegisterCallback creates a register callback that does not record any metrics.
func (noopMeter) RegisterCallback([]instrument.Asynchronous, func(context.Context)) error {
	return nil
}

type nonrecordingAsyncFloat64Instrument struct {
	instrument.Asynchronous
}

var (
	_ instrument.AsyncInstrumentProvider[float64] = nonrecordingAsyncFloat64Instrument{}
	_ instrument.AsyncCounter[float64]            = nonrecordingAsyncFloat64Instrument{}
	_ instrument.AsyncUpDownCounter[float64]      = nonrecordingAsyncFloat64Instrument{}
	_ instrument.AsyncGauge[float64]              = nonrecordingAsyncFloat64Instrument{}
)

func (n nonrecordingAsyncFloat64Instrument) Counter(string, ...instrument.Option) (instrument.AsyncCounter[float64], error) {
	return n, nil
}

func (n nonrecordingAsyncFloat64Instrument) UpDownCounter(string, ...instrument.Option) (instrument.AsyncUpDownCounter[float64], error) {
	return n, nil
}

func (n nonrecordingAsyncFloat64Instrument) Gauge(string, ...instrument.Option) (instrument.AsyncGauge[float64], error) {
	return n, nil
}

func (nonrecordingAsyncFloat64Instrument) Observe(context.Context, float64, ...attribute.KeyValue) {

}

type nonrecordingAsyncInt64Instrument struct {
	instrument.Asynchronous
}

var (
	_ instrument.AsyncInstrumentProvider[int64] = nonrecordingAsyncInt64Instrument{}
	_ instrument.AsyncCounter[int64]            = nonrecordingAsyncInt64Instrument{}
	_ instrument.AsyncUpDownCounter[int64]      = nonrecordingAsyncInt64Instrument{}
	_ instrument.AsyncGauge[int64]              = nonrecordingAsyncInt64Instrument{}
)

func (n nonrecordingAsyncInt64Instrument) Counter(string, ...instrument.Option) (instrument.AsyncCounter[int64], error) {
	return n, nil
}

func (n nonrecordingAsyncInt64Instrument) UpDownCounter(string, ...instrument.Option) (instrument.AsyncUpDownCounter[int64], error) {
	return n, nil
}

func (n nonrecordingAsyncInt64Instrument) Gauge(string, ...instrument.Option) (instrument.AsyncGauge[int64], error) {
	return n, nil
}

func (nonrecordingAsyncInt64Instrument) Observe(context.Context, int64, ...attribute.KeyValue) {
}

type nonrecordingSyncFloat64Instrument struct {
	instrument.Synchronous
}

var (
	_ instrument.SyncInstrumentProvider[float64] = nonrecordingSyncFloat64Instrument{}
	_ instrument.SyncCounter[float64]            = nonrecordingSyncFloat64Instrument{}
	_ instrument.SyncUpDownCounter[float64]      = nonrecordingSyncFloat64Instrument{}
	_ instrument.SyncHistogram[float64]          = nonrecordingSyncFloat64Instrument{}
)

func (n nonrecordingSyncFloat64Instrument) Counter(string, ...instrument.Option) (instrument.SyncCounter[float64], error) {
	return n, nil
}

func (n nonrecordingSyncFloat64Instrument) UpDownCounter(string, ...instrument.Option) (instrument.SyncUpDownCounter[float64], error) {
	return n, nil
}

func (n nonrecordingSyncFloat64Instrument) Histogram(string, ...instrument.Option) (instrument.SyncHistogram[float64], error) {
	return n, nil
}

func (nonrecordingSyncFloat64Instrument) Add(context.Context, float64, ...attribute.KeyValue) {

}

func (nonrecordingSyncFloat64Instrument) Record(context.Context, float64, ...attribute.KeyValue) {

}

type nonrecordingSyncInt64Instrument struct {
	instrument.Synchronous
}

var (
	_ instrument.SyncInstrumentProvider[int64] = nonrecordingSyncInt64Instrument{}
	_ instrument.SyncCounter[int64]            = nonrecordingSyncInt64Instrument{}
	_ instrument.SyncUpDownCounter[int64]      = nonrecordingSyncInt64Instrument{}
	_ instrument.SyncHistogram[int64]          = nonrecordingSyncInt64Instrument{}
)

func (n nonrecordingSyncInt64Instrument) Counter(string, ...instrument.Option) (instrument.SyncCounter[int64], error) {
	return n, nil
}

func (n nonrecordingSyncInt64Instrument) UpDownCounter(string, ...instrument.Option) (instrument.SyncUpDownCounter[int64], error) {
	return n, nil
}

func (n nonrecordingSyncInt64Instrument) Histogram(string, ...instrument.Option) (instrument.SyncHistogram[int64], error) {
	return n, nil
}

func (nonrecordingSyncInt64Instrument) Add(context.Context, int64, ...attribute.KeyValue) {
}
func (nonrecordingSyncInt64Instrument) Record(context.Context, int64, ...attribute.KeyValue) {
}
