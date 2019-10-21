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
	"sync/atomic"

	"go.opentelemetry.io/api/core"
)

var global atomic.Value

// GlobalMeter returns a meter registered as a global meter. If no
// meter is registered then an instance of noop Meter is returned.
func GlobalMeter() Meter {
	if t := global.Load(); t != nil {
		return t.(Meter)
	}
	return noopMeter{}
}

// SetGlobalMeter sets provided meter as a global meter.
func SetGlobalMeter(t Meter) {
	global.Store(t)
}

// Labels gets a LabelSet from the global Meter.
func Labels(ctx context.Context, kv ...core.KeyValue) LabelSet {
	return GlobalMeter().Labels(ctx, kv...)
}

// NewInt64Counter creates a new integral counter with the global
// Meter.
func NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	return GlobalMeter().NewInt64Counter(name, cos...)
}

// NewFloat64Counter creates a new floating point counter with the
// global Meter.
func NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	return GlobalMeter().NewFloat64Counter(name, cos...)
}

// NewInt64Gauge creates a new integral gauge with the global Meter.
func NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	return GlobalMeter().NewInt64Gauge(name, gos...)
}

// NewFloat64Gauge creates a new floating point gauge with the global
// Meter.
func NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	return GlobalMeter().NewFloat64Gauge(name, gos...)
}

// NewInt64Measure creates a new integral measure with the global
// Meter.
func NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	return GlobalMeter().NewInt64Measure(name, mos...)
}

// NewFloat64Measure creates a new floating point measure with the
// global Meter.
func NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	return GlobalMeter().NewFloat64Measure(name, mos...)
}

// RecordBatch reports to the global Meter.
func RecordBatch(ctx context.Context, labels LabelSet, batch ...Measurement) {
	GlobalMeter().RecordBatch(ctx, labels, batch...)
}
