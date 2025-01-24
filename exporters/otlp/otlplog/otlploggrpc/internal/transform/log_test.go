// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlplog/transform/log_test.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package transform

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	lpb "go.opentelemetry.io/proto/otlp/logs/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	ts  = time.Date(2000, time.January, 0o1, 0, 0, 0, 0, time.FixedZone("GMT", 0))
	obs = ts.Add(30 * time.Second)

	tom   = api.String("user", "tom")
	jerry = api.String("user", "jerry")
	// A time before unix 0.
	negativeTs = time.Date(1969, 7, 20, 20, 17, 0, 0, time.UTC)

	pbTom = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "tom"},
	}}
	pbJerry = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "jerry"},
	}}

	sevC = api.SeverityInfo
	sevD = api.SeverityError

	pbSevC = lpb.SeverityNumber_SEVERITY_NUMBER_INFO
	pbSevD = lpb.SeverityNumber_SEVERITY_NUMBER_ERROR

	bodyC = api.StringValue("c")
	bodyD = api.StringValue("d")

	pbBodyC = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "c",
		},
	}
	pbBodyD = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "d",
		},
	}

	spanIDC  = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	spanIDD  = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	traceIDC = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	traceIDD = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	flagsC   = byte(1)
	flagsD   = byte(0)

	scope = instrumentation.Scope{
		Name:       "otel/test/code/path1",
		Version:    "v0.1.1",
		SchemaURL:  semconv.SchemaURL,
		Attributes: attribute.NewSet(attribute.String("foo", "bar")),
	}
	scope2 = instrumentation.Scope{
		Name:      "otel/test/code/path2",
		Version:   "v0.2.2",
		SchemaURL: semconv.SchemaURL,
	}
	scopeList = []instrumentation.Scope{scope, scope2}

	pbScope = &cpb.InstrumentationScope{
		Name:    "otel/test/code/path1",
		Version: "v0.1.1",
		Attributes: []*cpb.KeyValue{
			{
				Key: "foo",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "bar"},
				},
			},
		},
	}
	pbScope2 = &cpb.InstrumentationScope{
		Name:    "otel/test/code/path2",
		Version: "v0.2.2",
	}

	res = resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service1"),
		semconv.ServiceVersion("v0.1.1"),
	)
	res2 = resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service2"),
		semconv.ServiceVersion("v0.2.2"),
	)
	resList = []*resource.Resource{res, res2}

	pbRes = &rpb.Resource{
		Attributes: []*cpb.KeyValue{
			{
				Key: "service.name",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "service1"},
				},
			},
			{
				Key: "service.version",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "v0.1.1"},
				},
			},
		},
	}
	pbRes2 = &rpb.Resource{
		Attributes: []*cpb.KeyValue{
			{
				Key: "service.name",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "service2"},
				},
			},
			{
				Key: "service.version",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "v0.2.2"},
				},
			},
		},
	}

	records = func() []log.Record {
		var out []log.Record

		for _, r := range resList {
			for _, s := range scopeList {
				out = append(out, logtest.RecordFactory{
					Timestamp:            ts,
					ObservedTimestamp:    obs,
					EventName:            "evnt",
					Severity:             sevC,
					SeverityText:         "C",
					Body:                 bodyC,
					Attributes:           []api.KeyValue{tom},
					TraceID:              trace.TraceID(traceIDC),
					SpanID:               trace.SpanID(spanIDC),
					TraceFlags:           trace.TraceFlags(flagsC),
					InstrumentationScope: &s,
					Resource:             r,
				}.NewRecord())

				out = append(out, logtest.RecordFactory{
					Timestamp:            ts,
					ObservedTimestamp:    obs,
					Severity:             sevC,
					SeverityText:         "C",
					Body:                 bodyC,
					Attributes:           []api.KeyValue{jerry},
					TraceID:              trace.TraceID(traceIDC),
					SpanID:               trace.SpanID(spanIDC),
					TraceFlags:           trace.TraceFlags(flagsC),
					InstrumentationScope: &s,
					Resource:             r,
				}.NewRecord())

				out = append(out, logtest.RecordFactory{
					Timestamp:            ts,
					ObservedTimestamp:    obs,
					Severity:             sevD,
					SeverityText:         "D",
					Body:                 bodyD,
					Attributes:           []api.KeyValue{tom},
					TraceID:              trace.TraceID(traceIDD),
					SpanID:               trace.SpanID(spanIDD),
					TraceFlags:           trace.TraceFlags(flagsD),
					InstrumentationScope: &s,
					Resource:             r,
				}.NewRecord())

				out = append(out, logtest.RecordFactory{
					Timestamp:            ts,
					ObservedTimestamp:    obs,
					Severity:             sevD,
					SeverityText:         "D",
					Body:                 bodyD,
					Attributes:           []api.KeyValue{jerry},
					TraceID:              trace.TraceID(traceIDD),
					SpanID:               trace.SpanID(spanIDD),
					TraceFlags:           trace.TraceFlags(flagsD),
					InstrumentationScope: &s,
					Resource:             r,
				}.NewRecord())

				out = append(out, logtest.RecordFactory{
					Timestamp:            negativeTs,
					ObservedTimestamp:    obs,
					Severity:             sevD,
					SeverityText:         "D",
					Body:                 bodyD,
					Attributes:           []api.KeyValue{jerry},
					TraceID:              trace.TraceID(traceIDD),
					SpanID:               trace.SpanID(spanIDD),
					TraceFlags:           trace.TraceFlags(flagsD),
					InstrumentationScope: &s,
					Resource:             r,
				}.NewRecord())
			}
		}

		return out
	}()

	pbLogRecords = []*lpb.LogRecord{
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			EventName:            "evnt",
			SeverityNumber:       pbSevC,
			SeverityText:         "C",
			Body:                 pbBodyC,
			Attributes:           []*cpb.KeyValue{pbTom},
			Flags:                uint32(flagsC),
			TraceId:              traceIDC,
			SpanId:               spanIDC,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevC,
			SeverityText:         "C",
			Body:                 pbBodyC,
			Attributes:           []*cpb.KeyValue{pbJerry},
			Flags:                uint32(flagsC),
			TraceId:              traceIDC,
			SpanId:               spanIDC,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevD,
			SeverityText:         "D",
			Body:                 pbBodyD,
			Attributes:           []*cpb.KeyValue{pbTom},
			Flags:                uint32(flagsD),
			TraceId:              traceIDD,
			SpanId:               spanIDD,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevD,
			SeverityText:         "D",
			Body:                 pbBodyD,
			Attributes:           []*cpb.KeyValue{pbJerry},
			Flags:                uint32(flagsD),
			TraceId:              traceIDD,
			SpanId:               spanIDD,
		},
		{
			TimeUnixNano:         0,
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevD,
			SeverityText:         "D",
			Body:                 pbBodyD,
			Attributes:           []*cpb.KeyValue{pbJerry},
			Flags:                uint32(flagsD),
			TraceId:              traceIDD,
			SpanId:               spanIDD,
		},
	}

	pbScopeLogsList = []*lpb.ScopeLogs{
		{
			Scope:      pbScope,
			SchemaUrl:  semconv.SchemaURL,
			LogRecords: pbLogRecords,
		},
		{
			Scope:      pbScope2,
			SchemaUrl:  semconv.SchemaURL,
			LogRecords: pbLogRecords,
		},
	}

	pbResourceLogsList = []*lpb.ResourceLogs{
		{
			Resource:  pbRes,
			SchemaUrl: semconv.SchemaURL,
			ScopeLogs: pbScopeLogsList,
		},
		{
			Resource:  pbRes2,
			SchemaUrl: semconv.SchemaURL,
			ScopeLogs: pbScopeLogsList,
		},
	}
)

func TestResourceLogs(t *testing.T) {
	want := pbResourceLogsList
	assert.ElementsMatch(t, want, ResourceLogs(records))
}

func TestSeverityNumber(t *testing.T) {
	for i := 0; i <= int(api.SeverityFatal4); i++ {
		want := lpb.SeverityNumber(i)
		want += lpb.SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED
		assert.Equal(t, want, SeverityNumber(api.Severity(i)))
	}
}

func BenchmarkResourceLogs(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var out []*lpb.ResourceLogs
		for pb.Next() {
			out = ResourceLogs(records)
		}
		_ = out
	})
}
