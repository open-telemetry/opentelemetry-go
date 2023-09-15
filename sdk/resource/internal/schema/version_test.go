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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionErrors(t *testing.T) {
	schemaURLs := []string{
		"",
		"not a valid URL: 🌭",
		"https://invalid.host/schemas/1.21.0",
		"https://opentelemetry.io/invalid/path/1.21.0",
		"https://opentelemetry.io/schemas/invalid_version",
	}

	for _, u := range schemaURLs {
		_, err := Version(u)
		assert.Errorf(t, err, "schema URL: %q", u)
	}
}

func TestVersion(t *testing.T) {
	schemaURL := "https://opentelemetry.io/schemas/1.21.0"
	got, err := Version(schemaURL)
	require.NoError(t, err)
	assert.Equal(t, "1.21.0", got.Original())
}
