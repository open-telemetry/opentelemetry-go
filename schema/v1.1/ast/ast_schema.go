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

package ast // import "go.opentelemetry.io/otel/schema/v1.1/ast"

import (
	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
)

// Schema represents a Schema file in FileFormat 1.1.0 as defined in
// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md
type Schema struct {
	// Schema file format. SHOULD be 1.1.0 for the current specification version.
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
// Note that most of the fields are the same as in ast 1.0 package, only Metrics are defined
// differently, since only that field has changed from 1.0 to 1.1 of schema file format.
type VersionDef struct {
	All        ast10.Attributes
	Resources  ast10.Attributes
	Spans      ast10.Spans
	SpanEvents ast10.SpanEvents `yaml:"span_events"`
	Logs       ast10.Logs
	Metrics    Metrics
}
