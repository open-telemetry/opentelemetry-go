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

// Spans corresponds to a section representing a list of changes that happened
// to spans schema in a particular version.
type Spans struct {
	Changes []SpansChange
}

// SpanEvents corresponds to a section representing a list of changes that happened
// to span events schema in a particular version.
type SpanEvents struct {
	Changes []SpanEventsChange
}

// SpansChange corresponds to a section representing spans change.
type SpansChange struct {
	RenameAttributes *AttributeMapForSpans `yaml:"rename_attributes"`
}

// AttributeMapForSpans corresponds to a section representing a translation of
// attributes for specific spans.
type AttributeMapForSpans struct {
	ApplyToSpans []types.SpanName `yaml:"apply_to_spans"`
	AttributeMap AttributeMap     `yaml:"attribute_map"`
}

// SpanEventsChange corresponds to a section representing span events change.
type SpanEventsChange struct {
	RenameEvents     *RenameSpanEvents          `yaml:"rename_events"`
	RenameAttributes *RenameSpanEventAttributes `yaml:"rename_attributes"`
}

// RenameSpanEvents corresponds to section representing a renaming of span events.
type RenameSpanEvents struct {
	EventNameMap map[string]string `yaml:"name_map"`
}

// RenameSpanEventAttributes corresponds to section representing a renaming of
// attributes of span events.
type RenameSpanEventAttributes struct {
	ApplyToSpans  []types.SpanName  `yaml:"apply_to_spans"`
	ApplyToEvents []types.EventName `yaml:"apply_to_events"`
	AttributeMap  AttributeMap      `yaml:"attribute_map"`
}
