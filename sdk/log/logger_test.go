// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
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

	contextWithSpanContext := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0o1},
		SpanID:     trace.SpanID{0o2},
		TraceFlags: 0x1,
	}))

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

	emptyResource := resource.Empty()
	res := resource.NewSchemaless(attribute.String("key", "value"))

	testCases := []struct {
		name             string
		logger           *logger
		ctx              context.Context
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
				WithResource(res),
			), instrumentation.Scope{Name: "scope"}),
			ctx:      context.Background(),
			expected: true,
			expectedP0Params: []EnabledParameters{{
				Resource:             *res,
				InstrumentationScope: instrumentation.Scope{Name: "scope"},
			}},
			expectedP1Params: nil,
		},
		{
			name: "WithDisabledProcessors",
			logger: newLogger(NewLoggerProvider(
				WithProcessor(p2WithDisabled),
				WithResource(emptyResource),
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
				WithResource(emptyResource),
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
				WithResource(emptyResource),
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

			assert.Equal(t, tc.expected, tc.logger.Enabled(tc.ctx, log.EnabledParameters{}))
			assert.Equal(t, tc.expectedP0Params, p0.params)
			assert.Equal(t, tc.expectedP1Params, p1.params)
			assert.Equal(t, tc.expectedP2Params, p2WithDisabled.params)
		})
	}
}
