// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc

import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"

// Version is the current release version of the OpenTelemetry OTLP over gRPC logs exporter in use.
func Version() string {
	return internal.Version
}
