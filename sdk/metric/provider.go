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
	"go.opentelemetry.io/otel/metric"
)

// MeterProvider handles the creation and coordination of Meters. All Meters
// created by a MeterProvider will be associated with the same Resource, have
// the same Views applied to them, and have their produced metric telemetry
// passed to the configured Readers.
type MeterProvider struct {
	// TODO (#2820): implement.
}

// Compile-time check MeterProvider implements metric.MeterProvider.
var _ metric.MeterProvider = (*MeterProvider)(nil)

// NewMeterProvider returns a new and configured MeterProvider.
//
// By default, the returned MeterProvider is configured with the default
// Resource and no Readers.
func NewMeterProvider(options ...Option) *MeterProvider {
	// TODO (#2820): implement.
	return &MeterProvider{}
}

// Meter returns a Meter with the given name and configured with options.
//
// The name should be the name of the instrumentation library creating
// telemetry. This name may be the same as the instrumented code only if that
// code provides built-in instrumentation.
//
// If name is empty, the default (go.opentelemetry.io/otel/sdk/meter) will be
// used.
//
// This method is safe to call concurrently.
func (mp *MeterProvider) Meter(name string, options ...metric.MeterOption) metric.Meter {
	// TODO (#2821): ensure this is concurrent safe.
	// TODO: test this is concurrent safe.
	// TODO (#2821): register and track the created Meter.
	return &meter{}
}
