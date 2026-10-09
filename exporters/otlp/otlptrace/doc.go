// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package otlptrace contains abstractions for OTLP span exporters.
See the official OTLP span exporter implementations:
  - [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc],
  - [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp].

The SyncExporter and SyncClient provide a synchronous export path with an
explicit ownership contract. SyncClient.UploadTracesSync guarantees the
supplied trace data is fully consumed before returning and is not retained
afterwards, allowing SyncExporter to safely reuse arena-backed allocations.
*/
package otlptrace
