package spandata

import (
	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/reader"
)

type (
	Reader interface {
		Read(*Span)
	}

	Span struct {
		Events []reader.Event
	}

	spanReader struct {
		spans   map[core.SpanContext]*Span
		readers []Reader
	}
)

func NewReaderObserver(readers ...Reader) observer.Observer {
	return reader.NewReaderObserver(&spanReader{
		spans:   map[core.SpanContext]*Span{},
		readers: readers,
	})
}

func (s *spanReader) Read(data reader.Event) {
	if !data.SpanContext.HasSpanID() {
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
