package install

import (
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/spanlog"
)

// Use this import:
//
//   import _ "github.com/open-telemetry/opentelemetry-go/exporter/spanlog/install"
//
// to include the spanlog exporter by default.

func init() {
	observer.RegisterObserver(spanlog.New())
}
