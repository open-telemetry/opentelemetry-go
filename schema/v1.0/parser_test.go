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
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSchemaFile(t *testing.T) {
	ts, err := ParseFile("testdata/valid-example.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, ts)
}

func TestFailParseSchemaFile(t *testing.T) {
	ts, err := ParseFile("testdata/unsupported-file-format.yaml")
	assert.Error(t, err)
	assert.Nil(t, ts)

	ts, err = ParseFile("testdata/invalid-schema-url.yaml")
	assert.Error(t, err)
	assert.Nil(t, ts)
}

func TestFailParseSchema(t *testing.T) {
	_, err := Parse(bytes.NewReader([]byte("")))
	assert.Error(t, err)

	_, err = Parse(bytes.NewReader([]byte("invalid yaml")))
	assert.Error(t, err)

	_, err = Parse(bytes.NewReader([]byte("file_format: 1.0.0")))
	assert.Error(t, err)
}

func TestCheckFileFormatField(t *testing.T) {
	// Invalid file format version numbers.
	assert.Error(t, checkFileFormatField("not a semver"))
	assert.Error(t, checkFileFormatField("2.0.0"))
	assert.Error(t, checkFileFormatField("1.1.0"))

	// Valid cases.
	assert.NoError(t, checkFileFormatField("1.0.0"))
	assert.NoError(t, checkFileFormatField("1.0.1"))
	assert.NoError(t, checkFileFormatField("1.0.10000-alpha+4857"))
}
