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

package kv_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/kv/value"
)

func TestKeyValueConstructors(t *testing.T) {
	tt := []struct {
		name     string
		actual   kv.KeyValue
		expected kv.KeyValue
	}{
		{
			name:   "Bool",
			actual: kv.Bool("k1", true),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Bool(true),
			},
		},
		{
			name:   "Int64",
			actual: kv.Int64("k1", 123),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Int64(123),
			},
		},
		{
			name:   "Uint64",
			actual: kv.Uint64("k1", 1),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Uint64(1),
			},
		},
		{
			name:   "Float64",
			actual: kv.Float64("k1", 123.5),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Float64(123.5),
			},
		},
		{
			name:   "Int32",
			actual: kv.Int32("k1", 123),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Int32(123),
			},
		},
		{
			name:   "Uint32",
			actual: kv.Uint32("k1", 123),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Uint32(123),
			},
		},
		{
			name:   "Float32",
			actual: kv.Float32("k1", 123.5),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Float32(123.5),
			},
		},
		{
			name:   "String",
			actual: kv.String("k1", "123.5"),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.String("123.5"),
			},
		},
		{
			name:   "Int",
			actual: kv.Int("k1", 123),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Int(123),
			},
		},
		{
			name:   "Uint",
			actual: kv.Uint("k1", 123),
			expected: kv.KeyValue{
				Key:   "k1",
				Value: value.Uint(123),
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.actual, test.expected, cmp.AllowUnexported(value.Value{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInfer(t *testing.T) {
	builder := &strings.Builder{}
	builder.WriteString("foo")
	for _, testcase := range []struct {
		key       string
		value     interface{}
		wantType  value.Type
		wantValue interface{}
	}{
		{
			key:       "bool type inferred",
			value:     true,
			wantType:  value.BOOL,
			wantValue: true,
		},
		{
			key:       "int64 type inferred",
			value:     int64(42),
			wantType:  value.INT64,
			wantValue: int64(42),
		},
		{
			key:       "uint64 type inferred",
			value:     uint64(42),
			wantType:  value.UINT64,
			wantValue: uint64(42),
		},
		{
			key:       "float64 type inferred",
			value:     float64(42.1),
			wantType:  value.FLOAT64,
			wantValue: 42.1,
		},
		{
			key:       "int32 type inferred",
			value:     int32(42),
			wantType:  value.INT32,
			wantValue: int32(42),
		},
		{
			key:       "uint32 type inferred",
			value:     uint32(42),
			wantType:  value.UINT32,
			wantValue: uint32(42),
		},
		{
			key:       "float32 type inferred",
			value:     float32(42.1),
			wantType:  value.FLOAT32,
			wantValue: float32(42.1),
		},
		{
			key:       "string type inferred",
			value:     "foo",
			wantType:  value.STRING,
			wantValue: "foo",
		},
		{
			key:       "stringer type inferred",
			value:     builder,
			wantType:  value.STRING,
			wantValue: "foo",
		},
		{
			key:       "unknown value serialized as %v",
			value:     nil,
			wantType:  value.STRING,
			wantValue: "<nil>",
		},
	} {
		t.Logf("Running test case %s", testcase.key)
		keyValue := kv.Infer(testcase.key, testcase.value)
		if keyValue.Value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", keyValue.Value.Type(), testcase.wantType)
		}
		got := keyValue.Value.AsInterface()
		if diff := cmp.Diff(testcase.wantValue, got); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
	}
}
