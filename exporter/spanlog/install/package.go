package install

import (
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/spanlog"
)

// Use this import:
//
//   import _ "github.com/lightstep/opentelemetry-golang-prototype/exporter/spanlog/install"
//
// to include the spanlog exporter by default.

func init() {
	observer.RegisterObserver(spanlog.New())
}
