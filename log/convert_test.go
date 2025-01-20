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
			desc: "",
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
