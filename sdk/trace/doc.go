// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package trace contains support for OpenTelemetry distributed tracing.

The following assumes a basic familiarity with OpenTelemetry concepts.
See https://opentelemetry.io.

See [go.opentelemetry.io/otel/sdk/internal/x] for information about
the experimental features.

# Environment Variables

The following environment variables can be used for configuration.

OTEL_TRACES_SAMPLER (default: "parentbased_always_on") - the sampler used by
the TracerProvider returned from [NewTracerProvider]. Supported values are
"always_on", "always_off", "traceidratio", "parentbased_always_on",
"parentbased_always_off", and "parentbased_traceidratio". The value is
case-insensitive. Unsupported values are reported to the global error handler,
and the default is used instead.
The configuration can be overridden by the [WithSampler] option.

OTEL_TRACES_SAMPLER_ARG (default: "1.0") - the sampling probability, in the
range [0, 1], used by the "traceidratio" and "parentbased_traceidratio"
samplers. It is ignored if OTEL_TRACES_SAMPLER is not set. Invalid values are
reported to the global error handler, and the default is used instead.
The configuration can be overridden by the [WithSampler] option.

OTEL_BSP_SCHEDULE_DELAY (default: "5000") - the maximum time (in milliseconds)
a span processor returned from [NewBatchSpanProcessor] waits before exporting
the spans it holds. Invalid values are ignored, and the default is used
instead.
The configuration can be overridden by the [WithBatchTimeout] option.

OTEL_BSP_EXPORT_TIMEOUT (default: "30000") - the maximum time (in
milliseconds) a span processor returned from [NewBatchSpanProcessor] waits for
an export to complete. Invalid values are ignored, and the default is used
instead.
The configuration can be overridden by the [WithExportTimeout] option.

OTEL_BSP_MAX_QUEUE_SIZE (default: "2048") - the maximum number of spans a span
processor returned from [NewBatchSpanProcessor] queues for export. Valid values
are positive. Invalid values are ignored, and the default is used instead.
The configuration can be overridden by the [WithMaxQueueSize] option.

OTEL_BSP_MAX_EXPORT_BATCH_SIZE (default: "512") - the maximum number of spans a
span processor returned from [NewBatchSpanProcessor] exports in a single batch.
Valid values are positive. Invalid values are ignored, and the default is used
instead. If the value is greater than OTEL_BSP_MAX_QUEUE_SIZE, the smaller of
512 and OTEL_BSP_MAX_QUEUE_SIZE is used instead.
The configuration can be overridden by the [WithMaxExportBatchSize] option.

OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: no limit) - the maximum
allowed attribute value length, used by [NewSpanLimits] for the
AttributeValueLengthLimit field of [SpanLimits]. If it is not set,
OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT is used. Invalid values are ignored, and the
default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: no limit) - the maximum allowed
attribute value length, used when OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT is not
set. Invalid values are ignored, and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT (default: "128") - the maximum allowed span
attribute count, used by [NewSpanLimits] for the AttributeCountLimit field of
[SpanLimits]. If it is not set, OTEL_ATTRIBUTE_COUNT_LIMIT is used. Invalid
values are ignored, and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_ATTRIBUTE_COUNT_LIMIT (default: "128") - the maximum allowed span
attribute count, used when OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT is not set. Invalid
values are ignored, and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_SPAN_EVENT_COUNT_LIMIT (default: "128") - the maximum allowed span event
count, used by [NewSpanLimits] for the EventCountLimit field of [SpanLimits].
Invalid values are ignored, and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT (default: "128") - the maximum allowed
attribute count per span event, used by [NewSpanLimits] for the
AttributePerEventCountLimit field of [SpanLimits]. Invalid values are ignored,
and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_SPAN_LINK_COUNT_LIMIT (default: "128") - the maximum allowed span link
count, used by [NewSpanLimits] for the LinkCountLimit field of [SpanLimits].
Invalid values are ignored, and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.

OTEL_LINK_ATTRIBUTE_COUNT_LIMIT (default: "128") - the maximum allowed
attribute count per span link, used by [NewSpanLimits] for the
AttributePerLinkCountLimit field of [SpanLimits]. Invalid values are ignored,
and the default is used instead.
The configuration can be overridden by the [WithRawSpanLimits] option.
*/
package trace
