// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opentracing

import (
	oteltrace "go.opentelemetry.io/api/trace"
)

// NewTracerPair is a utility function that creates a BridgeTracer
// that forwards the calls to the WrapperTracer that wraps the passed
// tracer.
func NewTracerPair(tracer oteltrace.Tracer) (*BridgeTracer, *WrapperTracer) {
	bridgeTracer := NewBridgeTracer()
	wrapperTracer := NewWrapperTracer(bridgeTracer, tracer)
	bridgeTracer.SetOpenTelemetryTracer(wrapperTracer)
	return bridgeTracer, wrapperTracer
}
