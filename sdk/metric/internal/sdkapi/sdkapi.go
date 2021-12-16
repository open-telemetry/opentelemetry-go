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

package sdkapi // import "go.opentelemetry.io/otel/sdk/metric/internal/sdkapi"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/number"
)

// MeterImpl is the interface an SDK must implement to supply a Meter
// implementation.
type MeterImpl interface {
	// NewInstrument returns a newly constructed instrument
	// implementation or an error, should one occur.
	NewInstrument(descriptor Descriptor) (Instrument, error)

	NewCallback(insts []Instrument, callback func(context.Context) error) (Callback, error)
}

type Callback interface {
	Instruments() []Instrument
}

// Instrument is a common interface for synchronous and
// asynchronous instruments.
type Instrument interface {
	// Descriptor returns a copy of the instrument's Descriptor.
	Descriptor() Descriptor

	// Capture captures a single metric event.
	Capture(ctx context.Context, number number.Number, attrs []attribute.KeyValue)
}
