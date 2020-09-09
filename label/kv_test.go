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

package label_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/label"
)

func TestKeyValueConstructors(t *testing.T) {
	tt := []struct {
		name     string
		actual   label.KeyValue
		expected label.KeyValue
	}{
		{
			name:   "Bool",
			actual: label.Bool("k1", true),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.BoolValue(true),
			},
		},
		{
			name:   "Int64",
			actual: label.Int64("k1", 123),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.Int64Value(123),
			},
		},
		{
			name:   "Uint64",
			actual: label.Uint64("k1", 1),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.Uint64Value(1),
			},
		},
		{
			name:   "Float64",
			actual: label.Float64("k1", 123.5),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.Float64Value(123.5),
			},
		},
		{
			name:   "Int32",
			actual: label.Int32("k1", 123),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.Int32Value(123),
			},
		},
		{
			name:   "Uint32",
			actual: label.Uint32("k1", 123),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.Uint32Value(123),
			},
		},
		{
			name:   "Float32",
			actual: label.Float32("k1", 123.5),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.Float32Value(123.5),
			},
		},
		{
			name:   "String",
			actual: label.String("k1", "123.5"),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.StringValue("123.5"),
			},
		},
		{
			name:   "Int",
			actual: label.Int("k1", 123),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.IntValue(123),
			},
		},
		{
			name:   "Uint",
			actual: label.Uint("k1", 123),
			expected: label.KeyValue{
				Key:   "k1",
				Value: label.UintValue(123),
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.actual, test.expected, cmp.AllowUnexported(label.Value{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestAny(t *testing.T) {
	builder := &strings.Builder{}
	builder.WriteString("foo")
	jsonifyStruct := struct {
		Public    string
		private   string
		Tagged    string `json:"tagName"`
		Empty     string
		OmitEmpty string `json:",omitempty"`
		Omit      string `json:"-"`
	}{"foo", "bar", "baz", "", "", "omitted"}
	invalidStruct := struct {
		N complex64
	}{complex(0, 0)}
	for _, testcase := range []struct {
		key       string
		value     interface{}
		wantType  label.Type
		wantValue interface{}
	}{
		{
			key:       "bool type inferred",
			value:     true,
			wantType:  label.BOOL,
			wantValue: true,
		},
		{
			key:       "int64 type inferred",
			value:     int64(42),
			wantType:  label.INT64,
			wantValue: int64(42),
		},
		{
			key:       "uint64 type inferred",
			value:     uint64(42),
			wantType:  label.UINT64,
			wantValue: uint64(42),
		},
		{
			key:       "float64 type inferred",
			value:     float64(42.1),
			wantType:  label.FLOAT64,
			wantValue: 42.1,
		},
		{
			key:       "int32 type inferred",
			value:     int32(42),
			wantType:  label.INT32,
			wantValue: int32(42),
		},
		{
			key:       "uint32 type inferred",
			value:     uint32(42),
			wantType:  label.UINT32,
			wantValue: uint32(42),
		},
		{
			key:       "float32 type inferred",
			value:     float32(42.1),
			wantType:  label.FLOAT32,
			wantValue: float32(42.1),
		},
		{
			key:       "string type inferred",
			value:     "foo",
			wantType:  label.STRING,
			wantValue: "foo",
		},
		{
			key:       "stringer type inferred",
			value:     builder,
			wantType:  label.STRING,
			wantValue: "foo",
		},
		{
			key:       "unknown value serialized as %v",
			value:     nil,
			wantType:  label.STRING,
			wantValue: "<nil>",
		},
		{
			key:       "JSON struct serialized correctly",
			value:     &jsonifyStruct,
			wantType:  label.STRING,
			wantValue: `{"Public":"foo","tagName":"baz","Empty":""}`,
		},
		{
			key:       "Invalid JSON struct falls back to string",
			value:     &invalidStruct,
			wantType:  label.STRING,
			wantValue: "&{(0+0i)}",
		},
	} {
		t.Logf("Running test case %s", testcase.key)
		keyValue := label.Any(testcase.key, testcase.value)
		if keyValue.Value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", keyValue.Value.Type(), testcase.wantType)
		}
		got := keyValue.Value.AsInterface()
		if diff := cmp.Diff(testcase.wantValue, got); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
	}
}
