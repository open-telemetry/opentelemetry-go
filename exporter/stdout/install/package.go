package install

import (
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/stdout"
)

// Use this import:
//
//   import _ "github.com/open-telemetry/opentelemetry-go/exporter/stdout/install"
//
// to include the stderr exporter by default.

func init() {
	observer.RegisterObserver(stdout.New())
}
