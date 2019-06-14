package format

import (
	"fmt"
	"strings"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/exporter/reader"
)

func AppendEvent(buf *strings.Builder, data reader.Event) {

	f := func(skipIf bool) func(kv core.KeyValue) bool {
		return func(kv core.KeyValue) bool {
			if skipIf && data.Attributes.HasValue(kv.Key) {
				return true
			}
			buf.WriteString(" " + kv.Key.Name() + "=" + kv.Value.Emit())
			return true
		}
	}

	buf.WriteString(data.Time.Format("2006/01/02 15-04-05.000000"))
	buf.WriteString(" ")

	switch data.Type {
	case reader.START_SPAN:
		buf.WriteString("start ")
		buf.WriteString(data.Name)

		if !data.Parent.HasSpanID() {
			buf.WriteString(", a root span")
		} else {
			buf.WriteString(" <")
			if data.Parent.HasSpanID() {
				f(false)(trace.ParentSpanIDKey.String(data.SpanContext.SpanIDString()))
			}
			if data.ParentAttributes != nil {
				data.ParentAttributes.Foreach(f(false))
			}
			buf.WriteString(" >")
		}

	case reader.FINISH_SPAN:
		buf.WriteString("finish ")
		buf.WriteString(data.Name)

		buf.WriteString(" (")
		buf.WriteString(data.Duration.String())
		buf.WriteString(")")

	case reader.LOG_EVENT:
		buf.WriteString(data.Message)

	case reader.LOGF_EVENT:
		buf.WriteString(data.Message)

	case reader.MODIFY_ATTR:
		buf.WriteString("modify attr")
	case reader.RECORD_STATS:
		buf.WriteString("record")

		for _, s := range data.Stats {
			f(false)(s.Measure.V(s.Value))

			buf.WriteString(" {")
			i := 0
			s.Tags.Foreach(func(kv core.KeyValue) bool {
				if i != 0 {
					buf.WriteString(",")
				}
				i++
				buf.WriteString(kv.Key.Name())
				buf.WriteString("=")
				buf.WriteString(kv.Value.Emit())
				return true
			})
			buf.WriteString("}")
		}
	default:
		buf.WriteString(fmt.Sprintf("WAT? %d", data.Type))
	}

	// Attach the scope (span) attributes and context tags.
	buf.WriteString(" [")
	if data.Attributes != nil {
		data.Attributes.Foreach(f(false))
	}
	if data.Tags != nil {
		data.Tags.Foreach(f(true))
	}
	if data.SpanContext.HasSpanID() {
		f(false)(trace.SpanIDKey.String(data.SpanContext.SpanIDString()))
	}
	if data.SpanContext.HasTraceID() {
		f(false)(trace.TraceIDKey.String(data.SpanContext.TraceIDString()))
	}

	buf.WriteString(" ]\n")
}

func EventToString(data reader.Event) string {
	var buf strings.Builder
	AppendEvent(&buf, data)
	return buf.String()
}
