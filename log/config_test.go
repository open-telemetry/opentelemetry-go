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

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerConfig(t *testing.T) {
	v1 := "semver:0.0.1"
	v2 := "semver:1.0.0"
	schemaURL := "https://opentelemetry.io/schemas/1.2.0"
	tests := []struct {
		options  []LoggerOption
		expected LoggerConfig
	}{
		{
			// No non-zero-values should be set.
			[]LoggerOption{},
			LoggerConfig{},
		},
		{
			[]LoggerOption{
				WithInstrumentationVersion(v1),
			},
			LoggerConfig{
				instrumentationVersion: v1,
			},
		},
		{
			[]LoggerOption{
				// Multiple calls should overwrite.
				WithInstrumentationVersion(v1),
				WithInstrumentationVersion(v2),
			},
			LoggerConfig{
				instrumentationVersion: v2,
			},
		},
		{
			[]LoggerOption{
				WithSchemaURL(schemaURL),
			},
			LoggerConfig{
				schemaURL: schemaURL,
			},
		},
	}
	for _, test := range tests {
		config := NewLoggerConfig(test.options...)
		assert.Equal(t, test.expected, config)
	}
}

// Save benchmark results to a file level var to avoid the compiler optimizing
// away the actual work.
var (
	loggerConfig LoggerConfig
)

func BenchmarkNewTracerConfig(b *testing.B) {
	opts := []LoggerOption{
		WithInstrumentationVersion("testing verion"),
		WithSchemaURL("testing URL"),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerConfig = NewLoggerConfig(opts...)
	}
}
