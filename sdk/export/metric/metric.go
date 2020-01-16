// Copyright 2019, OpenTelemetry Authors
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

//go:generate stringer -type=Kind

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
)

// Batcher is responsible for deciding which kind of aggregation to
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
// checkpointed, allowing the batcher to build the set of metrics
// currently being exported.
//
// The `CheckpointSet` method is called during collection in a
// single-threaded context from the Exporter, giving the exporter
// access to a producer for iterating over the complete checkpoint.
type Batcher interface {
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
	// Labels, and the checkpointed Aggregator).  The Batcher
	// should be prepared to process duplicate (Descriptor,
	// Labels) pairs during this pass due to race conditions, but
	// this will usually be the ordinary course of events, as
	// Aggregators are typically merged according the output set
	// of labels.
	//
	// The Context argument originates from the controller that
	// orchestrates collection.
	Process(ctx context.Context, record Record) error

	// CheckpointSet is the interface used by the controller to
	// access the fully aggregated checkpoint after collection.
	//
	// The returned CheckpointSet is passed to the Exporter.
	CheckpointSet() CheckpointSet

	// FinishedCollection informs the Batcher that a complete
	// collection round was completed.  Stateless batchers might
	// reset state in this method, for example.
	FinishedCollection()
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
	AggregatorFor(*Descriptor) Aggregator
}

// Aggregator implements a specific aggregation behavior, e.g., a
// behavior to track a sequence of updates to a counter, a gauge, or a
// measure instrument.  For the most part, counter and gauge semantics
// are fixed and the provided implementations should be used.  Measure
// metrics offer a wide range of potential tradeoffs and several
// implementations are provided.
//
// Aggregators are meant to compute the change (i.e., delta) in state
// from one checkpoint to the next, with the exception of gauge
// aggregators.  Gauge aggregators are required to maintain the last
// value across checkpoints to implement montonic gauge support.
//
// Note that any Aggregator may be attached to any instrument--this is
// the result of the OpenTelemetry API/SDK separation.  It is possible
// to attach a counter aggregator to a measure instrument (to compute
// a simple sum) or a gauge instrument to a measure instrument (to
// compute the last value).
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
	Update(context.Context, core.Number, *Descriptor) error

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
	Checkpoint(context.Context, *Descriptor)

	// Merge combines the checkpointed state from the argument
	// aggregator into this aggregator's checkpointed state.
	// Merge() is called in a single-threaded context, no locking
	// is required.
	Merge(Aggregator, *Descriptor) error
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
	// The CheckpointSet interface refers to the Batcher that just
	// completed collection.
	Export(context.Context, CheckpointSet) error
}

// LabelEncoder enables an optimization for export pipelines that use
// text to encode their label sets.
//
// This interface allows configuring the encoder used in the SDK
// and/or the Batcher so that by the time the exporter is called, the
// same encoding may be used.
//
// If none is provided, a default will be used.
type LabelEncoder interface {
	// Encode is called (concurrently) in instrumentation context.
	// It should return a unique representation of the labels
	// suitable for the SDK to use as a map key.
	//
	// The exported Labels object retains a reference to its
	// LabelEncoder to determine which encoding was used.
	//
	// The expectation is that Exporters with a pre-determined to
	// syntax for serialized label sets should implement
	// LabelEncoder, thus avoiding duplicate computation in the
	// export path.
	Encode([]core.KeyValue) string
}

// CheckpointSet allows a controller to access a complete checkpoint of
// aggregated metrics from the Batcher.  This is passed to the
// Exporter which may then use ForEach to iterate over the collection
// of aggregated metrics.
type CheckpointSet interface {
	// ForEach iterates over aggregated checkpoints for all
	// metrics that were updated during the last collection
	// period.
	ForEach(func(Record))
}

// Record contains the exported data for a single metric instrument
// and label set.
type Record struct {
	descriptor *Descriptor
	labels     Labels
	aggregator Aggregator
}

// Labels stores complete information about a computed label set,
// including the labels in an appropriate order (as defined by the
// Batcher).  If the batcher does not re-order labels, they are
// presented in sorted order by the SDK.
type Labels struct {
	ordered []core.KeyValue
	encoded string
	encoder LabelEncoder
}

// NewLabels builds a Labels object, consisting of an ordered set of
// labels, a unique encoded representation, and the encoder that
// produced it.
func NewLabels(ordered []core.KeyValue, encoded string, encoder LabelEncoder) Labels {
	return Labels{
		ordered: ordered,
		encoded: encoded,
		encoder: encoder,
	}
}

// Ordered returns the labels in a specified order, according to the
// Batcher.
func (l Labels) Ordered() []core.KeyValue {
	return l.ordered
}

// Encoded is a pre-encoded form of the ordered labels.
func (l Labels) Encoded() string {
	return l.encoded
}

// Encoder is the encoder that computed the Encoded() representation.
func (l Labels) Encoder() LabelEncoder {
	return l.encoder
}

// Len returns the number of labels.
func (l Labels) Len() int {
	return len(l.ordered)
}

// NewRecord allows Batcher implementations to construct export
// records.  The Descriptor, Labels, and Aggregator represent
// aggregate metric events received over a single collection period.
func NewRecord(descriptor *Descriptor, labels Labels, aggregator Aggregator) Record {
	return Record{
		descriptor: descriptor,
		labels:     labels,
		aggregator: aggregator,
	}
}

// Aggregator returns the checkpointed aggregator. It is safe to
// access the checkpointed state without locking.
func (r Record) Aggregator() Aggregator {
	return r.aggregator
}

// Descriptor describes the metric instrument being exported.
func (r Record) Descriptor() *Descriptor {
	return r.descriptor
}

// Labels describes the labels associated with the instrument and the
// aggregated data.
func (r Record) Labels() Labels {
	return r.labels
}

// Kind describes the kind of instrument.
type Kind int8

const (
	// Counter kind indicates a counter instrument.
	CounterKind Kind = iota

	// Gauge kind indicates a gauge instrument.
	GaugeKind

	// Measure kind indicates a measure instrument.
	MeasureKind
)

// Descriptor describes a metric instrument to the exporter.
//
// Descriptors are created once per instrument and a pointer to the
// descriptor may be used to uniquely identify the instrument in an
// exporter.
type Descriptor struct {
	name        string
	metricKind  Kind
	keys        []core.Key
	description string
	unit        unit.Unit
	numberKind  core.NumberKind
	alternate   bool
}

// NewDescriptor builds a new descriptor, for use by `Meter`
// implementations in constructing new metric instruments.
//
// Descriptors are created once per instrument and a pointer to the
// descriptor may be used to uniquely identify the instrument in an
// exporter.
func NewDescriptor(
	name string,
	metricKind Kind,
	keys []core.Key,
	description string,
	unit unit.Unit,
	numberKind core.NumberKind,
	alternate bool,
) *Descriptor {
	return &Descriptor{
		name:        name,
		metricKind:  metricKind,
		keys:        keys,
		description: description,
		unit:        unit,
		numberKind:  numberKind,
		alternate:   alternate,
	}
}

// Name returns the metric instrument's name.
func (d *Descriptor) Name() string {
	return d.name
}

// MetricKind returns the kind of instrument: counter, gauge, or
// measure.
func (d *Descriptor) MetricKind() Kind {
	return d.metricKind
}

// Keys returns the recommended keys included in the metric
// definition.  These keys may be used by a Batcher as a default set
// of grouping keys for the metric instrument.
func (d *Descriptor) Keys() []core.Key {
	return d.keys
}

// Description provides a human-readable description of the metric
// instrument.
func (d *Descriptor) Description() string {
	return d.description
}

// Unit describes the units of the metric instrument.  Unitless
// metrics return the empty string.
func (d *Descriptor) Unit() unit.Unit {
	return d.unit
}

// NumberKind returns whether this instrument is declared over int64
// or a float64 values.
func (d *Descriptor) NumberKind() core.NumberKind {
	return d.numberKind
}

// Alternate returns true when the non-default behavior of the
// instrument was selected.  It returns true if:
//
//   - A counter instrument is non-monotonic
//   - A gauge instrument is monotonic
//   - A measure instrument is non-absolute
//
// TODO: Consider renaming this method, or expanding to provide
// kind-specific tests (e.g., Monotonic(), Absolute()).
func (d *Descriptor) Alternate() bool {
	return d.alternate
}
