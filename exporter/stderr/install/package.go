package install

import (
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/stderr"
)

// Use this import:
//
//   import _ "github.com/open-telemetry/opentelemetry-go/exporter/stderr/install"
//
// to include the stderr exporter by default.

func init() {
	observer.RegisterObserver(stderr.New())
}
