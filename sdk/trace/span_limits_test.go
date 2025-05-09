// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
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
		name   string
		env    map[string]string
		opt    *SpanLimits
		rawOpt *SpanLimits
		want   SpanLimits
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
			name:   "raw-opt",
			rawOpt: limits(42),
			want:   *(limits(42)),
		},
		{
			name: "opt-override",
			env:  envLimits("-2"),
			// Option take priority.
			opt:  limits(43),
			want: *(limits(43)),
		},
		{
			name: "raw-opt-override",
			env:  envLimits("-2"),
			// Option take priority.
			rawOpt: limits(43),
			want:   *(limits(43)),
		},
		{
			name:   "last-opt-wins",
			opt:    limits(-2),
			rawOpt: limits(-3),
			want:   *(limits(-3)),
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
			// Corrects to defaults.
			opt:  limits(-1),
			want: NewSpanLimits(),
		},
		{
			name:   "raw-opt(unlimited)",
			rawOpt: limits(-1),
			want:   *(limits(-1)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.env != nil {
				for k, v := range test.env {
					t.Setenv(k, v)
				}
			}

			var opts []TracerProviderOption
			if test.opt != nil {
				opts = append(opts, WithSpanLimits(*test.opt))
			}
			if test.rawOpt != nil {
				opts = append(opts, WithRawSpanLimits(*test.rawOpt))
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
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("testSpanLimits")

	ctx := context.Background()
	a := []attribute.KeyValue{attribute.Bool("one", true), attribute.Bool("two", true)}
	l := trace.Link{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: [16]byte{0x01},
			SpanID:  [8]byte{0x01},
		}),
		Attributes: a,
	}
	_, span := tracer.Start(ctx, "span-name", trace.WithLinks(l, l))
	span.SetAttributes(
		attribute.String("string", "abc"),
		attribute.StringSlice("stringSlice", []string{"abc", "def"}),
		attribute.String("euro", "€"), // this is a 3-byte rune
	)
	span.AddEvent("event 1", trace.WithAttributes(a...))
	span.AddEvent("event 2", trace.WithAttributes(a...))
	span.End()
	require.NoError(t, tp.Shutdown(ctx))

	require.Len(t, *rec, 1, "exported spans")
	return (*rec)[0]
}

func TestSpanLimits(t *testing.T) {
	t.Run("AttributeValueLengthLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		// Unlimited.
		limits.AttributeValueLengthLimit = -1
		attrs := testSpanLimits(t, limits).Attributes()
		assert.Contains(t, attrs, attribute.String("string", "abc"))
		assert.Contains(t, attrs, attribute.StringSlice("stringSlice", []string{"abc", "def"}))
		assert.Contains(t, attrs, attribute.String("euro", "€"))

		limits.AttributeValueLengthLimit = 2
		attrs = testSpanLimits(t, limits).Attributes()
		// Ensure string and string slice attributes are truncated.
		assert.Contains(t, attrs, attribute.String("string", "ab"))
		assert.Contains(t, attrs, attribute.StringSlice("stringSlice", []string{"ab", "de"}))
		assert.Contains(t, attrs, attribute.String("euro", "€"))

		limits.AttributeValueLengthLimit = 0
		attrs = testSpanLimits(t, limits).Attributes()
		assert.Contains(t, attrs, attribute.String("string", ""))
		assert.Contains(t, attrs, attribute.StringSlice("stringSlice", []string{"", ""}))
		assert.Contains(t, attrs, attribute.String("euro", ""))
	})

	t.Run("AttributeCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		// Unlimited.
		limits.AttributeCountLimit = -1
		assert.Len(t, testSpanLimits(t, limits).Attributes(), 3)

		limits.AttributeCountLimit = 1
		assert.Len(t, testSpanLimits(t, limits).Attributes(), 1)

		// Ensure this can be disabled.
		limits.AttributeCountLimit = 0
		assert.Empty(t, testSpanLimits(t, limits).Attributes())
	})

	t.Run("EventCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		// Unlimited.
		limits.EventCountLimit = -1
		assert.Len(t, testSpanLimits(t, limits).Events(), 2)

		limits.EventCountLimit = 1
		assert.Len(t, testSpanLimits(t, limits).Events(), 1)

		// Ensure this can be disabled.
		limits.EventCountLimit = 0
		assert.Empty(t, testSpanLimits(t, limits).Events())
	})

	t.Run("AttributePerEventCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		// Unlimited.
		limits.AttributePerEventCountLimit = -1
		for _, e := range testSpanLimits(t, limits).Events() {
			assert.Len(t, e.Attributes, 2)
		}

		limits.AttributePerEventCountLimit = 1
		for _, e := range testSpanLimits(t, limits).Events() {
			assert.Len(t, e.Attributes, 1)
		}

		// Ensure this can be disabled.
		limits.AttributePerEventCountLimit = 0
		for _, e := range testSpanLimits(t, limits).Events() {
			assert.Empty(t, e.Attributes)
		}
	})

	t.Run("LinkCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		// Unlimited.
		limits.LinkCountLimit = -1
		assert.Len(t, testSpanLimits(t, limits).Links(), 2)

		limits.LinkCountLimit = 1
		assert.Len(t, testSpanLimits(t, limits).Links(), 1)

		// Ensure this can be disabled.
		limits.LinkCountLimit = 0
		assert.Empty(t, testSpanLimits(t, limits).Links())
	})

	t.Run("AttributePerLinkCountLimit", func(t *testing.T) {
		limits := NewSpanLimits()
		// Unlimited.
		limits.AttributePerLinkCountLimit = -1
		for _, l := range testSpanLimits(t, limits).Links() {
			assert.Len(t, l.Attributes, 2)
		}

		limits.AttributePerLinkCountLimit = 1
		for _, l := range testSpanLimits(t, limits).Links() {
			assert.Len(t, l.Attributes, 1)
		}

		// Ensure this can be disabled.
		limits.AttributePerLinkCountLimit = 0
		for _, l := range testSpanLimits(t, limits).Links() {
			assert.Empty(t, l.Attributes)
		}
	})
}
