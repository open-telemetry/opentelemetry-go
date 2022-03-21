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

package metrictest // import "go.opentelemetry.io/otel/sdk/metric/metrictest"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (

	// Library is the same as "sdk/instrumentation".Library but there is
	// a package cycle to use it.
	Library struct {
		InstrumentationName    string
		InstrumentationVersion string
		SchemaURL              string
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Ctx          context.Context
		Labels       []attribute.KeyValue
		Library      Library
	}

	Measurement struct {
		// Number needs to be aligned for 64-bit atomic operations.
		Number     number.Number
		Instrument sdkapi.InstrumentImpl
	}
)

// NewDescriptor is a test helper for constructing test metric
// descriptors using standard options.
func NewDescriptor(name string, ikind sdkapi.InstrumentKind, nkind number.Kind, opts ...instrument.Option) sdkapi.Descriptor {
	cfg := instrument.NewConfig(opts...)
	return sdkapi.NewDescriptor(name, ikind, nkind, cfg.Description(), cfg.Unit())
}
