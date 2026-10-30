// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdoutmetric provides an exporter for OpenTelemetry metric
// telemetry.
//
// The default JSON output represents non-finite floating-point values as
// the strings "NaN", "Infinity", and "-Infinity", including values nested in
// attributes. Encoders supplied with [WithEncoder] receive the original metric
// data and determine their own representation.
//
// The exporter is intended to be used for testing and debugging, it is not
// meant for production use. Additionally, it does not provide an interchange
// format for OpenTelemetry that is supported with any stability or
// compatibility guarantees. If these are needed features, please use the OTLP
// exporter instead.
package stdoutmetric
