// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import "go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal"

// Version is the current release version of the OpenTelemetry OTLP over gRPC
// profiles exporter in use.
func Version() string {
	return internal.Version
}
