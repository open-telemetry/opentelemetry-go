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

package reader

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/distributedcontext"
	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

type Reader interface {
	Read(Event)
}

type Event struct {
	Type         exporter.EventType
	Time         time.Time
	Sequence     exporter.EventID
	SpanContext  core.SpanContext
	Entries      distributedcontext.Map // context entries
	Attributes   distributedcontext.Map // span attributes, metric labels
	Measurement  metric.Measurement
	Measurements []metric.Measurement

	Parent           core.SpanContext
	ParentAttributes distributedcontext.Map

	Duration time.Duration
	Name     string
	Message  string
	Status   codes.Code
}

type readerObserver struct {
	readers []Reader

	// core.EventID -> *readerSpan or *readerScope
	scopes sync.Map
}

type readerSpan struct {
	name         string
	start        time.Time
	startEntries distributedcontext.Map
	spanContext  core.SpanContext
	status       codes.Code

	*readerScope
}

type readerScope struct {
	span       *readerSpan
	parent     exporter.EventID
	attributes distributedcontext.Map
}

// NewReaderObserver returns an implementation that computes the
// necessary state needed by a reader to process events in memory.
// Practically, this means tracking live metric handles and scope
// attribute sets.
func NewReaderObserver(readers ...Reader) exporter.Observer {
	return &readerObserver{
		readers: readers,
	}
}

func (ro *readerObserver) Observe(event exporter.Event) {
	// TODO this should check for out-of-order events and buffer.
	ro.orderedObserve(event)
}

func (ro *readerObserver) orderedObserve(event exporter.Event) {
	read := Event{
		Time:       event.Time,
		Sequence:   event.Sequence,
		Attributes: distributedcontext.NewEmptyMap(),
		Entries:    distributedcontext.NewEmptyMap(),
	}

	if event.Context != nil {
		read.Entries = distributedcontext.FromContext(event.Context)
	}

	switch event.Type {
	case exporter.START_SPAN:
		// Save the span context entries, initial attributes, start time, and name.
		span := &readerSpan{
			name:         event.String,
			start:        event.Time,
			startEntries: distributedcontext.FromContext(event.Context),
			spanContext:  event.Scope.SpanContext,
			readerScope:  &readerScope{},
		}

		rattrs, _ := ro.readScope(event.Scope)

		span.readerScope.span = span
		span.readerScope.attributes = rattrs

		read.Name = span.name
		read.Type = exporter.START_SPAN
		read.SpanContext = span.spanContext
		read.Attributes = rattrs

		if event.Parent.EventID == 0 && event.Parent.HasTraceID() {
			// Remote parent
			read.Parent = event.Parent.SpanContext

			// Note: No parent attributes in the event for remote parents.
		} else {
			pattrs, pspan := ro.readScope(event.Parent)

			if pspan != nil {
				// Local parent
				read.Parent = pspan.spanContext
				read.ParentAttributes = pattrs
			}
		}

		ro.scopes.Store(event.Sequence, span)

	case exporter.END_SPAN:
		attrs, span := ro.readScope(event.Scope)
		if span == nil {
			panic(fmt.Sprint("span not found", event.Scope))
		}

		read.Name = span.name
		read.Type = exporter.END_SPAN

		read.Attributes = attrs
		read.Duration = event.Time.Sub(span.start)
		read.Entries = span.startEntries
		read.SpanContext = span.spanContext

		// TODO: recovered

	case exporter.NEW_SCOPE, exporter.MODIFY_ATTR:
		var span *readerSpan
		var m distributedcontext.Map

		sid := event.Scope

		if sid.EventID == 0 {
			m = distributedcontext.NewEmptyMap()
		} else {
			parentI, has := ro.scopes.Load(sid.EventID)
			if !has {
				panic("parent scope not found")
			}
			if parent, ok := parentI.(*readerScope); ok {
				m = parent.attributes
				span = parent.span
			} else if parent, ok := parentI.(*readerSpan); ok {
				m = parent.attributes
				span = parent
			}
		}

		sc := &readerScope{
			span:   span,
			parent: sid.EventID,
			attributes: m.Apply(
				distributedcontext.MapUpdate{
					SingleKV:      event.Attribute,
					MultiKV:       event.Attributes,
					SingleMutator: event.Mutator,
					MultiMutator:  event.Mutators,
				},
			),
		}

		ro.scopes.Store(event.Sequence, sc)

		if event.Type == exporter.NEW_SCOPE {
			return
		}

		read.Type = exporter.MODIFY_ATTR
		read.Attributes = sc.attributes

		if span != nil {
			read.SpanContext = span.spanContext
			read.Entries = span.startEntries
		}

	case exporter.ADD_EVENT:
		read.Type = exporter.ADD_EVENT
		read.Message = event.String

		attrs, span := ro.readScope(event.Scope)
		read.Attributes = attrs.Apply(distributedcontext.MapUpdate{
			MultiKV: event.Attributes,
		})
		if span != nil {
			read.SpanContext = span.spanContext
		}

	case exporter.SINGLE_METRIC:
		read.Type = exporter.SINGLE_METRIC

		if event.Context != nil {
			span := trace.CurrentSpan(event.Context)
			if span != nil {
				read.SpanContext = span.SpanContext()
			}
		}
		attrs, _ := ro.readScope(event.Scope)
		read.Attributes = attrs
		read.Measurement = event.Measurement

	case exporter.BATCH_METRIC:
		read.Type = event.Type

		if event.Context != nil {
			span := trace.CurrentSpan(event.Context)
			if span != nil {
				read.SpanContext = span.SpanContext()
			}
		}

		attrs, _ := ro.readScope(event.Scope)
		read.Attributes = attrs
		read.Measurements = make([]metric.Measurement, len(event.Measurements))
		copy(read.Measurements, event.Measurements)

	case exporter.SET_STATUS:
		read.Type = exporter.SET_STATUS
		read.Status = event.Status
		_, span := ro.readScope(event.Scope)
		if span != nil {
			span.status = event.Status
			read.SpanContext = span.spanContext
		}

	case exporter.SET_NAME:
		read.Type = exporter.SET_NAME
		read.Name = event.String

	default:
		panic(fmt.Sprint("Unhandled case: ", event.Type))
	}

	for _, reader := range ro.readers {
		reader.Read(read)
	}

	if event.Type == exporter.END_SPAN {
		ro.cleanupSpan(event.Scope.EventID)
	}
}

func (ro *readerObserver) readScope(id exporter.ScopeID) (distributedcontext.Map, *readerSpan) {
	if id.EventID == 0 {
		return distributedcontext.NewEmptyMap(), nil
	}
	ev, has := ro.scopes.Load(id.EventID)
	if !has {
		panic(fmt.Sprintln("scope not found", id.EventID))
	}
	if sp, ok := ev.(*readerScope); ok {
		return sp.attributes, sp.span
	} else if sp, ok := ev.(*readerSpan); ok {
		return sp.attributes, sp
	}
	return distributedcontext.NewEmptyMap(), nil
}

func (ro *readerObserver) cleanupSpan(id exporter.EventID) {
	for id != 0 {
		ev, has := ro.scopes.Load(id)
		if !has {
			panic(fmt.Sprintln("scope not found", id))
		}
		ro.scopes.Delete(id)

		if sp, ok := ev.(*readerScope); ok {
			id = sp.parent
		} else if sp, ok := ev.(*readerSpan); ok {
			id = sp.parent
		}
	}
}
