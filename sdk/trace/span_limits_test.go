// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/internal/env"
	semconv "go.opentelemetry.io/otel/semconv/v1.45.0"
	"go.opentelemetry.io/otel/trace"
)

func TestSettingSpanLimits(t *testing.T) {
	envLimits := func(val string) map[string]string {
		return map[string]string{
			env.SpanAttributeValueLengthKey:         val,
			"OTEL_SPAN_ATTRIBUTE_VALUE_DEPTH_LIMIT": val,
			env.SpanEventCountKey:                   val,
			env.SpanAttributeCountKey:               val,
			env.SpanLinkCountKey:                    val,
			env.SpanEventAttributeCountKey:          val,
			env.SpanLinkAttributeCountKey:           val,
		}
	}

	limits := func(n int) *SpanLimits {
		lims := NewSpanLimits()
		lims.AttributeValueLengthLimit = n
		lims.AttributeValueDepthLimit = n
		lims.AttributeCountLimit = n
		lims.EventCountLimit = n
		lims.LinkCountLimit = n
		lims.AttributePerEventCountLimit = n
		lims.AttributePerLinkCountLimit = n
		return &lims
	}
	envWant := func(n int) SpanLimits {
		lims := *limits(n)
		lims.AttributeValueDepthLimit = DefaultAttributeValueDepthLimit
		return lims
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
			want: envWant(42),
		},
		{
			name: "opt",
			opt:  limits(42),
			want: *limits(42),
		},
		{
			name:   "raw-opt",
			rawOpt: limits(42),
			want:   *limits(42),
		},
		{
			name:   "raw-opt-zero-depth",
			rawOpt: limits(0),
			want: SpanLimits{
				AttributeValueDepthLimit: DefaultAttributeValueDepthLimit,
			},
		},
		{
			name: "opt-override",
			env:  envLimits("-2"),
			// Option take priority.
			opt:  limits(43),
			want: *limits(43),
		},
		{
			name: "raw-opt-override",
			env:  envLimits("-2"),
			// Option take priority.
			rawOpt: limits(43),
			want:   *limits(43),
		},
		{
			name:   "last-opt-wins",
			opt:    limits(-2),
			rawOpt: limits(-3),
			want:   *limits(-3),
		},
		{
			name: "env(unlimited)",
			// OTel spec says negative SpanLinkAttributeCountKey is invalid,
			// but since we will revert to the default (unlimited) which uses
			// negative values to signal this than this value is expected to
			// pass through.
			env:  envLimits("-1"),
			want: envWant(-1),
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
			want:   *limits(-1),
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

func TestAttributeValueDepthLimitOptionPrecedence(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeValueDepthLimit = 7

	assert.Equal(t, 3, NewTracerProvider(
		WithRawSpanLimits(limits),
		WithAttributeValueDepthLimit(3),
	).spanLimits.AttributeValueDepthLimit)

	assert.Equal(t, 7, NewTracerProvider(
		WithAttributeValueDepthLimit(3),
		WithRawSpanLimits(limits),
	).spanLimits.AttributeValueDepthLimit)

	assert.Equal(t, DefaultAttributeValueDepthLimit, NewTracerProvider(
		WithAttributeValueDepthLimit(0),
	).spanLimits.AttributeValueDepthLimit)

	limits.AttributeValueDepthLimit = 0
	assert.Equal(t, DefaultAttributeValueDepthLimit, NewTracerProvider(
		WithRawSpanLimits(limits),
	).spanLimits.AttributeValueDepthLimit)
}

type recorder []ReadOnlySpan

func (*recorder) OnStart(context.Context, ReadWriteSpan) {}
func (r *recorder) OnEnd(s ReadOnlySpan)                 { *r = append(*r, s) }
func (*recorder) ForceFlush(context.Context) error       { return nil }
func (*recorder) Shutdown(context.Context) error         { return nil }

func testSpanLimits(t *testing.T, limits SpanLimits) ReadOnlySpan {
	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("testSpanLimits")

	ctx := t.Context()
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
		attribute.String("string", "a\uFFFDb"),
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
		assert.Contains(t, attrs, attribute.String("string", "a\uFFFDb"))
		assert.Contains(t, attrs, attribute.StringSlice("stringSlice", []string{"abc", "def"}))
		assert.Contains(t, attrs, attribute.String("euro", "€"))

		limits.AttributeValueLengthLimit = 2
		attrs = testSpanLimits(t, limits).Attributes()
		// Ensure string and string slice attributes are truncated.
		assert.Contains(t, attrs, attribute.String("string", "a\uFFFD"))
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

func TestSnapshotDroppedCountsWithZeroEventAndLinkLimits(t *testing.T) {
	// With a zero event/link limit the evicted queues never store
	// anything, but the dropped counts must still be reported by the
	// snapshot handed to processors and exporters on span end.
	limits := NewSpanLimits()
	limits.EventCountLimit = 0
	limits.LinkCountLimit = 0

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("testSnapshotDroppedCounts")

	ctx := t.Context()
	_, span := tracer.Start(ctx, "span-name")
	rws, ok := span.(ReadWriteSpan)
	require.True(t, ok)

	link := trace.Link{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: [16]byte{0x01},
			SpanID:  [8]byte{0x01},
		}),
	}
	span.AddEvent("event 1")
	span.AddEvent("event 2")
	span.AddEvent("event 3")
	rws.AddLink(link)
	rws.AddLink(link)

	assert.Equal(t, 3, rws.DroppedEvents())
	assert.Equal(t, 2, rws.DroppedLinks())

	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)

	snap := (*rec)[0]
	assert.Empty(t, snap.Events())
	assert.Empty(t, snap.Links())
	assert.Equal(t, 3, snap.DroppedEvents())
	assert.Equal(t, 2, snap.DroppedLinks())
}

func testLinkSpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{0x01},
		SpanID:  [8]byte{0x01},
	})
}

// assertEventAndLinkAttributes checks that applying limit to in yields want
// for a span-start link, an added event, and an added link.
func assertEventAndLinkAttributes(t *testing.T, limit int, in, want []attribute.KeyValue) {
	t.Helper()

	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = limit

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("eventLinkLength")
	ctx := t.Context()
	sc := testLinkSpanContext()

	_, span := tracer.Start(ctx, "span",
		trace.WithLinks(trace.Link{SpanContext: sc, Attributes: in}),
	)
	span.AddEvent("event", trace.WithAttributes(in...))
	span.AddLink(trace.Link{SpanContext: sc, Attributes: in})
	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)
	got := (*rec)[0]

	require.Len(t, got.Events(), 1)
	assert.Equal(t, want, got.Events()[0].Attributes)
	assert.Zero(t, got.Events()[0].DroppedAttributeCount)

	require.Len(t, got.Links(), 2)
	assert.Equal(t, want, got.Links()[0].Attributes)
	assert.Equal(t, want, got.Links()[1].Attributes)
	assert.Zero(t, got.Links()[0].DroppedAttributeCount)
	assert.Zero(t, got.Links()[1].DroppedAttributeCount)
}

func TestEventAndLinkAttributeValueLengthLimit(t *testing.T) {
	tests := []struct {
		name string
		in   attribute.KeyValue
		want map[int]attribute.KeyValue // keyed by limit: -1, 0, 3
	}{
		{
			name: "ascii",
			in:   attribute.String("ascii", "abcdef"),
			want: map[int]attribute.KeyValue{
				-1: attribute.String("ascii", "abcdef"),
				0:  attribute.String("ascii", ""),
				3:  attribute.String("ascii", "abc"),
			},
		},
		{
			name: "euro",
			in:   attribute.String("euro", "€€€€"),
			want: map[int]attribute.KeyValue{
				-1: attribute.String("euro", "€€€€"),
				0:  attribute.String("euro", ""),
				3:  attribute.String("euro", "€€€"),
			},
		},
		{
			// "日本語" is 3 runes and 9 bytes, so a limit of 3 keeps it intact.
			name: "multi-byte under limit",
			in:   attribute.String("cjk", "日本語"),
			want: map[int]attribute.KeyValue{
				-1: attribute.String("cjk", "日本語"),
				0:  attribute.String("cjk", ""),
				3:  attribute.String("cjk", "日本語"),
			},
		},
		{
			// The invalid byte is dropped, then the remaining runes are truncated.
			name: "invalid utf-8",
			in:   attribute.String("invalid", "ab\x80cd"),
			want: map[int]attribute.KeyValue{
				-1: attribute.String("invalid", "ab\x80cd"),
				0:  attribute.String("invalid", ""),
				3:  attribute.String("invalid", "abc"),
			},
		},
		{
			name: "exact",
			in:   attribute.String("exact", "abc"),
			want: map[int]attribute.KeyValue{
				-1: attribute.String("exact", "abc"),
				0:  attribute.String("exact", ""),
				3:  attribute.String("exact", "abc"),
			},
		},
		{
			name: "string slice",
			in:   attribute.StringSlice("strs", []string{"abcdef", "xy", "€€€€"}),
			want: map[int]attribute.KeyValue{
				-1: attribute.StringSlice("strs", []string{"abcdef", "xy", "€€€€"}),
				0:  attribute.StringSlice("strs", []string{"", "", ""}),
				3:  attribute.StringSlice("strs", []string{"abc", "xy", "€€€"}),
			},
		},
		{
			name: "byte slice",
			in:   attribute.ByteSlice("bytes", []byte{1, 2, 3, 4}),
			want: map[int]attribute.KeyValue{
				-1: attribute.ByteSlice("bytes", []byte{1, 2, 3, 4}),
				0:  attribute.ByteSlice("bytes", []byte{}),
				3:  attribute.ByteSlice("bytes", []byte{1, 2, 3}),
			},
		},
		{
			name: "bool",
			in:   attribute.Bool("flag", true),
			want: map[int]attribute.KeyValue{
				-1: attribute.Bool("flag", true),
				0:  attribute.Bool("flag", true),
				3:  attribute.Bool("flag", true),
			},
		},
		{
			name: "slice",
			in: attribute.Slice("slice",
				attribute.StringValue("abcdef"),
				attribute.StringSliceValue([]string{"abcdef", "xy", "€€€€"}),
				attribute.ByteSliceValue([]byte{1, 2, 3, 4}),
				attribute.BoolValue(true),
			),
			want: map[int]attribute.KeyValue{
				-1: attribute.Slice("slice",
					attribute.StringValue("abcdef"),
					attribute.StringSliceValue([]string{"abcdef", "xy", "€€€€"}),
					attribute.ByteSliceValue([]byte{1, 2, 3, 4}),
					attribute.BoolValue(true),
				),
				0: attribute.Slice("slice",
					attribute.StringValue(""),
					attribute.StringSliceValue([]string{"", "", ""}),
					attribute.ByteSliceValue([]byte{}),
					attribute.BoolValue(true),
				),
				3: attribute.Slice("slice",
					attribute.StringValue("abc"),
					attribute.StringSliceValue([]string{"abc", "xy", "€€€"}),
					attribute.ByteSliceValue([]byte{1, 2, 3}),
					attribute.BoolValue(true),
				),
			},
		},
		{
			name: "map",
			in: attribute.Map("map",
				attribute.String("s", "abcdef"),
				attribute.StringSlice("ss", []string{"abcdef"}),
				attribute.Map("nested", attribute.String("s", "abcdef")),
				attribute.Bool("b", true),
			),
			want: map[int]attribute.KeyValue{
				-1: attribute.Map("map",
					attribute.String("s", "abcdef"),
					attribute.StringSlice("ss", []string{"abcdef"}),
					attribute.Map("nested", attribute.String("s", "abcdef")),
					attribute.Bool("b", true),
				),
				0: attribute.Map("map",
					attribute.String("s", ""),
					attribute.StringSlice("ss", []string{""}),
					attribute.Map("nested", attribute.String("s", "")),
					attribute.Bool("b", true),
				),
				3: attribute.Map("map",
					attribute.String("s", "abc"),
					attribute.StringSlice("ss", []string{"abc"}),
					attribute.Map("nested", attribute.String("s", "abc")),
					attribute.Bool("b", true),
				),
			},
		},
		{
			// Nested map keys collapse to the last value, then that value is truncated.
			name: "nested map duplicates",
			in: attribute.Map("dups",
				attribute.String("k", "abcdef"),
				attribute.String("k", "zzzzzz"),
			),
			want: map[int]attribute.KeyValue{
				-1: attribute.Map("dups", attribute.String("k", "zzzzzz")),
				0:  attribute.Map("dups", attribute.String("k", "")),
				3:  attribute.Map("dups", attribute.String("k", "zzz")),
			},
		},
	}

	for _, tt := range tests {
		for _, limit := range []int{-1, 0, 3} {
			t.Run(fmt.Sprintf("%s/limit=%d", tt.name, limit), func(t *testing.T) {
				assertEventAndLinkAttributes(t, limit,
					[]attribute.KeyValue{tt.in},
					[]attribute.KeyValue{tt.want[limit]},
				)
			})
		}
	}

	// Top-level keys are not deduplicated. Each value is still truncated.
	in := []attribute.KeyValue{
		attribute.String("dup", "abcdef"),
		attribute.String("dup", "zzzzzz"),
	}
	topLevel := map[int][]attribute.KeyValue{
		-1: {
			attribute.String("dup", "abcdef"),
			attribute.String("dup", "zzzzzz"),
		},
		0: {
			attribute.String("dup", ""),
			attribute.String("dup", ""),
		},
		3: {
			attribute.String("dup", "abc"),
			attribute.String("dup", "zzz"),
		},
	}
	for _, limit := range []int{-1, 0, 3} {
		t.Run(fmt.Sprintf("top-level duplicates/limit=%d", limit), func(t *testing.T) {
			assertEventAndLinkAttributes(t, limit, in, topLevel[limit])
		})
	}
}

func TestSpanAttributeAndRecordErrorValueLengthLimit(t *testing.T) {
	errType := "*errors.errorString"
	errMsg := "€€€€"

	tests := []struct {
		name      string
		limit     int
		spanValue string
		user      string
		userSlice []string
		errType   string
		errMsg    string
	}{
		{
			name:      "Unlimited",
			limit:     -1,
			spanValue: "abcdef",
			user:      "abcdef",
			userSlice: []string{"abcdef", "€€€€"},
			errType:   errType,
			errMsg:    errMsg,
		},
		{
			name:      "Zero",
			limit:     0,
			spanValue: "",
			user:      "",
			userSlice: []string{"", ""},
			errType:   "",
			errMsg:    "",
		},
		{
			name:      "Truncated",
			limit:     3,
			spanValue: "abc",
			user:      "abc",
			userSlice: []string{"abc", "€€€"},
			errType:   "*er",
			errMsg:    "€€€",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limits := NewSpanLimits()
			limits.AttributeValueLengthLimit = tt.limit

			rec := new(recorder)
			tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
			tracer := tp.Tracer("spanAndExceptionLength")
			ctx := t.Context()

			_, span := tracer.Start(ctx, "span",
				trace.WithAttributes(attribute.String("span", "abcdef")),
			)
			span.RecordError(errors.New(errMsg),
				trace.WithAttributes(
					attribute.String("user", "abcdef"),
					attribute.StringSlice("userSlice", []string{"abcdef", "€€€€"}),
				),
				trace.WithStackTrace(true),
			)
			span.End()
			require.NoError(t, tp.Shutdown(ctx))
			require.Len(t, *rec, 1)
			got := (*rec)[0]

			assert.Contains(t, got.Attributes(), attribute.String("span", tt.spanValue))

			require.Len(t, got.Events(), 1)
			exc := got.Events()[0]
			require.GreaterOrEqual(t, len(exc.Attributes), 5)
			assert.Equal(t, []attribute.KeyValue{
				attribute.String("user", tt.user),
				attribute.StringSlice("userSlice", tt.userSlice),
				semconv.ExceptionType(tt.errType),
				semconv.ExceptionMessage(tt.errMsg),
			}, exc.Attributes[:4])
			assert.Equal(t, semconv.ExceptionStacktraceKey, exc.Attributes[4].Key)
			stack := exc.Attributes[4].Value.AsString()
			if tt.limit < 0 {
				assert.Greater(t, utf8.RuneCountInString(stack), 3)
			} else {
				assert.Equal(t, tt.limit, utf8.RuneCountInString(stack))
			}
			assert.Zero(t, exc.DroppedAttributeCount)
		})
	}
}

func TestEventAndLinkAttributeValueLengthLimitWithCountLimit(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = 3
	limits.AttributePerEventCountLimit = 1
	limits.AttributePerLinkCountLimit = 1

	input := []attribute.KeyValue{
		attribute.String("keep", "abcdef"),
		attribute.String("drop-str", "€€€€"),
		attribute.Bool("drop-bool", true),
	}
	caller := append([]attribute.KeyValue(nil), input...)

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("eventLinkLengthCount")
	ctx := t.Context()
	sc := testLinkSpanContext()

	_, span := tracer.Start(ctx, "span",
		trace.WithLinks(trace.Link{SpanContext: sc, Attributes: caller}),
	)
	span.AddEvent("event", trace.WithAttributes(caller...))
	span.AddLink(trace.Link{SpanContext: sc, Attributes: caller})
	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)
	got := (*rec)[0]

	want := []attribute.KeyValue{attribute.String("keep", "abc")}
	require.Len(t, got.Events(), 1)
	assert.Equal(t, want, got.Events()[0].Attributes)
	assert.Equal(t, 2, got.Events()[0].DroppedAttributeCount)
	require.Len(t, got.Links(), 2)
	assert.Equal(t, want, got.Links()[0].Attributes)
	assert.Equal(t, want, got.Links()[1].Attributes)
	assert.Equal(t, 2, got.Links()[0].DroppedAttributeCount)
	assert.Equal(t, 2, got.Links()[1].DroppedAttributeCount)
	assert.Equal(t, input, caller)
}

func TestEventAndLinkAttributeValueLengthLimitDepth(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = 3
	limits.AttributeValueDepthLimit = 1

	input := []attribute.KeyValue{
		attribute.Map("nested",
			attribute.String("keep", "abcdef"),
			attribute.Map("child", attribute.String("leaf", "abcdef")),
		),
	}
	want := []attribute.KeyValue{
		attribute.Map("nested",
			attribute.String("keep", "abc"),
			attribute.KeyValue{Key: "child"},
		),
	}

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("eventLinkLengthDepth")
	ctx := t.Context()
	sc := testLinkSpanContext()

	_, span := tracer.Start(ctx, "span",
		trace.WithLinks(trace.Link{SpanContext: sc, Attributes: input}),
	)
	span.AddEvent("event", trace.WithAttributes(input...))
	span.AddLink(trace.Link{SpanContext: sc, Attributes: input})
	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)
	got := (*rec)[0]

	require.Len(t, got.Events(), 1)
	assert.Equal(t, want, got.Events()[0].Attributes)
	require.Len(t, got.Links(), 2)
	assert.Equal(t, want, got.Links()[0].Attributes)
	assert.Equal(t, want, got.Links()[1].Attributes)
}

func TestEventAndLinkAttributeValueLengthLimitPartial(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = 3

	// The truncated value is not first, so normalization has to keep the
	// unchanged prefix when it copies.
	input := []attribute.KeyValue{
		attribute.Bool("ok", true),
		attribute.String("exact", "abc"),
		attribute.String("long", "abcdef"),
		attribute.StringSlice("strs", []string{"xy", "abcdef"}),
	}
	want := []attribute.KeyValue{
		attribute.Bool("ok", true),
		attribute.String("exact", "abc"),
		attribute.String("long", "abc"),
		attribute.StringSlice("strs", []string{"xy", "abc"}),
	}
	within := []attribute.KeyValue{
		attribute.Bool("ok", true),
		attribute.String("exact", "abc"),
		attribute.StringSlice("strs", []string{"xy", "z"}),
		attribute.Int("n", 7),
	}
	wantWithin := []attribute.KeyValue{
		attribute.Bool("ok", true),
		attribute.String("exact", "abc"),
		attribute.StringSlice("strs", []string{"xy", "z"}),
		attribute.Int("n", 7),
	}
	wantInput := []attribute.KeyValue{
		attribute.Bool("ok", true),
		attribute.String("exact", "abc"),
		attribute.String("long", "abcdef"),
		attribute.StringSlice("strs", []string{"xy", "abcdef"}),
	}

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("eventLinkLengthPartial")
	ctx := t.Context()
	sc := testLinkSpanContext()

	_, span := tracer.Start(ctx, "span",
		trace.WithLinks(trace.Link{SpanContext: sc, Attributes: input}),
	)
	span.AddEvent("partial", trace.WithAttributes(input...))
	span.AddEvent("within", trace.WithAttributes(within...))
	span.AddLink(trace.Link{SpanContext: sc, Attributes: within})
	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)
	got := (*rec)[0]

	require.Len(t, got.Events(), 2)
	assert.Equal(t, want, got.Events()[0].Attributes)
	assert.Equal(t, wantWithin, got.Events()[1].Attributes)
	require.Len(t, got.Links(), 2)
	assert.Equal(t, want, got.Links()[0].Attributes)
	assert.Equal(t, wantWithin, got.Links()[1].Attributes)
	assert.Equal(t, wantInput, input)
	assert.Equal(t, wantWithin, within)
}

func TestNormAttrsWithinLimitDoesNotAllocate(t *testing.T) {
	kvs := []attribute.KeyValue{
		attribute.Bool("ok", true),
		attribute.String("s", "ab"),
		attribute.StringSlice("ss", []string{"xy", "z"}),
		attribute.Int("n", 7),
	}
	for _, limit := range []int{-1, 3} {
		allocs := testing.AllocsPerRun(1000, func() {
			out := normAttrs(kvs, DefaultAttributeValueDepthLimit, limit)
			if len(out) != len(kvs) || out[1].Value.AsString() != "ab" {
				t.Fatalf("normAttrs(%d) = %#v", limit, out)
			}
		})
		assert.Zero(t, allocs, "limit %d", limit)
	}
}

func TestEventAndLinkAttributeValueLengthLimitDoesNotMutateCaller(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = 3

	caller := []attribute.KeyValue{
		attribute.String("k", "abcdef"),
		attribute.StringSlice("ss", []string{"abcdef", "wxyz"}),
		attribute.Map("m", attribute.String("s", "abcdef")),
	}
	wantCaller := []attribute.KeyValue{
		attribute.String("k", "abcdef"),
		attribute.StringSlice("ss", []string{"abcdef", "wxyz"}),
		attribute.Map("m", attribute.String("s", "abcdef")),
	}
	wantExported := []attribute.KeyValue{
		attribute.String("k", "abc"),
		attribute.StringSlice("ss", []string{"abc", "wxy"}),
		attribute.Map("m", attribute.String("s", "abc")),
	}

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer("eventLinkLengthAlias")
	ctx := t.Context()
	sc := testLinkSpanContext()

	_, span := tracer.Start(ctx, "span",
		trace.WithLinks(trace.Link{SpanContext: sc, Attributes: caller}),
	)
	assert.Equal(t, wantCaller, caller)
	span.AddEvent("event", trace.WithAttributes(caller...))
	assert.Equal(t, wantCaller, caller)
	span.AddLink(trace.Link{SpanContext: sc, Attributes: caller})
	assert.Equal(t, wantCaller, caller)
	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)
	got := (*rec)[0]

	require.Len(t, got.Events(), 1)
	assert.Equal(t, wantExported, got.Events()[0].Attributes)
	require.Len(t, got.Links(), 2)
	assert.Equal(t, wantExported, got.Links()[0].Attributes)
	assert.Equal(t, wantExported, got.Links()[1].Attributes)
	assert.Equal(t, wantCaller, caller)
}

// oversizedStringSlice is a 1,024-element array of strings that exceed a
// length limit of 3. The depth-only normalizer does not walk it; value-length
// truncation does.
func oversizedStringSlice() (long string, values []string, attr attribute.KeyValue) {
	long = strings.Repeat("abcdefghij", 5)
	values = make([]string, 1024)
	for i := range values {
		values[i] = long
	}
	return long, values, attribute.StringSlice("arr", values)
}

func spanLimitsForAttrCounts(length, count int) SpanLimits {
	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = length
	limits.AttributePerEventCountLimit = count
	limits.AttributePerLinkCountLimit = count
	return limits
}

func TestEventAndLinkAttributeCountLimitBeforeNormalization(t *testing.T) {
	long, values, sliceAttr := oversizedStringSlice()
	sc := testLinkSpanContext()

	tests := []struct {
		name    string
		count   int
		attrs   []attribute.KeyValue
		want    []attribute.KeyValue
		dropped int
	}{
		{
			name:    "zero count drops every attribute",
			count:   0,
			attrs:   []attribute.KeyValue{attribute.String("a", long), sliceAttr, attribute.Bool("ok", true)},
			want:    nil,
			dropped: 3,
		},
		{
			name:  "positive count keeps a normalized prefix",
			count: 2,
			attrs: []attribute.KeyValue{
				attribute.String("keep", long),
				attribute.Bool("also", true),
				sliceAttr,
				attribute.String("drop", long),
			},
			want: []attribute.KeyValue{
				attribute.String("keep", "abc"),
				attribute.Bool("also", true),
			},
			dropped: 2,
		},
		{
			name:  "negative count is unlimited",
			count: -1,
			attrs: []attribute.KeyValue{
				attribute.String("s", long),
				attribute.StringSlice("ss", []string{long, "xy"}),
				attribute.Bool("ok", true),
			},
			want: []attribute.KeyValue{
				attribute.String("s", "abc"),
				attribute.StringSlice("ss", []string{"abc", "xy"}),
				attribute.Bool("ok", true),
			},
		},
		{
			name:  "count above length keeps every attribute",
			count: 8,
			attrs: []attribute.KeyValue{
				attribute.String("s", long),
				attribute.Bool("ok", true),
			},
			want: []attribute.KeyValue{
				attribute.String("s", "abc"),
				attribute.Bool("ok", true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caller := append([]attribute.KeyValue(nil), tt.attrs...)
			before := append([]attribute.KeyValue(nil), caller...)

			limits := spanLimitsForAttrCounts(3, tt.count)
			rec := new(recorder)
			tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
			tracer := tp.Tracer(t.Name())
			ctx := t.Context()

			_, span := tracer.Start(ctx, "span",
				trace.WithLinks(trace.Link{SpanContext: sc, Attributes: caller}),
			)
			assert.Equal(t, before, caller)
			span.AddEvent("event", trace.WithAttributes(caller...))
			assert.Equal(t, before, caller)
			span.AddLink(trace.Link{SpanContext: sc, Attributes: caller})
			assert.Equal(t, before, caller)
			span.End()
			require.NoError(t, tp.Shutdown(ctx))
			require.Len(t, *rec, 1)
			got := (*rec)[0]

			require.Len(t, got.Events(), 1)
			assert.Equal(t, tt.want, got.Events()[0].Attributes)
			assert.Equal(t, tt.dropped, got.Events()[0].DroppedAttributeCount)
			require.Len(t, got.Links(), 2)
			for _, link := range got.Links() {
				assert.Equal(t, tt.want, link.Attributes)
				assert.Equal(t, tt.dropped, link.DroppedAttributeCount)
			}
			assert.Equal(t, before, caller)
			// The discarded string slice is a shared value. Truncating it in
			// place would replace these elements with "abc".
			assert.Equal(t, values, sliceAttr.Value.AsStringSlice())
		})
	}

	t.Run("zero count with no attributes", func(t *testing.T) {
		limits := spanLimitsForAttrCounts(3, 0)
		rec := new(recorder)
		tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
		tracer := tp.Tracer(t.Name())
		ctx := t.Context()
		_, span := tracer.Start(ctx, "span")
		span.AddEvent("event")
		span.AddLink(trace.Link{SpanContext: sc})
		span.End()
		require.NoError(t, tp.Shutdown(ctx))
		require.Len(t, *rec, 1)
		got := (*rec)[0]
		require.Len(t, got.Events(), 1)
		assert.Empty(t, got.Events()[0].Attributes)
		assert.Zero(t, got.Events()[0].DroppedAttributeCount)
		require.Len(t, got.Links(), 1)
		assert.Empty(t, got.Links()[0].Attributes)
		assert.Zero(t, got.Links()[0].DroppedAttributeCount)
	})

	// Retained map values are still deduplicated, depth-limited, and
	// truncated. The discarded suffix is not.
	t.Run("retained prefix still normalized", func(t *testing.T) {
		limits := spanLimitsForAttrCounts(3, 1)
		limits.AttributeValueDepthLimit = 1
		caller := []attribute.KeyValue{
			attribute.Map("keep",
				attribute.String("s", long),
				attribute.String("s", "zzzzzz"),
				attribute.Map("child", attribute.String("leaf", long)),
			),
			sliceAttr,
		}
		want := []attribute.KeyValue{
			attribute.Map("keep",
				attribute.KeyValue{Key: "child"},
				attribute.String("s", "zzz"),
			),
		}

		rec := new(recorder)
		tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
		tracer := tp.Tracer(t.Name())
		ctx := t.Context()
		_, span := tracer.Start(ctx, "span",
			trace.WithLinks(trace.Link{SpanContext: sc, Attributes: caller}),
		)
		span.AddEvent("event", trace.WithAttributes(caller...))
		span.AddLink(trace.Link{SpanContext: sc, Attributes: caller})
		span.End()
		require.NoError(t, tp.Shutdown(ctx))
		require.Len(t, *rec, 1)
		got := (*rec)[0]

		require.Len(t, got.Events(), 1)
		assert.Equal(t, want, got.Events()[0].Attributes)
		assert.Equal(t, 1, got.Events()[0].DroppedAttributeCount)
		require.Len(t, got.Links(), 2)
		assert.Equal(t, want, got.Links()[0].Attributes)
		assert.Equal(t, want, got.Links()[1].Attributes)
		assert.Equal(t, values, caller[1].Value.AsStringSlice())
		kept := caller[0].Value.AsMap()
		require.Len(t, kept, 3)
		assert.Equal(t, attribute.Key("child"), kept[0].Key)
		child := kept[0].Value.AsMap()
		require.Len(t, child, 1)
		assert.Equal(t, long, child[0].Value.AsString())
		assert.Equal(t, long, kept[1].Value.AsString())
		assert.Equal(t, "zzzzzz", kept[2].Value.AsString())
	})
}

func TestEventAndLinkUnchangedPrefixAliasesCallerSlice(t *testing.T) {
	_, _, sliceAttr := oversizedStringSlice()
	caller := []attribute.KeyValue{
		attribute.String("keep", "ab"),
		attribute.Bool("also", true),
		sliceAttr,
	}
	before := append([]attribute.KeyValue(nil), caller...)

	limits := spanLimitsForAttrCounts(3, 2)
	tp := NewTracerProvider(WithRawSpanLimits(limits))
	tracer := tp.Tracer(t.Name())
	ctx := t.Context()
	sc := testLinkSpanContext()

	_, span := tracer.Start(ctx, "span",
		trace.WithLinks(trace.Link{SpanContext: sc, Attributes: caller}),
	)
	span.AddEvent("event", trace.WithAttributes(caller...))
	span.AddLink(trace.Link{SpanContext: sc, Attributes: caller})

	rs := span.(*recordingSpan)
	require.Len(t, rs.links.queue, 2)
	require.Len(t, rs.events.queue, 1)
	// AddLink and links passed to Tracer.Start store the retained prefix in
	// place when normalization changes nothing. AddEvent copies attributes
	// through EventConfig first, so it must not alias the caller either.
	assert.Same(t, &caller[0], &rs.links.queue[0].Attributes[0])
	assert.Same(t, &caller[0], &rs.links.queue[1].Attributes[0])
	assert.Equal(t, []attribute.KeyValue{
		attribute.String("keep", "ab"),
		attribute.Bool("also", true),
	}, rs.links.queue[0].Attributes)
	assert.Equal(t, rs.links.queue[0].Attributes, rs.links.queue[1].Attributes)
	assert.Equal(t, 1, rs.links.queue[0].DroppedAttributeCount)
	assert.Equal(t, 1, rs.links.queue[1].DroppedAttributeCount)
	assert.Equal(t, []attribute.KeyValue{
		attribute.String("keep", "ab"),
		attribute.Bool("also", true),
	}, rs.events.queue[0].Attributes)
	assert.Equal(t, 1, rs.events.queue[0].DroppedAttributeCount)
	assert.NotSame(t, &caller[0], &rs.events.queue[0].Attributes[0])
	assert.Equal(t, before, caller)

	// Replacing a caller element is visible through the aliased link prefix.
	// The SDK itself does not write into that backing array.
	caller[0] = attribute.String("keep", "zz")
	assert.Equal(t, "zz", rs.links.queue[0].Attributes[0].Value.AsString())
	assert.Equal(t, "ab", rs.events.queue[0].Attributes[0].Value.AsString())

	span.End()
	require.NoError(t, tp.Shutdown(ctx))
}

func TestEventAndLinkCountLimitSkipsDiscardedAllocations(t *testing.T) {
	long, _, sliceAttr := oversizedStringSlice()
	dropped := []attribute.KeyValue{attribute.String("a", long), sliceAttr, attribute.Bool("ok", true)}
	unchangedPrefix := []attribute.KeyValue{
		attribute.String("keep", "ab"),
		attribute.Bool("also", true),
		sliceAttr,
	}
	truncatedPrefix := []attribute.KeyValue{
		attribute.String("keep", long),
		sliceAttr,
	}

	// A finite length limit must not allocate more than an unlimited one when
	// the count limit drops the values truncation would otherwise walk.
	eventZeroUnlimited, linkZeroUnlimited := measureEventLinkAllocs(t, spanLimitsForAttrCounts(-1, 0), dropped)
	eventZeroFinite, linkZeroFinite := measureEventLinkAllocs(t, spanLimitsForAttrCounts(3, 0), dropped)
	assert.Equal(t, eventZeroUnlimited, eventZeroFinite, "AddEvent count 0")
	assert.Equal(t, linkZeroUnlimited, linkZeroFinite, "AddLink count 0")
	assert.Zero(t, linkZeroUnlimited, "AddLink measurement allocated")
	assert.Zero(t, linkZeroFinite, "AddLink count 0")

	_, linkUnchanged := measureEventLinkAllocs(t, spanLimitsForAttrCounts(3, 2), unchangedPrefix)
	assert.Zero(t, linkUnchanged, "AddLink unchanged retained prefix")

	_, linkTruncated := measureEventLinkAllocs(t, spanLimitsForAttrCounts(3, 1), truncatedPrefix)
	_, linkKept := measureEventLinkAllocs(t, spanLimitsForAttrCounts(3, -1), truncatedPrefix)
	assert.Less(t, linkTruncated, linkKept, "discarded string slice was normalized")
}

func measureEventLinkAllocs(
	t *testing.T,
	limits SpanLimits,
	attrs []attribute.KeyValue,
) (eventAllocs, linkAllocs float64) {
	t.Helper()

	limits.EventCountLimit = 1
	limits.LinkCountLimit = 1
	tp := NewTracerProvider(WithRawSpanLimits(limits))
	ctx := t.Context()
	_, span := tp.Tracer(t.Name()).Start(ctx, "span")
	t.Cleanup(func() {
		span.End()
		//nolint:usetesting // The test context is canceled before Cleanup.
		require.NoError(t, tp.Shutdown(context.Background()))
	})

	link := trace.Link{SpanContext: testLinkSpanContext(), Attributes: attrs}
	span.AddEvent("prime", trace.WithAttributes(attrs...))
	span.AddEvent("prime", trace.WithAttributes(attrs...))
	span.AddLink(link)
	span.AddLink(link)

	eventAllocs = testing.AllocsPerRun(100, func() {
		span.AddEvent("event", trace.WithAttributes(attrs...))
	})
	linkAllocs = testing.AllocsPerRun(100, func() {
		span.AddLink(link)
	})
	return eventAllocs, linkAllocs
}
