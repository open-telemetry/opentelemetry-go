package stdout

import (
	"os"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/reader"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/reader/format"
)

type (
	stdoutLog struct{}
)

func New() observer.Observer {
	return reader.NewReaderObserver(&stdoutLog{})
}

func (s *stdoutLog) Read(data reader.Event) {
	os.Stdout.WriteString(format.EventToString(data))
}
