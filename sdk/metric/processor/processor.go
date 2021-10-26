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

package processor // import "go.opentelemetry.io/otel/sdk/metric/processor"

// aggregator import is temporary means to allow for soft deprecation
// of `sdk/export/metric` and will converted to `sdk/metric/aggregator`
// once the work has been done.
import aggregator "go.opentelemetry.io/otel/sdk/export/metric"

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
	aggregator.AggregatorSelector

	// Process is called by the SDK once per internal record,
	// passing the export Accumulation (a Descriptor, the corresponding
	// Labels, and the checkpointed Aggregator). This call has no
	// Context argument because it is expected to perform only
	// computation. An SDK is not expected to call exporters from
	// with Process, use a controller for that (see
	// ./controllers/{pull,push}.
	Process(accum aggregator.Accumulation) error
}
