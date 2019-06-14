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

	"github.com/open-telemetry/opentelemetry-go/api/core"
)

type (
	EventType int

	Event struct {
		// Automatic fields
		Sequence core.EventID // Auto-filled
		Time     time.Time    // Auto-filled

		// Type, Scope, Context
		Type    EventType       // All events
		Scope   core.ScopeID    // All events
		Context context.Context // core.FromContext() and scope.Active()

		// Arguments (type-specific)
		Attribute  core.KeyValue   // SET_ATTRIBUTE
		Attributes []core.KeyValue // SET_ATTRIBUTES, LOG_EVENT
		Mutator    core.Mutator    // SET_ATTRIBUTE
		Mutators   []core.Mutator  // SET_ATTRIBUTES
		Arguments  []interface{}   // LOGF_EVENT
		Recovered  interface{}     // FINISH_SPAN

		// Values
		String  string // START_SPAN, EVENT, ...
		Float64 float64
		Parent  core.ScopeID // START_SPAN
		Stats   []core.Measurement
		Stat    core.Measurement
	}

	Observer interface {
		Observe(data Event)
	}

	observersMap map[Observer]struct{}
)

//go:generate stringer -type=EventType
const (
	// TODO: rename these NOUN_VERB
	INVALID EventType = iota
	START_SPAN
	FINISH_SPAN
	LOG_EVENT
	LOGF_EVENT
	NEW_SCOPE
	NEW_MEASURE
	NEW_METRIC
	MODIFY_ATTR
	RECORD_STATS
)

var (
	observerMu sync.Mutex
	observers  atomic.Value

	sequenceNum uint64
)

func NextEventID() core.EventID {
	return core.EventID(atomic.AddUint64(&sequenceNum, 1))
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

func Record(event Event) core.EventID {
	if event.Sequence == 0 {
		event.Sequence = NextEventID()
	}
	if event.Time.IsZero() {
		event.Time = time.Now()
	}

	observers, _ := observers.Load().(observersMap)
	for observer, _ := range observers {
		observer.Observe(event)
	}
	return event.Sequence
}

func Foreach(f func(Observer)) {
	observers, _ := observers.Load().(observersMap)
	for observer, _ := range observers {
		f(observer)
	}
}
