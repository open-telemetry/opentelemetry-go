// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

// EnabledParameters represents payload for [Logger]'s Enabled method.
type EnabledParameters struct {
	severity Severity
}

// Severity returns the [Severity] level.
func (r *EnabledParameters) Severity() Severity {
	return r.severity
}

// SetSeverity sets the [Severity] level.
func (r *EnabledParameters) SetSeverity(level Severity) {
	r.severity = level
}
