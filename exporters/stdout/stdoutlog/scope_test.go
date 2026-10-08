// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func TestScopeAttributes(t *testing.T) {
	testCases := []struct {
		name       string
		attributes []attribute.KeyValue
		want       string
	}{
		{
			name: "empty scope",
			want: `{}`,
		},
		{
			name: "primitives",
			attributes: []attribute.KeyValue{
				attribute.Bool("bool", true),
				attribute.Int64("int", 42),
				attribute.Int64("min", math.MinInt64),
				attribute.Int64("max", math.MaxInt64),
				attribute.Float64("float", 1.25),
				attribute.String("string", "hello \"world\"\n\\<>&\u2028\u2029"),
				attribute.ByteSlice("bytes", []byte{0, 1, 2, 255}),
			},
			want: `{
				"bool": true,
				"int": 42,
				"min": -9223372036854775808,
				"max": 9223372036854775807,
				"float": 1.25,
				"string": "hello \"world\"\n\\<>&\u2028\u2029",
				"bytes": "AAEC/w=="
			}`,
		},
		{
			name: "zero and empty values",
			attributes: []attribute.KeyValue{
				{Key: "null"},
				attribute.Bool("bool", false),
				attribute.Int64("int", 0),
				attribute.Float64("float", 0),
				attribute.String("string", ""),
				attribute.ByteSlice("bytes", nil),
				attribute.ByteSlice("empty bytes", []byte{}),
				attribute.BoolSlice("bools", nil),
				attribute.Int64Slice("ints", nil),
				attribute.Float64Slice("floats", nil),
				attribute.StringSlice("strings", nil),
				attribute.Slice("slice"),
				attribute.Map("map"),
			},
			want: `{
				"null": null,
				"bool": false,
				"int": 0,
				"float": 0,
				"string": "",
				"bytes": "",
				"empty bytes": "",
				"bools": [],
				"ints": [],
				"floats": [],
				"strings": [],
				"slice": [],
				"map": {}
			}`,
		},
		{
			name: "collections",
			attributes: []attribute.KeyValue{
				attribute.BoolSlice("bools", []bool{true, false}),
				attribute.Int64Slice("ints", []int64{math.MinInt64, 0, math.MaxInt64}),
				attribute.Float64Slice("floats", []float64{1.25, -2.5}),
				attribute.StringSlice("strings", []string{"first", ""}),
				attribute.Slice("slice",
					attribute.BoolValue(false),
					attribute.Int64Value(math.MaxInt64),
					attribute.Float64Value(1.25),
					attribute.StringValue("value"),
					attribute.ByteSliceValue([]byte{0, 1, 2, 255}),
					attribute.Value{},
					attribute.MapValue(attribute.String("nested", "map")),
					attribute.SliceValue(attribute.StringValue("nested"), attribute.Int64Value(2)),
				),
				attribute.Map("map",
					attribute.String("string", "value"),
					attribute.Map("nested", attribute.Slice("array",
						attribute.MapValue(attribute.Bool("bool", true)),
						attribute.Int64SliceValue([]int64{1, 2}),
					)),
				),
			},
			want: `{
				"bools": [true, false],
				"ints": [-9223372036854775808, 0, 9223372036854775807],
				"floats": [1.25, -2.5],
				"strings": ["first", ""],
				"slice": [false, 9223372036854775807, 1.25, "value", "AAEC/w==", null, {"nested": "map"}, ["nested", 2]],
				"map": {"string": "value", "nested": {"array": [{"bool": true}, [1, 2]]}}
			}`,
		},
		{
			name: "case and duplicate keys",
			attributes: []attribute.KeyValue{
				attribute.String("Key", "upper"),
				attribute.String("key", "first"),
				attribute.String("key", "last"),
				attribute.String("escaped\n\"key", "value"),
			},
			want: `{"Key": "upper", "key": "last", "escaped\n\"key": "value"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testScopeOutput(t, tc.attributes, tc.want)
		})
	}
}

func TestScopeNonFiniteAttributes(t *testing.T) {
	testCases := []struct {
		name       string
		attributes []attribute.KeyValue
		want       string
	}{
		{
			name: "scalars",
			attributes: []attribute.KeyValue{
				attribute.Float64("nan", math.NaN()),
				attribute.Float64("positive", math.Inf(1)),
				attribute.Float64("negative", math.Inf(-1)),
			},
			want: `{"nan": "NaN", "positive": "Infinity", "negative": "-Infinity"}`,
		},
		{
			name: "float slice",
			attributes: []attribute.KeyValue{
				attribute.Float64Slice("floats", []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1.25}),
			},
			want: `{"floats": ["NaN", "Infinity", "-Infinity", 1.25]}`,
		},
		{
			name: "heterogeneous array",
			attributes: []attribute.KeyValue{
				attribute.Slice("slice",
					attribute.Float64Value(math.NaN()),
					attribute.StringValue("value"),
					attribute.Float64Value(math.Inf(1)),
					attribute.Value{},
					attribute.Float64Value(math.Inf(-1)),
				),
			},
			want: `{"slice": ["NaN", "value", "Infinity", null, "-Infinity"]}`,
		},
		{
			name: "recursive maps and arrays",
			attributes: []attribute.KeyValue{
				attribute.Map("map", attribute.Slice("array",
					attribute.MapValue(
						attribute.Float64("nan", math.NaN()),
						attribute.Float64("positive", math.Inf(1)),
						attribute.Float64("negative", math.Inf(-1)),
						attribute.Float64Slice("floats", []float64{math.NaN(), math.Inf(1), math.Inf(-1)}),
						attribute.Slice("slice", attribute.Float64Value(math.NaN()),
							attribute.Float64Value(math.Inf(1)), attribute.Float64Value(math.Inf(-1))),
					),
				)),
			},
			want: `{"map": {"array": [{"nan": "NaN", "positive": "Infinity", "negative": "-Infinity", "floats": ["NaN", "Infinity", "-Infinity"], "slice": ["NaN", "Infinity", "-Infinity"]}]}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testScopeOutput(t, tc.attributes, tc.want)
		})
	}
}

func testScopeOutput(t *testing.T, attributes []attribute.KeyValue, wantAttributes string) {
	t.Helper()
	for _, pretty := range []bool{false, true} {
		name := "compact"
		if pretty {
			name = "pretty"
		}
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			options := []stdoutlog.Option{stdoutlog.WithWriter(&buf), stdoutlog.WithoutTimestamps()}
			if pretty {
				options = append(options, stdoutlog.WithPrettyPrint())
			}
			exporter, err := stdoutlog.New(options...)
			require.NoError(t, err)
			provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exporter)))
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
			logger := provider.Logger("scope-name",
				log.WithInstrumentationVersion("scope-version"),
				log.WithSchemaURL("https://example.com/scope-schema"),
				log.WithInstrumentationAttributes(attributes...),
			)

			bodies := []string{"first", "second"}
			for _, body := range bodies {
				var record log.Record
				record.SetBody(attribute.StringValue(body))
				logger.Emit(t.Context(), record)
			}
			assert.Equal(t, pretty, strings.Contains(buf.String(), "\n\t"))

			decoder := json.NewDecoder(&buf)
			decoder.UseNumber()
			var records []map[string]any
			for {
				var record map[string]any
				err := decoder.Decode(&record)
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				records = append(records, record)
			}
			require.Len(t, records, len(bodies))

			var wantScope any
			wantDecoder := json.NewDecoder(
				strings.NewReader(
					`{"Name": "scope-name", "Version": "scope-version", "SchemaURL": "https://example.com/scope-schema", "Attributes": ` + wantAttributes + `}`,
				),
			)
			wantDecoder.UseNumber()
			require.NoError(t, wantDecoder.Decode(&wantScope))
			for i, record := range records {
				assert.Equal(t, wantScope, record["Scope"])
				assert.Equal(t, map[string]any{"Type": "STRING", "Value": bodies[i]}, record["Body"])
				assert.NotContains(t, record, "Timestamp")
				assert.NotContains(t, record, "ObservedTimestamp")
			}
		})
	}
}
