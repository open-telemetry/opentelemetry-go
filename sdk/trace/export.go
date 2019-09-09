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
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
)

// Exporter is a type for functions that receive sampled trace spans.
//
// The ExportSpan method should be safe for concurrent use and should return
// quickly; if an Exporter takes a significant amount of time to process a
// SpanData, that work should be done on another goroutine.
//
// The SpanData should not be modified, but a pointer to it can be kept.
type Exporter interface {
	ExportSpan(s *SpanData)
}

type exportersMap map[Exporter]struct{}

var (
	exporterMu sync.Mutex
	exporters  atomic.Value
)

// RegisterExporter adds to the list of Exporters that will receive sampled
// trace spans.
//
// Binaries can register exporters, libraries shouldn't register exporters.
func RegisterExporter(e Exporter) {
	exporterMu.Lock()
	defer exporterMu.Unlock()
	new := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	new[e] = struct{}{}
	exporters.Store(new)
}

// UnregisterExporter removes from the list of Exporters the Exporter that was
// registered with the given name.
func UnregisterExporter(e Exporter) {
	exporterMu.Lock()
	defer exporterMu.Unlock()
	new := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	delete(new, e)
	exporters.Store(new)
}

// SpanData contains all the information collected by a span.
type SpanData struct {
	SpanContext  core.SpanContext
	ParentSpanID uint64
	SpanKind     int
	Name         string
	StartTime    time.Time
	// The wall clock time of EndTime will be adjusted to always be offset
	// from StartTime by the duration of the span.
	EndTime time.Time
	// The values of Attributes each have type string, bool, or int64.
	Attributes               map[string]interface{}
	MessageEvents            []Event
	Status                   codes.Code
	HasRemoteParent          bool
	DroppedAttributeCount    int
	DroppedMessageEventCount int
	DroppedLinkCount         int

	// ChildSpanCount holds the number of child span created for this span.
	ChildSpanCount int
}

// Event is used to describe an Event with a message string and set of
// Attributes.
type Event struct {
	// Message describes the Event.
	Message string

	// Attributes contains a list of keyvalue pairs.
	Attributes []core.KeyValue

	// Time is the time at which this event was recorded.
	Time time.Time
}
