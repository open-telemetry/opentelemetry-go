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

	alice = api.String("user", "alice")
	bob   = api.String("user", "bob")

	pbAlice = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "alice"},
	}}
	pbBob = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "bob"},
	}}

	sevA = api.SeverityInfo
	sevB = api.SeverityError

	pbSevA = lpb.SeverityNumber_SEVERITY_NUMBER_INFO
	pbSevB = lpb.SeverityNumber_SEVERITY_NUMBER_ERROR

	bodyA = api.StringValue("a")
	bodyB = api.StringValue("b")

	pbBodyA = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "a",
		},
	}
	pbBodyB = &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{
			StringValue: "b",
		},
	}

	spanIDA  = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	spanIDB  = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	traceIDA = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	traceIDB = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	flagsA   = byte(1)
	flagsB   = byte(0)

	scope = instrumentation.Scope{
		Name:      "test/code/path",
		Version:   "v0.1.0",
		SchemaURL: semconv.SchemaURL,
	}
	pbScope = &cpb.InstrumentationScope{
		Name:    "test/code/path",
		Version: "v0.1.0",
	}

	res = resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("test server"),
		semconv.ServiceVersion("v0.1.0"),
	)
	pbRes = &rpb.Resource{
		Attributes: []*cpb.KeyValue{
			{
				Key: "service.name",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "test server"},
				},
			},
			{
				Key: "service.version",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "v0.1.0"},
				},
			},
		},
	}

	records = func() []log.Record {
		var out []log.Record

		out = append(out, logtest.RecordFactory{
			Timestamp:            ts,
			ObservedTimestamp:    obs,
			Severity:             sevA,
			SeverityText:         "A",
			Body:                 bodyA,
			Attributes:           []api.KeyValue{alice},
			TraceID:              trace.TraceID(traceIDA),
			SpanID:               trace.SpanID(spanIDA),
			TraceFlags:           trace.TraceFlags(flagsA),
			InstrumentationScope: &scope,
			Resource:             res,
		}.NewRecord())

		out = append(out, logtest.RecordFactory{
			Timestamp:            ts,
			ObservedTimestamp:    obs,
			Severity:             sevA,
			SeverityText:         "A",
			Body:                 bodyA,
			Attributes:           []api.KeyValue{bob},
			TraceID:              trace.TraceID(traceIDA),
			SpanID:               trace.SpanID(spanIDA),
			TraceFlags:           trace.TraceFlags(flagsA),
			InstrumentationScope: &scope,
			Resource:             res,
		}.NewRecord())

		out = append(out, logtest.RecordFactory{
			Timestamp:            ts,
			ObservedTimestamp:    obs,
			Severity:             sevB,
			SeverityText:         "B",
			Body:                 bodyB,
			Attributes:           []api.KeyValue{alice},
			TraceID:              trace.TraceID(traceIDB),
			SpanID:               trace.SpanID(spanIDB),
			TraceFlags:           trace.TraceFlags(flagsB),
			InstrumentationScope: &scope,
			Resource:             res,
		}.NewRecord())

		out = append(out, logtest.RecordFactory{
			Timestamp:            ts,
			ObservedTimestamp:    obs,
			Severity:             sevB,
			SeverityText:         "B",
			Body:                 bodyB,
			Attributes:           []api.KeyValue{bob},
			TraceID:              trace.TraceID(traceIDB),
			SpanID:               trace.SpanID(spanIDB),
			TraceFlags:           trace.TraceFlags(flagsB),
			InstrumentationScope: &scope,
			Resource:             res,
		}.NewRecord())

		return out
	}()

	pbLogRecords = []*lpb.LogRecord{
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevA,
			SeverityText:         "A",
			Body:                 pbBodyA,
			Attributes:           []*cpb.KeyValue{pbAlice},
			Flags:                uint32(flagsA),
			TraceId:              traceIDA,
			SpanId:               spanIDA,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevA,
			SeverityText:         "A",
			Body:                 pbBodyA,
			Attributes:           []*cpb.KeyValue{pbBob},
			Flags:                uint32(flagsA),
			TraceId:              traceIDA,
			SpanId:               spanIDA,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevB,
			SeverityText:         "B",
			Body:                 pbBodyB,
			Attributes:           []*cpb.KeyValue{pbAlice},
			Flags:                uint32(flagsB),
			TraceId:              traceIDB,
			SpanId:               spanIDB,
		},
		{
			TimeUnixNano:         uint64(ts.UnixNano()),
			ObservedTimeUnixNano: uint64(obs.UnixNano()),
			SeverityNumber:       pbSevB,
			SeverityText:         "B",
			Body:                 pbBodyB,
			Attributes:           []*cpb.KeyValue{pbBob},
			Flags:                uint32(flagsB),
			TraceId:              traceIDB,
			SpanId:               spanIDB,
		},
	}

	pbScopeLogs = &lpb.ScopeLogs{
		Scope:      pbScope,
		SchemaUrl:  semconv.SchemaURL,
		LogRecords: pbLogRecords,
	}

	pbResourceLogs = &lpb.ResourceLogs{
		Resource:  pbRes,
		SchemaUrl: semconv.SchemaURL,
		ScopeLogs: []*lpb.ScopeLogs{pbScopeLogs},
	}
)

func TestResourceLogs(t *testing.T) {
	want := []*lpb.ResourceLogs{pbResourceLogs}
	assert.Equal(t, want, ResourceLogs(records))
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
