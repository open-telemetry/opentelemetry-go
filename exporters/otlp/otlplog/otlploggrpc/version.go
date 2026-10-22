// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc

import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"

// Version is the current release version of the OpenTelemetry OTLP over gRPC logs exporter in use.
//
// It reports the module's registered version reference rather than a second
// literal. This module's version is published through
// `./internal/version.go` in versions.yaml, which is what the release tooling
// updates; a literal here is not version-refs entry, so it is never bumped and
// silently reports whatever release it was last written at.
func Version() string {
	return internal.Version
}
