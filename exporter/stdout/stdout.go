package stdout

import (
	"os"

	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/reader"
	"github.com/open-telemetry/opentelemetry-go/exporter/reader/format"
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
