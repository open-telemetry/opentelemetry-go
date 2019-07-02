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

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/event"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type Reader interface {
	Read(Event)
}

type EventType int

type Event struct {
	Type        EventType
	Time        time.Time
	Sequence    core.EventID
	SpanContext core.SpanContext
	Tags        tag.Map
	Attributes  tag.Map
	Event       event.Event
	Stats       []Measurement

	Parent           core.SpanContext
	ParentAttributes tag.Map

	Duration time.Duration
	Name     string
	Message  string
	Status   codes.Code
}

type Measurement struct {
	Measure core.Measure
	Value   float64
	Tags    tag.Map
}

type readerObserver struct {
	readers []Reader

	// core.EventID -> *readerSpan or *readerScope
	scopes sync.Map

	// core.EventID -> *readerMeasure
	measures sync.Map

	// core.EventID -> *readerMetric
	metrics sync.Map
}

type readerSpan struct {
	name        string
	start       time.Time
	startTags   tag.Map
	spanContext core.SpanContext
	status      codes.Code

	*readerScope
}

type readerMeasure struct {
	name string
	// TODO[rghetia]: comment to avoid compile errors. Remove it if not required
	// desc string
	// unit unit.Unit
}

type readerMetric struct {
	*readerMeasure
	// TODO[rghetia]: comment to avoid compile errors. Remove it if not required
	// mtype  metric.MetricType
	// fields []core.Measure
}

type readerScope struct {
	span       *readerSpan
	parent     core.EventID
	attributes tag.Map
}

const (
	INVALID EventType = iota
	START_SPAN
	FINISH_SPAN
	ADD_EVENT
	MODIFY_ATTR
	RECORD_STATS
	SET_STATUS
)

// NewReaderObserver returns an implementation that computes the
// necessary state needed by a reader to process events in memory.
// Practically, this means tracking live metric handles and scope
// attribute sets.
func NewReaderObserver(readers ...Reader) observer.Observer {
	return &readerObserver{
		readers: readers,
	}
}

func (ro *readerObserver) Observe(event observer.Event) {
	read := Event{
		Time:       event.Time,
		Sequence:   event.Sequence,
		Attributes: tag.NewEmptyMap(),
		Tags:       tag.NewEmptyMap(),
	}

	if event.Context != nil {
		read.Tags = tag.FromContext(event.Context)
	}

	switch event.Type {
	case observer.START_SPAN:
		// Save the span context tags, initial attributes, start time, and name.
		span := &readerSpan{
			name:        event.String,
			start:       event.Time,
			startTags:   tag.FromContext(event.Context),
			spanContext: event.Scope.SpanContext,
			readerScope: &readerScope{},
		}

		rattrs, _ := ro.readScope(event.Scope)

		span.readerScope.span = span
		span.readerScope.attributes = rattrs

		read.Name = span.name
		read.Type = START_SPAN
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

	case observer.FINISH_SPAN:
		attrs, span := ro.readScope(event.Scope)
		if span == nil {
			panic("span not found")
		}

		read.Name = span.name
		read.Type = FINISH_SPAN

		read.Attributes = attrs
		read.Duration = event.Time.Sub(span.start)
		read.Tags = span.startTags
		read.SpanContext = span.spanContext

		// TODO: recovered

	case observer.NEW_SCOPE, observer.MODIFY_ATTR:
		var span *readerSpan
		var m tag.Map
		var sid core.ScopeID

		if event.Scope.EventID == 0 {
			// TODO: This is racey. Do this at the call
			// site via Resources.
			sid = trace.GlobalTracer().ScopeID()
		} else {
			sid = event.Scope
		}
		if sid.EventID == 0 {
			m = tag.NewEmptyMap()
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
				event.Attribute,
				event.Attributes,
				event.Mutator,
				event.Mutators,
			),
		}

		ro.scopes.Store(event.Sequence, sc)

		if event.Type == observer.NEW_SCOPE {
			return
		}

		read.Type = MODIFY_ATTR
		read.Attributes = sc.attributes

		if span != nil {
			read.SpanContext = span.spanContext
			read.Tags = span.startTags
		}

	case observer.NEW_MEASURE:
		measure := &readerMeasure{
			name: event.String,
		}
		ro.measures.Store(event.Sequence, measure)
		return

	case observer.NEW_METRIC:
		measureI, has := ro.measures.Load(event.Scope.EventID)
		if !has {
			panic("metric measure not found")
		}
		metric := &readerMetric{
			readerMeasure: measureI.(*readerMeasure),
		}
		ro.metrics.Store(event.Sequence, metric)
		return

	case observer.ADD_EVENT:
		read.Type = ADD_EVENT
		read.Event = event.Event

		attrs, span := ro.readScope(event.Scope)
		read.Attributes = attrs.Apply(core.KeyValue{}, event.Attributes, core.Mutator{}, nil)
		if span != nil {
			read.SpanContext = span.spanContext
		}

	case observer.RECORD_STATS:
		read.Type = RECORD_STATS

		_, span := ro.readScope(event.Scope)
		if span != nil {
			read.SpanContext = span.spanContext
		}
		for _, es := range event.Stats {
			ro.addMeasurement(&read, es)
		}
		if event.Stat.Measure != nil {
			ro.addMeasurement(&read, event.Stat)
		}

	case observer.SET_STATUS:
		read.Type = SET_STATUS
		read.Status = event.Status
		_, span := ro.readScope(event.Scope)
		if span != nil {
			span.status = event.Status
			read.SpanContext = span.spanContext
		}

	default:
		panic(fmt.Sprint("Unhandled case: ", event.Type))
	}

	for _, reader := range ro.readers {
		reader.Read(read)
	}

	if event.Type == observer.FINISH_SPAN {
		ro.cleanupSpan(event.Scope.EventID)
	}
}

func (ro *readerObserver) addMeasurement(e *Event, m core.Measurement) {
	attrs, _ := ro.readScope(m.ScopeID)
	e.Stats = append(e.Stats, Measurement{
		Measure: m.Measure,
		Value:   m.Value,
		Tags:    attrs,
	})
}

func (ro *readerObserver) readScope(id core.ScopeID) (tag.Map, *readerSpan) {
	if id.EventID == 0 {
		return tag.NewEmptyMap(), nil
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
	return tag.NewEmptyMap(), nil
}

func (ro *readerObserver) cleanupSpan(id core.EventID) {
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
