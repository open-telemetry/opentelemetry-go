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

	"go.opentelemetry.io/otel/metric/aggregation"
	"go.opentelemetry.io/otel/metric/number"
)

// Aggregator implements a specific aggregation behavior, e.g., a
// behavior to track a sequence of updates to an instrument.  Sum-only
// instruments commonly use a simple Sum aggregator, but for the
// distribution instruments (ValueRecorder, ValueObserver) there are a
// number of possible aggregators with different cost and accuracy
// tradeoffs.
//
// Note that any Aggregator may be attached to any instrument--this is
// the result of the OpenTelemetry API/SDK separation.  It is possible
// to attach a Sum aggregator to a ValueRecorder instrument or a
// MinMaxSumCount aggregator to a Counter instrument.
type Aggregator interface {
	// Aggregation returns an Aggregation interface to access the
	// current state of this Aggregator.  The caller is
	// responsible for synchronization and must not call any the
	// other methods in this interface concurrently while using
	// the Aggregation.
	Aggregation() aggregation.Aggregation

	// Update receives a new measured value and incorporates it
	// into the aggregation.  Update() calls may be called
	// concurrently.
	//
	// Descriptor.NumberKind() should be consulted to determine
	// whether the provided number is an int64 or float64.
	//
	// The Context argument comes from user-level code and could be
	// inspected for a `correlation.Map` or `trace.SpanContext`.
	Update(ctx context.Context, number number.Number, descriptor *Descriptor) error

	// SynchronizedMove is called during collection to finish one
	// period of aggregation by atomically saving the
	// currently-updating state into the argument Aggregator AND
	// resetting the current value to the zero state.
	//
	// SynchronizedMove() is called concurrently with Update().  These
	// two methods must be synchronized with respect to each
	// other, for correctness.
	//
	// After saving a synchronized copy, the Aggregator can be converted
	// into one or more of the interfaces in the `aggregation` sub-package,
	// according to kind of Aggregator that was selected.
	//
	// This method will return an InconsistentAggregatorError if
	// this Aggregator cannot be copied into the destination due
	// to an incompatible type.
	//
	// This call has no Context argument because it is expected to
	// perform only computation.
	//
	// When called with a nil `destination`, this Aggregator is reset
	// and the current value is discarded.
	SynchronizedMove(destination Aggregator, descriptor *Descriptor) error

	// Merge combines the checkpointed state from the argument
	// Aggregator into this Aggregator.  Merge is not synchronized
	// with respect to Update or SynchronizedMove.
	//
	// The owner of an Aggregator being merged is responsible for
	// synchronization of both Aggregator states.
	Merge(aggregator Aggregator, descriptor *Descriptor) error
}

// AggregatorSelector supports selecting the kind of Aggregator to
// use at runtime for a specific metric instrument.
type AggregatorSelector interface {
	// AggregatorFor allocates a variable number of aggregators of
	// a kind suitable for the requested export.  This method
	// initializes a `...*Aggregator`, to support making a single
	// allocation.
	//
	// When the call returns without initializing the *Aggregator
	// to a non-nil value, the metric instrument is explicitly
	// disabled.
	//
	// This must return a consistent type to avoid confusion in
	// later stages of the metrics export process, i.e., when
	// Merging or Checkpointing aggregators for a specific
	// instrument.
	//
	// Note: This is context-free because the aggregator should
	// not relate to the incoming context.  This call should not
	// block.
	AggregatorFor(descriptor *Descriptor, aggregator ...*Aggregator)
}
