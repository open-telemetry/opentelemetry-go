// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := log.ConvertAttributeValue(tc.v)
			assert.True(t, got.Equal(tc.want), "%v.Equal(%v)", got, tc.want)
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
			assert.True(t, got.Equal(tc.want), "%v.Equal(%v)", got, tc.want)
		})
	}
}
