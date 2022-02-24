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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/sdk/internal/env"
)

func TestSettingSpanLimits(t *testing.T) {
	envLimits := func(val string) map[string]string {
		return map[string]string{
			env.SpanAttributeValueLengthKey: val,
			env.SpanEventCountKey:           val,
			env.SpanAttributeCountKey:       val,
			env.SpanLinkCountKey:            val,
			env.SpanEventAttributeCountKey:  val,
			env.SpanLinkAttributeCountKey:   val,
		}
	}

	limits := func(n int) *SpanLimits {
		lims := NewSpanLimits()
		lims.AttributeValueLengthLimit = n
		lims.AttributeCountLimit = n
		lims.EventCountLimit = n
		lims.LinkCountLimit = n
		lims.AttributePerEventCountLimit = n
		lims.AttributePerLinkCountLimit = n
		return &lims
	}

	tests := []struct {
		name string
		env  map[string]string
		opt  *SpanLimits
		want SpanLimits
	}{
		{
			name: "defaults",
			want: NewSpanLimits(),
		},
		{
			name: "env",
			env:  envLimits("42"),
			want: *(limits(42)),
		},
		{
			name: "opt",
			opt:  limits(42),
			want: *(limits(42)),
		},
		{
			name: "opt-override",
			env:  envLimits("-2"),
			// Option take priority.
			opt:  limits(43),
			want: *(limits(43)),
		},
		{
			name: "env(unlimited)",
			// OTel spec says negative SpanLinkAttributeCountKey is invalid,
			// but since we will revert to the default (unlimited) which uses
			// negative values to signal this than this value is expected to
			// pass through.
			env:  envLimits("-1"),
			want: *(limits(-1)),
		},
		{
			name: "opt(unlimited)",
			opt:  limits(-1),
			want: *(limits(-1)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.env != nil {
				es := ottest.NewEnvStore()
				t.Cleanup(func() { require.NoError(t, es.Restore()) })
				for k, v := range test.env {
					es.Record(k)
					require.NoError(t, os.Setenv(k, v))
				}
			}

			var tp *TracerProvider
			if test.opt != nil {
				tp = NewTracerProvider(WithSpanLimits(*test.opt))
			} else {
				tp = NewTracerProvider()
			}

			assert.Equal(t, test.want, tp.spanLimits)
		})
	}
}
