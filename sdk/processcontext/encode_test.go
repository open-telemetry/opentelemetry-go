// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package processcontext

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	otlpcommon "go.opentelemetry.io/proto/otlp/common/v1"
	processcontextpb "go.opentelemetry.io/proto/otlp/processcontext/v1development"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

// ---- encodeProcessContext ------------------------------------------------

func TestEncodeProcessContextEmpty(t *testing.T) {
	b, err := encodeProcessContext(resource.Empty())
	require.NoError(t, err)

	var msg processcontextpb.ProcessContext
	require.NoError(t, proto.Unmarshal(b, &msg))
	require.NotNil(t, msg.Resource)
	assert.Empty(t, msg.Resource.Attributes)
}

func TestEncodeProcessContextRoundTrip(t *testing.T) {
	r, err := resource.New(t.Context(), resource.WithAttributes(
		attribute.String("service.name", "my-service"),
		attribute.String("service.version", "1.0.0"),
	))
	require.NoError(t, err)

	b, err := encodeProcessContext(r)
	require.NoError(t, err)
	assert.NotEmpty(t, b)
	assert.LessOrEqual(t, len(b), MaxPayloadSize)

	var msg processcontextpb.ProcessContext
	require.NoError(t, proto.Unmarshal(b, &msg))
	require.NotNil(t, msg.Resource)
	assert.NotEmpty(t, msg.Resource.Attributes)
}

func TestEncodeProcessContextAllTypes(t *testing.T) {
	r, err := resource.New(t.Context(), resource.WithAttributes(
		attribute.String("s", "hello"),
		attribute.Bool("b", true),
		attribute.Int64("i", -42),
		attribute.Float64("f", 3.14),
		attribute.StringSlice("ss", []string{"a", "b"}),
		attribute.BoolSlice("bs", []bool{true, false}),
		attribute.Int64Slice("is", []int64{1, 2, 3}),
		attribute.Float64Slice("fs", []float64{1.1, 2.2}),
	))
	require.NoError(t, err)

	b, err := encodeProcessContext(r)
	require.NoError(t, err)
	assert.NotEmpty(t, b)
	assert.LessOrEqual(t, len(b), MaxPayloadSize)

	var msg processcontextpb.ProcessContext
	require.NoError(t, proto.Unmarshal(b, &msg))
	require.NotNil(t, msg.Resource)
	assert.NotEmpty(t, msg.Resource.Attributes)
}

// ---- resourceToOTLPAttrs -------------------------------------------------

func TestResourceToOTLPAttrsMultiple(t *testing.T) {
	r, err := resource.New(t.Context(), resource.WithAttributes(
		attribute.String("a", "1"),
		attribute.String("b", "2"),
	))
	require.NoError(t, err)

	kvs := resourceToOTLPAttrs(r)

	m := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		m[kv.Key] = kv.Value.GetStringValue()
	}
	assert.Equal(t, "1", m["a"])
	assert.Equal(t, "2", m["b"])
}

// ---- attrValueToOTLP -----------------------------------------------------

func TestAttrValueToOTLPString(t *testing.T) {
	v := attrValueToOTLP(attribute.StringValue("hello"))
	assert.Equal(t, "hello", v.GetStringValue())
}

func TestAttrValueToOTLPBoolTrue(t *testing.T) {
	v := attrValueToOTLP(attribute.BoolValue(true))
	assert.True(t, v.GetBoolValue())
}

func TestAttrValueToOTLPBoolFalse(t *testing.T) {
	v := attrValueToOTLP(attribute.BoolValue(false))
	assert.NotNil(t, v.GetValue())
	assert.False(t, v.GetBoolValue())
}

func TestAttrValueToOTLPInt64Positive(t *testing.T) {
	v := attrValueToOTLP(attribute.Int64Value(42))
	assert.Equal(t, int64(42), v.GetIntValue())
}

func TestAttrValueToOTLPInt64Negative(t *testing.T) {
	v := attrValueToOTLP(attribute.Int64Value(-1))
	assert.Equal(t, int64(-1), v.GetIntValue())
}

func TestAttrValueToOTLPFloat64(t *testing.T) {
	v := attrValueToOTLP(attribute.Float64Value(3.14))
	assert.InDelta(t, 3.14, v.GetDoubleValue(), 1e-10)
}

func TestAttrValueToOTLPStringSlice(t *testing.T) {
	v := attrValueToOTLP(attribute.StringSliceValue([]string{"x", "y"}))
	arr := v.GetArrayValue()
	require.NotNil(t, arr)
	require.Len(t, arr.Values, 2)
	assert.Equal(t, "x", arr.Values[0].GetStringValue())
	assert.Equal(t, "y", arr.Values[1].GetStringValue())
}

func TestAttrValueToOTLPBoolSlice(t *testing.T) {
	v := attrValueToOTLP(attribute.BoolSliceValue([]bool{true, false}))
	arr := v.GetArrayValue()
	require.NotNil(t, arr)
	require.Len(t, arr.Values, 2)
	assert.True(t, arr.Values[0].GetBoolValue())
	assert.False(t, arr.Values[1].GetBoolValue())
}

func TestAttrValueToOTLPInt64Slice(t *testing.T) {
	v := attrValueToOTLP(attribute.Int64SliceValue([]int64{1, -2, 3}))
	arr := v.GetArrayValue()
	require.NotNil(t, arr)
	require.Len(t, arr.Values, 3)
	assert.Equal(t, int64(1), arr.Values[0].GetIntValue())
	assert.Equal(t, int64(-2), arr.Values[1].GetIntValue())
	assert.Equal(t, int64(3), arr.Values[2].GetIntValue())
}

func TestAttrValueToOTLPFloat64Slice(t *testing.T) {
	v := attrValueToOTLP(attribute.Float64SliceValue([]float64{1.1, 2.2}))
	arr := v.GetArrayValue()
	require.NotNil(t, arr)
	require.Len(t, arr.Values, 2)
	assert.InDelta(t, 1.1, arr.Values[0].GetDoubleValue(), 1e-10)
	assert.InDelta(t, 2.2, arr.Values[1].GetDoubleValue(), 1e-10)
}

func TestAttrValueToOTLPEmptySlices(t *testing.T) {
	for _, v := range []attribute.Value{
		attribute.StringSliceValue(nil),
		attribute.BoolSliceValue(nil),
		attribute.Int64SliceValue(nil),
		attribute.Float64SliceValue(nil),
	} {
		av := attrValueToOTLP(v)
		arr := av.GetArrayValue()
		require.NotNil(t, arr)
		assert.Empty(t, arr.Values)
	}
}

func TestAttrValueToOTLPByteSlice(t *testing.T) {
	v := attrValueToOTLP(attribute.ByteSliceValue([]byte{1, 2, 3}))
	assert.Equal(t, []byte{1, 2, 3}, v.GetBytesValue())
}

func TestAttrValueToOTLPSlice(t *testing.T) {
	v := attrValueToOTLP(attribute.SliceValue(
		attribute.StringValue("a"),
		attribute.Int64Value(7),
	))
	arr := v.GetArrayValue()
	require.NotNil(t, arr)
	require.Len(t, arr.Values, 2)
	assert.Equal(t, "a", arr.Values[0].GetStringValue())
	assert.Equal(t, int64(7), arr.Values[1].GetIntValue())
}

func TestAttrValueToOTLPMap(t *testing.T) {
	v := attrValueToOTLP(attribute.MapValue(
		attribute.String("k1", "v1"),
		attribute.Int64("k2", 42),
	))
	kvlist := v.GetKvlistValue()
	require.NotNil(t, kvlist)
	require.Len(t, kvlist.Values, 2)
	m := make(map[string]*otlpcommon.AnyValue, len(kvlist.Values))
	for _, kv := range kvlist.Values {
		m[kv.Key] = kv.Value
	}
	assert.Equal(t, "v1", m["k1"].GetStringValue())
	assert.Equal(t, int64(42), m["k2"].GetIntValue())
}

func TestAttrValueToOTLPEmpty(t *testing.T) {
	// EMPTY is a valid value; Value field must be nil (not "INVALID").
	v := attrValueToOTLP(attribute.Value{})
	assert.NotNil(t, v)
	assert.Nil(t, v.GetValue())
}
