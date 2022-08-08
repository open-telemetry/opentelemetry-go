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

package schema // import "go.opentelemetry.io/otel/schema/v1.1"

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"

	"go.opentelemetry.io/otel/schema/internal"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
)

// Major file version number that this library supports.
const supportedFormatMajor = 1

// Maximum minor version number that this library supports.
const supportedFormatMinor = 1

// ParseFile a schema file. schemaFilePath is the file path.
func ParseFile(schemaFilePath string) (*ast.Schema, error) {
	file, err := os.Open(schemaFilePath)
	if err != nil {
		return nil, err
	}
	return Parse(file)
}

// Parse a schema file. schemaFileContent is the readable content of the schema file.
func Parse(schemaFileContent io.Reader) (*ast.Schema, error) {
	var ts ast.Schema
	d := yaml.NewDecoder(schemaFileContent)
	err := d.Decode(&ts)
	if err != nil {
		return nil, err
	}

	err = internal.CheckFileFormatField(ts.FileFormat, supportedFormatMajor, supportedFormatMinor)
	if err != nil {
		return nil, err
	}

	err = internal.CheckSchemaURL(ts.SchemaURL)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}
