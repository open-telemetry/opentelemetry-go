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

package asyncint64 // import "go.opentelemetry.io/otel/metric/instrument/asyncint64"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

// Observer is a recorder of measurement values.
//
// Warning: methods may be added to this interface in minor releases.
type Observer interface {
	instrument.Asynchronous

	// Observe records the measurement value for a set of attributes.
	//
	// It is only valid to call this within a callback. If called outside of
	// the registered callback it should have no effect on the instrument, and
	// an error will be reported via the error handler.
	Observe(ctx context.Context, value int64, attributes ...attribute.KeyValue)
}

// Callback is a function that returns observations for an Asynchronous
// instrument.
//
// The function needs to complete in a finite amount of time and the deadline
// of the passed context is expected to be honored.
//
// The function needs to be concurrent safe.
type Callback func(context.Context, Observer) error

// InstrumentProvider provides access to individual instruments.
//
// Warning: methods may be added to this interface in minor releases.
type InstrumentProvider interface {
	// Counter creates an instrument for recording increasing values.
	Counter(name string, opts ...Option) (Counter, error)

	// UpDownCounter creates an instrument for recording changes of a value.
	UpDownCounter(name string, opts ...Option) (UpDownCounter, error)

	// Gauge creates an instrument for recording the current value.
	Gauge(name string, opts ...Option) (Gauge, error)
}

// Counter is an instrument used to asynchronously record increasing int64
// measurements once per a measurement collection cycle. The Observe method is
// used to record the measured state of the instrument when it is called.
// Implementations will assume the observed value to be the cumulative sum of
// the count.
//
// Warning: methods may be added to this interface in minor releases.
type Counter interface{ Observer }

// UpDownCounter is an instrument used to asynchronously record int64
// measurements once per a measurement collection cycle. The Observe method is
// used to record the measured state of the instrument when it is called.
// Implementations will assume the observed value to be the cumulative sum of
// the count.
//
// Warning: methods may be added to this interface in minor releases.
type UpDownCounter interface{ Observer }

// Gauge is an instrument used to asynchronously record instantaneous int64
// measurements once per a measurement collection cycle.
//
// Warning: methods may be added to this interface in minor releases.
type Gauge interface{ Observer }
