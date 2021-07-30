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

package tracetest // import "go.opentelemetry.io/otel/sdk/trace/tracetest"

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SpanRecorder records started and ended spans.
type SpanRecorder struct {
	started []sdktrace.ReadWriteSpan
	ended   []sdktrace.ReadOnlySpan
}

var _ sdktrace.SpanProcessor = (*SpanRecorder)(nil)

func NewSpanRecorder() *SpanRecorder {
	return new(SpanRecorder)
}

// OnStart records started spans.
func (sr *SpanRecorder) OnStart(_ context.Context, s sdktrace.ReadWriteSpan) {
	sr.started = append(sr.started, s)
}

// OnEnd records completed spans.
func (sr *SpanRecorder) OnEnd(s sdktrace.ReadOnlySpan) {
	sr.ended = append(sr.ended, s)
}

// Shutdown does nothing.
func (sr *SpanRecorder) Shutdown(context.Context) error {
	return nil
}

// ForceFlush does nothing.
func (sr *SpanRecorder) ForceFlush(context.Context) error {
	return nil
}

// Started returns a copy of all started spans that have been recorded.
func (sr *SpanRecorder) Started() []sdktrace.ReadWriteSpan {
	dst := make([]sdktrace.ReadWriteSpan, len(sr.started))
	copy(dst, sr.started)
	return dst
}

// Ended returns a copy of all ended spans that have been recorded.
func (sr *SpanRecorder) Ended() []sdktrace.ReadOnlySpan {
	dst := make([]sdktrace.ReadOnlySpan, len(sr.ended))
	copy(dst, sr.ended)
	return dst
}
