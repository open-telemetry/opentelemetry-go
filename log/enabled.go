// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

// EnabledParam represents payload for [Logger]'s Enabled method.
type EnabledParam struct {
	severity Severity
}

// Severity returns the [Severity] level.
func (r *EnabledParam) Severity() Severity {
	return r.severity
}

// SetSeverity sets the [Severity] level.
func (r *EnabledParam) SetSeverity(level Severity) {
	r.severity = level
}
