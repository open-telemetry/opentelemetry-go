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

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
)

var _ trace.Span = (*Span)(nil)

type Span struct {
	lock         *sync.Mutex
	tracer       *Tracer
	spanContext  core.SpanContext
	parentSpanID core.SpanID
	ended        bool
	name         string
	startTime    time.Time
	endTime      time.Time
	status       codes.Code
	attributes   map[core.Key]core.Value
	events       []Event
	links        map[core.SpanContext][]core.KeyValue
}

func (s *Span) Tracer() trace.Tracer {
	return s.tracer
}

func (s *Span) End(opts ...trace.EndOption) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		var c trace.EndOptions

		for _, opt := range opts {
			opt(&c)
		}

		s.endTime = time.Now()

		if endTime := c.EndTime; !endTime.IsZero() {
			s.endTime = endTime
		}

		s.ended = true
	}
}

func (s *Span) AddEvent(ctx context.Context, msg string, attrs ...core.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), msg, attrs...)
}

func (s *Span) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...core.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		attributes := make(map[core.Key]core.Value)

		for _, attr := range attrs {
			attributes[attr.Key] = attr.Value
		}

		s.events = append(s.events, Event{
			Timestamp:  timestamp,
			Message:    msg,
			Attributes: attributes,
		})
	}
}

func (s *Span) IsRecording() bool {
	return true
}

func (s *Span) AddLink(link trace.Link) {
	s.Link(link.SpanContext, link.Attributes...)
}

func (s *Span) Link(sc core.SpanContext, attrs ...core.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		s.links[sc] = attrs
	}
}

func (s *Span) SpanContext() core.SpanContext {
	return s.spanContext
}

func (s *Span) SetStatus(status codes.Code) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		s.status = status
	}
}

func (s *Span) SetName(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		s.name = name
	}
}

func (s *Span) SetAttribute(attr core.KeyValue) {
	s.SetAttributes(attr)
}

func (s *Span) SetAttributes(attrs ...core.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.ended {
		for _, attr := range attrs {
			s.attributes[attr.Key] = attr.Value
		}
	}
}

func (s *Span) Name() string {
	return s.name
}

func (s *Span) ParentSpanID() core.SpanID {
	return s.parentSpanID
}

func (s *Span) Attributes() map[core.Key]core.Value {
	return s.attributes
}

func (s *Span) Events() []Event {
	return s.events
}

func (s *Span) Links() map[core.SpanContext][]core.KeyValue {
	links := make(map[core.SpanContext][]core.KeyValue)

	for sc, attributes := range s.links {
		links[sc] = append([]core.KeyValue{}, attributes...)
	}

	return links
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
	return s.status
}
