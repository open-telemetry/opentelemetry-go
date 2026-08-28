// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x provides a stable-SDK-backed experimental facade with bound
// Int64Counter and exact-attribute Finish support for Sum aggregation.
//
// Configure the provider with options, readers, views, and exporters from
// go.opentelemetry.io/otel/sdk/metric. Only counters obtained from the provider
// returned by NewMeterProvider implement the experimental lifecycle methods.
package x // import "go.opentelemetry.io/otel/sdk/metric/x"
