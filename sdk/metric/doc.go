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

Package metric implements the OpenTelemetry `Meter` API.  The SDK
supports configurable metrics export behavior through a
`export.MetricBatcher` API.  Most metrics behavior is controlled
by the `MetricBatcher`, including:

1. Selecting the concrete type of aggregation to use
2. Receiving exported data during SDK.Collect()

The call to SDK.Collect() initiates collection.  The SDK calls the
`MetricBatcher` for each current record, asking the aggregator to
export itself.  Aggregators, found in `./aggregators`, are responsible
for receiving updates and exporting their current state.

The SDK.Collect() API should be called by an exporter.  During the
call to Collect(), the exporter receives calls in a single-threaded
context.  No locking is required because the SDK.Collect() call
prevents concurrency.

The SDK uses lock-free algorithms to maintain its internal state.
There are three central data structures at work:

1. A sync.Map maps unique (InstrumentID, LabelSet) to records
2. A "primary" atomic list of records
3. A "reclaim" atomic list of records

Collection is oriented around epochs.  The SDK internally has a
notion of the "current" epoch, which is incremented each time
Collect() is called.  Records contain two atomic counter values,
the epoch in which it was last modified and the epoch in which it
was last collected.  Records may be garbage collected when the
epoch in which they were last updated is less than the epoch in
which they were last collected.

Collect() performs a record-by-record scan of all active records
and exports their current state, before incrementing the current
epoch.  Collection events happen at a point in time during
`Collect()`, but all records are not collected in the same instant.

The purpose of the two lists: the primary list is appended-to when
new handles are created and atomically cleared during collect.  The
reclaim list is used as a second chance, in case there is a race
between looking up a record and record deletion.
*/
package metric
