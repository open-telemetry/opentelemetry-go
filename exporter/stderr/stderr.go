package stderr

import (
	"os"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/reader"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/reader/format"
)

type (
	stderrLog struct{}
)

func New() observer.Observer {
	return reader.NewReaderObserver(&stderrLog{})
}

func (s *stderrLog) Read(data reader.Event) {
	os.Stderr.WriteString(format.EventToString(data))
}
