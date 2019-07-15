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

package event

import (
	"fmt"

	"go.opentelemetry.io/api/core"
	apievent "go.opentelemetry.io/api/event"
)

type event struct {
	message    string
	attributes []core.KeyValue
}

var _ apievent.Event = (*event)(nil)

// WithAttr creates an Event with Attributes and a message.
// Attributes are immutable.
func WithAttr(msg string, attributes ...core.KeyValue) apievent.Event {
	return event{message: msg, attributes: attributes}
}

// WithString creates an Event with formatted string.
func WithString(f string, args ...interface{}) apievent.Event {
	return event{message: fmt.Sprint(f, args), attributes: nil}
}

func (e event) Message() string {
	return e.message
}

func (e event) Attributes() []core.KeyValue {
	return append(e.attributes[:0:0], e.attributes...)
}
