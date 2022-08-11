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

package internal // import "go.opentelemetry.io/otel/schema/internal"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFileFormatField(t *testing.T) {
	// Invalid file format version numbers.
	assert.Error(t, CheckFileFormatField("not a semver", 1, 0))
	assert.Error(t, CheckFileFormatField("2.0.0", 1, 0))
	assert.Error(t, CheckFileFormatField("1.1.0", 1, 0))

	assert.Error(t, CheckFileFormatField("1.2.0", 1, 1))

	// Valid cases.
	assert.NoError(t, CheckFileFormatField("1.0.0", 1, 0))
	assert.NoError(t, CheckFileFormatField("1.0.1", 1, 0))
	assert.NoError(t, CheckFileFormatField("1.0.10000-alpha+4857", 1, 0))

	assert.NoError(t, CheckFileFormatField("1.0.0", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.0.1", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.0.10000-alpha+4857", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.1.0", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.1.1", 1, 1))
}
