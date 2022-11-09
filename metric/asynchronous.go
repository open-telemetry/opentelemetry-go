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

// Observable is a kind of instrument that record measurements within a
// callback function. These instruments can created with callbacks to register
// for them, or they can be passed to the RegisterCallback method of the Meter
// that created them with a callback to register.
type Observable interface {
	observable()
}

// Callback is a function that records an observation for an Observable
// instrument by calling that instruments Observe method during its execution.
// The Callback needs to only record observations for the instruments it is
// registered with.
//
// The function needs to complete in a finite amount of time and the deadline
// of the passed context is expected to be honored.
//
// The function needs to be concurrent safe.
type Callback func(context.Context) error

// Float64ObservableCounter is an observable instrument that records increasing
// incremental float64 values within a callback.
//
// Warning: methods may be added to this interface in minor releases.
type Float64ObservableCounter interface {
	Observable

	// Observe records the state of the instrument to be x for attrs.
	// Implementations will assume x to be the cumulative sum of the count.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
}

// Int64ObservableCounter is an observable instrument that records increasing
// incremental int64 values within a callback.
//
// Warning: methods may be added to this interface in minor releases.
type Int64ObservableCounter interface {
	Observable

	// Observe records the state of the instrument to be x for attrs.
	// Implementations will assume x to be the cumulative sum of the count.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}

// Float64ObservableCounter is an observable instrument that records increasing
// or decreasing incremental float64 values within a callback.
//
// Warning: methods may be added to this interface in minor releases.
type Float64ObservableUpDownCounter interface {
	Observable

	// Observe records the state of the instrument to be x for attrs.
	// Implementations will assume x to be the cumulative sum of the count.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
}

// Int64ObservableUpDownCounter is an observable instrument that records
// increasing or decreasing incremental int64 values within a callback.
//
// Warning: methods may be added to this interface in minor releases.
type Int64ObservableUpDownCounter interface {
	Observable

	// Observe records the state of the instrument to be x for attrs.
	// Implementations will assume x to be the cumulative sum of the count.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}

// Float64ObservableGauge is an observable instrument that records independent
// float64 values within a callback.
//
// Warning: methods may be added to this interface in minor releases.
type Float64ObservableGauge interface {
	Observable

	// Observe records the state of the instrument to be x for attrs.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x float64, attrs ...attribute.KeyValue)
}

// Int64ObservableGauge is an observable instrument that records independent
// int64 values within a callback.
//
// Warning: methods may be added to this interface in minor releases.
type Int64ObservableGauge interface {
	Observable

	// Observe records the state of the instrument to be x for attrs.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)
}
