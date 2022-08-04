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

//go:build go1.18
// +build go1.18

package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/transform"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

type attributeTest struct {
	attrs    []attribute.KeyValue
	expected []*commonpb.KeyValue
}

func TestAttributes(t *testing.T) {
	for _, test := range []attributeTest{
		{nil, nil},
		{
			[]attribute.KeyValue{
				attribute.Int("int to int", 123),
				attribute.Int64("int64 to int64", 1234567),
				attribute.Float64("float64 to double", 1.61),
				attribute.String("string to string", "string"),
				attribute.Bool("bool to bool", true),
			},
			[]*commonpb.KeyValue{
				{
					Key: "int to int",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_IntValue{
							IntValue: 123,
						},
					},
				},
				{
					Key: "int64 to int64",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_IntValue{
							IntValue: 1234567,
						},
					},
				},
				{
					Key: "float64 to double",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_DoubleValue{
							DoubleValue: 1.61,
						},
					},
				},
				{
					Key: "string to string",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_StringValue{
							StringValue: "string",
						},
					},
				},
				{
					Key: "bool to bool",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_BoolValue{
							BoolValue: true,
						},
					},
				},
			},
		},
	} {
		got := KeyValues(test.attrs)
		if !assert.Len(t, got, len(test.expected)) {
			continue
		}
		for i, actual := range got {
			if a, ok := actual.Value.Value.(*commonpb.AnyValue_DoubleValue); ok {
				e, ok := test.expected[i].Value.Value.(*commonpb.AnyValue_DoubleValue)
				if !ok {
					t.Errorf("expected AnyValue_DoubleValue, got %T", test.expected[i].Value.Value)
					continue
				}
				if !assert.InDelta(t, e.DoubleValue, a.DoubleValue, 0.01) {
					continue
				}
				e.DoubleValue = a.DoubleValue
			}
			assert.Equal(t, test.expected[i], actual)
		}
	}
}

func TestArrayAttributes(t *testing.T) {
	// Array KeyValue supports only arrays of primitive types:
	// "bool", "int", "int64",
	// "float64", "string",
	for _, test := range []attributeTest{
		{nil, nil},
		{
			[]attribute.KeyValue{
				{
					Key:   attribute.Key("invalid"),
					Value: attribute.Value{},
				},
			},
			[]*commonpb.KeyValue{
				{
					Key: "invalid",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_StringValue{
							StringValue: "INVALID",
						},
					},
				},
			},
		},
		{
			[]attribute.KeyValue{
				attribute.BoolSlice("bool slice to bool array", []bool{true, false}),
				attribute.IntSlice("int slice to int64 array", []int{1, 2, 3}),
				attribute.Int64Slice("int64 slice to int64 array", []int64{1, 2, 3}),
				attribute.Float64Slice("float64 slice to double array", []float64{1.11, 2.22, 3.33}),
				attribute.StringSlice("string slice to string array", []string{"foo", "bar", "baz"}),
			},
			[]*commonpb.KeyValue{
				newOTelBoolArray("bool slice to bool array", []bool{true, false}),
				newOTelIntArray("int slice to int64 array", []int64{1, 2, 3}),
				newOTelIntArray("int64 slice to int64 array", []int64{1, 2, 3}),
				newOTelDoubleArray("float64 slice to double array", []float64{1.11, 2.22, 3.33}),
				newOTelStringArray("string slice to string array", []string{"foo", "bar", "baz"}),
			},
		},
	} {
		actualArrayAttributes := KeyValues(test.attrs)
		expectedArrayAttributes := test.expected
		if !assert.Len(t, actualArrayAttributes, len(expectedArrayAttributes)) {
			continue
		}

		for i, actualArrayAttr := range actualArrayAttributes {
			expectedArrayAttr := expectedArrayAttributes[i]
			expectedKey, actualKey := expectedArrayAttr.Key, actualArrayAttr.Key
			if !assert.Equal(t, expectedKey, actualKey) {
				continue
			}

			expected := expectedArrayAttr.Value.GetArrayValue()
			actual := actualArrayAttr.Value.GetArrayValue()
			if expected == nil {
				assert.Nil(t, actual)
				continue
			}
			if assert.NotNil(t, actual, "expected not nil for %s", actualKey) {
				assertExpectedArrayValues(t, expected.Values, actual.Values)
			}
		}

	}
}

func assertExpectedArrayValues(t *testing.T, expectedValues, actualValues []*commonpb.AnyValue) {
	for i, actual := range actualValues {
		expected := expectedValues[i]
		if a, ok := actual.Value.(*commonpb.AnyValue_DoubleValue); ok {
			e, ok := expected.Value.(*commonpb.AnyValue_DoubleValue)
			if !ok {
				t.Errorf("expected AnyValue_DoubleValue, got %T", expected.Value)
				continue
			}
			if !assert.InDelta(t, e.DoubleValue, a.DoubleValue, 0.01) {
				continue
			}
			e.DoubleValue = a.DoubleValue
		}
		assert.Equal(t, expected, actual)
	}
}

func newOTelBoolArray(key string, values []bool) *commonpb.KeyValue {
	arrayValues := []*commonpb.AnyValue{}
	for _, b := range values {
		arrayValues = append(arrayValues, &commonpb.AnyValue{
			Value: &commonpb.AnyValue_BoolValue{
				BoolValue: b,
			},
		})
	}

	return newOTelArray(key, arrayValues)
}

func newOTelIntArray(key string, values []int64) *commonpb.KeyValue {
	arrayValues := []*commonpb.AnyValue{}

	for _, i := range values {
		arrayValues = append(arrayValues, &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: i,
			},
		})
	}

	return newOTelArray(key, arrayValues)
}

func newOTelDoubleArray(key string, values []float64) *commonpb.KeyValue {
	arrayValues := []*commonpb.AnyValue{}

	for _, d := range values {
		arrayValues = append(arrayValues, &commonpb.AnyValue{
			Value: &commonpb.AnyValue_DoubleValue{
				DoubleValue: d,
			},
		})
	}

	return newOTelArray(key, arrayValues)
}

func newOTelStringArray(key string, values []string) *commonpb.KeyValue {
	arrayValues := []*commonpb.AnyValue{}

	for _, s := range values {
		arrayValues = append(arrayValues, &commonpb.AnyValue{
			Value: &commonpb.AnyValue_StringValue{
				StringValue: s,
			},
		})
	}

	return newOTelArray(key, arrayValues)
}

func newOTelArray(key string, arrayValues []*commonpb.AnyValue) *commonpb.KeyValue {
	return &commonpb.KeyValue{
		Key: key,
		Value: &commonpb.AnyValue{
			Value: &commonpb.AnyValue_ArrayValue{
				ArrayValue: &commonpb.ArrayValue{
					Values: arrayValues,
				},
			},
		},
	}
}

func TestAttrIter(t *testing.T) {
	tests := []struct {
		kvs      []attribute.KeyValue
		expected []*commonpb.KeyValue
	}{
		{
			nil,
			nil,
		},
		{
			[]attribute.KeyValue{},
			nil,
		},
		{
			[]attribute.KeyValue{
				attribute.Bool("true", true),
				attribute.Int64("one", 1),
				attribute.Int64("two", 2),
				attribute.Float64("three", 3),
				attribute.Int("four", 4),
				attribute.Int("five", 5),
				attribute.Float64("six", 6),
				attribute.Int("seven", 7),
				attribute.Int("eight", 8),
				attribute.String("the", "final word"),
			},
			[]*commonpb.KeyValue{
				{Key: "eight", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 8}}},
				{Key: "five", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 5}}},
				{Key: "four", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 4}}},
				{Key: "one", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 1}}},
				{Key: "seven", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 7}}},
				{Key: "six", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 6.0}}},
				{Key: "the", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "final word"}}},
				{Key: "three", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 3.0}}},
				{Key: "true", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_BoolValue{BoolValue: true}}},
				{Key: "two", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 2}}},
			},
		},
	}

	for _, test := range tests {
		labels := attribute.NewSet(test.kvs...)
		assert.Equal(t, test.expected, AttrIter(labels.Iter()))
	}
}
