// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdoutlog provides an exporter for OpenTelemetry logs.
//
// Values and attributes use the attribute package's non-OTLP representation.
//
// The exporter is intended to be used for testing and debugging; it is not
// meant for production use. Additionally, it does not provide an interchange
// format for OpenTelemetry that is supported with any stability or
// compatibility guarantees. If such guarantees are needed, use the OTLP
// exporter instead.
package stdoutlog
