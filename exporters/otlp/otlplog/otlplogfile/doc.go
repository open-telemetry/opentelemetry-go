// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package otlplogfile provides an OTLP log exporter that outputs log records to a JSON line file. The exporter uses a buffered
file writer to write log records to file to reduce I/O and improve performance.

All Exporters must be created with [New].

See: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/file-exporter.md
*/
package otlplogfile // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile"
