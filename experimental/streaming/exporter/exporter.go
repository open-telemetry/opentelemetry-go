package exporter

import (
	"context"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/stats"
	"go.opentelemetry.io/api/tag"
)

type EventType int

type EventID uint64

type ScopeID struct {
	EventID
	core.SpanContext
}

type Event struct {
	// Automatic fields
	Sequence EventID   // Auto-filled
	Time     time.Time // Auto-filled

	// Type, Scope, Context
	Type    EventType       // All events
	Scope   ScopeID         // All events
	Context context.Context // core.FromContext() and scope.Active()

	// Arguments (type-specific)
	Attribute  core.KeyValue   // SET_ATTRIBUTE
	Attributes []core.KeyValue // SET_ATTRIBUTES
	Mutator    tag.Mutator     // SET_ATTRIBUTE
	Mutators   []tag.Mutator   // SET_ATTRIBUTES
	Recovered  interface{}     // END_SPAN
	Status     codes.Code      // SET_STATUS

	// Values
	String  string // START_SPAN, EVENT, SET_NAME, ...
	Float64 float64
	Parent  ScopeID // START_SPAN
	Stats   []stats.Measurement
	Stat    stats.Measurement
}

type Observer interface {
	Observe(data Event)
}

//go:generate stringer -type=EventType
const (
	INVALID EventType = iota
	START_SPAN
	END_SPAN
	ADD_EVENT
	NEW_SCOPE
	NEW_MEASURE
	NEW_METRIC
	MODIFY_ATTR
	RECORD_STATS
	SET_STATUS
	SET_NAME
)

type Exporter struct {
	sequence  uint64
	observers []Observer
}

func NewExporter(observers ...Observer) *Exporter {
	return &Exporter{
		observers: observers,
	}
}

func (e *Exporter) NextEventID() EventID {
	return EventID(atomic.AddUint64(&e.sequence, 1))
}

func (e *Exporter) Record(event Event) EventID {
	if event.Sequence == 0 {
		event.Sequence = e.NextEventID()
	}
	if event.Time.IsZero() {
		event.Time = time.Now()
	}
	for _, observer := range e.observers {
		observer.Observe(event)
	}
	return event.Sequence
}

func (e *Exporter) Foreach(f func(Observer)) {
	for _, observer := range e.observers {
		f(observer)
	}
}

func (e *Exporter) NewScope(parent ScopeID, attributes ...core.KeyValue) ScopeID {
	eventID := e.Record(Event{
		Type:       NEW_SCOPE,
		Scope:      parent,
		Attributes: attributes,
	})
	return ScopeID{
		EventID:     eventID,
		SpanContext: parent.SpanContext,
	}
}
