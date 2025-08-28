// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
	"go.opentelemetry.io/otel/trace"
)

func TestLoggerEmit(t *testing.T) {
	nowDate := time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

	nowSwap := now
	t.Cleanup(func() {
		now = nowSwap
	})
	now = func() time.Time {
		return nowDate
	}

	p0, p1, p2WithError := newProcessor("0"), newProcessor("1"), newProcessor("2")
	p2WithError.Err = errors.New("error")

	r := log.Record{}
	r.SetEventName("testing.name")
	r.SetTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetBody(log.StringValue("testing body value"))
	r.SetSeverity(log.SeverityInfo)
	r.SetSeverityText("testing text")
	r.AddAttributes(
		log.String("k1", "str"),
		log.Float64("k2", 1.0),
	)
	r.SetObservedTimestamp(time.Date(2001, time.January, 1, 0, 0, 0, 0, time.UTC))

	rWithNoObservedTimestamp := r
	rWithNoObservedTimestamp.SetObservedTimestamp(time.Time{})

	rWithAllowKeyDuplication := r
	rWithAllowKeyDuplication.AddAttributes(
		log.String("k1", "str1"),
	)
	rWithAllowKeyDuplication.SetBody(log.MapValue(
		log.Int64("1", 2),
		log.Int64("1", 3),
	))

	rWithDuplicatesInBody := r
	rWithDuplicatesInBody.SetBody(log.MapValue(
		log.Int64("1", 2),
		log.Int64("1", 3),
	))

	contextWithSpanContext := trace.ContextWithSpanContext(
		context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    trace.TraceID{0o1},
			SpanID:     trace.SpanID{0o2},
			TraceFlags: 0x1,
		}),
	)

	testCases := []struct {
		name            string
		logger          *logger
		ctx             context.Context
		record          log.Record
		expectedRecords []Record
	}{
		{
			name:   "NoProcessors",
			logger: newLogger(NewLoggerProvider(), instrumentation.Scope{}),
			ctx:    context.Background(),
			record: r,
		},
		{
			name: "WithProcessors",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeValueLengthLimit(3),
				WithAttributeCountLimit(2),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			), instrumentation.Scope{Name: "scope"}),
			ctx:    context.Background(),
			record: r,
			expectedRecords: []Record{
				{
					eventName:                 r.EventName(),
					timestamp:                 r.Timestamp(),
					body:                      r.Body(),
					severity:                  r.Severity(),
					severityText:              r.SeverityText(),
					observedTimestamp:         r.ObservedTimestamp(),
					resource:                  resource.NewSchemaless(attribute.String("key", "value")),
					attributeValueLengthLimit: 3,
					attributeCountLimit:       2,
					scope:                     &instrumentation.Scope{Name: "scope"},
					front: [attributesInlineCount]log.KeyValue{
						log.String("k1", "str"),
						log.Float64("k2", 1.0),
					},
					nFront: 2,
				},
			},
		},
		{
			name: "WithProcessorsWithError",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p2WithError),
				WithAttributeValueLengthLimit(3),
				WithAttributeCountLimit(2),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			), instrumentation.Scope{Name: "scope"}),
			ctx: context.Background(),
		},
		{
			name: "WithTraceSpanInContext",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeValueLengthLimit(3),
				WithAttributeCountLimit(2),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			), instrumentation.Scope{Name: "scope"}),
			ctx:    contextWithSpanContext,
			record: r,
			expectedRecords: []Record{
				{
					eventName:                 r.EventName(),
					timestamp:                 r.Timestamp(),
					body:                      r.Body(),
					severity:                  r.Severity(),
					severityText:              r.SeverityText(),
					observedTimestamp:         r.ObservedTimestamp(),
					resource:                  resource.NewSchemaless(attribute.String("key", "value")),
					attributeValueLengthLimit: 3,
					attributeCountLimit:       2,
					scope:                     &instrumentation.Scope{Name: "scope"},
					front: [attributesInlineCount]log.KeyValue{
						log.String("k1", "str"),
						log.Float64("k2", 1.0),
					},
					nFront:     2,
					traceID:    trace.TraceID{0o1},
					spanID:     trace.SpanID{0o2},
					traceFlags: 0x1,
				},
			},
		},
		{
			name: "WithNilContext",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeValueLengthLimit(3),
				WithAttributeCountLimit(2),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			), instrumentation.Scope{Name: "scope"}),
			ctx:    context.Background(),
			record: r,
			expectedRecords: []Record{
				{
					eventName:                 r.EventName(),
					timestamp:                 r.Timestamp(),
					body:                      r.Body(),
					severity:                  r.Severity(),
					severityText:              r.SeverityText(),
					observedTimestamp:         r.ObservedTimestamp(),
					resource:                  resource.NewSchemaless(attribute.String("key", "value")),
					attributeValueLengthLimit: 3,
					attributeCountLimit:       2,
					scope:                     &instrumentation.Scope{Name: "scope"},
					front: [attributesInlineCount]log.KeyValue{
						log.String("k1", "str"),
						log.Float64("k2", 1.0),
					},
					nFront: 2,
				},
			},
		},
		{
			name: "NoObservedTimestamp",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeValueLengthLimit(3),
				WithAttributeCountLimit(2),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			), instrumentation.Scope{Name: "scope"}),
			ctx:    context.Background(),
			record: rWithNoObservedTimestamp,
			expectedRecords: []Record{
				{
					eventName:                 rWithNoObservedTimestamp.EventName(),
					timestamp:                 rWithNoObservedTimestamp.Timestamp(),
					body:                      rWithNoObservedTimestamp.Body(),
					severity:                  rWithNoObservedTimestamp.Severity(),
					severityText:              rWithNoObservedTimestamp.SeverityText(),
					observedTimestamp:         nowDate,
					resource:                  resource.NewSchemaless(attribute.String("key", "value")),
					attributeValueLengthLimit: 3,
					attributeCountLimit:       2,
					scope:                     &instrumentation.Scope{Name: "scope"},
					front: [attributesInlineCount]log.KeyValue{
						log.String("k1", "str"),
						log.Float64("k2", 1.0),
					},
					nFront: 2,
				},
			},
		},
		{
			name: "WithAllowKeyDuplication",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeValueLengthLimit(5),
				WithAttributeCountLimit(5),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
				WithAllowKeyDuplication(),
			), instrumentation.Scope{Name: "scope"}),
			ctx:    context.Background(),
			record: rWithAllowKeyDuplication,
			expectedRecords: []Record{
				{
					eventName:                 rWithAllowKeyDuplication.EventName(),
					timestamp:                 rWithAllowKeyDuplication.Timestamp(),
					body:                      rWithAllowKeyDuplication.Body(),
					severity:                  rWithAllowKeyDuplication.Severity(),
					severityText:              rWithAllowKeyDuplication.SeverityText(),
					observedTimestamp:         rWithAllowKeyDuplication.ObservedTimestamp(),
					resource:                  resource.NewSchemaless(attribute.String("key", "value")),
					attributeValueLengthLimit: 5,
					attributeCountLimit:       5,
					scope:                     &instrumentation.Scope{Name: "scope"},
					front: [attributesInlineCount]log.KeyValue{
						log.String("k1", "str"),
						log.Float64("k2", 1.0),
						log.String("k1", "str1"),
					},
					nFront:       3,
					allowDupKeys: true,
				},
			},
		},
		{
			name: "WithDuplicatesInBody",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeValueLengthLimit(5),
				WithAttributeCountLimit(5),
				WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			), instrumentation.Scope{Name: "scope"}),
			ctx:    context.Background(),
			record: rWithDuplicatesInBody,
			expectedRecords: []Record{
				{
					eventName: rWithDuplicatesInBody.EventName(),
					timestamp: rWithDuplicatesInBody.Timestamp(),
					body: log.MapValue(
						log.Int64("1", 3),
					),
					severity:                  rWithDuplicatesInBody.Severity(),
					severityText:              rWithDuplicatesInBody.SeverityText(),
					observedTimestamp:         rWithDuplicatesInBody.ObservedTimestamp(),
					resource:                  resource.NewSchemaless(attribute.String("key", "value")),
					attributeValueLengthLimit: 5,
					attributeCountLimit:       5,
					scope:                     &instrumentation.Scope{Name: "scope"},
					front: [attributesInlineCount]log.KeyValue{
						log.String("k1", "str"),
						log.Float64("k2", 1.0),
					},
					nFront: 2,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up the records before the test.
			p0.records = nil
			p1.records = nil

			tc.logger.Emit(tc.ctx, tc.record)

			assert.Equal(t, tc.expectedRecords, p0.records)
			assert.Equal(t, tc.expectedRecords, p1.records)
		})
	}
}

func TestLoggerEnabled(t *testing.T) {
	p0 := newFltrProcessor("0", true)
	p1 := newFltrProcessor("1", true)
	p2WithDisabled := newFltrProcessor("2", false)

	testCases := []struct {
		name             string
		logger           *logger
		ctx              context.Context
		param            log.EnabledParameters
		expected         bool
		expectedP0Params []EnabledParameters
		expectedP1Params []EnabledParameters
		expectedP2Params []EnabledParameters
	}{
		{
			name:     "NoProcessors",
			logger:   newLogger(NewLoggerProvider(), instrumentation.Scope{}),
			ctx:      context.Background(),
			expected: false,
		},
		{
			name: "WithProcessors",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
			), instrumentation.Scope{Name: "scope"}),
			ctx: context.Background(),
			param: log.EnabledParameters{
				Severity:  log.SeverityInfo,
				EventName: "test_event",
			},
			expected: true,
			expectedP0Params: []EnabledParameters{{
				InstrumentationScope: instrumentation.Scope{Name: "scope"},
				Severity:             log.SeverityInfo,
				EventName:            "test_event",
			}},
			expectedP1Params: nil,
		},
		{
			name: "WithDisabledProcessors",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p2WithDisabled),
			), instrumentation.Scope{}),
			ctx:              context.Background(),
			expected:         false,
			expectedP2Params: []EnabledParameters{{}},
		},
		{
			name: "ContainsDisabledProcessor",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p2WithDisabled),
				WithProcessor(p0),
			), instrumentation.Scope{}),
			ctx:              context.Background(),
			expected:         true,
			expectedP2Params: []EnabledParameters{{}},
			expectedP0Params: []EnabledParameters{{}},
		},
		{
			name: "WithNilContext",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
			), instrumentation.Scope{}),
			ctx:              nil,
			expected:         true,
			expectedP0Params: []EnabledParameters{{}},
			expectedP1Params: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up the records before the test.
			p0.params = nil
			p1.params = nil
			p2WithDisabled.params = nil

			assert.Equal(t, tc.expected, tc.logger.Enabled(tc.ctx, tc.param))
			assert.Equal(t, tc.expectedP0Params, p0.params)
			assert.Equal(t, tc.expectedP1Params, p1.params)
			assert.Equal(t, tc.expectedP2Params, p2WithDisabled.params)
		})
	}
}

func TestLoggerSelfObservability(t *testing.T) {
	testCases := []struct {
		name                     string
		selfObservabilityEnabled bool
		records                  []log.Record
		wantLogRecordCount       int64
	}{
		{
			name:                     "Disabled",
			selfObservabilityEnabled: false,
			records:                  []log.Record{{}, {}},
			wantLogRecordCount:       0,
		},
		{
			name:                     "Enabled",
			selfObservabilityEnabled: true,
			records:                  []log.Record{{}, {}, {}, {}, {}},
			wantLogRecordCount:       5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", strconv.FormatBool(tc.selfObservabilityEnabled))
			prev := otel.GetMeterProvider()
			t.Cleanup(func() {
				otel.SetMeterProvider(prev)
			})
			r := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(r))
			otel.SetMeterProvider(mp)
			l := newLogger(NewLoggerProvider(), instrumentation.Scope{})

			for _, record := range tc.records {
				l.Emit(context.Background(), record)
			}

			gotMetrics := new(metricdata.ResourceMetrics)
			assert.NoError(t, r.Collect(context.Background(), gotMetrics))
			if tc.wantLogRecordCount == 0 {
				assert.Empty(t, gotMetrics.ScopeMetrics)
				return
			}

			require.Len(t, gotMetrics.ScopeMetrics, 1)
			sm := gotMetrics.ScopeMetrics[0]
			assert.Equal(t, instrumentation.Scope{
				Name:      "go.opentelemetry.io/otel/sdk/log",
				Version:   sdk.Version(),
				SchemaURL: semconv.SchemaURL,
			}, sm.Scope)

			wantMetric := metricdata.Metrics{
				Name:        otelconv.SDKLogCreated{}.Name(),
				Description: otelconv.SDKLogCreated{}.Description(),
				Unit:        otelconv.SDKLogCreated{}.Unit(),
				Data: metricdata.Sum[int64]{
					DataPoints:  []metricdata.DataPoint[int64]{{Value: tc.wantLogRecordCount}},
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
				},
			}
			metricdatatest.AssertEqual(t, wantMetric, sm.Metrics[0], metricdatatest.IgnoreTimestamp())
		})
	}
}

func TestNewLoggerSelfObservabilityErrorHandled(t *testing.T) {
	errHandler := otel.GetErrorHandler()
	t.Cleanup(func() {
		otel.SetErrorHandler(errHandler)
	})

	var errs []error
	eh := otel.ErrorHandlerFunc(func(e error) { errs = append(errs, e) })
	otel.SetErrorHandler(eh)

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })
	otel.SetMeterProvider(&errMeterProvider{err: assert.AnError})

	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
	l := newLogger(NewLoggerProvider(), instrumentation.Scope{})
	_ = l
	require.Len(t, errs, 1)
	assert.ErrorIs(t, errs[0], assert.AnError)
}

type errMeterProvider struct {
	metric.MeterProvider

	err error
}

func (mp *errMeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	return &errMeter{err: mp.err}
}

type errMeter struct {
	metric.Meter

	err error
}

func (m *errMeter) Int64Counter(string, ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Int64UpDownCounter(string, ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return nil, m.err
}
