// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"errors"
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

// eventLinkLengthAttrs is the attribute set applied to events and links.
// Nested map keys are deduplicated (last value wins) before length truncation,
// matching span attributes. Top-level keys are not deduplicated.
func eventLinkLengthAttrs(limit int) []attribute.KeyValue {
	ascii, euro, cjk, invalid, exact := "abcdef", "€€€€", "日本語", "ab\x80cd", "abc"
	strs := []string{"abcdef", "xy", "€€€€"}
	raw := []byte{1, 2, 3, 4}
	dupA, dupB := "abcdef", "zzzzzz"

	switch limit {
	case -1:
	case 0:
		ascii, euro, cjk, invalid, exact = "", "", "", "", ""
		strs = []string{"", "", ""}
		raw = []byte{}
		dupA, dupB = "", ""
	case 3:
		// "日本語" is 3 runes and 9 bytes, so a limit of 3 keeps it intact.
		// "abc" is exactly at the limit. "xy" is shorter than the limit.
		ascii, euro, invalid = "abc", "€€€", "abc"
		strs = []string{"abc", "xy", "€€€"}
		raw = []byte{1, 2, 3}
		dupA, dupB = "abc", "zzz"
	default:
		panic("eventLinkLengthAttrs limit must be -1, 0, or 3")
	}

	return []attribute.KeyValue{
		attribute.String("ascii", ascii),
		attribute.String("euro", euro),
		attribute.String("cjk", cjk),
		attribute.String("invalid", invalid),
		attribute.String("exact", exact),
		attribute.StringSlice("strs", strs),
		attribute.ByteSlice("bytes", raw),
		attribute.Bool("flag", true),
		attribute.Slice("slice",
			attribute.StringValue(ascii),
			attribute.StringSliceValue(strs),
			attribute.ByteSliceValue(raw),
			attribute.BoolValue(true),
		),
		attribute.Map("map",
			attribute.String("s", ascii),
			attribute.StringSlice("ss", []string{ascii}),
			attribute.Map("nested", attribute.String("s", ascii)),
			attribute.Bool("b", true),
		),
		// Nested duplicates collapse to the last value, then that value is
		// truncated. Both top-level duplicates are kept.
		attribute.Map("dups",
			attribute.String("k", dupB),
		),
		attribute.String("dup", dupA),
		attribute.String("dup", dupB),
	}
}

func eventLinkLengthInput() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("ascii", "abcdef"),
		attribute.String("euro", "€€€€"),
		attribute.String("cjk", "日本語"),
		attribute.String("invalid", "ab\x80cd"),
		attribute.String("exact", "abc"),
		attribute.StringSlice("strs", []string{"abcdef", "xy", "€€€€"}),
		attribute.ByteSlice("bytes", []byte{1, 2, 3, 4}),
		attribute.Bool("flag", true),
		attribute.Slice("slice",
			attribute.StringValue("abcdef"),
			attribute.StringSliceValue([]string{"abcdef", "xy", "€€€€"}),
			attribute.ByteSliceValue([]byte{1, 2, 3, 4}),
			attribute.BoolValue(true),
		),
		attribute.Map("map",
			attribute.String("s", "abcdef"),
			attribute.StringSlice("ss", []string{"abcdef"}),
			attribute.Map("nested", attribute.String("s", "abcdef")),
			attribute.Bool("b", true),
		),
		attribute.Map("dups",
			attribute.String("k", "abcdef"),
			attribute.String("k", "zzzzzz"),
		),
		attribute.String("dup", "abcdef"),
		attribute.String("dup", "zzzzzz"),
	}
}

func TestEventAndLinkAttributeValueLengthLimit(t *testing.T) {
	errType := "*errors.errorString"
	errMsg := "€€€€"

	tests := []struct {
		name    string
		limit   int
		errType string
		errMsg  string
	}{
		{name: "Unlimited", limit: -1, errType: errType, errMsg: errMsg},
		{name: "Zero", limit: 0, errType: "", errMsg: ""},
		{name: "Truncated", limit: 3, errType: "*er", errMsg: "€€€"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limits := NewSpanLimits()
			limits.AttributeValueLengthLimit = tt.limit
			input := eventLinkLengthInput()
			want := eventLinkLengthAttrs(tt.limit)

			rec := new(recorder)
			tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
			tracer := tp.Tracer("eventLinkLength")
			ctx := t.Context()
			sc := testLinkSpanContext()

			_, span := tracer.Start(ctx, "span",
				trace.WithAttributes(attribute.String("span", "abcdef")),
				trace.WithLinks(trace.Link{SpanContext: sc, Attributes: input}),
			)
			span.AddEvent("event", trace.WithAttributes(input...))
			span.AddLink(trace.Link{SpanContext: sc, Attributes: input})
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

			spanValue := "abcdef"
			if tt.limit == 0 {
				spanValue = ""
			} else if tt.limit > 0 {
				spanValue = "abc"
			}
			assert.Contains(t, got.Attributes(), attribute.String("span", spanValue))

			require.Len(t, got.Events(), 2)
			assert.Equal(t, want, got.Events()[0].Attributes)
			assert.Zero(t, got.Events()[0].DroppedAttributeCount)

			require.Len(t, got.Links(), 2)
			assert.Equal(t, want, got.Links()[0].Attributes)
			assert.Equal(t, want, got.Links()[1].Attributes)
			assert.Zero(t, got.Links()[0].DroppedAttributeCount)
			assert.Zero(t, got.Links()[1].DroppedAttributeCount)

			user, userSlice := "abcdef", []string{"abcdef", "€€€€"}
			if tt.limit == 0 {
				user, userSlice = "", []string{"", ""}
			} else if tt.limit > 0 {
				user, userSlice = "abc", []string{"abc", "€€€"}
			}
			exc := got.Events()[1]
			require.GreaterOrEqual(t, len(exc.Attributes), 5)
			assert.Equal(t, []attribute.KeyValue{
				attribute.String("user", user),
				attribute.StringSlice("userSlice", userSlice),
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
