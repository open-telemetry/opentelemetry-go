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

package attribute_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/attribute"
)

func TestKeyValueConstructors(t *testing.T) {
	tt := []struct {
		name     string
		actual   attribute.KeyValue
		expected attribute.KeyValue
	}{
		{
			name:   "Bool",
			actual: attribute.Bool("k1", true),
			expected: attribute.KeyValue{
				Key:   "k1",
				Value: attribute.BoolValue(true),
			},
		},
		{
			name:   "Int64",
			actual: attribute.Int64("k1", 123),
			expected: attribute.KeyValue{
				Key:   "k1",
				Value: attribute.Int64Value(123),
			},
		},
		{
			name:   "Float64",
			actual: attribute.Float64("k1", 123.5),
			expected: attribute.KeyValue{
				Key:   "k1",
				Value: attribute.Float64Value(123.5),
			},
		},
		{
			name:   "String",
			actual: attribute.String("k1", "123.5"),
			expected: attribute.KeyValue{
				Key:   "k1",
				Value: attribute.StringValue("123.5"),
			},
		},
		{
			name:   "Int",
			actual: attribute.Int("k1", 123),
			expected: attribute.KeyValue{
				Key:   "k1",
				Value: attribute.IntValue(123),
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.actual, test.expected, cmp.AllowUnexported(attribute.Value{})); diff != "" {
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
		wantType  attribute.Type
		wantValue interface{}
	}{
		{
			key:       "bool type inferred",
			value:     true,
			wantType:  attribute.BOOL,
			wantValue: true,
		},
		{
			key:       "int64 type inferred",
			value:     int64(42),
			wantType:  attribute.INT64,
			wantValue: int64(42),
		},
		{
			key:       "float64 type inferred",
			value:     float64(42.1),
			wantType:  attribute.FLOAT64,
			wantValue: 42.1,
		},
		{
			key:       "string type inferred",
			value:     "foo",
			wantType:  attribute.STRING,
			wantValue: "foo",
		},
		{
			key:       "stringer type inferred",
			value:     builder,
			wantType:  attribute.STRING,
			wantValue: "foo",
		},
		{
			key:       "unknown value serialized as %v",
			value:     nil,
			wantType:  attribute.STRING,
			wantValue: "<nil>",
		},
		{
			key:       "JSON struct serialized correctly",
			value:     &jsonifyStruct,
			wantType:  attribute.STRING,
			wantValue: `{"Public":"foo","tagName":"baz","Empty":""}`,
		},
		{
			key:       "Invalid JSON struct falls back to string",
			value:     &invalidStruct,
			wantType:  attribute.STRING,
			wantValue: "&{(0+0i)}",
		},
	} {
		t.Logf("Running test case %s", testcase.key)
		keyValue := attribute.Any(testcase.key, testcase.value)
		if keyValue.Value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", keyValue.Value.Type(), testcase.wantType)
		}
		got := keyValue.Value.AsInterface()
		if diff := cmp.Diff(testcase.wantValue, got); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
	}
}

func TestKeyValueValid(t *testing.T) {
	tests := []struct {
		desc  string
		valid bool
		kv    attribute.KeyValue
	}{
		{
			desc:  "uninitialized KeyValue should be invalid",
			valid: false,
			kv:    attribute.KeyValue{},
		},
		{
			desc:  "empty key value should be invalid",
			valid: false,
			kv:    attribute.Key("").Bool(true),
		},
		{
			desc:  "INVALID value type should be invalid",
			valid: false,
			kv: attribute.KeyValue{
				Key: attribute.Key("valid key"),
				// Default type is INVALID.
				Value: attribute.Value{},
			},
		},
		{
			desc:  "non-empty key with BOOL type Value should be valid",
			valid: true,
			kv:    attribute.Bool("bool", true),
		},
		{
			desc:  "non-empty key with INT64 type Value should be valid",
			valid: true,
			kv:    attribute.Int64("int64", 0),
		},
		{
			desc:  "non-empty key with FLOAT64 type Value should be valid",
			valid: true,
			kv:    attribute.Float64("float64", 0),
		},
		{
			desc:  "non-empty key with STRING type Value should be valid",
			valid: true,
			kv:    attribute.String("string", ""),
		},
		{
			desc:  "non-empty key with ARRAY type Value should be valid",
			valid: true,
			kv:    attribute.Array("array", []int{}),
		},
	}

	for _, test := range tests {
		if got, want := test.kv.Valid(), test.valid; got != want {
			t.Error(test.desc)
		}
	}
}
