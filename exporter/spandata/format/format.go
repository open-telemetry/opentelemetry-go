package format

import (
	"strings"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/reader/format"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/spandata"
)

func AppendSpan(buf *strings.Builder, data *spandata.Span) {
	for _, event := range data.Events {
		format.AppendEvent(buf, event)
	}
}

func SpanToString(data *spandata.Span) string {
	var buf strings.Builder
	AppendSpan(&buf, data)
	return buf.String()
}
