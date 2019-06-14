package impl

import (
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr"
)

// Use this import:
//
//   import _ "github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr/impl"
//
// to include the stderr exporter by default.

func init() {
	observer.RegisterObserver(stderr.New())
}
