// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

var y2k = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

type mockTestingT struct {
	errors []string
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.errors = append(m.errors, format)
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

	got := assertEqual(t, a, b)
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
			result := assertEqual(mockT, tc.a, tc.b, tc.opts...)
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assertEqual(mockT, tc.a, tc.b, tc.opts...)
			if result != tc.want {
				t.Errorf("AssertEqual() = %v, want %v", result, tc.want)
			}
			if !tc.want && len(mockT.errors) == 0 {
				t.Errorf("expected Errorf call but got none")
			}
		})
	}
}
