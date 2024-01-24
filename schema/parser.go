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

// Package schema provides types and utilities used to interact with
// [OpenTelemetry schema files].
//
// [OpenTelemetry schema files]: https://github.com/open-telemetry/opentelemetry-specification/blob/007f415120090972e22a90afd499640321f160f3/specification/schemas/file_format_v1.1.0.md
package schema // import "go.opentelemetry.io/otel/schema"

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseFile parses a Schema from the schema file found at path.
func ParseFile(path string) (*Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Parse(file)
}

// Parse parses a Schema read from r.
//
// If r contains a Schema with a file format version outside FileFormatRange,
// an error will be returned.
//
// If r contains an invalid schema URL an error will be returned.
func Parse(r io.Reader) (*Schema, error) {
	d := yaml.NewDecoder(r)
	d.KnownFields(true)

	var s Schema
	if err := d.Decode(&s); err != nil {
		return nil, err
	}

	if err := s.validate(); err != nil {
		return nil, err
	}
	return &s, nil
}
