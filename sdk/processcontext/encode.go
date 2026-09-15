// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package processcontext

import (
	otlpcommon "go.opentelemetry.io/proto/otlp/common/v1"
	processcontextpb "go.opentelemetry.io/proto/otlp/processcontext/v1development"
	otlpresource "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

// encodeProcessContext serializes r as a ProcessContext protobuf message.
func encodeProcessContext(r *resource.Resource) ([]byte, error) {
	msg := &processcontextpb.ProcessContext{
		Resource: &otlpresource.Resource{
			Attributes: resourceToOTLPAttrs(r),
		},
	}
	return proto.Marshal(msg)
}

func resourceToOTLPAttrs(r *resource.Resource) []*otlpcommon.KeyValue {
	iter := r.Iter()
	kvs := make([]*otlpcommon.KeyValue, 0, iter.Len())
	for iter.Next() {
		a := iter.Attribute()
		kvs = append(kvs, &otlpcommon.KeyValue{
			Key:   string(a.Key),
			Value: attrValueToOTLP(a.Value),
		})
	}
	return kvs
}

func attrValueToOTLP(v attribute.Value) *otlpcommon.AnyValue {
	av := new(otlpcommon.AnyValue)
	switch v.Type() {
	case attribute.STRING:
		av.Value = &otlpcommon.AnyValue_StringValue{StringValue: v.AsString()}
	case attribute.BOOL:
		av.Value = &otlpcommon.AnyValue_BoolValue{BoolValue: v.AsBool()}
	case attribute.INT64:
		av.Value = &otlpcommon.AnyValue_IntValue{IntValue: v.AsInt64()}
	case attribute.FLOAT64:
		av.Value = &otlpcommon.AnyValue_DoubleValue{DoubleValue: v.AsFloat64()}
	case attribute.STRINGSLICE:
		ss := v.AsStringSlice()
		vals := make([]*otlpcommon.AnyValue, len(ss))
		for i, s := range ss {
			vals[i] = &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_StringValue{StringValue: s}}
		}
		av.Value = &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{Values: vals}}
	case attribute.BOOLSLICE:
		bs := v.AsBoolSlice()
		vals := make([]*otlpcommon.AnyValue, len(bs))
		for i, b := range bs {
			vals[i] = &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_BoolValue{BoolValue: b}}
		}
		av.Value = &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{Values: vals}}
	case attribute.INT64SLICE:
		is := v.AsInt64Slice()
		vals := make([]*otlpcommon.AnyValue, len(is))
		for i, n := range is {
			vals[i] = &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_IntValue{IntValue: n}}
		}
		av.Value = &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{Values: vals}}
	case attribute.FLOAT64SLICE:
		fs := v.AsFloat64Slice()
		vals := make([]*otlpcommon.AnyValue, len(fs))
		for i, f := range fs {
			vals[i] = &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_DoubleValue{DoubleValue: f}}
		}
		av.Value = &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{Values: vals}}
	case attribute.BYTESLICE:
		av.Value = &otlpcommon.AnyValue_BytesValue{BytesValue: v.AsByteSlice()}
	case attribute.SLICE:
		elems := v.AsSlice()
		vals := make([]*otlpcommon.AnyValue, len(elems))
		for i, e := range elems {
			vals[i] = attrValueToOTLP(e)
		}
		av.Value = &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{Values: vals}}
	case attribute.MAP:
		kvs := v.AsMap()
		pairs := make([]*otlpcommon.KeyValue, len(kvs))
		for i, kv := range kvs {
			pairs[i] = &otlpcommon.KeyValue{Key: string(kv.Key), Value: attrValueToOTLP(kv.Value)}
		}
		av.Value = &otlpcommon.AnyValue_KvlistValue{
			KvlistValue: &otlpcommon.KeyValueList{Values: pairs},
		}
	case attribute.EMPTY:
		// Value field stays nil; an empty value is a valid value.
	}
	return av
}
