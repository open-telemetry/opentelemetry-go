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
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TestSetStatus(t *testing.T) {
	tests := []struct {
		name        string
		span        recordingSpan
		code        codes.Code
		description string
		expected    Status
	}{
		{
			"Error and description should overwrite Unset",
			recordingSpan{},
			codes.Error,
			"description",
			Status{Code: codes.Error, Description: "description"},
		},
		{
			"Ok should overwrite Unset and ignore description",
			recordingSpan{},
			codes.Ok,
			"description",
			Status{Code: codes.Ok},
		},
		{
			"Error and description should return error and overwrite description",
			recordingSpan{status: Status{Code: codes.Error, Description: "d1"}},
			codes.Error,
			"d2",
			Status{Code: codes.Error, Description: "d2"},
		},
		{
			"Ok should overwrite error and remove description",
			recordingSpan{status: Status{Code: codes.Error, Description: "d1"}},
			codes.Ok,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Error and description should be ignored when already Ok",
			recordingSpan{status: Status{Code: codes.Ok}},
			codes.Error,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Ok should be noop when already Ok",
			recordingSpan{status: Status{Code: codes.Ok}},
			codes.Ok,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Unset should be noop when already Ok",
			recordingSpan{status: Status{Code: codes.Ok}},
			codes.Unset,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Unset should be noop when already Error",
			recordingSpan{status: Status{Code: codes.Error, Description: "d1"}},
			codes.Unset,
			"d2",
			Status{Code: codes.Error, Description: "d1"},
		},
	}

	for i := range tests {
		tc := &tests[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.span.SetStatus(tc.code, tc.description)
			assert.Equal(t, tc.expected, tc.span.status)
		})
	}
}

func TestTruncateAttr(t *testing.T) {
	const key = "key"

	strAttr := attribute.String(key, "value")
	strSliceAttr := attribute.StringSlice(key, []string{"value-0", "value-1"})

	tests := []struct {
		limit      int
		attr, want attribute.KeyValue
	}{
		{
			limit: -1,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: -1,
			attr:  strSliceAttr,
			want:  strSliceAttr,
		},
		{
			limit: 0,
			attr:  attribute.Bool(key, true),
			want:  attribute.Bool(key, true),
		},
		{
			limit: 0,
			attr:  attribute.BoolSlice(key, []bool{true, false}),
			want:  attribute.BoolSlice(key, []bool{true, false}),
		},
		{
			limit: 0,
			attr:  attribute.Int(key, 42),
			want:  attribute.Int(key, 42),
		},
		{
			limit: 0,
			attr:  attribute.IntSlice(key, []int{42, -1}),
			want:  attribute.IntSlice(key, []int{42, -1}),
		},
		{
			limit: 0,
			attr:  attribute.Int64(key, 42),
			want:  attribute.Int64(key, 42),
		},
		{
			limit: 0,
			attr:  attribute.Int64Slice(key, []int64{42, -1}),
			want:  attribute.Int64Slice(key, []int64{42, -1}),
		},
		{
			limit: 0,
			attr:  attribute.Float64(key, 42),
			want:  attribute.Float64(key, 42),
		},
		{
			limit: 0,
			attr:  attribute.Float64Slice(key, []float64{42, -1}),
			want:  attribute.Float64Slice(key, []float64{42, -1}),
		},
		{
			limit: 0,
			attr:  strAttr,
			want:  attribute.String(key, ""),
		},
		{
			limit: 0,
			attr:  strSliceAttr,
			want:  attribute.StringSlice(key, []string{"", ""}),
		},
		{
			limit: 0,
			attr:  attribute.Stringer(key, bytes.NewBufferString("value")),
			want:  attribute.String(key, ""),
		},
		{
			limit: 1,
			attr:  strAttr,
			want:  attribute.String(key, "v"),
		},
		{
			limit: 1,
			attr:  strSliceAttr,
			want:  attribute.StringSlice(key, []string{"v", "v"}),
		},
		{
			limit: 5,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: 7,
			attr:  strSliceAttr,
			want:  strSliceAttr,
		},
		{
			limit: 6,
			attr:  attribute.StringSlice(key, []string{"value", "value-1"}),
			want:  attribute.StringSlice(key, []string{"value", "value-"}),
		},
		{
			limit: 128,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: 128,
			attr:  strSliceAttr,
			want:  strSliceAttr,
		},
		{
			// This tests the ordinary safeTruncate().
			limit: 10,
			attr:  attribute.String(key, "€€€€"), // 3 bytes each
			want:  attribute.String(key, "€€€"),
		},
		{
			// This tests truncation with an invalid UTF-8 input.
			//
			// Note that after removing the invalid rune,
			// the string is over length and still has to
			// be truncated on a code point boundary.
			limit: 10,
			attr:  attribute.String(key, "€"[0:2]+"hello€€"), // corrupted first rune, then over limit
			want:  attribute.String(key, "hello€"),
		},
		{
			// This tests the fallback to invalidTruncate()
			// where after validation the string does not require
			// truncation.
			limit: 6,
			attr:  attribute.String(key, "€"[0:2]+"hello"), // corrupted first rune, then not over limit
			want:  attribute.String(key, "hello"),
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s->%s(limit:%d)", test.attr.Key, test.attr.Value.Emit(), test.limit)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, truncateAttr(test.limit, test.attr))
		})
	}
}
