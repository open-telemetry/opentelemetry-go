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

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

func NewNoopMeterProvider() MeterProvider {
	return noopMeterProvider{}
}

type noopMeterProvider struct{}

var _ MeterProvider = noopMeterProvider{}

func (noopMeterProvider) Meter(instrumentationName string, opts ...MeterOption) Meter {
	return noopMeter{}
}

func NewNoopMeter() Meter {
	return noopMeter{}
}

type noopMeter struct{}

var _ Meter = noopMeter{}

func (noopMeter) AsyncInt64() asyncint64.Instruments {
	return asyncint64.NewNoopInstruments()
}
func (noopMeter) AsyncFloat64() asyncfloat64.Instruments {
	return asyncfloat64.NewNoopInstruments()
}
func (noopMeter) SyncInt64() syncint64.Instruments {
	return syncint64.NewNoopInstruments()
}
func (noopMeter) SyncFloat64() syncfloat64.Instruments {
	return syncfloat64.NewNoopInstruments()
}
func (noopMeter) RegisterCallback([]instrument.Asynchronous, func(context.Context)) error {
	return nil
}
