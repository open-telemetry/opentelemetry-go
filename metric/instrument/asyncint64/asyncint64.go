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

// Package asyncint64 provides asynchronous int64 instrument types.
//
// Deprecated: Use go.opentelemetry.io/otel/metric instead.
package asyncint64 // import "go.opentelemetry.io/otel/metric/instrument/asyncint64"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument" //nolint:staticcheck  // Known deprecation.
)

// InstrumentProvider provides access to individual instruments.
//
// Warning: methods may be added to this interface in minor releases.
//
// Deprecated: use the go.opentelemetry.io/otel/metric.Meter methods instead.
type InstrumentProvider interface {
	// Counter creates an instrument for recording increasing values.
	Counter(name string, opts ...instrument.Option) (Counter, error) //nolint:staticcheck  // Known deprecation.

	// UpDownCounter creates an instrument for recording changes of a value.
	UpDownCounter(name string, opts ...instrument.Option) (UpDownCounter, error) //nolint:staticcheck  // Known deprecation.

	// Gauge creates an instrument for recording the current value.
	Gauge(name string, opts ...instrument.Option) (Gauge, error) //nolint:staticcheck  // Known deprecation.
}

// Counter is an instrument that records increasing values.
//
// Warning: methods may be added to this interface in minor releases.
//
// Deprecated: use go.opentelemetry.io/otel/metric.Int64ObservableCounter
// instead.
type Counter interface {
	// Observe records the state of the instrument to be x. Implementations
	// will assume x to be the cumulative sum of the count.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)

	instrument.Asynchronous //nolint:staticcheck  // Known deprecation.
}

// UpDownCounter is an instrument that records increasing or decreasing values.
//
// Warning: methods may be added to this interface in minor releases.
//
// Deprecated: use go.opentelemetry.io/otel/metric.Int64ObservableUpDownCounter
// instead.
type UpDownCounter interface {
	// Observe records the state of the instrument to be x. Implementations
	// will assume x to be the cumulative sum of the count.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)

	instrument.Asynchronous //nolint:staticcheck  // Known deprecation.
}

// Gauge is an instrument that records independent readings.
//
// Warning: methods may be added to this interface in minor releases.
//
// Deprecated: use go.opentelemetry.io/otel/metric.Int64ObservableGauge
// instead.
type Gauge interface {
	// Observe records the state of the instrument to be x.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x int64, attrs ...attribute.KeyValue)

	instrument.Asynchronous //nolint:staticcheck  // Known deprecation.
}
