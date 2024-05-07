// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutout"

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

func TestExporter(t *testing.T) {
	var buf bytes.Buffer
	now := time.Now()

	testCases := []struct {
		name     string
		exporter *Exporter
		want     string
	}{
		{
			name:     "zero value",
			exporter: &Exporter{},
			want:     "",
		},
		{
			name: "new",
			exporter: func() *Exporter {
				defaultWriterSwap := defaultWriter
				defer func() {
					defaultWriter = defaultWriterSwap
				}()
				defaultWriter = &buf

				exporter, err := New()
				require.NoError(t, err)
				require.NotNil(t, exporter)

				return exporter
			}(),
			want: getJSON(&now),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write to buffer for testing
			defaultWriterSwap := defaultWriter
			defer func() {
				defaultWriter = defaultWriterSwap
			}()
			defaultWriter = &buf

			buf.Reset()

			var err error

			exporter := tc.exporter

			record := getRecord(now)

			// Export a record
			err = exporter.Export(context.Background(), []sdklog.Record{record})
			assert.NoError(t, err)

			// Check the writer
			assert.Equal(t, tc.want, buf.String())

			// Flush the exporter
			err = exporter.ForceFlush(context.Background())
			assert.NoError(t, err)

			// Shutdown the exporter
			err = exporter.Shutdown(context.Background())
			assert.NoError(t, err)

			// Export a record after shutdown, this should not be written
			err = exporter.Export(context.Background(), []sdklog.Record{record})
			assert.NoError(t, err)

			// Check the writer
			assert.Equal(t, tc.want, buf.String())
		})
	}
}

func TestExporterExport(t *testing.T) {
	now := time.Now()

	record := getRecord(now)
	records := []sdklog.Record{record, record}

	testCases := []struct {
		name       string
		options    []Option
		ctx        context.Context
		records    []sdklog.Record
		wantResult string
		wantError  error
	}{
		{
			name:       "default",
			options:    []Option{},
			ctx:        context.Background(),
			records:    records,
			wantResult: getJSONs(&now),
		},
		{
			name:       "NoRecords",
			options:    []Option{},
			ctx:        context.Background(),
			records:    nil,
			wantResult: "",
		},
		{
			name:       "WithPrettyPrint",
			options:    []Option{WithPrettyPrint()},
			ctx:        context.Background(),
			records:    records,
			wantResult: getPrettyJSONs(&now),
		},
		{
			name:       "WithoutTimestamps",
			options:    []Option{WithoutTimestamps()},
			ctx:        context.Background(),
			records:    records,
			wantResult: getJSONs(nil),
		},
		{
			name:       "WithoutTimestamps and WithPrettyPrint",
			options:    []Option{WithoutTimestamps(), WithPrettyPrint()},
			ctx:        context.Background(),
			records:    records,
			wantResult: getPrettyJSONs(nil),
		},
		{
			name: "WithCanceledContext",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			records:    records,
			wantResult: "",
			wantError:  context.Canceled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write to buffer for testing
			var buf bytes.Buffer

			exporter, err := New(append(tc.options, WithWriter(&buf))...)
			assert.NoError(t, err)

			err = exporter.Export(tc.ctx, tc.records)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantResult, buf.String())
		})
	}
}

func getJSON(now *time.Time) string {
	var timestamps string
	if now != nil {
		serializedNow, _ := json.Marshal(now)
		timestamps = "\"Timestamp\":" + string(serializedNow) + ",\"ObservedTimestamp\":" + string(serializedNow) + ","
	}

	return "{" + timestamps + "\"Severity\":9,\"SeverityText\":\"INFO\",\"Body\":{\"Type\":\"String\",\"Value\":\"test\"},\"Attributes\":[{\"Key\":\"key\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key2\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key3\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key4\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key5\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"bool\",\"Value\":{\"Type\":\"Bool\",\"Value\":true}}],\"TraceID\":\"0102030405060708090a0b0c0d0e0f10\",\"SpanID\":\"0102030405060708\",\"TraceFlags\":\"01\",\"Resource\":[{\"Key\":\"foo\",\"Value\":{\"Type\":\"STRING\",\"Value\":\"bar\"}}],\"Scope\":{\"Name\":\"name\",\"Version\":\"version\",\"SchemaURL\":\"https://example.com/custom-schema\"},\"DroppedAttributes\":10}\n"
}

func getJSONs(now *time.Time) string {
	return getJSON(now) + getJSON(now)
}

func getPrettyJSON(now *time.Time) string {
	var timestamps string
	if now != nil {
		serializedNow, _ := json.Marshal(now)
		timestamps = "\n\t\"Timestamp\": " + string(serializedNow) + ",\n\t\"ObservedTimestamp\": " + string(serializedNow) + ","
	}

	return `{` + timestamps + `
	"Severity": 9,
	"SeverityText": "INFO",
	"Body": {
		"Type": "String",
		"Value": "test"
	},
	"Attributes": [
		{
			"Key": "key",
			"Value": {
				"Type": "String",
				"Value": "value"
			}
		},
		{
			"Key": "key2",
			"Value": {
				"Type": "String",
				"Value": "value"
			}
		},
		{
			"Key": "key3",
			"Value": {
				"Type": "String",
				"Value": "value"
			}
		},
		{
			"Key": "key4",
			"Value": {
				"Type": "String",
				"Value": "value"
			}
		},
		{
			"Key": "key5",
			"Value": {
				"Type": "String",
				"Value": "value"
			}
		},
		{
			"Key": "bool",
			"Value": {
				"Type": "Bool",
				"Value": true
			}
		}
	],
	"TraceID": "0102030405060708090a0b0c0d0e0f10",
	"SpanID": "0102030405060708",
	"TraceFlags": "01",
	"Resource": [
		{
			"Key": "foo",
			"Value": {
				"Type": "STRING",
				"Value": "bar"
			}
		}
	],
	"Scope": {
		"Name": "name",
		"Version": "version",
		"SchemaURL": "https://example.com/custom-schema"
	},
	"DroppedAttributes": 10
}
`
}

func getPrettyJSONs(now *time.Time) string {
	return getPrettyJSON(now) + getPrettyJSON(now)
}

func TestExporterShutdown(t *testing.T) {
	exporter, err := New()
	assert.NoError(t, err)

	assert.NoError(t, exporter.Shutdown(context.Background()))
}

func TestExporterForceFlush(t *testing.T) {
	exporter, err := New()
	assert.NoError(t, err)

	assert.NoError(t, exporter.ForceFlush(context.Background()))
}

func getRecord(now time.Time) sdklog.Record {
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")

	rf := logtest.RecordFactory{
		Timestamp:         now,
		ObservedTimestamp: now,
		Severity:          log.SeverityInfo1,
		SeverityText:      "INFO",
		Body:              log.StringValue("test"),
		Attributes: []log.KeyValue{
			// More than 5 attributes to test back slice
			log.String("key", "value"),
			log.String("key2", "value"),
			log.String("key3", "value"),
			log.String("key4", "value"),
			log.String("key5", "value"),
			log.Bool("bool", true),
		},
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,

		Resource: resource.NewWithAttributes(
			"https://example.com/custom-resource-schema",
			attribute.String("foo", "bar"),
		),
		InstrumentationScope: &instrumentation.Scope{Name: "name", Version: "version", SchemaURL: "https://example.com/custom-schema"},
		DroppedAttributes:    10,
	}

	return rf.NewRecord()
}

func TestExporterConcurrentSafe(t *testing.T) {
	testCases := []struct {
		name     string
		exporter *Exporter
	}{
		{
			name:     "zero value",
			exporter: &Exporter{},
		},
		{
			name: "new",
			exporter: func() *Exporter {
				exporter, err := New()
				require.NoError(t, err)
				require.NotNil(t, exporter)

				return exporter
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exporter := tc.exporter

			const goroutines = 10
			var wg sync.WaitGroup
			wg.Add(goroutines)
			for i := 0; i < goroutines; i++ {
				go func() {
					defer wg.Done()
					err := exporter.Export(context.Background(), []sdklog.Record{{}})
					assert.NoError(t, err)
					err = exporter.ForceFlush(context.Background())
					assert.NoError(t, err)
					err = exporter.Shutdown(context.Background())
					assert.NoError(t, err)
				}()
			}
			wg.Wait()
		})
	}
}
