// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdouttrace contains an OpenTelemetry exporter for tracing
// telemetry to be written to an output destination as JSON.
//
// Span attributes use the
// [OpenTelemetry representation for non-OTLP protocols].
//
// See [go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/x] for information about
// the experimental features.
//
// [OpenTelemetry representation for non-OTLP protocols]: https://opentelemetry.io/docs/specs/otel/common/#anyvalue-representation-for-non-otlp-protocols
package stdouttrace
