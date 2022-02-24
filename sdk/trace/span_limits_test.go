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
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/sdk/internal/env"
	"go.opentelemetry.io/otel/trace"
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

			var opts []TracerProviderOption
			if test.opt != nil {
				opts = append(opts, WithSpanLimits(*test.opt))
			}

			assert.Equal(t, test.want, NewTracerProvider(opts...).spanLimits)
		})
	}
}

type recorder []ReadOnlySpan

func (r *recorder) OnStart(context.Context, ReadWriteSpan) {}
func (r *recorder) OnEnd(s ReadOnlySpan)                   { *r = append(*r, s) }
func (r *recorder) ForceFlush(context.Context) error       { return nil }
func (r *recorder) Shutdown(context.Context) error         { return nil }

func testSpanLimits(t *testing.T, limits SpanLimits) ReadOnlySpan {
	rec := new(recorder)
	tp := NewTracerProvider(WithSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("testSpanLimits")

	ctx := context.Background()
	_, span := tracer.Start(ctx, "span-name", trace.WithLinks(
		trace.Link{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: [16]byte{0x01},
				SpanID:  [8]byte{0x01},
			}),
			Attributes: []attribute.KeyValue{
				attribute.Bool("one", true),
				attribute.Bool("two", true),
			},
		},
		trace.Link{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: [16]byte{0x01},
				SpanID:  [8]byte{0x01},
			}),
			Attributes: []attribute.KeyValue{
				attribute.Bool("one", true),
				attribute.Bool("two", true),
			},
		},
	))
	span.SetAttributes(
		attribute.String("string", "abc"),
		attribute.StringSlice("stringSlice", []string{"abc", "def"}),
	)
	span.AddEvent("event 1", trace.WithAttributes(
		attribute.Bool("one", true),
		attribute.Bool("two", true),
	))
	span.AddEvent("event 2", trace.WithAttributes(
		attribute.Bool("one", true),
		attribute.Bool("two", true),
	))
	span.End()
	tp.Shutdown(ctx)

	require.Len(t, *rec, 1, "exported spans")
	return (*rec)[0]
}

func TestSpanLimits(t *testing.T) {
	t.Run("AttributeValueLengthLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		attrs := testSpanLimits(t, limits).Attributes()
		require.Contains(t, attrs, attribute.String("string", "abc"))
		require.Contains(t, attrs, attribute.StringSlice("stringSlice", []string{"abc", "def"}))

		limits.AttributeValueLengthLimit = 2
		attrs = testSpanLimits(t, limits).Attributes()
		// Ensure string and string slice attributes are truncated.
		assert.Contains(t, attrs, attribute.String("string", "ab"))
		assert.Contains(t, attrs, attribute.StringSlice("stringSlice", []string{"ab", "de"}))
	})

	t.Run("AttributeCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		require.Len(t, testSpanLimits(t, limits).Attributes(), 2)

		limits.AttributeCountLimit = 1
		assert.Len(t, testSpanLimits(t, limits).Attributes(), 1)
	})

	t.Run("EventCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		require.Len(t, testSpanLimits(t, limits).Events(), 2)

		limits.EventCountLimit = 1
		assert.Len(t, testSpanLimits(t, limits).Events(), 1)
	})

	t.Run("AttributePerEventCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		for _, e := range testSpanLimits(t, limits).Events() {
			require.Len(t, e.Attributes, 2)
		}

		limits.AttributePerEventCountLimit = 1
		for _, e := range testSpanLimits(t, limits).Events() {
			require.Len(t, e.Attributes, 1)
		}
	})

	t.Run("LinkCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		require.Len(t, testSpanLimits(t, limits).Links(), 2)

		limits.LinkCountLimit = 1
		assert.Len(t, testSpanLimits(t, limits).Links(), 1)
	})

	t.Run("AttributePerLinkCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		for _, l := range testSpanLimits(t, limits).Links() {
			require.Len(t, l.Attributes, 2)
		}

		limits.AttributePerLinkCountLimit = 1
		for _, l := range testSpanLimits(t, limits).Links() {
			require.Len(t, l.Attributes, 1)
		}
	})
}
