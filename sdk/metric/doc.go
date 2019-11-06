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

/*
Package metric implements the OpenTelemetry metric.Meter API.  The SDK
supports configurable metrics export behavior through a collection of
export interfaces that support various export strategies, described below.

The metric.Meter API consists of methods for constructing each of the
basic kinds of metric instrument.  There are six types of instrument
available to the end user, comprised of three basic kinds of metric
instrument (Counter, Gauge, Measure) crossed with two kinds of number
(int64, float64).

The API assists the SDK by consolidating the variety of metric instruments
into a narrower interface, allowing the SDK to avoid repetition of
boilerplate.  The API and SDK are separated such that an event reacheing
the SDK has a uniform structure: an instrument, a label set, and a
numerical value.

To this end, the API uses a core.Number type to represent either an int64
or a float64, depending on the instrument's definition.  A single
implementation interface is used for instruments, metric.InstrumentImpl,
and a single implementation interface is used for handles,
metric.HandleImpl.

There are three entry points for events in the Metrics API: via instrument
handles, via direct instrument calls, and via BatchRecord.  The SDK is
designed with handles as the primary entry point, the other two entry
points are implemented in terms of short-lived handles.  For example, the
implementation of a direct call allocates a handle, operates on the
handle, and releases the handle. Similarly, the implementation of
RecordBatch uses a short-lived handle for each measurement in the batch.

Internal Structure

The SDK is designed with minimal use of locking, to avoid adding
contention for user-level code.  For each handle, whether it is held by
user-level code or a short-lived device, there exists an internal record
managed by the SDK.  Each internal record corresponds to a specific
instrument and label set combination.

A sync.Map maintains the mapping of current instruments and label sets to
internal records.  To create a new handle, the SDK consults the Map to
locate an existing record, otherwise it constructs a new record.  The SDK
maintains a count of the number of references to each record, ensuring
that records are not reclaimed from the Map while they are still active
from the user's perspective.

Metric collection is performed via a single-threaded call to Collect that
sweeps through all records in the SDK, checkpointing their state.  When a
record is discovered that has no references and has not been updated since
the prior collection pass, it is marked for reclamation and removed from
the Map.  There exists, at this moment, a race condition since another
goroutine could, in the same instant, obtain a reference to the handle.

The SDK is designed to tolerate this sort of race condition, in the name
of reducing lock contention.  It is possible for more than one record with
identical instrument and label set to exist simultaneously, though only
one can be linked from the Map at a time.  To avoid lost updates, the SDK
maintains two additional linked lists of records, one managed by the
collection code path and one managed by the instrumentation code path.

The SDK maintains a current epoch number, corresponding to the number of
completed collections.  Each record contains the last epoch during which
it was collected and updated.  These variables allow the collection code
path to detect stale records while allowing the instrumentation code path
to detect potential reclamations.  When the instrumentation code path
detects a potential reclamation, it adds itself to the second linked list,
where records are saved from reclamation.

Each record has an associated aggregator, which maintains the current
state resulting from all metric events since its last checkpoint.
Aggregators may be lock-free or they may use locking, but they should
expect to be called concurrently.  Because of the tolerated race condition
described above, aggregators must be capable of merging with another
aggregator of the same type.

Export Pipeline

While the SDK serves to maintain a current set of records and
coordinate collection, the behavior of a metrics export pipeline is
configured through the export types in
go.opentelemetry.io/otel/sdk/export/metric.  It is important to keep
in mind the context these interfaces are called from.  There are two
contexts, instrumentation context, where a user-level goroutine that
enters the SDK resulting in a new record, and collection context,
where a system-level thread performs a collection pass through the
SDK.

Descriptor is a struct that describes the metric instrument to the
export pipeline, containing the name, recommended aggregation keys,
units, description, metric kind (counter, gauge, or measure), number
kind (int64 or float64), and whether the instrument has alternate
semantics or not (i.e., monotonic=false counter, monotonic=true gauge,
absolute=false measure).  A Descriptor accompanies metric data as it
passes through the export pipeline.

The AggregationSelector interface supports choosing the method of
aggregation to apply to a particular instrument.  Given the
Descriptor, this AggregatorFor method returns an implementation of
Aggregator.  If this interface returns nil, the metric will be
disabled.  The aggregator should be matched to the capabilities of the
exporter.  Selecting the aggregator for counter and gauge instruments
is relatively straightforward, but for measure instruments there are
numerous choices with different cost and quality tradeoffs.

Aggregator is an interface which implements a concrete strategy for
aggregating metric updates.  Several Aggregator implementations are
provided by the SDK.  Aggregators may be lock-free or use locking,
depending on their structure and semantics.  Aggregators implement an
Update method, called in instrumentation context, to receive a single
metric event.  Aggregators implement a Checkpoint method, called in
collection context, to save a checkpoint of the current state.
Aggregators implement a Merge method, also called in collection
context, that combines state from two aggregators into one.  Each SDK
record has an associated aggregator.

Batcher

LabelEncoder

Producer

ProducedRecord

Exporter

Controller

metric.MeterProvider


TODO: think about name for Producer/ProducedRecord


*/
package metric // import "go.opentelemetry.io/otel/sdk/metric"
