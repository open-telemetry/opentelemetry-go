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
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Processor is responsible for deciding which kind of aggregation to
// use (via AggregatorSelector), gathering exported results from the
// SDK during collection, and deciding over which dimensions to group
// the exported data.
//
// The SDK supports binding only one of these interfaces, as it has
// the sole responsibility of determining which Aggregator to use for
// each record.
//
// The embedded AggregatorSelector interface is called (concurrently)
// in instrumentation context to select the appropriate Aggregator for
// an instrument.
//
// The `Process` method is called during collection in a
// single-threaded context from the SDK, after the aggregator is
// checkpointed, allowing the processor to build the set of metrics
// currently being exported.
type Processor interface {
	// AggregatorSelector is responsible for selecting the
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
	AggregatorSelector

	// Process is called by the SDK once per internal record,
	// passing the export Accumulation (a Descriptor, the corresponding
	// Labels, and the checkpointed Aggregator). This call has no
	// Context argument because it is expected to perform only
	// computation. An SDK is not expected to call exporters from
	// with Process, use a controller for that (see
	// ./controllers/{pull,push}.
	Process(accum Accumulation) error
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
	AggregatorFor(descriptor *sdkapi.Descriptor, aggregator ...*Aggregator)
}

// Checkpointer is the interface used by a Controller to coordinate
// the Processor with Accumulator(s) and Exporter(s).  The
// StartCollection() and FinishCollection() methods start and finish a
// collection interval.  Controllers call the Accumulator(s) during
// collection to process Accumulations.
type Checkpointer interface {
	// Processor processes metric data for export.  The Process
	// method is bracketed by StartCollection and FinishCollection
	// calls.  The embedded AggregatorSelector can be called at
	// any time.
	Processor

	// Reader returns the current data set.  This may be
	// called before and after collection.  The
	// implementation is required to return the same value
	// throughout its lifetime, since Reader exposes a
	// sync.Locker interface.  The caller is responsible for
	// locking the Reader before initiating collection.
	Reader() Reader

	// StartCollection begins a collection interval.
	StartCollection()

	// FinishCollection ends a collection interval.
	FinishCollection() error
}

// CheckpointerFactory is an interface for producing configured
// Checkpointer instances.
type CheckpointerFactory interface {
	NewCheckpointer() Checkpointer
}

// Aggregator implements a specific aggregation behavior, e.g., a
// behavior to track a sequence of updates to an instrument.  Counter
// instruments commonly use a simple Sum aggregator, but for the
// distribution instruments (Histogram, GaugeObserver) there are a
// number of possible aggregators with different cost and accuracy
// tradeoffs.
//
// Note that any Aggregator may be attached to any instrument--this is
// the result of the OpenTelemetry API/SDK separation.  It is possible
// to attach a Sum aggregator to a Histogram instrument.
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
	Update(ctx context.Context, number number.Number, descriptor *sdkapi.Descriptor) error

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
	SynchronizedMove(destination Aggregator, descriptor *sdkapi.Descriptor) error

	// Merge combines the checkpointed state from the argument
	// Aggregator into this Aggregator.  Merge is not synchronized
	// with respect to Update or SynchronizedMove.
	//
	// The owner of an Aggregator being merged is responsible for
	// synchronization of both Aggregator states.
	Merge(aggregator Aggregator, descriptor *sdkapi.Descriptor) error
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
	// The InstrumentationLibraryReader interface refers to the
	// Processor that just completed collection.
	Export(ctx context.Context, resource *resource.Resource, reader InstrumentationLibraryReader) error

	// TemporalitySelector is an interface used by the Processor
	// in deciding whether to compute Delta or Cumulative
	// Aggregations when passing Records to this Exporter.
	aggregation.TemporalitySelector
}

// InstrumentationLibraryReader is an interface for exporters to iterate
// over one instrumentation library of metric data at a time.
type InstrumentationLibraryReader interface {
	// ForEach calls the passed function once per instrumentation library,
	// allowing the caller to emit metrics grouped by the library that
	// produced them.
	ForEach(readerFunc func(instrumentation.Library, Reader) error) error
}

// Reader allows a controller to access a complete checkpoint of
// aggregated metrics from the Processor for a single library of
// metric data.  This is passed to the Exporter which may then use
// ForEach to iterate over the collection of aggregated metrics.
type Reader interface {
	// ForEach iterates over aggregated checkpoints for all
	// metrics that were updated during the last collection
	// period. Each aggregated checkpoint returned by the
	// function parameter may return an error.
	//
	// The TemporalitySelector argument is used to determine
	// whether the Record is computed using Delta or Cumulative
	// aggregation.
	//
	// ForEach tolerates ErrNoData silently, as this is
	// expected from the Meter implementation. Any other kind
	// of error will immediately halt ForEach and return
	// the error to the caller.
	ForEach(tempSelector aggregation.TemporalitySelector, recordFunc func(Record) error) error

	// Locker supports locking the checkpoint set.  Collection
	// into the checkpoint set cannot take place (in case of a
	// stateful processor) while it is locked.
	//
	// The Processor attached to the Accumulator MUST be called
	// with the lock held.
	sync.Locker

	// RLock acquires a read lock corresponding to this Locker.
	RLock()
	// RUnlock releases a read lock corresponding to this Locker.
	RUnlock()
}

// Metadata contains the common elements for exported metric data that
// are shared by the Accumulator->Processor and Processor->Exporter
// steps.
type Metadata struct {
	descriptor *sdkapi.Descriptor
	labels     *attribute.Set
}

// Accumulation contains the exported data for a single metric instrument
// and label set, as prepared by an Accumulator for the Processor.
type Accumulation struct {
	Metadata
	aggregator Aggregator
}

// Record contains the exported data for a single metric instrument
// and label set, as prepared by the Processor for the Exporter.
// This includes the effective start and end time for the aggregation.
type Record struct {
	Metadata
	aggregation aggregation.Aggregation
	start       time.Time
	end         time.Time
}

// Descriptor describes the metric instrument being exported.
func (m Metadata) Descriptor() *sdkapi.Descriptor {
	return m.descriptor
}

// Labels describes the labels associated with the instrument and the
// aggregated data.
func (m Metadata) Labels() *attribute.Set {
	return m.labels
}

// NewAccumulation allows Accumulator implementations to construct new
// Accumulations to send to Processors. The Descriptor, Labels,
// and Aggregator represent aggregate metric events received over a single
// collection period.
func NewAccumulation(descriptor *sdkapi.Descriptor, labels *attribute.Set, aggregator Aggregator) Accumulation {
	return Accumulation{
		Metadata: Metadata{
			descriptor: descriptor,
			labels:     labels,
		},
		aggregator: aggregator,
	}
}

// Aggregator returns the checkpointed aggregator. It is safe to
// access the checkpointed state without locking.
func (r Accumulation) Aggregator() Aggregator {
	return r.aggregator
}

// NewRecord allows Processor implementations to construct export
// records.  The Descriptor, Labels, and Aggregator represent
// aggregate metric events received over a single collection period.
func NewRecord(descriptor *sdkapi.Descriptor, labels *attribute.Set, aggregation aggregation.Aggregation, start, end time.Time) Record {
	return Record{
		Metadata: Metadata{
			descriptor: descriptor,
			labels:     labels,
		},
		aggregation: aggregation,
		start:       start,
		end:         end,
	}
}

// Aggregation returns the aggregation, an interface to the record and
// its aggregator, dependent on the kind of both the input and exporter.
func (r Record) Aggregation() aggregation.Aggregation {
	return r.aggregation
}

// StartTime is the start time of the interval covered by this aggregation.
func (r Record) StartTime() time.Time {
	return r.start
}

// EndTime is the end time of the interval covered by this aggregation.
func (r Record) EndTime() time.Time {
	return r.end
}
