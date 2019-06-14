package impl

import (
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
	"github.com/lightstep/opentelemetry-golang-prototype/exporter/stdout"
)

// Use this import:
//
//   import _ "github.com/lightstep/opentelemetry-golang-prototype/exporter/stdout/impl"
//
// to include the stderr exporter by default.

func init() {
	observer.RegisterObserver(stdout.New())
}
