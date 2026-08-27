// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x provides a proof-of-concept experimental metric SDK with bound
// Int64Counter and exact-attribute Finish support.
//
// This package is an alternative to go.opentelemetry.io/otel/sdk/metric. Its
// Reader and MeterProvider types cannot be mixed with the stable SDK types.
// Only Int64Counter and Sum aggregation are implemented by this proof of
// concept; other instrument methods return no-op instruments.
package x // import "go.opentelemetry.io/otel/sdk/metric/x"
