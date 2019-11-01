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
	oteltrace "go.opentelemetry.io/otel/api/trace"
)

// NewTracerPair is a utility function that creates a BridgeTracer
// and a WrapperProvider. WrapperProvider creates a single instance of
// WrapperTracer. The BridgeTracer forwards the calls to the WrapperTracer
// that wraps the passed tracer. BridgeTracer and WrapperProvider are returned to
// the caller and the caller is expected to register BridgeTracer with opentracing and
// WrapperProvider with opentelemetry.
func NewTracerPair(tracer oteltrace.Tracer) (*BridgeTracer, *WrapperProvider) {
	bridgeTracer := NewBridgeTracer()
	wrapperProvider := NewWrappedProvider(bridgeTracer, tracer)
	bridgeTracer.SetOpenTelemetryTracer(wrapperProvider.GetTracer(""))
	return bridgeTracer, wrapperProvider
}
