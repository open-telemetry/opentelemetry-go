// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x provides stable SDK options for experimental bound instruments
// and exact-attribute Finish support.
//
// Pass WithBinding, WithFinish, or both to
// go.opentelemetry.io/otel/sdk/metric.NewMeterProvider. The proof of concept
// currently supports Int64Counter with Sum aggregation.
package x // import "go.opentelemetry.io/otel/sdk/metric/x"
