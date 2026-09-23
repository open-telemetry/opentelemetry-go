// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package trace contains support for OpenTelemetry distributed tracing.

The following assumes a basic familiarity with OpenTelemetry concepts.
See https://opentelemetry.io.

See [go.opentelemetry.io/otel/sdk/internal/x] for information about
the experimental features.

# Environment Variables

The environment variables described below can be used for configuration.
Values of the OTEL_BSP_* and span limit environment variables that are not
valid integers are ignored and the default is used.

OTEL_TRACES_SAMPLER (default: "parentbased_always_on") -
the sampler used by the TracerProvider returned from [NewTracerProvider].
Supported values: "always_on", "always_off", "traceidratio",
"parentbased_always_on", "parentbased_always_off", "parentbased_traceidratio".
The value is case-insensitive.
An unsupported value is reported to the global error handler and the default is used.
The configuration can be overridden by [WithSampler] option.

OTEL_TRACES_SAMPLER_ARG (default: "1.0") -
the sampling probability, in the range [0, 1], used by the "traceidratio" and
"parentbased_traceidratio" samplers.
An invalid value is reported to the global error handler and 1.0 is used.
It is ignored if OTEL_TRACES_SAMPLER is not set.
The configuration can be overridden by [WithSampler] option.

OTEL_BSP_SCHEDULE_DELAY (default: "5000") -
maximum time in milliseconds the span processor returned from
[NewBatchSpanProcessor] waits before exporting the spans it holds.
The configuration can be overridden by [WithBatchTimeout] option.

OTEL_BSP_EXPORT_TIMEOUT (default: "30000") -
maximum time in milliseconds the span processor returned from
[NewBatchSpanProcessor] waits for an export to complete.
The configuration can be overridden by [WithExportTimeout] option.

OTEL_BSP_MAX_QUEUE_SIZE (default: "2048") -
maximum number of spans the span processor returned from
[NewBatchSpanProcessor] queues for export.
Valid values are positive.
The configuration can be overridden by [WithMaxQueueSize] option.

OTEL_BSP_MAX_EXPORT_BATCH_SIZE (default: "512") -
maximum number of spans the span processor returned from
[NewBatchSpanProcessor] exports in a single batch.
Valid values are positive.
If the value is greater than OTEL_BSP_MAX_QUEUE_SIZE, the smaller of 512 and
OTEL_BSP_MAX_QUEUE_SIZE is used.
The configuration can be overridden by [WithMaxExportBatchSize] option.

The span limit environment variables below are read by [NewSpanLimits], which
provides the default limits of the TracerProvider returned from [NewTracerProvider].
Their values are interpreted as described by the corresponding [SpanLimits] fields.
The configuration can be overridden by [WithRawSpanLimits] option.

OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT, OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: unlimited) -
maximum allowed attribute value length.
OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT takes precedence over OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT.

OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT, OTEL_ATTRIBUTE_COUNT_LIMIT (default: "128") -
maximum allowed span attribute count.
OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT takes precedence over OTEL_ATTRIBUTE_COUNT_LIMIT.

OTEL_SPAN_EVENT_COUNT_LIMIT (default: "128") -
maximum allowed span event count.

OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT (default: "128") -
maximum allowed attribute count per span event.

OTEL_SPAN_LINK_COUNT_LIMIT (default: "128") -
maximum allowed span link count.

OTEL_LINK_ATTRIBUTE_COUNT_LIMIT (default: "128") -
maximum allowed attribute count per span link.
*/
package trace
