// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetric // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric"

import (
	"context"

	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

// Client manages connections to the collector, handles the
// transformation of data into wire format, and the transmission of that
// data to the collector.
type Client interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// UploadMetrics sends protoMetrics to connected endpoint.
	//
	// Retryable errors from the server will be handled according to any
	// RetryConfig the client was created with.
	UploadMetrics(context.Context, *metricpb.ResourceMetrics) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// Shutdown shuts down the client, freeing all resources.
	Shutdown(context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}
