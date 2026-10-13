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

Use the transport sync constructors to avoid a type assertion:

	httpClient := otlptracehttp.NewSyncClient(otlptracehttp.WithEndpointURL(...))
	exp, err := NewSync(ctx, httpClient)

	grpcClient := otlptracegrpc.NewSyncClient(otlptracegrpc.WithEndpoint(...))
	exp, err := NewSync(ctx, grpcClient)
*/
package otlptrace
