// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/counter"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/observ"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/semconv/v1.40.0/otelconv"
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
			err = exporter.Export(t.Context(), []sdklog.Record{record})
			assert.NoError(t, err)

			// Check the writer
			assert.Equal(t, tc.want, buf.String())

			// Flush the exporter
			err = exporter.ForceFlush(t.Context())
			assert.NoError(t, err)

			// Shutdown the exporter
			err = exporter.Shutdown(t.Context())
			assert.NoError(t, err)

			// Export a record after shutdown, this should not be written
			err = exporter.Export(t.Context(), []sdklog.Record{record})
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
			ctx:        t.Context(),
			records:    records,
			wantResult: getJSONs(&now),
		},
		{
			name:       "NoRecords",
			options:    []Option{},
			ctx:        t.Context(),
			records:    nil,
			wantResult: "",
		},
		{
			name:       "WithPrettyPrint",
			options:    []Option{WithPrettyPrint()},
			ctx:        t.Context(),
			records:    records,
			wantResult: getPrettyJSONs(&now),
		},
		{
			name:       "WithoutTimestamps",
			options:    []Option{WithoutTimestamps()},
			ctx:        t.Context(),
			records:    records,
			wantResult: getJSONs(nil),
		},
		{
			name:       "WithoutTimestamps and WithPrettyPrint",
			options:    []Option{WithoutTimestamps(), WithPrettyPrint()},
			ctx:        t.Context(),
			records:    records,
			wantResult: getPrettyJSONs(nil),
		},
		{
			name: "WithCanceledContext",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(t.Context())
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

	return "{" + timestamps + "\"EventName\":\"testing.event\",\"Severity\":9,\"SeverityText\":\"INFO\",\"Body\":{\"Type\":\"String\",\"Value\":\"test\"},\"Attributes\":[{\"Key\":\"key\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key2\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key3\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key4\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"key5\",\"Value\":{\"Type\":\"String\",\"Value\":\"value\"}},{\"Key\":\"bool\",\"Value\":{\"Type\":\"Bool\",\"Value\":true}}],\"TraceID\":\"0102030405060708090a0b0c0d0e0f10\",\"SpanID\":\"0102030405060708\",\"TraceFlags\":\"01\",\"Resource\":[{\"Key\":\"foo\",\"Value\":{\"Type\":\"STRING\",\"Value\":\"bar\"}}],\"Scope\":{\"Name\":\"name\",\"Version\":\"version\",\"SchemaURL\":\"https://example.com/custom-schema\",\"Attributes\":{}},\"DroppedAttributes\":10}\n"
}

func getJSONs(now *time.Time) string {
	return getJSON(now) + getJSON(now)
}

func getPrettyJSON(now *time.Time) string {
	var timestamps string
	if now != nil {
		serializedNow, _ := json.Marshal(now)
		timestamps = "\n\t\"Timestamp\": " + string(
			serializedNow,
		) + ",\n\t\"ObservedTimestamp\": " + string(
			serializedNow,
		) + ","
	}

	return `{` + timestamps + `
	"EventName": "testing.event",
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
		"SchemaURL": "https://example.com/custom-schema",
		"Attributes": {}
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

	assert.NoError(t, exporter.Shutdown(t.Context()))
}

func TestExporterForceFlush(t *testing.T) {
	exporter, err := New()
	assert.NoError(t, err)

	assert.NoError(t, exporter.ForceFlush(t.Context()))
}

func getRecord(now time.Time) sdklog.Record {
	traceID, _ := trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")

	rf := logtest.RecordFactory{
		EventName:         "testing.event",
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
		InstrumentationScope: &instrumentation.Scope{
			Name:      "name",
			Version:   "version",
			SchemaURL: "https://example.com/custom-schema",
		},
		DroppedAttributes: 10,
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
			for range goroutines {
				go func() {
					defer wg.Done()
					err := exporter.Export(t.Context(), []sdklog.Record{{}})
					assert.NoError(t, err)
					err = exporter.ForceFlush(t.Context())
					assert.NoError(t, err)
					err = exporter.Shutdown(t.Context())
					assert.NoError(t, err)
				}()
			}
			wg.Wait()
		})
	}
}

func TestValueMarshalJSON(t *testing.T) {
	testCases := []struct {
		value log.Value
		want  string
	}{
		{
			value: log.Empty("test").Value,
			want:  `{"Type":"Empty","Value":null}`,
		},
		{
			value: log.BoolValue(true),
			want:  `{"Type":"Bool","Value":true}`,
		},
		{
			value: log.Float64Value(3.14),
			want:  `{"Type":"Float64","Value":3.14}`,
		},
		{
			value: log.Int64Value(42),
			want:  `{"Type":"Int64","Value":42}`,
		},
		{
			value: log.StringValue("hello"),
			want:  `{"Type":"String","Value":"hello"}`,
		},
		{
			value: log.BytesValue([]byte{1, 2, 3}),
			// The base64 encoding of []byte{1, 2, 3} is "AQID".
			want: `{"Type":"Bytes","Value":"AQID"}`,
		},
		{
			value: log.SliceValue(
				log.Empty("empty").Value,
				log.BoolValue(true),
				log.Float64Value(2.2),
				log.IntValue(3),
				log.StringValue("4"),
				log.BytesValue([]byte{5}),
				log.SliceValue(
					log.IntValue(6),
					log.MapValue(
						log.Int("seven", 7),
					),
				),
				log.MapValue(
					log.Int("nine", 9),
				),
			),
			want: `{"Type":"Slice","Value":[{"Type":"Empty","Value":null},{"Type":"Bool","Value":true},{"Type":"Float64","Value":2.2},{"Type":"Int64","Value":3},{"Type":"String","Value":"4"},{"Type":"Bytes","Value":"BQ=="},{"Type":"Slice","Value":[{"Type":"Int64","Value":6},{"Type":"Map","Value":[{"Key":"seven","Value":{"Type":"Int64","Value":7}}]}]},{"Type":"Map","Value":[{"Key":"nine","Value":{"Type":"Int64","Value":9}}]}]}`,
		},
		{
			value: log.MapValue(
				log.Empty("empty"),
				log.Bool("one", true),
				log.Float64("two", 2.2),
				log.Int("three", 3),
				log.String("four", "4"),
				log.Bytes("five", []byte{5}),
				log.Slice("six",
					log.IntValue(6),
					log.MapValue(
						log.Int("seven", 7),
					),
				),
				log.Map("eight",
					log.Int("nine", 9),
				),
			),
			want: `{"Type":"Map","Value":[{"Key":"empty","Value":{"Type":"Empty","Value":null}},{"Key":"one","Value":{"Type":"Bool","Value":true}},{"Key":"two","Value":{"Type":"Float64","Value":2.2}},{"Key":"three","Value":{"Type":"Int64","Value":3}},{"Key":"four","Value":{"Type":"String","Value":"4"}},{"Key":"five","Value":{"Type":"Bytes","Value":"BQ=="}},{"Key":"six","Value":{"Type":"Slice","Value":[{"Type":"Int64","Value":6},{"Type":"Map","Value":[{"Key":"seven","Value":{"Type":"Int64","Value":7}}]}]}},{"Key":"eight","Value":{"Type":"Map","Value":[{"Key":"nine","Value":{"Type":"Int64","Value":9}}]}}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.value.String(), func(t *testing.T) {
			got, err := json.Marshal(value{Value: tc.value})
			require.NoError(t, err)
			assert.JSONEq(t, tc.want, string(got))
		})
	}
}

func TestObservability(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		test    func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics)
	}{
		{
			name:    "Disabled",
			enabled: false,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				var buf bytes.Buffer
				exporter, err := New(WithWriter(&buf))
				require.NoError(t, err)
				assert.Nil(t, exporter.inst)
			},
		},
		{
			name:    "upload success",
			enabled: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				ctx := t.Context()
				var buf bytes.Buffer
				exporter, err := New(WithWriter(&buf))
				require.NoError(t, err)
				err = exporter.Export(ctx, []sdklog.Record{getRecord(time.Now())})
				require.NoError(t, err)
				assertStdoutLogObservabilityMetrics(t, scopeMetrics(), 1, 1, nil)
			},
		},
		{
			name:    "upload failed",
			enabled: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				ctx := t.Context()
				writeErr := errors.New("write failed")
				exporter, err := New(WithWriter(&failingWriter{err: writeErr}))
				require.NoError(t, err)
				err = exporter.Export(ctx, []sdklog.Record{getRecord(time.Now())})
				require.ErrorIs(t, err, writeErr)
				assertStdoutLogObservabilityMetrics(t, scopeMetrics(), 1, 0, writeErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.enabled {
				t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
				_ = counter.SetExporterID(0)
			}
			provider := otel.GetMeterProvider()
			t.Cleanup(func() {
				otel.SetMeterProvider(provider)
			})
			r := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(r))
			otel.SetMeterProvider(mp)

			scopeMetrics := func() metricdata.ScopeMetrics {
				var got metricdata.ResourceMetrics
				err := r.Collect(t.Context(), &got)
				require.NoError(t, err)
				require.Len(t, got.ScopeMetrics, 1)
				return got.ScopeMetrics[0]
			}
			tc.test(t, scopeMetrics)
		})
	}
}

func BenchmarkExporterObservability(b *testing.B) {
	ctx := b.Context()
	rec := getRecord(time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC))
	records := []sdklog.Record{rec}

	b.Run("Disabled", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
		var buf bytes.Buffer
		exporter, err := New(WithWriter(&buf))
		require.NoError(b, err)

		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			buf.Reset()
			_ = exporter.Export(ctx, records)
		}
	})

	setupObservability := func(b *testing.B) {
		b.Helper()
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		_ = counter.SetExporterID(0)
		provider := otel.GetMeterProvider()
		b.Cleanup(func() {
			otel.SetMeterProvider(provider)
		})
		r := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(r))
		otel.SetMeterProvider(mp)
	}

	b.Run("UploadSuccess", func(b *testing.B) {
		setupObservability(b)
		var buf bytes.Buffer
		exporter, err := New(WithWriter(&buf))
		require.NoError(b, err)
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			buf.Reset()
			_ = exporter.Export(ctx, records)
		}
	})
}

// failingWriter implements [io.Writer] and always returns an error.
type failingWriter struct {
	err error
}

func (w *failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func stdoutObservAttrSet(err error) attribute.Set {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(observ.GetComponentName(0)),
		semconv.OTelComponentNameKey.String(observ.ComponentType),
	}
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attribute.NewSet(attrs...)
}

func stdoutLogInflightMetric() metricdata.Metrics {
	inflight := otelconv.SDKExporterLogInflight{}
	return metricdata.Metrics{
		Name:        inflight.Name(),
		Description: inflight.Description(),
		Unit:        inflight.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: stdoutObservAttrSet(nil), Value: 0},
			},
		},
	}
}

func stdoutLogExportedMetric(success, total int64, err error) metricdata.Metrics {
	dp := []metricdata.DataPoint[int64]{
		{Attributes: stdoutObservAttrSet(nil), Value: success},
	}
	if err != nil {
		dp = append(dp, metricdata.DataPoint[int64]{
			Attributes: stdoutObservAttrSet(err),
			Value:      total - success,
		})
	}
	exported := otelconv.SDKExporterLogExported{}
	return metricdata.Metrics{
		Name:        exported.Name(),
		Description: exported.Description(),
		Unit:        exported.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dp,
		},
	}
}

func stdoutLogDurationMetric(err error) metricdata.Metrics {
	duration := otelconv.SDKExporterOperationDuration{}
	return metricdata.Metrics{
		Name:        duration.Name(),
		Description: duration.Description(),
		Unit:        duration.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: stdoutObservAttrSet(err)},
			},
		},
	}
}

func assertStdoutLogObservabilityMetrics(
	t *testing.T,
	got metricdata.ScopeMetrics,
	logs int64,
	success int64,
	err error,
) {
	t.Helper()
	wantScope := instrumentation.Scope{
		Name:      observ.ScopeName,
		Version:   internal.Version,
		SchemaURL: semconv.SchemaURL,
	}
	assert.Equal(t, wantScope, got.Scope)

	m := got.Metrics
	require.Len(t, m, 3)

	o := metricdatatest.IgnoreTimestamp()
	metricdatatest.AssertEqual(t, stdoutLogInflightMetric(), m[0], o)
	metricdatatest.AssertEqual(t, stdoutLogExportedMetric(success, logs, err), m[1], o)
	metricdatatest.AssertEqual(t, stdoutLogDurationMetric(err), m[2], metricdatatest.IgnoreValue(), o)
}
