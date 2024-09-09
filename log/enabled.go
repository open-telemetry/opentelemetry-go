// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

// EnabledOpts represents payload for [Logger]'s Enabled method.
type EnabledOpts struct {
	severity Severity
}

// Severity returns the [Severity] level.
func (r *EnabledOpts) Severity() Severity {
	return r.severity
}

// SetSeverity sets the [Severity] level.
func (r *EnabledOpts) SetSeverity(level Severity) {
	r.severity = level
}
