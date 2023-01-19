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
	"errors"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

var (
	tid, _ = trace.TraceIDFromHex("01020304050607080102040810203040")
	sid, _ = trace.SpanIDFromHex("0102040810203040")
)

type testExporter struct {
	logRecords []sdklog.ReadOnlyLogRecord
	shutdown   bool
}

func (t *testExporter) ExportLogRecords(ctx context.Context, logRecords []sdklog.ReadOnlyLogRecord) error {
	t.logRecords = append(t.logRecords, logRecords...)
	return nil
}

func (t *testExporter) Shutdown(ctx context.Context) error {
	t.shutdown = true
	select {
	case <-ctx.Done():
		// Ensure context deadline tests receive the expected error.
		return ctx.Err()
	default:
		return nil
	}
}

var _ sdklog.LogRecordExporter = (*testExporter)(nil)

func TestNewSimpleLogRecordProcessor(t *testing.T) {
	if ssp := sdklog.NewSimpleLogRecordProcessor(&testExporter{}); ssp == nil {
		t.Error("failed to create new SimpleLogRecordProcessor")
	}
}

func TestNewSimpleLogRecordProcessorWithNilExporter(t *testing.T) {
	if ssp := sdklog.NewSimpleLogRecordProcessor(nil); ssp == nil {
		t.Error("failed to create new SimpleLogRecordProcessor with nil exporter")
	}
}

func emitLogRecord(tp log.LoggerProvider) {
	tr := tp.Logger("SimpleLogRecordProcessor")
	sc := trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID:    tid,
			SpanID:     sid,
			TraceFlags: 0x1,
		},
	)
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)
	tr.Emit(ctx)
}

func TestSimpleLogRecordProcessorOnEnd(t *testing.T) {
	tp := basicLoggerProvider(t)
	te := testExporter{}
	ssp := sdklog.NewSimpleLogRecordProcessor(&te)

	tp.RegisterLogRecordProcessor(ssp)
	emitLogRecord(tp)

	wantTraceID := tid
	gotTraceID := te.logRecords[0].SpanContext().TraceID()
	if wantTraceID != gotTraceID {
		t.Errorf("SimpleLogRecordProcessor OnEmit() check: got %+v, want %+v\n", gotTraceID, wantTraceID)
	}
}

func TestSimpleLogRecordProcessorShutdown(t *testing.T) {
	exporter := &testExporter{}
	ssp := sdklog.NewSimpleLogRecordProcessor(exporter)

	// Ensure we can export a span before we test we cannot after shutdown.
	tp := basicLoggerProvider(t)
	tp.RegisterLogRecordProcessor(ssp)
	emitLogRecord(tp)
	nExported := len(exporter.logRecords)
	if nExported != 1 {
		t.Error("failed to verify span export")
	}

	if err := ssp.Shutdown(context.Background()); err != nil {
		t.Errorf("shutting the SimpleLogRecordProcessor down: %v", err)
	}
	if !exporter.shutdown {
		t.Error("SimpleLogRecordProcessor.Shutdown did not shut down exporter")
	}

	emitLogRecord(tp)
	if len(exporter.logRecords) > nExported {
		t.Error("exported span to shutdown exporter")
	}
}

func TestSimpleLogRecordProcessorShutdownOnEndConcurrency(t *testing.T) {
	exporter := &testExporter{}
	ssp := sdklog.NewSimpleLogRecordProcessor(exporter)
	tp := basicLoggerProvider(t)
	tp.RegisterLogRecordProcessor(ssp)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for {
			select {
			case <-stop:
				return
			default:
				emitLogRecord(tp)
			}
		}
	}()

	if err := ssp.Shutdown(context.Background()); err != nil {
		t.Errorf("shutting the SimpleLogRecordProcessor down: %v", err)
	}
	if !exporter.shutdown {
		t.Error("SimpleLogRecordProcessor.Shutdown did not shut down exporter")
	}

	stop <- struct{}{}
	<-done
}

func TestSimpleLogRecordProcessorShutdownHonorsContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	ssp := sdklog.NewSimpleLogRecordProcessor(&testExporter{})
	if got, want := ssp.Shutdown(ctx), context.DeadlineExceeded; !errors.Is(got, want) {
		t.Errorf("SimpleLogRecordProcessor.Shutdown did not return %v, got %v", want, got)
	}
}

func TestSimpleLogRecordProcessorShutdownHonorsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ssp := sdklog.NewSimpleLogRecordProcessor(&testExporter{})
	if got, want := ssp.Shutdown(ctx), context.Canceled; !errors.Is(got, want) {
		t.Errorf("SimpleLogRecordProcessor.Shutdown did not return %v, got %v", want, got)
	}
}
