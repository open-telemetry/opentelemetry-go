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

package schema // import "go.opentelemetry.io/otel/schema/v1.0"

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v2"

	"go.opentelemetry.io/otel/schema/v1.0/ast"
)

// Major file version number that this library supports.
const supportedFormatMajor = 1

// Maximum minor version number that this library supports.
const supportedFormatMinor = 0

// Maximum major+minor version number that this library supports, as a string.
var supportedFormatMajorMinor = strconv.Itoa(supportedFormatMajor) + "." +
	strconv.Itoa(supportedFormatMinor) // 1.0

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

	if err := checkFileFormatField(ts.FileFormat); err != nil {
		return nil, err
	}

	if strings.TrimSpace(ts.SchemaURL) == "" {
		return nil, fmt.Errorf("schema_url field is missing")
	}

	if _, err := url.Parse(ts.SchemaURL); err != nil {
		return nil, fmt.Errorf("invalid URL specified in schema_url field: %w", err)
	}

	return &ts, nil
}

// checkFileFormatField validates the file format field according to the rules here:
// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-file-format-number
func checkFileFormatField(fileFormat string) error {
	// Verify that the version number in the file is a semver.
	fileFormatParsed, err := semver.StrictNewVersion(fileFormat)
	if err != nil {
		return fmt.Errorf(
			"invalid schema file format version number %q (expected semver): %w",
			fileFormat, err,
		)
	}

	// Check that the major version number in the file is the same as what we expect.
	if fileFormatParsed.Major() != supportedFormatMajor {
		return fmt.Errorf(
			"this library cannot parse file formats with major version other than %v",
			supportedFormatMajor,
		)
	}

	// Check that the file minor version number is not greater than
	// what is requested supports.
	if fileFormatParsed.Minor() > supportedFormatMinor {
		return fmt.Errorf(
			"unsupported schema file format minor version number, expected no newer than %v, got %v",
			supportedFormatMajorMinor+".x", fileFormat,
		)
	}

	// Patch, prerelease and metadata version number does not matter, so we don't check it.

	return nil
}
