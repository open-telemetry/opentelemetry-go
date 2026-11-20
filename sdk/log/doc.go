// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package log provides the OpenTelemetry Logs SDK.

See https://opentelemetry.io/docs/concepts/signals/logs/ for information
about OpenTelemetry Logs and
https://opentelemetry.io/docs/concepts/components/ for more information
about OpenTelemetry SDKs.

The entry point for the log package is [NewLoggerProvider].
[LoggerProvider] is the object that all Bridge API calls use to create Loggers
and ultimately emit log records. It should also be used to control the
lifecycle (start, flush, and shutdown) of the Logs SDK.

A LoggerProvider needs to be configured to process log records. This is done
by configuring it with a [Processor] implementation using [WithProcessor].
The log package provides [BatchProcessor] and [SimpleProcessor], which are
configured with an [Exporter] implementation that exports log records to a
given destination. See
[go.opentelemetry.io/otel/exporters] for exporters that can be used with these
Processors.

A LoggerProvider needs to include information about the origin of the data it
generates. It needs to be configured with a Resource by using [WithResource]
to include this information. This Resource should describe the unique runtime
environment in which the instrumented code runs. That way, when telemetry from
multiple instances of the code is collected at a single endpoint, the origin
of each instance is decipherable.

See [go.opentelemetry.io/otel/sdk/log/internal/x] for information about
the experimental features.

See [go.opentelemetry.io/otel/log] for more information about
the OpenTelemetry Logs API.

# Environment Variables

The following environment variables can be used for configuration.

OTEL_BLRP_SCHEDULE_DELAY (default: "1000") - the delay interval (in
milliseconds) between two consecutive exports of a [BatchProcessor].
Non-positive or invalid values are ignored, and the default is used instead.
The configuration can be overridden by the [WithExportInterval] option.

OTEL_BLRP_EXPORT_TIMEOUT (default: "30000") - the maximum time (in
milliseconds) a [BatchProcessor] waits for an export to complete before
canceling it. Non-positive or invalid values are ignored, and the default is
used instead.
The configuration can be overridden by the [WithExportTimeout] option.

OTEL_BLRP_MAX_QUEUE_SIZE (default: "2048") - the maximum number of log
records held in the queue of a [BatchProcessor] before log records are
dropped. Non-positive or invalid values are ignored, and the default is used
instead.
The configuration can be overridden by the [WithMaxQueueSize] option.

OTEL_BLRP_MAX_EXPORT_BATCH_SIZE (default: "512") - the maximum number of log
records a [BatchProcessor] exports in a single batch. Non-positive or invalid
values are ignored, and the default is used instead.
The configuration can be overridden by the [WithExportMaxBatchSize] option.

OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT (default: "128") - the maximum allowed
log record attribute count. Invalid values are ignored, and the default is
used instead.
The configuration can be overridden by the [WithAttributeCountLimit] option.

OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT (default: no limit) - the maximum
allowed attribute value length. Invalid values are ignored, and the default
is used instead.
The configuration can be overridden by the [WithAttributeValueLengthLimit] option.
*/
package log
