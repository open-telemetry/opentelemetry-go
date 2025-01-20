// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestConvertAttributeValue(t *testing.T) {
	testCases := []struct {
		desc string
		v    attribute.Value
		want log.Value
	}{
		{
			desc: "Empty",
			v:    attribute.Value{},
			want: log.Value{},
		},
		{
			desc: "Bool",
			v:    attribute.BoolValue(true),
			want: log.BoolValue(true),
		},
		{
			desc: "BoolSlice",
			v:    attribute.BoolSliceValue([]bool{true, false}),
			want: log.SliceValue(log.BoolValue(true), log.BoolValue(false)),
		},
		{
			desc: "Int64",
			v:    attribute.Int64Value(13),
			want: log.Int64Value(13),
		},
		{
			desc: "Int64Slice",
			v:    attribute.Int64SliceValue([]int64{12, 34}),
			want: log.SliceValue(log.Int64Value(12), log.Int64Value(34)),
		},
		{
			desc: "Float64",
			v:    attribute.Float64Value(3.14),
			want: log.Float64Value(3.14),
		},
		{
			desc: "Float64Slice",
			v:    attribute.Float64SliceValue([]float64{3.14, 2.72}),
			want: log.SliceValue(log.Float64Value(3.14), log.Float64Value(2.72)),
		},
		{
			desc: "String",
			v:    attribute.StringValue("foo"),
			want: log.StringValue("foo"),
		},
		{
			desc: "StringSlice",
			v:    attribute.StringSliceValue([]string{"foo", "bar"}),
			want: log.SliceValue(log.StringValue("foo"), log.StringValue("bar")),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := log.ConvertAttributeValue(tc.v)
			if !got.Equal(tc.want) {
				t.Errorf("got: %v; want:%v", got, tc.want)
			}
		})
	}
}

func TestConvertAttributeKeyValue(t *testing.T) {
	testCases := []struct {
		desc string
		kv   attribute.KeyValue
		want log.KeyValue
	}{
		{
			desc: "Empty",
			kv:   attribute.KeyValue{},
			want: log.KeyValue{},
		},
		{
			desc: "Bool",
			kv:   attribute.Bool("k", true),
			want: log.Bool("k", true),
		},
		{
			desc: "BoolSlice",
			kv:   attribute.BoolSlice("k", []bool{true, false}),
			want: log.Slice("k", log.BoolValue(true), log.BoolValue(false)),
		},
		{
			desc: "Int64",
			kv:   attribute.Int64("k", 13),
			want: log.Int64("k", 13),
		},
		{
			desc: "Int64Slice",
			kv:   attribute.Int64Slice("k", []int64{12, 34}),
			want: log.Slice("k", log.Int64Value(12), log.Int64Value(34)),
		},
		{
			desc: "Float64",
			kv:   attribute.Float64("k", 3.14),
			want: log.Float64("k", 3.14),
		},
		{
			desc: "Float64Slice",
			kv:   attribute.Float64Slice("k", []float64{3.14, 2.72}),
			want: log.Slice("k", log.Float64Value(3.14), log.Float64Value(2.72)),
		},
		{
			desc: "String",
			kv:   attribute.String("k", "foo"),
			want: log.String("k", "foo"),
		},
		{
			desc: "StringSlice",
			kv:   attribute.StringSlice("k", []string{"foo", "bar"}),
			want: log.Slice("k", log.StringValue("foo"), log.StringValue("bar")),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := log.ConvertAttributeKeyValue(tc.kv)
			if !got.Equal(tc.want) {
				t.Errorf("got: %v; want:%v", got, tc.want)
			}
		})
	}
}
