// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
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

func TestNewSelfObservability(t *testing.T) {
	testCases := []struct {
		name   string
		enable bool
		test   func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics)
	}{
		{
			name:   "inflight",
			enable: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log 1"),
				}
				record1 := rf.NewRecord()

				rf.Body = log.StringValue("test log 2")
				record2 := rf.NewRecord()

				records := []sdklog.Record{record1, record2}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				assert.NotEmpty(t, got.Metrics)

				assert.Equal(t, "go.opentelemetry.io/otel/exporters/stdout/stdoutlog", got.Scope.Name)
				assert.NotEmpty(t, got.Scope.Version)

				var inflightMetric metricdata.Metrics
				inflightInstrument := otelconv.SDKExporterLogInflight{}
				for _, m := range got.Metrics {
					if m.Name == inflightInstrument.Name() {
						inflightMetric = m
						break
					}
				}
				require.NotEmpty(t, inflightMetric, "inflight metric not found")

				expected := metricdata.Metrics{
					Name:        inflightInstrument.Name(),
					Description: inflightInstrument.Description(),
					Unit:        inflightInstrument.Unit(),
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: false,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Value: 0,
								Attributes: attribute.NewSet(
									attribute.KeyValue{
										Key:   "otel.component.name",
										Value: attribute.StringValue(otelComponentType + "/0"),
									},
									attribute.KeyValue{
										Key:   "otel.component.type",
										Value: attribute.StringValue(otelComponentType),
									},
								),
							},
						},
					},
				}

				metricdatatest.AssertEqual(
					t,
					expected,
					inflightMetric,
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreValue(),
				)
			},
		},
		{
			name:   "exported",
			enable: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log 1"),
				}
				record1 := rf.NewRecord()

				rf.Body = log.StringValue("test log 2")
				record2 := rf.NewRecord()

				rf.Body = log.StringValue("test log 3")
				record3 := rf.NewRecord()

				records := []sdklog.Record{record1, record2, record3}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				assert.NotEmpty(t, got.Metrics)

				assert.Equal(t, "go.opentelemetry.io/otel/exporters/stdout/stdoutlog", got.Scope.Name)
				assert.NotEmpty(t, got.Scope.Version)

				var exportedMetric metricdata.Metrics
				exportedInstrument := otelconv.SDKExporterLogExported{}
				for _, m := range got.Metrics {
					if m.Name == exportedInstrument.Name() {
						exportedMetric = m
						break
					}
				}
				require.NotEmpty(t, exportedMetric, "exported metric not found")

				expected := metricdata.Metrics{
					Name:        exportedInstrument.Name(),
					Description: exportedInstrument.Description(),
					Unit:        exportedInstrument.Unit(),
					Data: metricdata.Sum[int64]{
						Temporality: metricdata.CumulativeTemporality,
						IsMonotonic: true,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Value: 3,
								Attributes: attribute.NewSet(
									attribute.KeyValue{
										Key:   "otel.component.name",
										Value: attribute.StringValue(otelComponentType + "/0"),
									},
									attribute.KeyValue{
										Key:   "otel.component.type",
										Value: attribute.StringValue(otelComponentType),
									},
								),
							},
						},
					},
				}

				metricdatatest.AssertEqual(
					t,
					expected,
					exportedMetric,
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreValue(),
				)
			},
		},
		{
			name:   "duration",
			enable: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log"),
				}
				record := rf.NewRecord()
				records := []sdklog.Record{record}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				assert.NotEmpty(t, got.Metrics)

				assert.Equal(t, "go.opentelemetry.io/otel/exporters/stdout/stdoutlog", got.Scope.Name)
				assert.NotEmpty(t, got.Scope.Version)

				var durationMetric metricdata.Metrics
				durationInstrument := otelconv.SDKExporterOperationDuration{}
				for _, m := range got.Metrics {
					if m.Name == durationInstrument.Name() {
						durationMetric = m
						break
					}
				}
				require.NotEmpty(t, durationMetric, "duration metric not found")

				expected := metricdata.Metrics{
					Name:        durationInstrument.Name(),
					Description: durationInstrument.Description(),
					Unit:        durationInstrument.Unit(),
					Data: metricdata.Histogram[float64]{
						Temporality: metricdata.CumulativeTemporality,
						DataPoints: []metricdata.HistogramDataPoint[float64]{
							{
								Count: 1,
								Attributes: attribute.NewSet(
									attribute.KeyValue{
										Key:   "otel.component.name",
										Value: attribute.StringValue(otelComponentType + "/0"),
									},
									attribute.KeyValue{
										Key:   "otel.component.type",
										Value: attribute.StringValue(otelComponentType),
									},
								),
							},
						},
					},
				}

				metricdatatest.AssertEqual(
					t,
					expected,
					durationMetric,
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreValue(),
				)
			},
		},
		{
			name:   "multiple_exports",
			enable: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				for i := 0; i < 3; i++ {
					rf := logtest.RecordFactory{
						Timestamp: time.Now(),
						Body:      log.StringValue("test log"),
					}
					record := rf.NewRecord()
					records := []sdklog.Record{record}
					err = exporter.Export(context.Background(), records)
					require.NoError(t, err)
				}

				got := scopeMetrics()
				assert.NotEmpty(t, got.Metrics)

				assert.Equal(t, "go.opentelemetry.io/otel/exporters/stdout/stdoutlog", got.Scope.Name)
				assert.NotEmpty(t, got.Scope.Version)

				expectedAttrs := attribute.NewSet(
					attribute.KeyValue{
						Key:   "otel.component.name",
						Value: attribute.StringValue(otelComponentType + "/0"),
					},
					attribute.KeyValue{Key: "otel.component.type", Value: attribute.StringValue(otelComponentType)},
				)

				expected := metricdata.ScopeMetrics{
					Scope: got.Scope,
					Metrics: []metricdata.Metrics{
						{
							Name:        otelconv.SDKExporterLogInflight{}.Name(),
							Description: otelconv.SDKExporterLogInflight{}.Description(),
							Unit:        otelconv.SDKExporterLogInflight{}.Unit(),
							Data: metricdata.Sum[int64]{
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: false,
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value:      0,
										Attributes: expectedAttrs,
									},
								},
							},
						},
						{
							Name:        otelconv.SDKExporterLogExported{}.Name(),
							Description: otelconv.SDKExporterLogExported{}.Description(),
							Unit:        otelconv.SDKExporterLogExported{}.Unit(),
							Data: metricdata.Sum[int64]{
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: true,
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value:      3,
										Attributes: expectedAttrs,
									},
								},
							},
						},
						{
							Name:        otelconv.SDKExporterOperationDuration{}.Name(),
							Description: otelconv.SDKExporterOperationDuration{}.Description(),
							Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
							Data: metricdata.Histogram[float64]{
								Temporality: metricdata.CumulativeTemporality,
								DataPoints: []metricdata.HistogramDataPoint[float64]{
									{
										Count:      3,
										Attributes: expectedAttrs,
									},
								},
							},
						},
					},
				}

				metricdatatest.AssertEqual(
					t,
					expected,
					got,
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreValue(),
				)
			},
		},
		{
			name:   "empty_records",
			enable: true,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)
				err = exporter.Export(context.Background(), []sdklog.Record{})
				require.NoError(t, err)

				got := scopeMetrics()
				assert.Equal(t, "go.opentelemetry.io/otel/exporters/stdout/stdoutlog", got.Scope.Name)
				assert.NotEmpty(t, got.Scope.Version)
				assert.NotEmpty(t, got.Metrics, "metrics should be recorded even for empty records")

				assert.Len(t, got.Metrics, 3, "should have 3 metrics for self-observability")
				metricNames := make(map[string]bool)
				for _, metric := range got.Metrics {
					metricNames[metric.Name] = true
				}
				assert.True(t, metricNames["otel.sdk.exporter.log.exported"], "exported metric should be present")
				assert.True(t, metricNames["otel.sdk.exporter.operation.duration"], "duration metric should be present")
				assert.True(t, metricNames["otel.sdk.exporter.log.inflight"], "inflight metric should be present")
			},
		},
		{
			name:   "self_observability_disabled",
			enable: false,
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log"),
				}
				record := rf.NewRecord()
				records := []sdklog.Record{record}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				assert.Empty(t, got.Metrics, "no metrics should be recorded when self-observability is disabled")
			},
		},
	}
	ranOnce := false
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.enable {
				t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
			}

			if ranOnce {
				// Reset the global exporter ID counter for deterministic tests
				exporterIDCounter.Store(0) // First call to nextExporterID() will return 0
			}

			prev := otel.GetMeterProvider()
			t.Cleanup(func() { otel.SetMeterProvider(prev) })
			r := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(r))
			otel.SetMeterProvider(mp)

			scopeMetrics := func() metricdata.ScopeMetrics {
				var got metricdata.ResourceMetrics
				err := r.Collect(context.Background(), &got)
				require.NoError(t, err)
				if len(got.ScopeMetrics) == 0 {
					return metricdata.ScopeMetrics{}
				}
				return got.ScopeMetrics[0]
			}
			tc.test(t, scopeMetrics)
		})
		ranOnce = true
	}
}
