package sdk

import (
	apitrace "github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/experimental/streaming/sdk/trace"
)

func init() {
	apitrace.SetGlobalTracer(trace.New())
}
