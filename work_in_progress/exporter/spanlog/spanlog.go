package spanlog

import (
	"os"
	"strings"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/buffer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/spandata"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/spandata/format"
)

type (
	spanLog struct{}
)

func New() observer.Observer {
	return buffer.NewBuffer(1000, spandata.NewReaderObserver(&spanLog{}))
}

func (s *spanLog) Read(data *spandata.Span) {
	var buf strings.Builder
	buf.WriteString("----------------------------------------------------------------------\n")
	format.AppendSpan(&buf, data)
	os.Stdout.WriteString(buf.String())
}
