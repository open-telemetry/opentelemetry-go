// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing // import "go.opentelemetry.io/otel/bridge/opentracing"

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// NewTracerPair is a utility function that creates a BridgeTracer and a
// WrapperTracerProvider. WrapperTracerProvider creates a single instance of
// WrapperTracer. The BridgeTracer forwards the calls to the WrapperTracer
// that wraps the passed tracer. BridgeTracer and WrapperTracerProvider are
// returned to the caller and the caller is expected to register BridgeTracer
// with opentracing and WrapperTracerProvider with opentelemetry.
func NewTracerPair(tracer trace.Tracer) (*BridgeTracer, *WrapperTracerProvider) {
	bridgeTracer := NewBridgeTracer()
	wrapperProvider := NewWrappedTracerProvider(bridgeTracer, tracer)
	bridgeTracer.SetOpenTelemetryTracer(wrapperProvider.Tracer(""))
	return bridgeTracer, wrapperProvider
}

// NewTracerPairWithContext is a convenience function. It calls NewTracerPair
// and returns a hooked version of ctx with the created BridgeTracer along
// with the BridgeTracer and WrapperTracerProvider.
func NewTracerPairWithContext(ctx context.Context, tracer trace.Tracer) (context.Context, *BridgeTracer, *WrapperTracerProvider) {
	bridgeTracer, wrapperProvider := NewTracerPair(tracer)
	ctx = bridgeTracer.NewHookedContext(ctx)
	return ctx, bridgeTracer, wrapperProvider
}
