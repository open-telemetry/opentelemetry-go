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
//
// Deprecated: Use NewNoopMeterProvider().Meter() instead.
func NewNoopMeter() Meter {
	return noopMeter{}
}

type noopMeter struct{}

func (noopMeter) Float64Counter(string, ...InstrumentOption) (Float64Counter, error) {
	return noopFloat64Inst{}, nil
}

func (noopMeter) Float64UpDownCounter(string, ...InstrumentOption) (Float64UpDownCounter, error) {
	return noopFloat64Inst{}, nil
}

func (noopMeter) Float64Histogram(string, ...InstrumentOption) (Float64Histogram, error) {
	return noopFloat64Inst{}, nil
}

func (noopMeter) Float64ObservableCounter(string, ...ObservableOption) (Float64ObservableCounter, error) {
	return noopFloat64ObservableInst{}, nil
}

func (noopMeter) Float64ObservableUpDownCounter(string, ...ObservableOption) (Float64ObservableUpDownCounter, error) {
	return noopFloat64ObservableInst{}, nil
}

func (noopMeter) Float64ObservableGauge(string, ...ObservableOption) (Float64ObservableGauge, error) {
	return noopFloat64ObservableInst{}, nil
}

func (noopMeter) Int64Counter(string, ...InstrumentOption) (Int64Counter, error) {
	return noopInt64Inst{}, nil
}

func (noopMeter) Int64UpDownCounter(string, ...InstrumentOption) (Int64UpDownCounter, error) {
	return noopInt64Inst{}, nil
}

func (noopMeter) Int64Histogram(string, ...InstrumentOption) (Int64Histogram, error) {
	return noopInt64Inst{}, nil
}

func (noopMeter) Int64ObservableCounter(string, ...ObservableOption) (Int64ObservableCounter, error) {
	return noopInt64ObservableInst{}, nil
}

func (noopMeter) Int64ObservableUpDownCounter(string, ...ObservableOption) (Int64ObservableUpDownCounter, error) {
	return noopInt64ObservableInst{}, nil
}

func (noopMeter) Int64ObservableGauge(string, ...ObservableOption) (Int64ObservableGauge, error) {
	return noopInt64ObservableInst{}, nil
}

func (noopMeter) RegisterCallback(Callback, Observable, ...Observable) (Unregisterer, error) {
	return unregisterer{}, nil
}

type unregisterer struct{}

func (unregisterer) Unregister() error { return nil }

type noopFloat64ObservableInst struct {
	Observable
}

var (
	_ Float64ObservableCounter       = noopFloat64ObservableInst{}
	_ Float64ObservableUpDownCounter = noopFloat64ObservableInst{}
	_ Float64ObservableGauge         = noopFloat64ObservableInst{}
)

func (noopFloat64ObservableInst) Observe(context.Context, float64, ...attribute.KeyValue) {}

type noopInt64ObservableInst struct {
	Observable
}

var (
	_ Int64ObservableCounter       = noopInt64ObservableInst{}
	_ Int64ObservableUpDownCounter = noopInt64ObservableInst{}
	_ Int64ObservableGauge         = noopInt64ObservableInst{}
)

func (noopInt64ObservableInst) Observe(context.Context, int64, ...attribute.KeyValue) {}

type noopFloat64Inst struct{}

var (
	_ Float64Counter       = noopFloat64Inst{}
	_ Float64UpDownCounter = noopFloat64Inst{}
	_ Float64Histogram     = noopFloat64Inst{}
)

func (noopFloat64Inst) Add(context.Context, float64, ...attribute.KeyValue)    {}
func (noopFloat64Inst) Record(context.Context, float64, ...attribute.KeyValue) {}

type noopInt64Inst struct{}

var (
	_ Int64Counter       = noopInt64Inst{}
	_ Int64UpDownCounter = noopInt64Inst{}
	_ Int64Histogram     = noopInt64Inst{}
)

func (noopInt64Inst) Add(context.Context, int64, ...attribute.KeyValue)    {}
func (noopInt64Inst) Record(context.Context, int64, ...attribute.KeyValue) {}
