// Copyright 2019, OpenTelemetry Authors
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

package testtrace

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/distributedcontext"
	"go.opentelemetry.io/api/trace"
)

var _ trace.Span = (*Span)(nil)

type Span struct {
	lock        *sync.Mutex
	tracer      *Tracer
	spanContext core.SpanContext
	ended       bool
	name        string
	attributes  []core.KeyValue
	startTime   time.Time
	endTime     time.Time
	code        codes.Code
}

func (s *Span) Tracer() trace.Tracer {
	return s.tracer
}

func (s *Span) End(opts ...trace.EndOption) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		s.endTime = time.Now()
		s.ended = true
	}
}

func (s *Span) AddEvent(ctx context.Context, msg string, attrs ...core.KeyValue) {}

func (s *Span) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...core.KeyValue) {
}

func (s *Span) IsRecording() bool {
	return true
}

func (s *Span) AddLink(link trace.Link) {}

func (s *Span) Link(sc core.SpanContext, attrs ...core.KeyValue) {}

func (s *Span) SpanContext() core.SpanContext {
	return s.spanContext
}

func (s *Span) SetStatus(code codes.Code) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.code = code
}

func (s *Span) SetName(name string) {}

func (s *Span) SetAttribute(attr core.KeyValue) {
	s.SetAttributes(attr)
}

func (s *Span) SetAttributes(attrs ...core.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.attributes = append(s.attributes, attrs...)
}

func (s *Span) ModifyAttribute(mutator distributedcontext.Mutator) {}

func (s *Span) ModifyAttributes(mutators ...distributedcontext.Mutator) {}

func (s *Span) Name() string {
	return s.name
}

func (s *Span) Attributes() []core.KeyValue {
	return append([]core.KeyValue{}, s.attributes...)
}

func (s *Span) StartTime() time.Time {
	return s.startTime
}

func (s *Span) EndTime() (time.Time, bool) {
	return s.endTime, s.ended
}

func (s *Span) Ended() bool {
	return s.ended
}

func (s *Span) Status() codes.Code {
	return s.code
}
