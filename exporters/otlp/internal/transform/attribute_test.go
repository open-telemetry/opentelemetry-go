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
	commonpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/common/v1"
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
				kv.Array("int array to int array", []int{1, 2, 3}),
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
				{
					Key: "int array to int array",
					Value: &commonpb.AnyValue{
						Value: &commonpb.AnyValue_ArrayValue{
							ArrayValue: &commonpb.ArrayValue{
								Values: []*commonpb.AnyValue{
									&commonpb.AnyValue{
										Value: &commonpb.AnyValue_IntValue{
											IntValue: 1,
										},
									},
									&commonpb.AnyValue{
										Value: &commonpb.AnyValue_IntValue{
											IntValue: 2,
										},
									},
									&commonpb.AnyValue{
										Value: &commonpb.AnyValue_IntValue{
											IntValue: 3,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		assertExpectedAttributes(t, test)
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
				kv.Array("int array to int array", []int{1, 2, 3}),
				kv.Array("int32 array to int array", []int32{1, 2, 3}),
				kv.Array("int64 array to int array", []int64{1, 2, 3}),
			},
			[]*commonpb.KeyValue{
				newOtelBoolArray("bool array to bool array", []bool{true, false}),
				newOtelIntArray("int array to int array", []int{1, 2, 3}),
				newOtelIntArray("int32 array to int array", []int{1, 2, 3}),
				newOtelIntArray("int64 array to int array", []int{1, 2, 3}),
			},
		},
	} {
		assertExpectedAttributes(t, test)
	}
}

func assertExpectedAttributes(t *testing.T, test attributeTest) {
	got := Attributes(test.attrs)
	if !assert.Len(t, got, len(test.expected)) {
		return
	}
	for i, actual := range got {
		if a, ok := actual.Value.Value.(*commonpb.AnyValue_DoubleValue); ok {
			e, ok := test.expected[i].Value.Value.(*commonpb.AnyValue_DoubleValue)
			if !ok {
				t.Errorf("expected AnyValue_DoubleValue, got %T", test.expected[i].Value.Value)
				return
			}
			if !assert.InDelta(t, e.DoubleValue, a.DoubleValue, 0.01) {
				return
			}
			e.DoubleValue = a.DoubleValue
		}
		assert.Equal(t, test.expected[i], actual)
	}
}

func newOtelBoolArray(key string, values []bool) *commonpb.KeyValue {

	arrayValues := []*commonpb.AnyValue{}

	for _, b := range values {
		arrayValues = append(arrayValues, &commonpb.AnyValue{
			Value: &commonpb.AnyValue_BoolValue{
				BoolValue: b,
			},
		})
	}

	result := &commonpb.KeyValue{
		Key: key,
		Value: &commonpb.AnyValue{
			Value: &commonpb.AnyValue_ArrayValue{
				ArrayValue: &commonpb.ArrayValue{
					Values: arrayValues,
				},
			},
		},
	}
	return result
}

func newOtelIntArray(key string, values []int) *commonpb.KeyValue {

	arrayValues := []*commonpb.AnyValue{}

	for _, i := range values {
		arrayValues = append(arrayValues, &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: int64(i),
			},
		})
	}

	result := &commonpb.KeyValue{
		Key: key,
		Value: &commonpb.AnyValue{
			Value: &commonpb.AnyValue_ArrayValue{
				ArrayValue: &commonpb.ArrayValue{
					Values: arrayValues,
				},
			},
		},
	}
	return result
}
