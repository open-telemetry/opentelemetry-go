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

package tracetransform

import (
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"

	"go.opentelemetry.io/otel/sdk/resource"
)

// KeyValues transforms a slice of attribute KeyValues into OTLP key-values.
func KeyValues(attrs []attribute.KeyValue) []*commonpb.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*commonpb.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, KeyValue(kv))
	}
	return out
}

// Iterator transforms an attribute iterator into OTLP key-values.
func Iterator(iter attribute.Iterator) []*commonpb.KeyValue {
	l := iter.Len()
	if l == 0 {
		return nil
	}

	out := make([]*commonpb.KeyValue, 0, l)
	for iter.Next() {
		out = append(out, KeyValue(iter.Attribute()))
	}
	return out
}

// ResourceAttributes transforms a Resource OTLP key-values.
func ResourceAttributes(resource *resource.Resource) []*commonpb.KeyValue {
	return Iterator(resource.Iter())
}

// KeyValue transforms an attribute KeyValue into an OTLP key-value.
func KeyValue(kv attribute.KeyValue) *commonpb.KeyValue {
	return &commonpb.KeyValue{Key: string(kv.Key), Value: Value(kv.Value)}
}

// Value transforms an attribute Value into an OTLP AnyValue.
func Value(v attribute.Value) *commonpb.AnyValue {
	av := new(commonpb.AnyValue)
	switch v.Type() {
	case attribute.BOOL:
		av.Value = &commonpb.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		}
	case attribute.INT64:
		av.Value = &commonpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case attribute.FLOAT64:
		av.Value = &commonpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case attribute.STRING:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case attribute.ARRAY:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: arrayValues(v),
			},
		}
	default:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

func arrayValues(v attribute.Value) []*commonpb.AnyValue {
	a := v.AsArray()
	aType := reflect.TypeOf(a)
	var valueFunc func(reflect.Value) *commonpb.AnyValue
	switch aType.Elem().Kind() {
	case reflect.Bool:
		valueFunc = func(v reflect.Value) *commonpb.AnyValue {
			return &commonpb.AnyValue{
				Value: &commonpb.AnyValue_BoolValue{
					BoolValue: v.Bool(),
				},
			}
		}
	case reflect.Int, reflect.Int64:
		valueFunc = func(v reflect.Value) *commonpb.AnyValue {
			return &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: v.Int(),
				},
			}
		}
	case reflect.Uintptr:
		valueFunc = func(v reflect.Value) *commonpb.AnyValue {
			return &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: int64(v.Uint()),
				},
			}
		}
	case reflect.Float64:
		valueFunc = func(v reflect.Value) *commonpb.AnyValue {
			return &commonpb.AnyValue{
				Value: &commonpb.AnyValue_DoubleValue{
					DoubleValue: v.Float(),
				},
			}
		}
	case reflect.String:
		valueFunc = func(v reflect.Value) *commonpb.AnyValue {
			return &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: v.String(),
				},
			}
		}
	}

	results := make([]*commonpb.AnyValue, aType.Len())
	for i, aValue := 0, reflect.ValueOf(a); i < aValue.Len(); i++ {
		results[i] = valueFunc(aValue.Index(i))
	}
	return results
}
