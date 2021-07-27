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

package internal

import (
	"testing"

	"go.opentelemetry.io/otel/trace"
)

type span struct {
	trace.Span

	recording bool
	ended     bool
	sc        trace.SpanContext
}

func (s *span) IsRecording() bool              { return s.recording }
func (s *span) End(...trace.SpanEndOption)     { s.ended = true }
func (s *span) SpanContext() trace.SpanContext { return s.sc }

func TestSpanIsRecordingEvents(t *testing.T) {
	s := &span{recording: true}
	ocS := NewSpan(s)
	if !ocS.IsRecordingEvents() {
		t.Errorf("span.IsRecordingEvents() = false, want true")
	}
	s.recording = false
	if ocS.IsRecordingEvents() {
		t.Errorf("span.IsRecordingEvents() = true, want false")
	}
}

func TestSpanEnd(t *testing.T) {
	s := new(span)
	ocS := NewSpan(s)
	if s.ended {
		t.Fatal("new span already ended")
	}

	ocS.End()
	if !s.ended {
		t.Error("span.End() did not end OpenTelemetry span")
	}
}

func TestSpanSpanContext(t *testing.T) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})
	converted := util
}
