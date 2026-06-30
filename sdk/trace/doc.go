// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package trace contains support for OpenTelemetry distributed tracing.

The following assumes a basic familiarity with OpenTelemetry concepts.
See https://opentelemetry.io.

# Environment Variables

The following environment variables are used by this package.

OTEL_TRACES_SAMPLER, OTEL_TRACES_SAMPLER_ARG -
configure the default sampler used by [NewTracerProvider]. Supported sampler
names are `always_on`, `always_off`, `traceidratio`,
`parentbased_always_on`, `parentbased_always_off`, and
`parentbased_traceidratio`. Invalid or unsupported values fall back to
`ParentBased(AlwaysSample)`.

OTEL_BSP_SCHEDULE_DELAY (default: `5000`), OTEL_BSP_EXPORT_TIMEOUT
(default: `30000`), OTEL_BSP_MAX_QUEUE_SIZE (default: `2048`), and
OTEL_BSP_MAX_EXPORT_BATCH_SIZE (default: `512`) -
configure the batch span processor created by [NewBatchSpanProcessor] or
[WithBatcher]. The duration values are interpreted as milliseconds.

OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: unlimited),
OTEL_ATTRIBUTE_COUNT_LIMIT (default: `128`),
OTEL_SPAN_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: unlimited),
OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT (default: `128`),
OTEL_SPAN_EVENT_COUNT_LIMIT (default: `128`),
OTEL_EVENT_ATTRIBUTE_COUNT_LIMIT (default: `128`),
OTEL_SPAN_LINK_COUNT_LIMIT (default: `128`), and
OTEL_LINK_ATTRIBUTE_COUNT_LIMIT (default: `128`) -
configure [NewSpanLimits]. The span-specific attribute limit variables take
precedence over their general `OTEL_ATTRIBUTE_*` counterparts.

Resource-related environment variables, including OTEL_RESOURCE_ATTRIBUTES and
OTEL_SERVICE_NAME, are documented in the `go.opentelemetry.io/otel/sdk/resource`
package and are applied when this package uses the default resource or
[WithResource].

See [go.opentelemetry.io/otel/sdk/internal/x] for information about
shared experimental SDK environment variables, including
OTEL_GO_X_OBSERVABILITY and OTEL_GO_X_SELF_OBSERVABILITY.
*/
package trace // import "go.opentelemetry.io/otel/sdk/trace"
