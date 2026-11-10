// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdoutlog provides an exporter for OpenTelemetry logs.
//
// Log bodies and attributes use the
// [OpenTelemetry representation for non-OTLP protocols].
//
// The exporter is intended to be used for testing and debugging; it is not
// meant for production use. Additionally, it does not provide an interchange
// format for OpenTelemetry that is supported with any stability or
// compatibility guarantees. If such guarantees are needed, use the OTLP
// exporter instead.
//
// [OpenTelemetry representation for non-OTLP protocols]: https://opentelemetry.io/docs/specs/otel/common/#anyvalue-representation-for-non-otlp-protocols
package stdoutlog
