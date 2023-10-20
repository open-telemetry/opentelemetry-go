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

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestNewTraceConfig(t *testing.T) {
	globalTP := noop.NewTracerProvider()
	customTP := noop.NewTracerProvider()
	otel.SetTracerProvider(globalTP)
	for _, tc := range []struct {
		desc     string
		opts     []TraceOption
		expected traceConfig
	}{
		{
			desc: "default",
			expected: traceConfig{
				tp: globalTP,
			},
		},
		{
			desc: "overridden",
			opts: []TraceOption{
				WithTracerProvider(customTP),
			},
			expected: traceConfig{
				tp: customTP,
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := newTraceConfig(tc.opts)
			assert.Equal(t, tc.expected, cfg)
		})
	}
}
