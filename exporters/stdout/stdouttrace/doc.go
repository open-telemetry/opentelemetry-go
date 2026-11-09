// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdouttrace contains an OpenTelemetry exporter for tracing
// telemetry to be written to an output destination as JSON.
//
// Attribute values use the OpenTelemetry representation for non-OTLP protocols.
// Span, event, and link attributes are lists of single-key JSON objects that
// preserve their original order. Resource and instrumentation scope attributes
// are JSON objects. Values use native JSON types, including nested arrays and
// maps; byte slices use Base64 strings and empty values use null. Non-finite
// floats use the JSON strings "NaN", "Infinity", and "-Infinity".
//
// See [go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/x] for information about
// the experimental features.
package stdouttrace
