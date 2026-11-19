// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp

import "go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal"

// Version is the current release version of the OpenTelemetry OTLP over
// HTTP/protobuf profiles exporter in use.
func Version() string {
	return internal.Version
}
