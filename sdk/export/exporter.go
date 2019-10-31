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

package export

import (
	"context"
)

// SpanSyncer is a type for functions that receive a single sampled trace span.
//
// The ExportSpan method is called synchronously. Therefore, it should not take
// forever to process the span.
//
// The SpanData should not be modified.
type SpanSyncer interface {
	ExportSpan(context.Context, *SpanData)
}

// SpanBatcher is a type for functions that receive batched of sampled trace
// spans.
//
// The ExportSpans method is called asynchronously. However its should not take
// forever to process the spans.
//
// The SpanData should not be modified.
type SpanBatcher interface {
	ExportSpans(context.Context, []*SpanData)
}

// MetricBatcher is responsible for deciding which kind of aggregation
// to use and gathering exported results from the SDK.  The standard SDK
// supports binding only one of these interfaces, i.e., a single exporter.
//
// Multiple-exporters could be implemented by implementing this interface
// for a group of MetricBatcher.
type MetricBatcher interface {
	// AggregatorFor should return the kind of aggregator
	// suited to the requested export.  Returning `nil`
	// indicates to ignore the metric update.
	//
	// Note: This is context-free because the handle should not be
	// bound to the incoming context.  This call should not block.
	AggregatorFor(MetricRecord) MetricAggregator

	// Export receives pairs of records and aggregators
	// during the SDK Collect().  Exporter implementations
	// must access the specific aggregator to receive the
	// exporter data, since the format of the data varies
	// by aggregation.
	Export(context.Context, MetricRecord, MetricAggregator)
}
