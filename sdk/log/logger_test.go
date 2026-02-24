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
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
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
		t.Context(),
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
			ctx:    t.Context(),
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
			ctx:    t.Context(),
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
			ctx: t.Context(),
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
			ctx:    t.Context(),
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
			ctx:    t.Context(),
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
			ctx:    t.Context(),
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
			ctx:    t.Context(),
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

func TestAddExceptionFromError(t *testing.T) {
	t.Run("AddsMissing", func(t *testing.T) {
		r := &Record{}
		r.attributeValueLengthLimit = -1
		addExceptionFromError(r, errors.New("boom"))

		var gotType, gotMessage, gotStack string
		r.WalkAttributes(func(kv log.KeyValue) bool {
			switch kv.Key {
			case string(semconv.ExceptionTypeKey):
				gotType = kv.Value.AsString()
			case string(semconv.ExceptionMessageKey):
				gotMessage = kv.Value.AsString()
			case string(semconv.ExceptionStacktraceKey):
				gotStack = kv.Value.AsString()
			}
			return true
		})

		assert.Equal(t, "*errors.errorString", gotType)
		assert.Equal(t, "boom", gotMessage)
		assert.Empty(t, gotStack)
	})

	t.Run("DoesNotOverwrite", func(t *testing.T) {
		r := &Record{}
		r.attributeValueLengthLimit = -1
		r.AddAttributes(
			log.String(string(semconv.ExceptionTypeKey), "existing.type"),
			log.String(string(semconv.ExceptionMessageKey), "existing.message"),
			log.String(string(semconv.ExceptionStacktraceKey), "existing.stack"),
		)

		addExceptionFromError(r, errors.New("boom"))

		var gotType, gotMessage, gotStack string
		r.WalkAttributes(func(kv log.KeyValue) bool {
			switch kv.Key {
			case string(semconv.ExceptionTypeKey):
				gotType = kv.Value.AsString()
			case string(semconv.ExceptionMessageKey):
				gotMessage = kv.Value.AsString()
			case string(semconv.ExceptionStacktraceKey):
				gotStack = kv.Value.AsString()
			}
			return true
		})

		assert.Equal(t, "existing.type", gotType)
		assert.Equal(t, "existing.message", gotMessage)
		assert.Equal(t, "existing.stack", gotStack)
	})

	t.Run("DoesNotMixPartial", func(t *testing.T) {
		r := &Record{}
		r.attributeValueLengthLimit = -1
		r.AddAttributes(
			log.String(string(semconv.ExceptionTypeKey), "existing.type"),
		)

		addExceptionFromError(r, errors.New("boom"))

		var gotType, gotMessage, gotStack string
		r.WalkAttributes(func(kv log.KeyValue) bool {
			switch kv.Key {
			case string(semconv.ExceptionTypeKey):
				gotType = kv.Value.AsString()
			case string(semconv.ExceptionMessageKey):
				gotMessage = kv.Value.AsString()
			case string(semconv.ExceptionStacktraceKey):
				gotStack = kv.Value.AsString()
			}
			return true
		})

		assert.Equal(t, "existing.type", gotType)
		assert.Empty(t, gotMessage)
		assert.Empty(t, gotStack)
	})

	t.Run("ShortCircuitsAtAttributeLimit", func(t *testing.T) {
		r := &Record{}
		r.attributeValueLengthLimit = -1
		r.attributeCountLimit = 1
		r.AddAttributes(log.String("k1", "v1"))

		addExceptionFromError(r, errors.New("boom"))

		var gotType, gotMessage string
		r.WalkAttributes(func(kv log.KeyValue) bool {
			switch kv.Key {
			case string(semconv.ExceptionTypeKey):
				gotType = kv.Value.AsString()
			case string(semconv.ExceptionMessageKey):
				gotMessage = kv.Value.AsString()
			}
			return true
		})

		assert.Empty(t, gotType)
		assert.Equal(t, "boom", gotMessage)
	})
}

func TestErrorType(t *testing.T) {
	t.Run("UsesErrorTypeMethod", func(t *testing.T) {
		err := errWithType{msg: "boom", typ: "custom.type"}
		assert.Equal(t, "custom.type", errorType(err))
	})

	t.Run("FallsBackWhenErrorTypeEmpty", func(t *testing.T) {
		err := errWithType{msg: "boom", typ: ""}
		assert.Equal(t, "go.opentelemetry.io/otel/sdk/log.errWithType", errorType(err))
	})

	t.Run("NilError", func(t *testing.T) {
		assert.Empty(t, errorType(nil))
	})

	t.Run("UnnamedType", func(t *testing.T) {
		var err error = struct{ baseErr }{}
		assert.Contains(t, errorType(err), "struct")
	})
}

type errWithType struct {
	msg string
	typ string
}

func (e errWithType) Error() string { return e.msg }

func (e errWithType) ErrorType() string { return e.typ }

type baseErr struct{}

func (baseErr) Error() string { return "boom" }

func TestNewRecordSkipsExceptionWhenPresent(t *testing.T) {
	l := newLogger(NewLoggerProvider(), instrumentation.Scope{})

	var r log.Record
	r.SetBody(log.StringValue("boom"))
	r.SetSeverity(log.SeverityError)
	r.SetErr(errors.New("boom"))
	r.AddAttributes(log.String(string(semconv.ExceptionMessageKey), "existing.message"))

	got := l.newRecord(t.Context(), r)

	var gotType, gotMessage string
	got.WalkAttributes(func(kv log.KeyValue) bool {
		switch kv.Key {
		case string(semconv.ExceptionTypeKey):
			gotType = kv.Value.AsString()
		case string(semconv.ExceptionMessageKey):
			gotMessage = kv.Value.AsString()
		}
		return true
	})

	assert.Equal(t, "existing.message", gotMessage)
	assert.Empty(t, gotType)
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
			ctx:      t.Context(),
			expected: false,
		},
		{
			name: "WithProcessors",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p0),
				WithProcessor(p1),
			), instrumentation.Scope{Name: "scope"}),
			ctx: t.Context(),
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
			ctx:              t.Context(),
			expected:         false,
			expectedP2Params: []EnabledParameters{{}},
		},
		{
			name: "ContainsDisabledProcessor",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p2WithDisabled),
				WithProcessor(p0),
			), instrumentation.Scope{}),
			ctx:              t.Context(),
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

func TestLoggerObservability(t *testing.T) {
	testCases := []struct {
		name               string
		enabled            bool
		records            []log.Record
		wantLogRecordCount int64
	}{
		{
			name:               "Disabled",
			enabled:            false,
			records:            []log.Record{{}, {}},
			wantLogRecordCount: 0,
		},
		{
			name:               "Enabled",
			enabled:            true,
			records:            []log.Record{{}, {}, {}, {}, {}},
			wantLogRecordCount: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_OBSERVABILITY", strconv.FormatBool(tc.enabled))
			prev := otel.GetMeterProvider()
			t.Cleanup(func() {
				otel.SetMeterProvider(prev)
			})
			r := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(r))
			otel.SetMeterProvider(mp)
			l := newLogger(NewLoggerProvider(), instrumentation.Scope{})

			for _, record := range tc.records {
				l.Emit(t.Context(), record)
			}

			gotMetrics := new(metricdata.ResourceMetrics)
			assert.NoError(t, r.Collect(t.Context(), gotMetrics))
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

func TestNewLoggerObservabilityErrorHandled(t *testing.T) {
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

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
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
