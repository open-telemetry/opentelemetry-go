package sdk

import (
	"go.opentelemetry.io/api/trace"
)

func init() {
	trace.SetGlobalTracer(New())
}
