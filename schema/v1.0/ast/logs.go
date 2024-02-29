// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ast // import "go.opentelemetry.io/otel/schema/v1.0/ast"

// Logs corresponds to a section representing a list of changes that happened
// to logs schema in a particular version.
type Logs struct {
	Changes []LogsChange
}

// LogsChange corresponds to a section representing logs change.
type LogsChange struct {
	RenameAttributes *RenameAttributes `yaml:"rename_attributes"`
}
