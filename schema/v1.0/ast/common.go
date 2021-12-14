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

// RenameAttributes corresponds to a section that describes attribute renaming.
type RenameAttributes struct {
	AttributeMap AttributeMap `yaml:"attribute_map"`
}

// AttributeMap corresponds to a section representing a mapping of attribute names.
// The keys are the old attribute name used in the previous version, the values are the
// new attribute name starting from this version.
type AttributeMap map[string]string
