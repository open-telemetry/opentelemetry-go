package spanlog

import (
	"os"
	"strings"

	"github.com/open-telemetry/opentelemetry-go/exporter/buffer"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/spandata"
	"github.com/open-telemetry/opentelemetry-go/exporter/spandata/format"
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
