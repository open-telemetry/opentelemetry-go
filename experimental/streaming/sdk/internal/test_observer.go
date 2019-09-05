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

package internal // import "go.opentelemetry.io/experimental/streaming/sdk/internal"

import (
	"go.opentelemetry.io/experimental/streaming/exporter/observer"
)

type eventsMap map[observer.EventType][]observer.Event

type TestObserver struct {
	events eventsMap
}

var _ observer.Observer = &TestObserver{}

func NewRegisteredObserver() *TestObserver {
	o := &TestObserver{}
	observer.RegisterObserver(o)
	return o
}

func (o *TestObserver) Unregister() {
	observer.UnregisterObserver(o)
}

func (o *TestObserver) Observe(e observer.Event) {
	if o.events == nil {
		o.events = make(eventsMap)
	}
	o.events[e.Type] = append(o.events[e.Type], e)
}

func (o *TestObserver) Clear() {
	o.events = nil
}

func (o *TestObserver) ClearAndUnregister() {
	o.Clear()
	o.Unregister()
}

func (o *TestObserver) Events(eType observer.EventType) []observer.Event {
	if o.events == nil {
		return nil
	}
	return o.events[eType]
}
