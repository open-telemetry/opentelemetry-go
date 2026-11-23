// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

// emitOne emits r through a Logger built with opts and returns the single
// Record the SDK produced.
func emitOne(t *testing.T, r log.Record, opts ...LoggerProviderOption) Record {
	t.Helper()
	p := newProcessor("capture")
	opts = append(opts, WithProcessor(p))
	NewLoggerProvider(opts...).Logger(t.Name()).Emit(t.Context(), r)
	require.Len(t, p.records, 1)
	return p.records[0]
}

func recordAttrs(r Record) []attribute.KeyValue {
	var out []attribute.KeyValue
	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		out = append(out, kv)
		return true
	})
	return out
}

func apiRecord(attrs ...attribute.KeyValue) log.Record {
	var r log.Record
	r.AddAttributes(attrs...)
	return r
}

// keys returns n distinct string attributes named k0..k(n-1).
func keys(n int) []attribute.KeyValue {
	out := make([]attribute.KeyValue, n)
	for i := range out {
		out[i] = attribute.String("k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
	}
	return out
}

func TestLoggerEmitAttributeDeduplication(t *testing.T) {
	a1, a2, a3 := attribute.String("a", "1"), attribute.String("a", "2"), attribute.String("a", "3")
	b := attribute.Int("b", 1)
	c := attribute.Bool("c", true)

	// 7 distinct keys spans the inline front array and the back slice.
	seven := keys(7)
	withDups := slices.Concat(seven, []attribute.KeyValue{
		attribute.String("k0", "front-updated"),
		attribute.String("k6", "back-updated"),
	})
	wantSeven := slices.Clone(seven)
	wantSeven[0] = attribute.String("k0", "front-updated")
	wantSeven[6] = attribute.String("k6", "back-updated")

	many := keys(40)
	manyDups := slices.Concat(many, many[:20])

	tests := []struct {
		name        string
		opts        []LoggerProviderOption
		in          []attribute.KeyValue
		want        []attribute.KeyValue
		wantDropped int
	}{
		{
			name: "NoAttributes",
		},
		{
			name: "Single",
			in:   []attribute.KeyValue{a1},
			want: []attribute.KeyValue{a1},
		},
		{
			name: "UniqueOrderPreserved",
			in:   keys(12),
			want: keys(12),
		},
		{
			name: "LastValueWinsAtFirstPosition",
			in:   []attribute.KeyValue{a1, b, a2, c, a3},
			want: []attribute.KeyValue{a3, b, c},
		},
		{
			name: "DuplicatesAcrossFrontAndBack",
			in:   withDups,
			want: wantSeven,
		},
		{
			name: "ManyDuplicates",
			in:   manyDups,
			want: many,
		},
		{
			name:        "CountLimitDuplicateOfRetainedKey",
			opts:        []LoggerProviderOption{WithAttributeCountLimit(2)},
			in:          []attribute.KeyValue{a1, b, c, a2},
			want:        []attribute.KeyValue{a2, b},
			wantDropped: 1,
		},
		{
			name: "CountLimitManyDuplicates",
			opts: []LoggerProviderOption{WithAttributeCountLimit(10)},
			in:   manyDups,
			want: many[:10],
			// 30 keys over the limit, then 10 repeats of keys k10-k19 that
			// were not retained.
			wantDropped: 40,
		},
		{
			name: "ValueLengthLimitAppliedToOverwrite",
			opts: []LoggerProviderOption{WithAttributeValueLengthLimit(3)},
			in: []attribute.KeyValue{
				attribute.String("a", "x"),
				b,
				attribute.String("a", "toolong"),
			},
			want: []attribute.KeyValue{attribute.String("a", "too"), b},
		},
		{
			name: "AllowKeyDuplication",
			opts: []LoggerProviderOption{WithAllowKeyDuplication()},
			in:   []attribute.KeyValue{a1, b, a2},
			want: []attribute.KeyValue{a1, b, a2},
		},
		{
			name: "InvalidKeysDropped",
			in:   []attribute.KeyValue{a1, attribute.String("", "x"), a2},
			want: []attribute.KeyValue{a2},
		},
		{
			name: "MapValueDeduplicated",
			in: []attribute.KeyValue{
				attribute.Map("m", attribute.String("x", "1"), attribute.String("x", "2")),
			},
			want: []attribute.KeyValue{
				attribute.Map("m", attribute.String("x", "2")),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := emitOne(t, apiRecord(tc.in...), tc.opts...)
			assert.Equal(t, tc.want, recordAttrs(got))
			assert.Equal(t, len(tc.want), got.AttributesLen())
			assert.Equal(t, tc.wantDropped, got.DroppedAttributes())
		})
	}
}

// TestLoggerEmitAttributeDeduplicationDroppedCount pins how duplicates of a
// key that was itself dropped by the count limit are counted: every
// occurrence is counted, as when the attributes are added to a Record one at
// a time.
func TestLoggerEmitAttributeDeduplicationDroppedCount(t *testing.T) {
	in := []attribute.KeyValue{
		attribute.String("a", "1"),
		attribute.String("b", "1"),
		attribute.String("b", "2"),
		attribute.String("b", "3"),
	}
	got := emitOne(t, apiRecord(in...), WithAttributeCountLimit(1))
	assert.Equal(t, []attribute.KeyValue{attribute.String("a", "1")}, recordAttrs(got))
	assert.Equal(t, 3, got.DroppedAttributes())
}

func TestLoggerEmitDuplicateExceptionAttributes(t *testing.T) {
	r := apiRecord(
		attribute.String(string(exceptionMessageKey), "first"),
		attribute.String(string(exceptionMessageKey), "second"),
	)
	r.SetErr(errors.New("derived"))

	got := recordAttrs(emitOne(t, r))
	require.Len(t, got, 2)
	assert.Equal(t, attribute.String(string(exceptionMessageKey), "second"), got[0])
	assert.Equal(t, exceptionTypeKey, got[1].Key)
}

func TestLoggerEmitDoesNotMutateInputAttributes(t *testing.T) {
	in := slices.Concat(keys(10), []attribute.KeyValue{
		attribute.String("k1", "dup"),
		attribute.String("k8", "dup"),
		attribute.String("k9", "toolong"),
	})
	r := apiRecord(in...)
	before := slices.Collect(func(yield func(attribute.KeyValue) bool) {
		r.WalkAttributes(yield)
	})

	_ = emitOne(t, r, WithAttributeValueLengthLimit(2))

	after := slices.Collect(func(yield func(attribute.KeyValue) bool) {
		r.WalkAttributes(yield)
	})
	assert.Equal(t, before, after)
	assert.Equal(t, in, after)
}

// referenceRecord returns the Record that results from adding attrs to an
// empty Record one at a time with AddAttributes. Logger.Emit must produce the
// same attributes and dropped count.
func referenceRecord(attrs []attribute.KeyValue, limit int) Record {
	r := Record{attributeCountLimit: limit, attributeValueLengthLimit: -1}
	for _, a := range attrs {
		r.AddAttributes(a)
	}
	return r
}

func TestLoggerEmitAttributeDeduplicationRandomized(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := range 500 {
		n := rng.IntN(60)
		nKeys := 1 + rng.IntN(30)
		limit := -1
		if rng.IntN(2) == 0 {
			limit = rng.IntN(20)
		}
		in := make([]attribute.KeyValue, n)
		for j := range in {
			in[j] = attribute.Int("k"+strconv.Itoa(rng.IntN(nKeys)), j)
		}

		name := fmt.Sprintf("%d/n=%d,keys=%d,limit=%d", i, n, nKeys, limit)
		t.Run(name, func(t *testing.T) {
			want := referenceRecord(in, limit)
			got := emitOne(t, apiRecord(in...), WithAttributeCountLimit(limit))
			assert.Equal(t, recordAttrs(want), recordAttrs(got))
			assert.Equal(t, want.DroppedAttributes(), got.DroppedAttributes())
		})
	}
}

func BenchmarkLoggerEmitAttributes(b *testing.B) {
	discard := newProcessor("discard")
	discard.onEmitFunc = func(context.Context, *Record) error { return nil }

	for _, n := range []int{5, 10, 20, 40, 80} {
		unique := apiRecord(keys(n)...)
		// Every fourth attribute repeats an earlier key.
		dupAttrs := keys(n)
		for i := 3; i < n; i += 4 {
			dupAttrs[i] = attribute.String("k"+strconv.Itoa(i/4), "dup")
		}
		dups := apiRecord(dupAttrs...)

		for _, bc := range []struct {
			name   string
			record log.Record
			opts   []LoggerProviderOption
		}{
			{"Unique", unique, nil},
			{"Duplicates", dups, nil},
			{"AllowKeyDuplication", unique, []LoggerProviderOption{WithAllowKeyDuplication()}},
		} {
			opts := append([]LoggerProviderOption{WithProcessor(discard)}, bc.opts...)
			logger := NewLoggerProvider(opts...).Logger(b.Name())
			b.Run(fmt.Sprintf("%s/%d", bc.name, n), func(b *testing.B) {
				ctx := b.Context()
				b.ReportAllocs()
				for b.Loop() {
					logger.Emit(ctx, bc.record)
				}
			})
		}
	}
}
