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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

// Environment variable names.
const (
	// batchSpanProcessorScheduleDelayKey is the delay interval between two
	// consecutive exports (i.e. 5000).
	batchSpanProcessorScheduleDelayKey = "OTEL_BSP_SCHEDULE_DELAY"
	// batchSpanProcessorExportTimeoutKey is the maximum allowed time to
	// export data (i.e. 3000).
	batchSpanProcessorExportTimeoutKey = "OTEL_BSP_EXPORT_TIMEOUT"
	// batchSpanProcessorMaxQueueSizeKey is the maximum queue size (i.e. 2048).
	batchSpanProcessorMaxQueueSizeKey = "OTEL_BSP_MAX_QUEUE_SIZE"
	// batchSpanProcessorMaxExportBatchSizeKey is the maximum batch size (i.e.
	// 512). Note: it must be less than or equal to
	// BatchSpanProcessorMaxQueueSize.
	batchSpanProcessorMaxExportBatchSizeKey = "OTEL_BSP_MAX_EXPORT_BATCH_SIZE"

	// attributeValueLengthKey is the maximum allowed attribute value size.
	attributeValueLengthKey = "OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT"

	// attributeCountKey is the maximum allowed span attribute count.
	attributeCountKey = "OTEL_ATTRIBUTE_COUNT_LIMIT"

	// spanAttributeValueLengthKey is the maximum allowed attribute value size
	// for a span.
	spanAttributeValueLengthKey = "OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT"

	// spanAttributeCountKey is the maximum allowed span attribute count for a
	// span.
	spanAttributeCountKey = "OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT"

	// spanEventCountKey is the maximum allowed span event count.
	spanEventCountKey = "OTEL_SPAN_EVENT_COUNT_LIMIT"

	// SpanEventAttributeCountKey is the maximum allowed attribute per span
	// event count.
	spanEventAttributeCountKey = "OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT"

	// spanLinkCountKey is the maximum allowed span link count.
	spanLinkCountKey = "OTEL_SPAN_LINK_COUNT_LIMIT"

	// spanLinkAttributeCountKey is the maximum allowed attribute per span
	// link count.
	spanLinkAttributeCountKey = "OTEL_LINK_ATTRIBUTE_COUNT_LIMIT"
)
