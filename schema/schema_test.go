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

package schema

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

var fileFormat = "1.1.0"

func TestSchemaInvalidFileFormat(t *testing.T) {
	s := &Schema{
		FileFormat: "not a semver",
		SchemaURL:  "http://localhost",
	}
	assert.ErrorContains(t, s.validate(), "invalid file format version")
}

func TestSchemaUnsupportedFileFormat(t *testing.T) {
	versions := []*semver.Version{
		semver.New(FileFormatRange.Min.Major()-1, 0, 0, "", ""),
		semver.New(FileFormatRange.Max.Major()+1, 0, 0, "", ""),
		semver.New(FileFormatRange.Max.Major(), FileFormatRange.Max.Minor()+1, 0, "", ""),
		semver.New(FileFormatRange.Max.Major(), FileFormatRange.Max.Minor(), FileFormatRange.Max.Patch()+1, "", ""),
	}
	for _, v := range versions {
		s := &Schema{FileFormat: v.String(), SchemaURL: "http://localhost"}
		assert.Error(t, s.validate(), "unsupported version: %s", v)
	}
}

func TestSchemaMissingSchemaURL(t *testing.T) {
	s := &Schema{FileFormat: fileFormat, SchemaURL: "  "}
	assert.ErrorIs(t, s.validate(), errMissingURL)
}

func TestSchemaInvalidSchemaURL(t *testing.T) {
	u := "\no\t \a valid URL"
	s := &Schema{FileFormat: fileFormat, SchemaURL: u}
	assert.ErrorContains(t, s.validate(), "invalid schema URL", u)
}
