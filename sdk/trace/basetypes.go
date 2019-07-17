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
	"time"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/event"
)

// LinkType specifies the relationship between the span that had the link
// added, and the linked span.
type LinkType int32

// LinkType values.
const (
	LinkTypeUnspecified LinkType = iota // The relationship of the two spans is unknown.
	LinkTypeChild                       // The linked span is a child of the current span.
	LinkTypeParent                      // The linked span is the parent of the current span.
)

// Link represents a reference from one span to another span.
type Link struct {
	Type LinkType
	// Attributes is a set of attributes on the link.
	Attributes map[string]interface{}
}

// MessageEventType specifies the type of message event.
type MessageEventType int32

// MessageEventType values.
const (
	MessageEventTypeUnspecified MessageEventType = iota // Unknown event type.
	MessageEventTypeSent                                // Indicates a sent RPC message.
	MessageEventTypeRecv                                // Indicates a received RPC message.
)

// MessageEvent represents an event describing a message sent or received on the network.
type MessageEvent struct {
	msg                  string
	attributes           []core.KeyValue
	Time                 time.Time
	EventType            MessageEventType
	MessageID            int64
	UncompressedByteSize int64
	CompressedByteSize   int64
}

var _ event.Event = &MessageEvent{}

func (me *MessageEvent) Message() string {
	return me.msg
}

func (me *MessageEvent) Attributes() []core.KeyValue {
	return me.attributes
}
