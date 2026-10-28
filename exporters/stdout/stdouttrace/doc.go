// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdouttrace contains an OpenTelemetry exporter for tracing
// telemetry to be written to an output destination as JSON.
//
// Non-finite floating-point attribute values are encoded as the JSON strings
// "NaN", "Infinity", and "-Infinity", including within arrays and maps.
// The output retains the attribute type information used for debugging.
//
// See [go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/x] for information about
// the experimental features.
package stdouttrace
