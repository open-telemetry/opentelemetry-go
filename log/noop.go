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

package log // import "go.opentelemetry.io/otel/log"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// NewNoopLoggerProvider returns an implementation of LoggerProvider that
// performs no operations. The Logger and Spans created from the returned
// LoggerProvider also perform no operations.
func NewNoopLoggerProvider() LoggerProvider {
	return noopLoggerProvider{}
}

type noopLoggerProvider struct{}

var _ LoggerProvider = noopLoggerProvider{}

// Logger returns noop implementation of Logger.
func (p noopLoggerProvider) Logger(string, ...LoggerOption) Logger {
	return noopLogger{}
}

// noopLogger is an implementation of Logger that preforms no operations.
type noopLogger struct{}

var _ Logger = noopLogger{}

// Emit carries forward a non-recording LogRecord, if one is present in the context, otherwise it
// creates a no-op LogRecord.
func (t noopLogger) Emit(ctx context.Context, _ ...LogRecordOption) {
}

// noopLogRecord is an implementation of LogRecord that preforms no operations.
type noopLogRecord struct{}

var _ LogRecord = noopLogRecord{}

// IsRecording always returns false.
func (noopLogRecord) IsRecording() bool { return false }

// SetAttributes does nothing.
func (noopLogRecord) SetAttributes(...attribute.KeyValue) {}

// LoggerProvider returns a no-op LoggerProvider.
func (noopLogRecord) LoggerProvider() LoggerProvider { return noopLoggerProvider{} }
