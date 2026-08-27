// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x provides a stable-SDK-backed experimental facade with bound
// Int64Counter and exact-attribute Finish support for Sum aggregation.
//
// Readers, views, exporters, and non-experimental instruments use
// go.opentelemetry.io/otel/sdk/metric directly. Only counters obtained from
// this package's MeterProvider implement the experimental lifecycle methods.
package x // import "go.opentelemetry.io/otel/sdk/metric/x"
