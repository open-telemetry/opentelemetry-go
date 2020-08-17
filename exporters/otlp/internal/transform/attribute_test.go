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

package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/api/kv"
	commonpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/common/v1"
)

type attributeTest struct {
	attrs    []kv.KeyValue
	expected []*commonpb.KeyValue
}

func TestAttributes(t *testing.T) {
	for _, test := range []attributeTest{
		{nil, nil},
		{
			[]kv.KeyValue{
				kv.Int("int to int", 123),
				kv.Uint("uint to int", 1234),
				kv.Int32("int32 to int", 12345),
				kv.Uint32("uint32 to int", 123456),
				kv.Int64("int64 to int64", 1234567),
				kv.Uint64("uint64 to int64", 12345678),
				kv.Float32("float32 to double", 3.14),
				kv.Float32("float64 to double", 1.61),
				kv.String("string to string", "string"),
				kv.Bool("bool to bool", true),
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
					Key: "uint to int",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_IntValue{
							IntValue: 1234,
						},
					},
				},
				{
					Key: "int32 to int",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_IntValue{
							IntValue: 12345,
						},
					},
				},
				{
					Key: "uint32 to int",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_IntValue{
							IntValue: 123456,
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
					Key: "uint64 to int64",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_IntValue{
							IntValue: 12345678,
						},
					},
				},
				{
					Key: "float32 to double",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_DoubleValue{
							DoubleValue: 3.14,
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
		got := Attributes(test.attrs)
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
	// "bool", "int", "int32", "int64",
	// "float32", "float64", "string",
	// "uint", "uint32", "uint64"
	for _, test := range []attributeTest{
		{nil, nil},
		{
			[]kv.KeyValue{
				kv.Array("bool array to bool array", []bool{true, false}),
				kv.Array("int array to int64 array", []int{1, 2, 3}),
				kv.Array("uint array to int64 array", []uint{1, 2, 3}),
				kv.Array("int32 array to int64 array", []int32{1, 2, 3}),
				kv.Array("uint32 array to int64 array", []uint32{1, 2, 3}),
				kv.Array("int64 array to int64 array", []int64{1, 2, 3}),
				kv.Array("uint64 array to int64 array", []uint64{1, 2, 3}),
				kv.Array("float32 array to double array", []float32{1.11, 2.22, 3.33}),
				kv.Array("float64 array to double array", []float64{1.11, 2.22, 3.33}),
				kv.Array("string array to string array", []string{"foo", "bar", "baz"}),
			},
			[]*commonpb.KeyValue{
				newOTelBoolArray("bool array to bool array", []bool{true, false}),
				newOTelIntArray("int array to int64 array", []int64{1, 2, 3}),
				newOTelIntArray("uint array to int64 array", []int64{1, 2, 3}),
				newOTelIntArray("int32 array to int64 array", []int64{1, 2, 3}),
				newOTelIntArray("uint32 array to int64 array", []int64{1, 2, 3}),
				newOTelIntArray("int64 array to int64 array", []int64{1, 2, 3}),
				newOTelIntArray("uint64 array to int64 array", []int64{1, 2, 3}),
				newOTelDoubleArray("float32 array to double array", []float64{1.11, 2.22, 3.33}),
				newOTelDoubleArray("float64 array to double array", []float64{1.11, 2.22, 3.33}),
				newOTelStringArray("string array to string array", []string{"foo", "bar", "baz"}),
			},
		},
	} {
		actualArrayAttributes := Attributes(test.attrs)
		expectedArrayAttributes := test.expected
		if !assert.Len(t, actualArrayAttributes, len(expectedArrayAttributes)) {
			continue
		}

		for i, actualArrayAttr := range actualArrayAttributes {
			expectedArrayAttr := expectedArrayAttributes[i]
			if !assert.Equal(t, expectedArrayAttr.Key, actualArrayAttr.Key) {
				continue
			}

			expectedArrayValue := expectedArrayAttr.Value.GetArrayValue()
			assert.NotNil(t, expectedArrayValue)

			actualArrayValue := actualArrayAttr.Value.GetArrayValue()
			assert.NotNil(t, actualArrayValue)

			assertExpectedArrayValues(t, expectedArrayValue.Values, actualArrayValue.Values)
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
