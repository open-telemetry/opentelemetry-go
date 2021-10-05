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

// Logs corresponds to a section representing a list of changes that happened
// to logs schema in a particular version.
type Logs struct {
	Changes []LogsChange
}

// LogsChange corresponds to a section representing logs change.
type LogsChange struct {
	RenameAttributes *RenameAttributes `yaml:"rename_attributes"`
}
