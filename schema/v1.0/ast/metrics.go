// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ast // import "go.opentelemetry.io/otel/schema/v1.0/ast"

import "go.opentelemetry.io/otel/schema/v1.0/types"

// Metrics corresponds to a section representing a list of changes that happened
// to metrics schema in a particular version.
type Metrics struct {
	Changes []MetricsChange
}

// MetricsChange corresponds to a section representing metrics change.
type MetricsChange struct {
	RenameMetrics    map[types.MetricName]types.MetricName `yaml:"rename_metrics"`
	RenameAttributes *AttributeMapForMetrics               `yaml:"rename_attributes"`
}

// AttributeMapForMetrics corresponds to a section representing a translation of
// attributes for specific metrics.
type AttributeMapForMetrics struct {
	ApplyToMetrics []types.MetricName `yaml:"apply_to_metrics"`
	AttributeMap   AttributeMap       `yaml:"attribute_map"`
}
