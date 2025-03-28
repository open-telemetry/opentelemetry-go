// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package types provides types for the OpenTelemetry schema.
package types // import "go.opentelemetry.io/otel/schema/v1.0/types"

// TelemetryVersion is a version number key in the schema file (e.g. "1.7.0").
type TelemetryVersion string

// SpanName is span name string.
type SpanName string

// EventName is an event name string.
type EventName string

// MetricName is a metric name string.
type MetricName string
