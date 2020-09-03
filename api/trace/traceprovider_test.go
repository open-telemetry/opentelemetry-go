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

package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracerConfigure(t *testing.T) {
	v1 := "semver:0.0.1"
	v2 := "semver:1.0.0"
	tests := []struct {
		options  []TracerOption
		expected *TracerConfig
	}{
		{
			// No non-zero-values should be set.
			[]TracerOption{},
			new(TracerConfig),
		},
		{
			[]TracerOption{
				WithInstrumentationVersion(v1),
			},
			&TracerConfig{
				InstrumentationVersion: v1,
			},
		},
		{
			[]TracerOption{
				// Multiple calls should overwrite.
				WithInstrumentationVersion(v1),
				WithInstrumentationVersion(v2),
			},
			&TracerConfig{
				InstrumentationVersion: v2,
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, TracerConfigure(test.options))
	}
}
