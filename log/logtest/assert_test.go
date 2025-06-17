// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
)

var y2k = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// Compile-time check to ensure testing structs implement TestingT.
var (
	_ TestingT = (*testing.T)(nil)
	_ TestingT = (*testing.B)(nil)
	_ TestingT = (*testing.F)(nil)
	_ TestingT = (*mockTestingT)(nil)
)

type mockTestingT struct {
	errors []string
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func TestAssertEqual(t *testing.T) {
	a := Recording{
		Scope{Name: t.Name()}: []Record{
			{Body: log.StringValue("msg"), Attributes: []log.KeyValue{log.String("foo", "bar"), log.Int("n", 1)}},
		},
	}
	b := Recording{
		Scope{Name: t.Name()}: []Record{
			{Body: log.StringValue("msg"), Attributes: []log.KeyValue{log.Int("n", 1), log.String("foo", "bar")}},
		},
	}

	got := AssertEqual(t, a, b)
	assert.True(t, got, "expected recordings to be equal")
}

func TestAssertEqualRecording(t *testing.T) {
	tests := []struct {
		name string
		a    Recording
		b    Recording
		opts []AssertOption
		want bool
	}{
		{
			name: "equal recordings",
			a: Recording{
				Scope{Name: t.Name()}: []Record{
					{
						Timestamp:  y2k,
						Context:    context.Background(),
						Attributes: []log.KeyValue{log.Int("n", 1), log.String("foo", "bar")},
					},
				},
			},
			b: Recording{
				Scope{Name: t.Name()}: []Record{
					{
						Timestamp:  y2k,
						Context:    context.Background(),
						Attributes: []log.KeyValue{log.String("foo", "bar"), log.Int("n", 1)},
					},
				},
			},
			want: true,
		},
		{
			name: "different recordings",
			a: Recording{
				Scope{Name: t.Name()}: []Record{
					{Attributes: []log.KeyValue{log.String("foo", "bar")}},
				},
			},
			b: Recording{
				Scope{Name: t.Name()}: []Record{
					{Attributes: []log.KeyValue{log.Int("n", 1)}},
				},
			},
			want: false,
		},
		{
			name: "equal empty scopes",
			a: Recording{
				Scope{Name: t.Name()}: nil,
			},
			b: Recording{
				Scope{Name: t.Name()}: []Record{},
			},
			want: true,
		},
		{
			name: "equal empty attributes",
			a: Recording{
				Scope{Name: t.Name()}: []Record{
					{Body: log.StringValue("msg"), Attributes: []log.KeyValue{}},
				},
			},
			b: Recording{
				Scope{Name: t.Name()}: []Record{
					{Body: log.StringValue("msg"), Attributes: nil},
				},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := AssertEqual(mockT, tc.a, tc.b, tc.opts...)
			if result != tc.want {
				t.Errorf("AssertEqual() = %v, want %v", result, tc.want)
			}
			if !tc.want && len(mockT.errors) == 0 {
				t.Errorf("expected Errorf call but got none")
			}
		})
	}
}

func TestAssertEqualRecord(t *testing.T) {
	tests := []struct {
		name string
		a    Record
		b    Record
		opts []AssertOption
		want bool
	}{
		{
			name: "equal records",
			a: Record{
				Timestamp:  y2k,
				Context:    context.Background(),
				Attributes: []log.KeyValue{log.Int("n", 1), log.String("foo", "bar")},
			},
			b: Record{
				Timestamp:  y2k,
				Context:    context.Background(),
				Attributes: []log.KeyValue{log.String("foo", "bar"), log.Int("n", 1)},
			},
			want: true,
		},
		{
			name: "different records",
			a: Record{
				Attributes: []log.KeyValue{log.String("foo", "bar")},
			},
			b: Record{
				Attributes: []log.KeyValue{log.Int("n", 1)},
			},
			want: false,
		},
		{
			name: "Transform to ignore timestamps",
			a: Record{
				Attributes: []log.KeyValue{log.Int("n", 1), log.String("foo", "bar")},
			},
			b: Record{
				Timestamp:  y2k,
				Attributes: []log.KeyValue{log.String("foo", "bar"), log.Int("n", 1)},
			},
			opts: []AssertOption{
				Transform(func(time.Time) time.Time {
					return time.Time{}
				}),
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := AssertEqual(mockT, tc.a, tc.b, tc.opts...)
			if result != tc.want {
				t.Errorf("AssertEqual() = %v, want %v", result, tc.want)
			}
			if !tc.want && len(mockT.errors) == 0 {
				t.Errorf("expected Errorf call but got none")
			}
		})
	}
}

func TestDesc(t *testing.T) {
	mockT := &mockTestingT{}
	a := Record{
		Attributes: []log.KeyValue{log.String("foo", "bar")},
	}
	b := Record{
		Attributes: []log.KeyValue{log.Int("n", 1)},
	}

	AssertEqual(mockT, a, b, Desc("custom message, %s", "test"))

	require.Len(t, mockT.errors, 1, "expected one error")
	assert.Contains(t, mockT.errors[0], "custom message, test\n", "expected custom message")
}
