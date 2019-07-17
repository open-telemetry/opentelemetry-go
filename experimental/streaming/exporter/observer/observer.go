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

package observer

import (
	"context"
	"sync"
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

// TODO: this Event is confusing with event.Event.
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
	Recovered  interface{}     // FINISH_SPAN
	Status     codes.Code      // SET_STATUS

	// Values
	String  string // START_SPAN, EVENT, ...
	Float64 float64
	Parent  ScopeID // START_SPAN
	Stats   []stats.Measurement
	Stat    stats.Measurement
}

type Observer interface {
	Observe(data Event)
}

type observersMap map[Observer]struct{}

//go:generate stringer -type=EventType
const (
	// TODO: rename these NOUN_VERB
	INVALID EventType = iota
	START_SPAN
	FINISH_SPAN
	ADD_EVENT
	ADD_EVENTF
	NEW_SCOPE
	NEW_MEASURE
	NEW_METRIC
	MODIFY_ATTR
	RECORD_STATS
	SET_STATUS
)

var (
	observerMu sync.Mutex
	observers  atomic.Value

	sequenceNum uint64
)

func NextEventID() EventID {
	return EventID(atomic.AddUint64(&sequenceNum, 1))
}

// RegisterObserver adds to the list of Observers that will receive sampled
// trace spans.
//
// Binaries can register observers, libraries shouldn't register observers.
func RegisterObserver(e Observer) {
	observerMu.Lock()
	new := make(observersMap)
	if old, ok := observers.Load().(observersMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	new[e] = struct{}{}
	observers.Store(new)
	observerMu.Unlock()
}

// UnregisterObserver removes from the list of Observers the Observer that was
// registered with the given name.
func UnregisterObserver(e Observer) {
	observerMu.Lock()
	new := make(observersMap)
	if old, ok := observers.Load().(observersMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	delete(new, e)
	observers.Store(new)
	observerMu.Unlock()
}

func Record(event Event) EventID {
	if event.Sequence == 0 {
		event.Sequence = NextEventID()
	}
	if event.Time.IsZero() {
		event.Time = time.Now()
	}

	observers, _ := observers.Load().(observersMap)
	for observer := range observers {
		observer.Observe(event)
	}
	return event.Sequence
}

func Foreach(f func(Observer)) {
	observers, _ := observers.Load().(observersMap)
	for observer := range observers {
		f(observer)
	}
}

func NewScope(parent ScopeID, attributes ...core.KeyValue) ScopeID {
	eventID := Record(Event{
		Type:       NEW_SCOPE,
		Scope:      parent,
		Attributes: attributes,
	})
	return ScopeID{
		EventID:     eventID,
		SpanContext: parent.SpanContext,
	}
}
