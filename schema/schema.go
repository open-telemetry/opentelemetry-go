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

package schema // import "go.opentelemetry.io/otel/schema"

// Schema represents an OpenTelemetry [Schema file].
//
// [Schema file]: https://github.com/open-telemetry/opentelemetry-specification/blob/007f415120090972e22a90afd499640321f160f3/specification/schemas/file_format_v1.1.0.md
type Schema struct {
	// FileFormat is the [schema file format].
	//
	// [schema file format]: https://github.com/open-telemetry/opentelemetry-specification/blob/007f415120090972e22a90afd499640321f160f3/specification/schemas/file_format_v1.1.0.md#schema-file-format-number
	FileFormat string `yaml:"file_format"`

	// SchemaURL is the [URL] for the Schema file.
	//
	// [URL]: https://github.com/open-telemetry/opentelemetry-specification/blob/007f415120090972e22a90afd499640321f160f3/specification/schemas/file_format_v1.1.0.md#schema-url
	SchemaURL string `yaml:"schema_url"`

	// Versions are the telemetry transforms that apply for each semantic
	// convention version.
	Versions map[SemConvVersion]Transform
}

// SemConvVersion is a semantic conventions version used by a schema file (e.g.
// "1.7.0").
type SemConvVersion string

// Transform is all the applicable telemetry changes for a particular semantic
// convention version.
type Transform struct {
	All        All
	Resources  Resources
	Spans      Spans
	SpanEvents SpanEvents `yaml:"span_events"`
	Logs       Logs
	Metrics    Metrics
}

// All defines the transforms that apply to all types of telemetry data.
type All struct {
	Changes []AllChange
}

// AllChange is a change that applies to all types of telemetry data.
type AllChange struct {
	RenameAttributes *RenameAttributes `yaml:"rename_attributes"`
}

// RenameAttributes defines a rename of attributes.
type RenameAttributes struct {
	// AttributeMap is a mapping of old attribute keys to the new attribute
	// keys. Attributes that have the same key as a key this map need to have
	// that their key renamed to the corresponding map value.
	AttributeMap map[string]string `yaml:"attribute_map"`
}

// Resources defines the transforms that apply to OpenTelemetry resources.
type Resources struct {
	Changes []ResourcesChange
}

// ResourcesChange is a change that applies to OpenTelemetry resources.
type ResourcesChange struct {
	RenameAttributes *RenameAttributes `yaml:"rename_attributes"`
}

// Spans defines the transforms that apply to OpenTelemetry spans.
type Spans struct {
	Changes []SpansChange
}

// SpansChange is a change that applies to OpenTelemetry spans.
type SpansChange struct {
	RenameAttributes *RenameSpansAttributes `yaml:"rename_attributes"`
}

// RenameSpansAttributes defines a rename of span attributes.
type RenameSpansAttributes struct {
	// ApplyToSpans is a slice of span names that this rename applies to. If no
	// span names are provided, the rename applies to all spans.
	ApplyToSpans []string `yaml:"apply_to_spans"`
	// AttributeMap is a mapping of old attribute keys to the new attribute
	// keys. Attributes that have the same key as a key this map need to have
	// that their key renamed to the corresponding map value.
	AttributeMap map[string]string `yaml:"attribute_map"`
}

// SpanEvents defines the transforms that apply to OpenTelemetry span
// events.
type SpanEvents struct {
	Changes []SpanEventsChange
}

// SpanEventsChange is a change that applies to OpenTelemetry span events.
type SpanEventsChange struct {
	RenameEvents     *RenameSpanEvents           `yaml:"rename_events"`
	RenameAttributes *RenameSpanEventsAttributes `yaml:"rename_attributes"`
}

// RenameSpanEvents defines a rename of span events.
type RenameSpanEvents struct {
	// EventNameMap is a mapping of old event names to the new event names.
	// Events that have the same name as a key in this map need to have their
	// name renamed to the corresponding map value.
	EventNameMap map[string]string `yaml:"name_map"`
}

// RenameSpansAttributes defines a rename of span event attributes.
type RenameSpanEventsAttributes struct {
	// ApplyToSpans is a slice of span names that this rename applies to. If no
	// span names are provided, the rename applies to all spans.
	ApplyToSpans []string `yaml:"apply_to_spans"`
	// ApplyToEvents is a slice of event names that this rename applies to. If
	// no event names are provided, the rename applies to all spans.
	ApplyToEvents []string `yaml:"apply_to_events"`
	// AttributeMap is a mapping of old attribute keys to the new attribute
	// keys. Attributes that have the same key as a key this map need to have
	// that their key renamed to the corresponding map value.
	AttributeMap map[string]string `yaml:"attribute_map"`
}

// Logs defines the transforms that apply to OpenTelemetry logs.
type Logs struct {
	Changes []LogsChange
}

// LogsChange is a change that applies to OpenTelemetry logs.
type LogsChange struct {
	RenameAttributes *RenameAttributes `yaml:"rename_attributes"`
}

// Metrics defines the transforms that apply to OpenTelemetry metrics.
type Metrics struct {
	Changes []MetricsChange
}

// MetricsChange is a change that applies to OpenTelemetry metrics.
type MetricsChange struct {
	// RenameMetrics is a mapping of old metric names to new metric names.
	// Metrics with names matching a key this map need to be renamed to the
	// corresponding map value.
	RenameMetrics    map[string]string        `yaml:"rename_metrics"`
	RenameAttributes *RenameMetricsAttributes `yaml:"rename_attributes"`
	Split            *SplitMetric             `yaml:"split"` // Added in schema file format 1.1.
}

// RenameMetricsAttributes defines a rename of metric attributes.
type RenameMetricsAttributes struct {
	// ApplyToMetrics is a slice of metric names that this rename applies to.
	// If no metric names are provided, the rename applies to all metrics.
	ApplyToMetrics []string `yaml:"apply_to_metrics"`
	// AttributeMap is a mapping of old attribute keys to the new attribute
	// keys. Attributes that have the same key as a key this map need to have
	// that their key renamed to the corresponding map value.
	AttributeMap map[string]string `yaml:"attribute_map"`
}

// SplitMetric defines how a metric should be split into multiple metrics by
// eliminating an attribute.
//
// This was introduced in [schema file format 1.1].
//
// [schema file format 1.1]: https://github.com/open-telemetry/opentelemetry-specification/pull/2653
type SplitMetric struct {
	// ApplyToMetric is the name of the old metric to split.
	ApplyToMetric string `yaml:"apply_to_metric"`

	// ByAttribute is the name of attribute in the old metric to use for
	// splitting. This attribute will be removed from all of the produced new
	// metrics.
	ByAttribute string `yaml:"by_attribute"`

	// MetricsFromAttributes is a mapping of new metric names to create based
	// on the value of the ByAttribute attribute.
	MetricsFromAttributes map[string]any `yaml:"metrics_from_attributes"`
}
