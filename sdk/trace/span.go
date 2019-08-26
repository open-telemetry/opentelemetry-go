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

package trace

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	apievent "go.opentelemetry.io/api/event"
	apitag "go.opentelemetry.io/api/tag"
	apitrace "go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/sdk/internal"
)

// span implements apitrace.Span interface.
type span struct {
	// data contains information recorded about the span.
	//
	// It will be non-nil if we are exporting the span or recording events for it.
	// Otherwise, data is nil, and the span is simply a carrier for the
	// SpanContext, so that the trace ID is propagated.
	data        *SpanData
	mu          sync.Mutex // protects the contents of *data (but not the pointer value.)
	spanContext core.SpanContext

	// lruAttributes are capped at configured limit. When the capacity is reached an oldest entry
	// is removed to create room for a new entry.
	lruAttributes *lruMap

	// messageEvents are stored in FIFO queue capped by configured limit.
	messageEvents *evictedQueue

	// links are stored in FIFO queue capped by configured limit.
	links *evictedQueue

	// spanStore is the spanStore this span belongs to, if any, otherwise it is nil.
	//*spanStore
	endOnce sync.Once

	executionTracerTaskEnd func()          // ends the execution tracer span
	tracer                 apitrace.Tracer // tracer used to create span.
}

var _ apitrace.Span = &span{}

func (s *span) SpanContext() core.SpanContext {
	if s == nil {
		return core.EmptySpanContext()
	}
	return s.spanContext
}

func (s *span) IsRecordingEvents() bool {
	if s == nil {
		return false
	}
	return s.data != nil
}

func (s *span) SetStatus(status codes.Code) {
	if s == nil {
		return
	}
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.data.Status = status
	s.mu.Unlock()
}

func (s *span) SetAttribute(attribute core.KeyValue) {
	if !s.IsRecordingEvents() {
		return
	}
	s.copyToCappedAttributes(attribute)
}

func (s *span) SetAttributes(attributes ...core.KeyValue) {
	if !s.IsRecordingEvents() {
		return
	}
	s.copyToCappedAttributes(attributes...)
}

// ModifyAttribute does nothing.
func (s *span) ModifyAttribute(mutator apitag.Mutator) {
}

// ModifyAttributes does nothing.
func (s *span) ModifyAttributes(mutators ...apitag.Mutator) {
}

func (s *span) Finish() {
	if s == nil {
		return
	}

	if s.executionTracerTaskEnd != nil {
		s.executionTracerTaskEnd()
	}
	if !s.IsRecordingEvents() {
		return
	}
	s.endOnce.Do(func() {
		exp, _ := exporters.Load().(exportersMap)
		mustExport := s.spanContext.IsSampled() && len(exp) > 0
		//if s.spanStore != nil || mustExport {
		if mustExport {
			sd := s.makeSpanData()
			sd.EndTime = internal.MonotonicEndTime(sd.StartTime)
			//if s.spanStore != nil {
			//	s.spanStore.finished(s, sd)
			//}
			if mustExport {
				for e := range exp {
					e.ExportSpan(sd)
				}
			}
		}
	})
}

func (s *span) Tracer() apitrace.Tracer {
	return s.tracer
}

func (s *span) AddEvent(ctx context.Context, event apievent.Event) {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageEvents.add(event)
}

func (s *span) Event(ctx context.Context, msg string, attrs ...core.KeyValue) {
	if !s.IsRecordingEvents() {
		return
	}
	now := time.Now()
	s.mu.Lock()
	s.messageEvents.add(event{
		msg:        msg,
		attributes: attrs,
		time:       now,
	})
	s.mu.Unlock()
}

// makeSpanData produces a SpanData representing the current state of the span.
// It requires that s.data is non-nil.
func (s *span) makeSpanData() *SpanData {
	var sd SpanData
	s.mu.Lock()
	defer s.mu.Unlock()
	sd = *s.data
	if s.lruAttributes.simpleLruMap.Len() > 0 {
		sd.Attributes = s.lruAttributesToAttributeMap()
		sd.DroppedAttributeCount = s.lruAttributes.droppedCount
	}
	if len(s.messageEvents.queue) > 0 {
		sd.MessageEvents = s.interfaceArrayToMessageEventArray()
		sd.DroppedMessageEventCount = s.messageEvents.droppedCount
	}
	return &sd
}

func (s *span) interfaceArrayToMessageEventArray() []event {
	messageEventArr := make([]event, 0)
	for _, value := range s.messageEvents.queue {
		messageEventArr = append(messageEventArr, value.(event))
	}
	return messageEventArr
}

func (s *span) lruAttributesToAttributeMap() map[string]interface{} {
	attributes := make(map[string]interface{})
	for _, key := range s.lruAttributes.simpleLruMap.Keys() {
		value, ok := s.lruAttributes.simpleLruMap.Get(key)
		if ok {
			key := key.(core.Key)
			attributes[key.Variable.Name] = value
		}
	}
	return attributes
}

func (s *span) copyToCappedAttributes(attributes ...core.KeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range attributes {
		s.lruAttributes.add(a.Key, a.Value)
	}
}

func (s *span) addChild() {
	if !s.IsRecordingEvents() {
		return
	}
	s.mu.Lock()
	s.data.ChildSpanCount++
	s.mu.Unlock()
}

func startSpanInternal(name string, parent core.SpanContext, remoteParent bool, o apitrace.SpanOptions) *span {
	var noParent bool
	span := &span{}
	span.spanContext = parent

	cfg := config.Load().(*Config)

	if parent == core.EmptySpanContext() {
		span.spanContext.TraceID = cfg.IDGenerator.NewTraceID()
		noParent = true
	}
	span.spanContext.SpanID = cfg.IDGenerator.NewSpanID()
	sampler := cfg.DefaultSampler

	// TODO: [rghetia] fix sampler
	//if !hasParent || remoteParent || o.Sampler != nil {
	if noParent || remoteParent {
		// If this span is the child of a local span and no Sampler is set in the
		// options, keep the parent's TraceOptions.
		//
		// Otherwise, consult the Sampler in the options if it is non-nil, otherwise
		// the default sampler.
		//if o.Sampler != nil {
		//	sampler = o.Sampler
		//}
		sampled := sampler(SamplingParameters{
			ParentContext:   parent,
			TraceID:         span.spanContext.TraceID,
			SpanID:          span.spanContext.SpanID,
			Name:            name,
			HasRemoteParent: remoteParent}).Sample
		if sampled {
			span.spanContext.TraceOptions = core.TraceOptionSampled
		}
	}

	// TODO: [rghetia] restore when spanstore is added.
	// if !internal.LocalSpanStoreEnabled && !span.spanContext.IsSampled() && !o.RecordEvent {
	if !span.spanContext.IsSampled() && !o.RecordEvent {
		return span
	}

	span.data = &SpanData{
		SpanContext: span.spanContext,
		StartTime:   time.Now(),
		// TODO;[rghetia] : fix spanKind
		//SpanKind:        o.SpanKind,
		Name:            name,
		HasRemoteParent: remoteParent,
	}
	span.lruAttributes = newLruMap(cfg.MaxAttributesPerSpan)
	span.messageEvents = newEvictedQueue(cfg.MaxEventsPerSpan)
	span.links = newEvictedQueue(cfg.MaxLinksPerSpan)

	if !noParent {
		span.data.ParentSpanID = parent.SpanID
	}
	// TODO: [rghetia] restore when spanstore is added.
	//if internal.LocalSpanStoreEnabled {
	//	ss := spanStoreForNameCreateIfNew(name)
	//	if ss != nil {
	//		span.spanStore = ss
	//		ss.add(span)
	//	}
	//}

	return span
}
