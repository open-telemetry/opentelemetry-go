// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdoutlog provides an exporter for OpenTelemetry logs.
//
// Values and attributes use the attribute package's non-OTLP representation.
// Nested maps and resource and scope attribute sets are encoded as JSON objects.
// Log attributes retain their order as an array of single-attribute objects.
// The log body retains its outer type envelope.
//
// Non-finite floating-point values are exported as the JSON strings "NaN",
// "Infinity", and "-Infinity".
//
// The exporter is intended to be used for testing and debugging; it is not
// meant for production use. Additionally, it does not provide an interchange
// format for OpenTelemetry that is supported with any stability or
// compatibility guarantees. If such guarantees are needed, use the OTLP
// exporter instead.
package stdoutlog
