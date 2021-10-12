// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
