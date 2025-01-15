// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlplog/transform/log_attr_test.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
)

var (
	logAttrBool    = log.Bool("bool", true)
	logAttrInt     = log.Int("int", 1)
	logAttrInt64   = log.Int64("int64", 1)
	logAttrFloat64 = log.Float64("float64", 1)
	logAttrString  = log.String("string", "o")
	logAttrBytes   = log.Bytes("bytes", []byte("test"))
	logAttrSlice   = log.Slice("slice", log.BoolValue(true))
	logAttrMap     = log.Map("map", logAttrString)
	logAttrEmpty   = log.Empty("")

	kvBytes = &cpb.KeyValue{
		Key: "bytes",
		Value: &cpb.AnyValue{
			Value: &cpb.AnyValue_BytesValue{
				BytesValue: []byte("test"),
			},
		},
	}
	kvSlice = &cpb.KeyValue{
		Key: "slice",
		Value: &cpb.AnyValue{
			Value: &cpb.AnyValue_ArrayValue{
				ArrayValue: &cpb.ArrayValue{
					Values: []*cpb.AnyValue{valBoolTrue},
				},
			},
		},
	}
	kvMap = &cpb.KeyValue{
		Key: "map",
		Value: &cpb.AnyValue{
			Value: &cpb.AnyValue_KvlistValue{
				KvlistValue: &cpb.KeyValueList{
					Values: []*cpb.KeyValue{kvString},
				},
			},
		},
	}
	kvEmpty = &cpb.KeyValue{
		Value: &cpb.AnyValue{
			Value: &cpb.AnyValue_StringValue{StringValue: "INVALID"},
		},
	}
)

func TestLogAttrs(t *testing.T) {
	type logAttrTest struct {
		name string
		in   []log.KeyValue
		want []*cpb.KeyValue
	}

	for _, test := range []logAttrTest{
		{"nil", nil, nil},
		{"len(0)", []log.KeyValue{}, nil},
		{
			"empty",
			[]log.KeyValue{logAttrEmpty},
			[]*cpb.KeyValue{kvEmpty},
		},
		{
			"bool",
			[]log.KeyValue{logAttrBool},
			[]*cpb.KeyValue{kvBool},
		},
		{
			"int",
			[]log.KeyValue{logAttrInt},
			[]*cpb.KeyValue{kvInt},
		},
		{
			"int64",
			[]log.KeyValue{logAttrInt64},
			[]*cpb.KeyValue{kvInt64},
		},
		{
			"float64",
			[]log.KeyValue{logAttrFloat64},
			[]*cpb.KeyValue{kvFloat64},
		},
		{
			"string",
			[]log.KeyValue{logAttrString},
			[]*cpb.KeyValue{kvString},
		},
		{
			"bytes",
			[]log.KeyValue{logAttrBytes},
			[]*cpb.KeyValue{kvBytes},
		},
		{
			"slice",
			[]log.KeyValue{logAttrSlice},
			[]*cpb.KeyValue{kvSlice},
		},
		{
			"map",
			[]log.KeyValue{logAttrMap},
			[]*cpb.KeyValue{kvMap},
		},
		{
			"all",
			[]log.KeyValue{
				logAttrBool,
				logAttrInt,
				logAttrInt64,
				logAttrFloat64,
				logAttrString,
				logAttrBytes,
				logAttrSlice,
				logAttrMap,
				logAttrEmpty,
			},
			[]*cpb.KeyValue{
				kvBool,
				kvInt,
				kvInt64,
				kvFloat64,
				kvString,
				kvBytes,
				kvSlice,
				kvMap,
				kvEmpty,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			assert.ElementsMatch(t, test.want, LogAttrs(test.in))
		})
	}
}
