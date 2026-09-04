// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package tracetransform provides conversion functionality for the otlptrace
// exporters.
package tracetransform

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"

	"go.opentelemetry.io/otel/attribute"
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
func ResourceAttributes(res *resource.Resource) []*commonpb.KeyValue {
	return Iterator(res.Iter())
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
	case attribute.BOOLSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: boolSliceValues(v.AsBoolSlice()),
			},
		}
	case attribute.INT64:
		av.Value = &commonpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case attribute.INT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: int64SliceValues(v.AsInt64Slice()),
			},
		}
	case attribute.FLOAT64:
		av.Value = &commonpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case attribute.FLOAT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: float64SliceValues(v.AsFloat64Slice()),
			},
		}
	case attribute.STRING:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case attribute.BYTESLICE:
		av.Value = &commonpb.AnyValue_BytesValue{
			BytesValue: v.AsByteSlice(),
		}
	case attribute.SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: values(v.AsSlice()),
			},
		}
	case attribute.MAP:
		av.Value = &commonpb.AnyValue_KvlistValue{
			KvlistValue: &commonpb.KeyValueList{
				Values: KeyValues(v.AsMap()),
			},
		}
	case attribute.STRINGSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: stringSliceValues(v.AsStringSlice()),
			},
		}
	case attribute.EMPTY:
	default:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

func boolSliceValues(vals []bool) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_BoolValue{
				BoolValue: v,
			},
		}
	}
	return converted
}

func int64SliceValues(vals []int64) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_IntValue{
				IntValue: v,
			},
		}
	}
	return converted
}

func float64SliceValues(vals []float64) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_DoubleValue{
				DoubleValue: v,
			},
		}
	}
	return converted
}

func stringSliceValues(vals []string) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &commonpb.AnyValue{
			Value: &commonpb.AnyValue_StringValue{
				StringValue: v,
			},
		}
	}
	return converted
}

func values(vals []attribute.Value) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = Value(v)
	}
	return converted
}

// KeyValuesWithArena transforms a slice of attribute KeyValues into OTLP key-values
// using the provided arena for allocation.
func KeyValuesWithArena(attrs []attribute.KeyValue, arena *Arena) []*commonpb.KeyValue {
	if len(attrs) == 0 {
		return nil
	}
	out := make([]*commonpb.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, KeyValueWithArena(kv, arena))
	}
	return out
}

// IteratorWithArena transforms an attribute iterator into OTLP key-values using arena.
func IteratorWithArena(iter attribute.Iterator, arena *Arena) []*commonpb.KeyValue {
	l := iter.Len()
	if l == 0 {
		return nil
	}
	out := make([]*commonpb.KeyValue, 0, l)
	for iter.Next() {
		out = append(out, KeyValueWithArena(iter.Attribute(), arena))
	}
	return out
}

// ResourceAttributesWithArena transforms a Resource into OTLP key-values using arena.
func ResourceAttributesWithArena(res *resource.Resource, arena *Arena) []*commonpb.KeyValue {
	return IteratorWithArena(res.Iter(), arena)
}

// KeyValueWithArena transforms an attribute KeyValue into an OTLP key-value using arena.
func KeyValueWithArena(kv attribute.KeyValue, arena *Arena) *commonpb.KeyValue {
	pbKV := arena.kvs.alloc()
	pbKV.Key = string(kv.Key)
	pbKV.Value = ValueWithArena(kv.Value, arena)
	return pbKV
}

// ValueWithArena transforms an attribute Value into an OTLP AnyValue using arena.
func ValueWithArena(v attribute.Value, arena *Arena) *commonpb.AnyValue {
	av := arena.avs.alloc()
	switch v.Type() {
	case attribute.BOOL:
		arena.avBoolValues = append(arena.avBoolValues, commonpb.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		})
		av.Value = &arena.avBoolValues[len(arena.avBoolValues)-1]
	case attribute.BOOLSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: boolSliceValues(v.AsBoolSlice()),
			},
		}
	case attribute.INT64:
		arena.avIntValues = append(arena.avIntValues, commonpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		})
		av.Value = &arena.avIntValues[len(arena.avIntValues)-1]
	case attribute.INT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: int64SliceValues(v.AsInt64Slice()),
			},
		}
	case attribute.FLOAT64:
		arena.avFloatValues = append(arena.avFloatValues, commonpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		})
		av.Value = &arena.avFloatValues[len(arena.avFloatValues)-1]
	case attribute.FLOAT64SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: float64SliceValues(v.AsFloat64Slice()),
			},
		}
	case attribute.STRING:
		arena.avStrValues = append(arena.avStrValues, commonpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		})
		av.Value = &arena.avStrValues[len(arena.avStrValues)-1]
	case attribute.BYTESLICE:
		av.Value = &commonpb.AnyValue_BytesValue{
			BytesValue: v.AsByteSlice(),
		}
	case attribute.SLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: valuesWithArena(v.AsSlice(), arena),
			},
		}
	case attribute.MAP:
		av.Value = &commonpb.AnyValue_KvlistValue{
			KvlistValue: &commonpb.KeyValueList{
				Values: KeyValuesWithArena(v.AsMap(), arena),
			},
		}
	case attribute.STRINGSLICE:
		av.Value = &commonpb.AnyValue_ArrayValue{
			ArrayValue: &commonpb.ArrayValue{
				Values: stringSliceValues(v.AsStringSlice()),
			},
		}
	case attribute.EMPTY:
	default:
		av.Value = &commonpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

func valuesWithArena(vals []attribute.Value, arena *Arena) []*commonpb.AnyValue {
	converted := make([]*commonpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = ValueWithArena(v, arena)
	}
	return converted
}
