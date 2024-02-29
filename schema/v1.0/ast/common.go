// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ast // import "go.opentelemetry.io/otel/schema/v1.0/ast"

// RenameAttributes corresponds to a section that describes attribute renaming.
type RenameAttributes struct {
	AttributeMap AttributeMap `yaml:"attribute_map"`
}

// AttributeMap corresponds to a section representing a mapping of attribute names.
// The keys are the old attribute name used in the previous version, the values are the
// new attribute name starting from this version.
type AttributeMap map[string]string
