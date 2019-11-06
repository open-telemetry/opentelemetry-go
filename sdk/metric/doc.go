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

While the SDK serves to maintain a current set of records and coordinate
collection, the behavior of a metrics export pipeline is configured
through the export types in go.opentelemetry.io/otel/sdk/export/metric.

AggregationSelector
LabelEncoder

They are briefly summarized here:

Aggregator


AggregationSelector: decides which aggregator to use
Batcher: determine the aggregation dimensions, group (and de-dup) records
Descriptor: summarizes an instrument and its metadata
Record: interface to the SDK-internal record
LabelEncoder: defines a unique mapping from label set to encoded string
Producer: interface to the batcher's checkpoint
ProducedRecord: result of the batcher's grouping
Exporter: output produced records to their final destination

One final type, a Controller, implements the metric.MeterProvider
interface and is responsible for initiating collection.

*/
package metric // import "go.opentelemetry.io/otel/sdk/metric"
