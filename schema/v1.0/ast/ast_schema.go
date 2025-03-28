// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package ast provides abstract syntax tree parsing for the OpenTelemetry
// schema.
package ast // import "go.opentelemetry.io/otel/schema/v1.0/ast"

import "go.opentelemetry.io/otel/schema/v1.0/types"

// Schema represents a Schema file in FileFormat 1.0.0 as defined in
// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md
type Schema struct {
	// Schema file format. SHOULD be 1.0.0 for the current specification version.
	// See https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-file-format-number
	FileFormat string `yaml:"file_format"`

	// Schema URL is an identifier of a Schema. The URL specifies a location of this
	// Schema File that can be retrieved (so it is a URL and not just a URI) using HTTP
	// or HTTPS protocol.
	// See https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-url
	SchemaURL string `yaml:"schema_url"`

	// Versions section that lists changes that happened in each particular version.
	Versions map[types.TelemetryVersion]VersionDef
}

// VersionDef corresponds to a section representing one version under the "versions"
// top-level key.
type VersionDef struct {
	All        Attributes
	Resources  Attributes
	Spans      Spans
	SpanEvents SpanEvents `yaml:"span_events"`
	Logs       Logs
	Metrics    Metrics
}

// Attributes corresponds to a section representing a list of changes that
// happened in a particular version.
type Attributes struct {
	Changes []AttributeChange
}

// AttributeChange corresponds to a section representing attribute changes.
type AttributeChange struct {
	RenameAttributes *RenameAttributes `yaml:"rename_attributes"`
}
