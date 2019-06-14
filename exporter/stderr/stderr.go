package stderr

import (
	"os"

	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/reader"
	"github.com/open-telemetry/opentelemetry-go/exporter/reader/format"
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
