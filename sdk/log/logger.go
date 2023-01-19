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

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/trace"
)

type logger struct {
	provider             *LoggerProvider
	instrumentationScope instrumentation.Scope
}

var _ log.Logger = &logger{}

// Emit starts a Span and returns it along with a context containing it.
//
// The Span is created with the provided name and as a child of any existing
// span context found in the passed context. The created Span will be
// configured appropriately by any SpanOption passed.
func (tr *logger) Emit(ctx context.Context, options ...log.LogRecordOption) {
	config := log.NewLogRecordConfig(options...)

	if ctx == nil {
		// Prevent log.ContextWithSpan from panicking.
		ctx = context.Background()
	}

	s := tr.newLogRecord(ctx, &config)
	if rw, ok := s.(ReadWriteLogRecord); ok {
		sps := tr.provider.logRecordProcessors.Load().(logRecordProcessorStates)
		for _, sp := range sps {
			sp.sp.OnEmit(ctx, rw)
		}
	}
}

// newLogRecord returns a new configured span.
func (tr *logger) newLogRecord(ctx context.Context, config *log.LogRecordConfig) log.LogRecord {
	// If told explicitly to make this a new root use a zero value SpanContext
	// as a parent which contains an invalid trace ID and is not remote.
	psc := trace.SpanContextFromContext(ctx)
	return tr.newRecordingLogRecord(psc, config)
}

// newRecordingLogRecord returns a new configured recordingLogRecord.
func (tr *logger) newRecordingLogRecord(
	psc trace.SpanContext, config *log.LogRecordConfig,
) *recordingLogRecord {
	timestamp := config.Timestamp()
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	s := &recordingLogRecord{
		// Do not pre-allocate the attributes slice here! Doing so will
		// allocate memory that is likely never going to be used, or if used,
		// will be over-sized. The default Go compiler has been tested to
		// dynamically allocate needed space very well. Benchmarking has shown
		// it to be more performant than what we can predetermine here,
		// especially for the common use case of few to no added
		// attributes.
		timestamp:   timestamp,
		logger:      tr,
		spanContext: psc,
	}

	s.SetAttributes(config.Attributes()...)

	return s
}
