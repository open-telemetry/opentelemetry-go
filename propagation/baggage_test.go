// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package propagation_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
)

type property struct {
	Key, Value string
}

type member struct {
	Key, Value string

	Properties []property
}

func (m member) Member(t *testing.T) baggage.Member {
	props := make([]baggage.Property, 0, len(m.Properties))
	for _, p := range m.Properties {
		p, err := baggage.NewKeyValuePropertyRaw(p.Key, p.Value)
		if err != nil {
			t.Fatal(err)
		}
		props = append(props, p)
	}
	bMember, err := baggage.NewMemberRaw(m.Key, m.Value, props...)
	if err != nil {
		t.Fatal(err)
	}
	return bMember
}

type members []member

func (m members) Baggage(t *testing.T) baggage.Baggage {
	bMembers := make([]baggage.Member, 0, len(m))
	for _, mem := range m {
		bMembers = append(bMembers, mem.Member(t))
	}
	bag, err := baggage.New(bMembers...)
	if err != nil {
		t.Fatal(err)
	}
	return bag
}

func TestExtractValidBaggage(t *testing.T) {
	prop := propagation.TextMapPropagator(propagation.Baggage{})
	tests := []struct {
		name   string
		header string
		want   members
	}{
		{
			name:   "valid w3cHeader",
			header: "key1=val1,key2=val2",
			want: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:   "valid w3cHeader with spaces",
			header: "key1 =   val1,  key2 =val2   ",
			want: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:   "valid w3cHeader with properties",
			header: "key1=val1,key2=val2;prop=1",
			want: members{
				{Key: "key1", Value: "val1"},
				{
					Key:   "key2",
					Value: "val2",
					Properties: []property{
						{Key: "prop", Value: "1"},
					},
				},
			},
		},
		{
			name:   "valid header with an invalid header",
			header: "key1=val1,key2=val2,a,val3",
			want: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:   "valid header with no value",
			header: "key1=,key2=val2",
			want: members{
				{Key: "key1", Value: ""},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:   "valid header with url encoded string",
			header: "key1=val%252",
			want: members{
				{Key: "key1", Value: "val%2"},
			},
		},
		{
			name:   "empty header",
			header: "",
			want:   members{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapCarr := propagation.MapCarrier{}
			mapCarr["baggage"] = tt.header
			req, _ := http.NewRequest(http.MethodGet, "http://example.com", http.NoBody)
			req.Header.Set("baggage", tt.header)

			// test with http header carrier (which implements ValuesGetter)
			ctx := prop.Extract(t.Context(), propagation.HeaderCarrier(req.Header))
			expected := tt.want.Baggage(t)
			assert.Equal(t, expected, baggage.FromContext(ctx), "should extract baggage for HeaderCarrier")

			// test with map carrier (which does not implement ValuesGetter)
			ctx = prop.Extract(t.Context(), mapCarr)
			expected = tt.want.Baggage(t)
			assert.Equal(t, expected, baggage.FromContext(ctx), "should extract baggage for MapCarrier")
		})
	}
}

// generateBaggageHeader creates a baggage header string with n members.
func generateBaggageHeader(n int, prefix string) string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = fmt.Sprintf("%s%d=v%d", prefix, i, i)
	}
	return strings.Join(parts, ",")
}

// generateMembers creates n members with keys like "prefix0", "prefix1", etc.
func generateMembers(n int, prefix string) members {
	m := make(members, n)
	for i := range m {
		m[i] = member{Key: fmt.Sprintf("%s%d", prefix, i), Value: fmt.Sprintf("v%d", i)}
	}
	return m
}

func TestExtractValidMultipleBaggageHeaders(t *testing.T) {
	// W3C Baggage spec limits: https://www.w3.org/TR/baggage/#limits
	const maxMembers = 64
	const maxBytesPerBaggageString = 8192

	prop := propagation.TextMapPropagator(propagation.Baggage{})
	tests := []struct {
		name         string
		headers      []string
		want         members
		wantCount    int // Used when want is nil and we only care about count.
		wantMaxBytes int // Used to check that baggage size doesn't exceed limit.
	}{
		{
			name:    "non conflicting headers",
			headers: []string{"key1=val1", "key2=val2"},
			want: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:    "conflicting keys, uses last val",
			headers: []string{"key1=val1", "key1=val2"},
			want: members{
				{Key: "key1", Value: "val2"},
			},
		},
		{
			name:    "single empty",
			headers: []string{"", "key1=val1"},
			want: members{
				{Key: "key1", Value: "val1"},
			},
		},
		{
			name:    "all empty",
			headers: []string{"", ""},
			want:    members{},
		},
		{
			name:    "none",
			headers: []string{},
			want:    members{},
		},
		{
			name:    "single header with one invalid skips invalid",
			headers: []string{"key1=val1,invalid-no-equals,key2=val2"},
			want: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name: "multiple headers with one invalid skips invalid and continues",
			headers: []string{
				"key1=val1",
				"invalid-no-equals",
				"key2=val2",
			},
			want: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:    "single header at max members limit",
			headers: []string{generateBaggageHeader(maxMembers, "k")},
			want:    generateMembers(maxMembers, "k"),
		},
		{
			name:    "single header exceeds max members limit keeps 64",
			headers: []string{generateBaggageHeader(maxMembers+1, "k")},
			want:    generateMembers(maxMembers, "k"),
		},
		{
			name: "multiple headers exceeds total max members limit keeps 64",
			headers: []string{
				generateBaggageHeader(maxMembers/2, "a"),
				generateBaggageHeader(maxMembers/2, "b"),
				generateBaggageHeader(1, "c"),
			},
			want:         nil, // Non-deterministic truncation by baggage.New()
			wantCount:    maxMembers,
			wantMaxBytes: maxBytesPerBaggageString,
		},
		{
			name:    "single header at max bytes limit",
			headers: []string{"k=" + strings.Repeat("v", maxBytesPerBaggageString-2)},
			want: members{
				{Key: "k", Value: strings.Repeat("v", maxBytesPerBaggageString-2)},
			},
		},
		{
			name:    "single header exceeds max bytes limit drops oversized member",
			headers: []string{"k=" + strings.Repeat("v", maxBytesPerBaggageString-1)},
			want:    members{},
		},
		{
			name: "multiple headers exceed total max bytes keeps one that fits",
			headers: []string{
				"k=" + strings.Repeat("v", maxBytesPerBaggageString-2),
				"y=" + strings.Repeat("v", maxBytesPerBaggageString-2),
			},
			want:         nil, // Non-deterministic: either k or y will be kept
			wantCount:    1,   // Only one member fits
			wantMaxBytes: maxBytesPerBaggageString,
		},
		{
			name: "multiple headers within total max bytes",
			headers: []string{
				"k=" + strings.Repeat("v", maxBytesPerBaggageString/2-2),
				// The comma as the separator of member would take 1 byte.
				"y=" + strings.Repeat("v", maxBytesPerBaggageString/2-2-1),
			},
			want: members{
				{Key: "k", Value: strings.Repeat("v", maxBytesPerBaggageString/2-2)},
				{Key: "y", Value: strings.Repeat("v", maxBytesPerBaggageString/2-2-1)},
			},
		},
		{
			name: "many headers exceeding member limit caps collection early",
			headers: func() []string {
				// 100 headers with 10 members each = 1000 total members.
				// The cap should stop collecting after ~maxMembers and
				// New() truncates to exactly maxMembers.
				h := make([]string, 100)
				for i := range h {
					h[i] = generateBaggageHeader(10, fmt.Sprintf("h%d_k", i))
				}
				return h
			}(),
			wantCount:    maxMembers,
			wantMaxBytes: maxBytesPerBaggageString,
		},
		{
			name: "skips large member that exceeds byte limit and continues",
			headers: []string{
				"small1=v1,small2=v2",
				"large=" + strings.Repeat("x", maxBytesPerBaggageString),
				"small3=v3",
			},
			want: members{
				{Key: "small1", Value: "v1"},
				{Key: "small2", Value: "v2"},
				{Key: "small3", Value: "v3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "http://example.com", http.NoBody)
			req.Header["Baggage"] = tt.headers

			ctx := t.Context()
			ctx = prop.Extract(ctx, propagation.HeaderCarrier(req.Header))
			got := baggage.FromContext(ctx)

			// If want is specified, check exact match
			if tt.want != nil {
				expected := tt.want.Baggage(t)
				assert.Equal(t, expected, got)
			} else if tt.wantCount > 0 {
				// If only count is specified, verify count and byte limit
				assert.Equal(t, tt.wantCount, got.Len(), "expected member count")
				assert.LessOrEqual(t, len(got.String()), tt.wantMaxBytes, "baggage size exceeds limit")
			}
		})
	}
}

func TestExtractInvalidDistributedContextFromHTTPReq(t *testing.T) {
	prop := propagation.TextMapPropagator(propagation.Baggage{})
	tests := []struct {
		name   string
		header string
		has    members
	}{
		{
			name:   "no key values",
			header: "header1",
		},
		{
			name:   "invalid header with existing context",
			header: "header2",
			has: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:   "empty header value",
			header: "",
			has: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name:   "with properties",
			header: "key1=val1,key2=val2;prop=1",
			has: members{
				{Key: "key1", Value: "val1"},
				{
					Key:   "key2",
					Value: "val2",
					Properties: []property{
						{Key: "prop", Value: "1"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "http://example.com", http.NoBody)
			req.Header.Set("baggage", tt.header)

			expected := tt.has.Baggage(t)
			ctx := baggage.ContextWithBaggage(t.Context(), expected)
			ctx = prop.Extract(ctx, propagation.HeaderCarrier(req.Header))
			assert.Equal(t, expected, baggage.FromContext(ctx))
		})
	}
}

func TestInjectBaggageToHTTPReq(t *testing.T) {
	propagator := propagation.Baggage{}
	tests := []struct {
		name         string
		mems         members
		wantInHeader []string
	}{
		{
			name: "two simple values",
			mems: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
			wantInHeader: []string{"key1=val1", "key2=val2"},
		},
		{
			name: "values with escaped chars",
			mems: members{
				{Key: "key2", Value: "val3,4"},
			},
			wantInHeader: []string{"key2=val3%2C4"},
		},
		{
			name: "with properties",
			mems: members{
				{Key: "key1", Value: "val1"},
				{
					Key:   "key2",
					Value: "val2",
					Properties: []property{
						{Key: "prop", Value: "1"},
					},
				},
			},
			wantInHeader: []string{
				"key1=val1",
				"key2=val2;prop=1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "http://example.com", http.NoBody)
			ctx := baggage.ContextWithBaggage(t.Context(), tt.mems.Baggage(t))
			propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

			got := strings.Split(req.Header.Get("baggage"), ",")
			assert.ElementsMatch(t, tt.wantInHeader, got)
		})
	}
}

func TestBaggageInjectExtractRoundtrip(t *testing.T) {
	propagator := propagation.Baggage{}
	tests := []struct {
		name string
		mems members
	}{
		{
			name: "two simple values",
			mems: members{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
			},
		},
		{
			name: "values with escaped chars",
			mems: members{
				{Key: "key1", Value: "val3=4"},
				{Key: "key2", Value: "mess,me%up"},
			},
		},
		{
			name: "with properties",
			mems: members{
				{Key: "key1", Value: "val1"},
				{
					Key:   "key2",
					Value: "val2",
					Properties: []property{
						{Key: "prop", Value: "1"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.mems.Baggage(t)
			req, _ := http.NewRequest(http.MethodGet, "http://example.com", http.NoBody)
			ctx := baggage.ContextWithBaggage(t.Context(), b)
			propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

			ctx = propagator.Extract(t.Context(), propagation.HeaderCarrier(req.Header))
			got := baggage.FromContext(ctx)

			assert.Equal(t, b, got)
		})
	}
}

func TestBaggagePropagatorGetAllKeys(t *testing.T) {
	var propagator propagation.Baggage
	want := []string{"baggage"}
	got := propagator.Fields()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
