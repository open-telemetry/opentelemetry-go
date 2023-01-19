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
	"sync"
)

// LogRecordProcessor is a processing pipeline for logRecords in the trace signal.
// LogRecordProcessors registered with a LoggerProvider and are called at the start
// and end of a Span's lifecycle, and are called in the order they are
// registered.
type LogRecordProcessor interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// OnEmit is called when a span is started. It is called synchronously
	// and should not block.
	OnEmit(parent context.Context, s ReadWriteLogRecord)
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the processor should be done in this call.
	//
	// Calls to OnStart, OnEnd, or ForceFlush after this has been called
	// should be ignored.
	//
	// All timeouts and cancellations contained in ctx must be honored, this
	// should not block indefinitely.
	Shutdown(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// ForceFlush exports all ended logRecords to the configured Exporter that have not yet
	// been exported.  It should only be called when absolutely necessary, such as when
	// using a FaaS provider that may suspend the process after an invocation, but before
	// the Processor can export the completed logRecords.
	ForceFlush(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}

type logRecordProcessorState struct {
	sp    LogRecordProcessor
	state *sync.Once
}

func newLogRecordProcessorState(sp LogRecordProcessor) *logRecordProcessorState {
	return &logRecordProcessorState{sp: sp, state: &sync.Once{}}
}

type logRecordProcessorStates []*logRecordProcessorState
