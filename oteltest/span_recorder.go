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

package oteltest

import (
	"context"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SpanRecorder performs operations to record a span as it starts and ends.
// It is designed to be concurrent safe and can by used by multiple goroutines.
type SpanRecorder struct {
	startedMu sync.RWMutex
	started   []sdktrace.ReadWriteSpan

	doneMu sync.RWMutex
	done   []sdktrace.ReadOnlySpan
}

func (ssr *SpanRecorder) Shutdown(ctx context.Context) error {
	return nil
}

func (ssr *SpanRecorder) ForceFlush(ctx context.Context) error {
	return nil
}

// OnStart records span as started.
func (ssr *SpanRecorder) OnStart(_ context.Context, s sdktrace.ReadWriteSpan) {
	ssr.startedMu.Lock()
	defer ssr.startedMu.Unlock()
	ssr.started = append(ssr.started, s)
}

// OnEnd records span as completed.
func (ssr *SpanRecorder) OnEnd(s sdktrace.ReadOnlySpan) {
	ssr.doneMu.Lock()
	defer ssr.doneMu.Unlock()
	ssr.done = append(ssr.done, s)
}

// Started returns a copy of all started Spans in the order they were started.
func (ssr *SpanRecorder) Started() []sdktrace.ReadOnlySpan {
	ssr.startedMu.RLock()
	defer ssr.startedMu.RUnlock()
	started := make([]sdktrace.ReadOnlySpan, len(ssr.started))
	for i := range ssr.started {
		started[i] = ssr.started[i]
	}
	return started
}

// Completed returns a copy of all ended Spans in the order they were ended.
func (ssr *SpanRecorder) Completed() []sdktrace.ReadOnlySpan {
	ssr.doneMu.RLock()
	defer ssr.doneMu.RUnlock()
	done := make([]sdktrace.ReadOnlySpan, len(ssr.done))
	for i := range ssr.done {
		done[i] = ssr.done[i]
	}
	return done
}
