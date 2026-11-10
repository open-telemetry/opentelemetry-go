// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdouttrace contains an OpenTelemetry exporter for tracing
// telemetry to be written to an output destination as JSON.
//
// Attributes use the
// [OpenTelemetry representation for non-OTLP protocols].
//
// The exporter is intended to be used for testing and debugging, it is not
// meant for production use. Additionally, it does not provide an interchange
// format for OpenTelemetry that is supported with any stability or
// compatibility guarantees. If these are needed features, please use the OTLP
// exporter instead.
//
// See [go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/x] for information about
// the experimental features.
//
// [OpenTelemetry representation for non-OTLP protocols]: https://opentelemetry.io/docs/specs/otel/common/#anyvalue-representation-for-non-otlp-protocols
package stdouttrace
