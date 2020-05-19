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

package metric // import "go.opentelemetry.io/otel/sdk/export/metric"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Integrator is responsible for deciding which kind of aggregation to
// use (via AggregationSelector), gathering exported results from the
// SDK during collection, and deciding over which dimensions to group
// the exported data.
//
// The SDK supports binding only one of these interfaces, as it has
// the sole responsibility of determining which Aggregator to use for
// each record.
//
// The embedded AggregationSelector interface is called (concurrently)
// in instrumentation context to select the appropriate Aggregator for
// an instrument.
//
// The `Process` method is called during collection in a
// single-threaded context from the SDK, after the aggregator is
// checkpointed, allowing the integrator to build the set of metrics
// currently being exported.
type Integrator interface {
	// AggregationSelector is responsible for selecting the
	// concrete type of Aggregator used for a metric in the SDK.
	//
	// This may be a static decision based on fields of the
	// Descriptor, or it could use an external configuration
	// source to customize the treatment of each metric
	// instrument.
	//
	// The result from AggregatorSelector.AggregatorFor should be
	// the same type for a given Descriptor or else nil.  The same
	// type should be returned for a given descriptor, because
	// Aggregators only know how to Merge with their own type.  If
	// the result is nil, the metric instrument will be disabled.
	//
	// Note that the SDK only calls AggregatorFor when new records
	// require an Aggregator. This does not provide a way to
	// disable metrics with active records.
	AggregationSelector

	// Process is called by the SDK once per internal record,
	// passing the export Record (a Descriptor, the corresponding
	// Labels, and the checkpointed Aggregator).
	//
	// The Context argument originates from the controller that
	// orchestrates collection.
	Process(ctx context.Context, record Record) error
}

// AggregationSelector supports selecting the kind of Aggregator to
// use at runtime for a specific metric instrument.
type AggregationSelector interface {
	// AggregatorFor returns the kind of aggregator suited to the
	// requested export.  Returning `nil` indicates to ignore this
	// metric instrument.  This must return a consistent type to
	// avoid confusion in later stages of the metrics export
	// process, i.e., when Merging multiple aggregators for a
	// specific instrument.
	//
	// Note: This is context-free because the aggregator should
	// not relate to the incoming context.  This call should not
	// block.
	AggregatorFor(*metric.Descriptor) Aggregator
}

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
	// Update receives a new measured value and incorporates it
	// into the aggregation.  Update() calls may arrive
	// concurrently as the SDK does not provide synchronization.
	//
	// Descriptor.NumberKind() should be consulted to determine
	// whether the provided number is an int64 or float64.
	//
	// The Context argument comes from user-level code and could be
	// inspected for distributed or span context.
	Update(context.Context, metric.Number, *metric.Descriptor) error

	// Checkpoint is called during collection to finish one period
	// of aggregation by atomically saving the current value.
	// Checkpoint() is called concurrently with Update().
	// Checkpoint should reset the current state to the empty
	// state, in order to begin computing a new delta for the next
	// collection period.
	//
	// After the checkpoint is taken, the current value may be
	// accessed using by converting to one a suitable interface
	// types in the `aggregator` sub-package.
	//
	// The Context argument originates from the controller that
	// orchestrates collection.
	Checkpoint(context.Context, *metric.Descriptor)

	// Merge combines the checkpointed state from the argument
	// aggregator into this aggregator's checkpointed state.
	// Merge() is called in a single-threaded context, no locking
	// is required.
	Merge(Aggregator, *metric.Descriptor) error
}

// Exporter handles presentation of the checkpoint of aggregate
// metrics.  This is the final stage of a metrics export pipeline,
// where metric data are formatted for a specific system.
type Exporter interface {
	// Export is called immediately after completing a collection
	// pass in the SDK.
	//
	// The Context comes from the controller that initiated
	// collection.
	//
	// The CheckpointSet interface refers to the Integrator that just
	// completed collection.
	Export(context.Context, CheckpointSet) error
}

// CheckpointSet allows a controller to access a complete checkpoint of
// aggregated metrics from the Integrator.  This is passed to the
// Exporter which may then use ForEach to iterate over the collection
// of aggregated metrics.
type CheckpointSet interface {
	// ForEach iterates over aggregated checkpoints for all
	// metrics that were updated during the last collection
	// period. Each aggregated checkpoint returned by the
	// function parameter may return an error.
	// ForEach tolerates ErrNoData silently, as this is
	// expected from the Meter implementation. Any other kind
	// of error will immediately halt ForEach and return
	// the error to the caller.
	ForEach(func(Record) error) error

	// Locker supports locking the checkpoint set.  Collection
	// into the checkpoint set cannot take place (in case of a
	// stateful integrator) while it is locked.
	//
	// The Integrator attached to the Accumulator MUST be called
	// with the lock held.
	sync.Locker

	// RLock acquires a read lock corresponding to this Locker.
	RLock()
	// RUnlock releases a read lock corresponding to this Locker.
	RUnlock()
}

// Record contains the exported data for a single metric instrument
// and label set.
type Record struct {
	descriptor *metric.Descriptor
	labels     *label.Set
	resource   *resource.Resource
	aggregator Aggregator
}

// NewRecord allows Integrator implementations to construct export
// records.  The Descriptor, Labels, and Aggregator represent
// aggregate metric events received over a single collection period.
func NewRecord(descriptor *metric.Descriptor, labels *label.Set, resource *resource.Resource, aggregator Aggregator) Record {
	return Record{
		descriptor: descriptor,
		labels:     labels,
		resource:   resource,
		aggregator: aggregator,
	}
}

// Aggregator returns the checkpointed aggregator. It is safe to
// access the checkpointed state without locking.
func (r Record) Aggregator() Aggregator {
	return r.aggregator
}

// Descriptor describes the metric instrument being exported.
func (r Record) Descriptor() *metric.Descriptor {
	return r.descriptor
}

// Labels describes the labels associated with the instrument and the
// aggregated data.
func (r Record) Labels() *label.Set {
	return r.labels
}

// Resource contains common attributes that apply to this metric event.
func (r Record) Resource() *resource.Resource {
	return r.resource
}
