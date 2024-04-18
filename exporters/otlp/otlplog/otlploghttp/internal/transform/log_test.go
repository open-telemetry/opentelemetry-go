// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package transform

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	cpb "go.opentelemetry.io/proto/slim/otlp/common/v1"
	lpb "go.opentelemetry.io/proto/slim/otlp/logs/v1"

	api "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
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

	records = func() []log.Record {
		r0 := new(log.Record)
		r0.SetTimestamp(ts)
		r0.SetObservedTimestamp(obs)
		r0.SetSeverity(sevA)
		r0.SetSeverityText("A")
		r0.SetBody(bodyA)
		r0.SetAttributes(alice)
		r0.SetTraceID(trace.TraceID(traceIDA))
		r0.SetSpanID(trace.SpanID(spanIDA))
		r0.SetTraceFlags(trace.TraceFlags(flagsA))

		r1 := new(log.Record)
		r1.SetTimestamp(ts)
		r1.SetObservedTimestamp(obs)
		r1.SetSeverity(sevA)
		r1.SetSeverityText("A")
		r1.SetBody(bodyA)
		r1.SetAttributes(bob)
		r1.SetTraceID(trace.TraceID(traceIDA))
		r1.SetSpanID(trace.SpanID(spanIDA))
		r1.SetTraceFlags(trace.TraceFlags(flagsA))

		r2 := new(log.Record)
		r2.SetTimestamp(ts)
		r2.SetObservedTimestamp(obs)
		r2.SetSeverity(sevB)
		r2.SetSeverityText("B")
		r2.SetBody(bodyB)
		r2.SetAttributes(alice)
		r2.SetTraceID(trace.TraceID(traceIDB))
		r2.SetSpanID(trace.SpanID(spanIDB))
		r2.SetTraceFlags(trace.TraceFlags(flagsB))

		r3 := new(log.Record)
		r3.SetTimestamp(ts)
		r3.SetObservedTimestamp(obs)
		r3.SetSeverity(sevB)
		r3.SetSeverityText("B")
		r3.SetBody(bodyB)
		r3.SetAttributes(bob)
		r3.SetTraceID(trace.TraceID(traceIDB))
		r3.SetSpanID(trace.SpanID(spanIDB))
		r3.SetTraceFlags(trace.TraceFlags(flagsB))

		return []log.Record{*r0, *r1, *r2, *r3}
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

	pbScopeLogs = &lpb.ScopeLogs{LogRecords: pbLogRecords}

	pbResourceLogs = &lpb.ResourceLogs{
		ScopeLogs: []*lpb.ScopeLogs{pbScopeLogs},
	}
)

func TestResourceLogs(t *testing.T) {
	want := []*lpb.ResourceLogs{pbResourceLogs}
	out, free := ResourceLogs(records)
	assert.Equal(t, want, out)
	free()
	want = []*lpb.ResourceLogs{{
		ScopeLogs: []*lpb.ScopeLogs{{
			LogRecords: pbLogRecords[2:],
		}},
	}}
	out, free = ResourceLogs(records[2:])
	assert.Equal(t, want, out)
}

func TestSeverityNumber(t *testing.T) {
	for i := 0; i <= int(api.SeverityFatal4); i++ {
		want := lpb.SeverityNumber(i)
		want += lpb.SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED
		assert.Equal(t, want, SeverityNumber(api.Severity(i)))
	}
}
