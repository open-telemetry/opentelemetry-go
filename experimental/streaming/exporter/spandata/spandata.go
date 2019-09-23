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

package spandata

import (
	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/experimental/streaming/exporter"
	"go.opentelemetry.io/experimental/streaming/exporter/reader"
)

type Reader interface {
	Read(*Span)
}

type Span struct {
	Events []reader.Event
}

type spanReader struct {
	spans   map[core.SpanContext]*Span
	readers []Reader
}

func NewReaderObserver(readers ...Reader) exporter.Observer {
	return reader.NewReaderObserver(&spanReader{
		spans:   map[core.SpanContext]*Span{},
		readers: readers,
	})
}

func (s *spanReader) Read(data reader.Event) {
	if !data.SpanContext.HasSpanID() {
		// @@@ This is happening, somehow span context is busted.
		return
	}
	var span *Span
	if data.Type == reader.START_SPAN {
		span = &Span{Events: make([]reader.Event, 0, 4)}
		s.spans[data.SpanContext] = span
	} else {
		span = s.spans[data.SpanContext]
		if span == nil {
			// TODO count and report this.
			return
		}
	}

	span.Events = append(span.Events, data)

	if data.Type == reader.FINISH_SPAN {
		for _, r := range s.readers {
			r.Read(span)
		}
		delete(s.spans, data.SpanContext)
	}
}
