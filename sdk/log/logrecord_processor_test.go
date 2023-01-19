// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

type testLogRecordProcessor struct {
	name          string
	logRecords    []sdklog.ReadWriteLogRecord
	shutdownCount int
}

func (t *testLogRecordProcessor) OnEmit(parent context.Context, s sdklog.ReadWriteLogRecord) {
	if t == nil {
		return
	}
	psc := trace.SpanContextFromContext(parent)
	kv := []attribute.KeyValue{
		{
			Key:   attribute.Key("SpanProcessorName" + t.name),
			Value: attribute.StringValue(t.name),
		},
		// Store parent trace ID and span ID as attributes to be read later in
		// tests so that we "do something" with the parent argument. Real
		// LogRecordProcessor implementations will likely use the parent argument in
		// a more meaningful way.
		{
			Key:   attribute.Key("ParentTraceID" + t.name),
			Value: attribute.StringValue(psc.TraceID().String()),
		},
		{
			Key:   attribute.Key("ParentSpanID" + t.name),
			Value: attribute.StringValue(psc.SpanID().String()),
		},
	}
	s.SetAttributes(kv...)
	t.logRecords = append(t.logRecords, s)
}

func (t *testLogRecordProcessor) Shutdown(context.Context) error {
	if t == nil {
		return nil
	}
	t.shutdownCount++
	return nil
}

func (t *testLogRecordProcessor) ForceFlush(context.Context) error {
	if t == nil {
		return nil
	}
	return nil
}

func TestRegisterLogRecordProcessor(t *testing.T) {
	name := "Register span processor before span starts"
	tp := basicLoggerProvider(t)
	spNames := []string{"sp1", "sp2", "sp3"}
	sps := NewNamedTestLogRecordProcessors(spNames)

	for _, sp := range sps {
		tp.RegisterLogRecordProcessor(sp)
	}

	tid, _ := trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ := trace.SpanIDFromHex("0102040810203040")
	parent := trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID: tid,
			SpanID:  sid,
		},
	)
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), parent)

	tr := tp.Logger("LogRecordProcessor")
	tr.Emit(ctx)
	wantCount := 1

	for _, sp := range sps {
		gotCount := len(sp.logRecords)
		if gotCount != wantCount {
			t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
		}

		nameOK := false
		tidOK := false
		sidOK := false
		for _, kv := range sp.logRecords[0].Attributes() {
			switch kv.Key {
			case attribute.Key("SpanProcessorName" + sp.name):
				gotValue := kv.Value.AsString()
				if gotValue != sp.name {
					t.Errorf("%s: attributes: got %s, want %s\n", name, gotValue, sp.name)
				}
				nameOK = true
			case attribute.Key("ParentTraceID" + sp.name):
				gotValue := kv.Value.AsString()
				if gotValue != parent.TraceID().String() {
					t.Errorf("%s: attributes: got %s, want %s\n", name, gotValue, parent.TraceID())
				}
				tidOK = true
			case attribute.Key("ParentSpanID" + sp.name):
				gotValue := kv.Value.AsString()
				if gotValue != parent.SpanID().String() {
					t.Errorf("%s: attributes: got %s, want %s\n", name, gotValue, parent.SpanID())
				}
				sidOK = true
			default:
				continue
			}
		}
		if !nameOK {
			t.Errorf("%s: expected attributes(SpanProcessorName)\n", name)
		}
		if !tidOK {
			t.Errorf("%s: expected attributes(ParentTraceID)\n", name)
		}
		if !sidOK {
			t.Errorf("%s: expected attributes(ParentSpanID)\n", name)
		}
	}
}

func TestUnregisterLogRecordProcessor(t *testing.T) {
	name := "Start span after unregistering span processor"
	tp := basicLoggerProvider(t)
	spNames := []string{"sp1", "sp2", "sp3"}
	sps := NewNamedTestLogRecordProcessors(spNames)

	for _, sp := range sps {
		tp.RegisterLogRecordProcessor(sp)
	}

	tr := tp.Logger("LogRecordProcessor")
	tr.Emit(context.Background())
	for _, sp := range sps {
		tp.UnregisterLogRecordProcessor(sp)
	}

	// start another span after unregistering span processor.
	tr.Emit(context.Background())

	for _, sp := range sps {
		wantCount := 1
		gotCount := len(sp.logRecords)
		if gotCount != wantCount {
			t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
		}
	}
}

func TestUnregisterLogRecordProcessorWhileSpanIsActive(t *testing.T) {
	name := "Unregister span processor while span is active"
	tp := basicLoggerProvider(t)
	sp := NewTestLogRecordProcessor("sp")
	tp.RegisterLogRecordProcessor(sp)

	tr := tp.Logger("LogRecordProcessor")
	tr.Emit(context.Background())
	tp.UnregisterLogRecordProcessor(sp)

	wantCount := 1
	gotCount := len(sp.logRecords)
	if gotCount != wantCount {
		t.Errorf("%s: started count: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestSpanProcessorShutdown(t *testing.T) {
	name := "Increment shutdown counter of a span processor"
	tp := basicLoggerProvider(t)
	sp := NewTestLogRecordProcessor("sp")
	tp.RegisterLogRecordProcessor(sp)

	wantCount := 1
	err := sp.Shutdown(context.Background())
	if err != nil {
		t.Error("Error shutting the testLogRecordProcessor down\n")
	}

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func TestMultipleUnregisterLogRecordProcessorCalls(t *testing.T) {
	name := "Increment shutdown counter after first UnregisterLogRecordProcessor call"
	tp := basicLoggerProvider(t)
	sp := NewTestLogRecordProcessor("sp")

	wantCount := 1

	tp.RegisterLogRecordProcessor(sp)
	tp.UnregisterLogRecordProcessor(sp)

	gotCount := sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}

	// Multiple UnregisterLogRecordProcessor should not trigger multiple Shutdown calls.
	tp.UnregisterLogRecordProcessor(sp)

	gotCount = sp.shutdownCount
	if wantCount != gotCount {
		t.Errorf("%s: wrong counter: got %d, want %d\n", name, gotCount, wantCount)
	}
}

func NewTestLogRecordProcessor(name string) *testLogRecordProcessor {
	return &testLogRecordProcessor{name: name}
}

func NewNamedTestLogRecordProcessors(names []string) []*testLogRecordProcessor {
	tsp := []*testLogRecordProcessor{}
	for _, n := range names {
		tsp = append(tsp, NewTestLogRecordProcessor(n))
	}
	return tsp
}
