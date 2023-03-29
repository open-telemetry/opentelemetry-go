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
)

func TestCompareVersions(t *testing.T) {
	invalid := `🌭`
	v13 := "https://opentelemetry.io/schemas/1.3.0"
	v12 := "https://opentelemetry.io/schemas/1.2.0"

	var targetErr *errInvalidVer
	t.Run("InvalidA", func(t *testing.T) {
		_, err := CompareVersions(invalid, v13)
		assert.ErrorAs(t, err, &targetErr)
	})
	t.Run("InvalidB", func(t *testing.T) {
		_, err := CompareVersions(v13, invalid)
		assert.ErrorAs(t, err, &targetErr)
	})
	t.Run("Equal", func(t *testing.T) {
		c, err := CompareVersions(v12, v12)
		assert.NoError(t, err)
		assert.Equal(t, EqualTo, c)
	})
	t.Run("LessThan", func(t *testing.T) {
		c, err := CompareVersions(v12, v13)
		assert.NoError(t, err)
		assert.Equal(t, LessThan, c)
	})
	t.Run("GreaterThan", func(t *testing.T) {
		c, err := CompareVersions(v13, v12)
		assert.NoError(t, err)
		assert.Equal(t, GreaterThan, c)
	})
}
