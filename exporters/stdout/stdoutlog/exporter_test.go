// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog_test // import "go.opentelemetry.io/otel/exporters/stdout/stdoutout"

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

func TestExporter(t *testing.T) {
	var buf bytes.Buffer

	exporter, err := stdoutlog.New(stdoutlog.WithWriter(&buf))
	assert.NoError(t, err)
	assert.NotNil(t, exporter)

	now := time.Now()
	record := getRecord(now)

	// Export a record
	err = exporter.Export(context.Background(), []sdklog.Record{record})
	assert.NoError(t, err)

	// Check the writer
	assert.Equal(t, expectedJSON(now, false), buf.String())

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
	assert.Equal(t, expectedJSON(now, false), buf.String())
}

func TestExporterExport(t *testing.T) {
	now := time.Now()

	record := getRecord(now)
	records := []sdklog.Record{record, record}

	// Get canceled context
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name           string
		options        []stdoutlog.Option
		ctx            context.Context
		records        []sdklog.Record
		expectedResult string
		expectedError  error
	}{
		{
			name:           "default",
			options:        []stdoutlog.Option{},
			ctx:            context.Background(),
			records:        records,
			expectedResult: expectedJSONs(now, false),
		},
		{
			name:           "NoRecords",
			options:        []stdoutlog.Option{},
			ctx:            context.Background(),
			records:        nil,
			expectedResult: "",
		},
		{
			name:           "WithPrettyPrint",
			options:        []stdoutlog.Option{stdoutlog.WithPrettyPrint()},
			ctx:            context.Background(),
			records:        records,
			expectedResult: expectedJSONs(now, true),
		},
		{
			name:           "WithoutTimestamps",
			options:        []stdoutlog.Option{stdoutlog.WithoutTimestamps()},
			ctx:            context.Background(),
			records:        records,
			expectedResult: expectedJSONs(time.Time{}, false),
		},
		{
			name:           "WithoutTimestamps and WithPrettyPrint",
			options:        []stdoutlog.Option{stdoutlog.WithoutTimestamps(), stdoutlog.WithPrettyPrint()},
			ctx:            context.Background(),
			records:        records,
			expectedResult: expectedJSONs(time.Time{}, true),
		},
		{
			name:           "WithCanceledContext",
			ctx:            canceledCtx,
			records:        records,
			expectedResult: "",
			expectedError:  context.Canceled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write to buffer for testing
			var buf bytes.Buffer

			exporter, err := stdoutlog.New(append(tc.options, stdoutlog.WithWriter(&buf))...)
			assert.NoError(t, err)

			err = exporter.Export(tc.ctx, tc.records)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, buf.String())
		})
	}
}

func expectedJSONs(now time.Time, prettyPrint bool) string {
	result := expectedJSON(now, prettyPrint)
	return result + result
}

// revive:disable-next-line:flag-parameter
func expectedJSON(now time.Time, prettyPrint bool) string {
	serializedNow, _ := json.Marshal(now)
	if prettyPrint {
		return `{
	"Timestamp": ` + string(serializedNow) + `,
	"ObservedTimestamp": ` + string(serializedNow) + `,
	"Severity": 9,
	"SeverityText": "INFO",
	"Body": {},
	"Attributes": [
		{
			"Key": "key",
			"Value": {}
		},
		{
			"Key": "key2",
			"Value": {}
		},
		{
			"Key": "key3",
			"Value": {}
		},
		{
			"Key": "key4",
			"Value": {}
		},
		{
			"Key": "key5",
			"Value": {}
		},
		{
			"Key": "bool",
			"Value": {}
		}
	],
	"TraceID": "0102030405060708090a0b0c0d0e0f10",
	"SpanID": "0102030405060708",
	"TraceFlags": "01",
	"Resource": {},
	"Scope": {
		"Name": "",
		"Version": "",
		"SchemaURL": ""
	},
	"AttributeValueLengthLimit": 0,
	"AttributeCountLimit": 0
}
`
	}
	return "{\"Timestamp\":" + string(serializedNow) + ",\"ObservedTimestamp\":" + string(serializedNow) + ",\"Severity\":9,\"SeverityText\":\"INFO\",\"Body\":{},\"Attributes\":[{\"Key\":\"key\",\"Value\":{}},{\"Key\":\"key2\",\"Value\":{}},{\"Key\":\"key3\",\"Value\":{}},{\"Key\":\"key4\",\"Value\":{}},{\"Key\":\"key5\",\"Value\":{}},{\"Key\":\"bool\",\"Value\":{}}],\"TraceID\":\"0102030405060708090a0b0c0d0e0f10\",\"SpanID\":\"0102030405060708\",\"TraceFlags\":\"01\",\"Resource\":{},\"Scope\":{\"Name\":\"\",\"Version\":\"\",\"SchemaURL\":\"\"},\"AttributeValueLengthLimit\":0,\"AttributeCountLimit\":0}\n"
}

func TestExporterShutdown(t *testing.T) {
	exporter, err := stdoutlog.New()
	assert.NoError(t, err)

	assert.NoError(t, exporter.Shutdown(context.Background()))
}

func TestExporterForceFlush(t *testing.T) {
	exporter, err := stdoutlog.New()
	assert.NoError(t, err)

	assert.NoError(t, exporter.ForceFlush(context.Background()))
}

func getRecord(now time.Time) sdklog.Record {
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")

	// Setup records
	record := sdklog.Record{}
	record.SetTimestamp(now)
	record.SetObservedTimestamp(now)
	record.SetSeverity(log.SeverityInfo1)
	record.SetSeverityText("INFO")
	record.SetBody(log.StringValue("test"))
	record.SetAttributes([]log.KeyValue{
		// More than 5 attributes to test back slice
		log.String("key", "value"),
		log.String("key2", "value"),
		log.String("key3", "value"),
		log.String("key4", "value"),
		log.String("key5", "value"),
		log.Bool("bool", true),
	}...)
	record.SetTraceID(traceID)
	record.SetSpanID(spanID)
	record.SetTraceFlags(trace.FlagsSampled)

	return record
}
