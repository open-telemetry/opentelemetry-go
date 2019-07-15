package sdk

import (
	"github.com/open-telemetry/opentelemetry-go/api/trace"
)

func init() {
	trace.SetGlobalTracer(New())
}
