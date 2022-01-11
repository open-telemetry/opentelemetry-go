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

package export // import "go.opentelemetry.io/otel/sdk/metric/export"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
)

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
type Aggregator[N number.Any, Agg, Opt any] interface {
	Init(opts ...Opt)

	Update(number N)

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
	SynchronizedMove(destination *Agg)

	// Merge combines the checkpointed state from the argument
	// Aggregator into this Aggregator.  Merge is not synchronized
	// with respect to Update or SynchronizedMove.
	//
	// The owner of an Aggregator being merged is responsible for
	// synchronization of both Aggregator states.
	Merge(aggregator *Agg)
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
	ForEach(recordFunc func(Record) error) error
}

// Metadata contains the common elements for exported metric data that
// are shared by the Accumulator->Processor and Processor->Exporter
// steps.
type Metadata struct {
	descriptor *sdkapi.Descriptor
	attributes attribute.Attributes
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

// Attributes describes the attribtes associated with the instrument
// and the aggregated data.
func (m Metadata) Attributes() attribute.Attributes {
	return m.attributes
}

// NewRecord allows Processor implementations to construct export
// records.  The Descriptor, Labels, and Aggregator represent
// aggregate metric events received over a single collection period.
func NewRecord(descriptor *sdkapi.Descriptor, attrs attribute.Attributes, aggregation aggregation.Aggregation, start, end time.Time) Record {
	return Record{
		Metadata: Metadata{
			descriptor: descriptor,
			attributes: attrs,
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
