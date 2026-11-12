// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// ContextFallbackMode controls how incoming span context is used for otherwise
// parentless spans.
type ContextFallbackMode int

const (
	// ContextFallbackParent uses the incoming span context as a parent.
	// It is the default mode of [WithContextFallback].
	ContextFallbackParent ContextFallbackMode = iota
	// ContextFallbackLink adds the incoming span context as a link, preserving
	// the new span's independent trace.
	ContextFallbackLink
	// ContextFallbackNone disables incoming span context and baggage fallback.
	ContextFallbackNone
)

type contextFallbackConfig struct {
	mode ContextFallbackMode
}

// ContextFallbackOption configures [WithContextFallback].
type ContextFallbackOption interface {
	apply(contextFallbackConfig) contextFallbackConfig
}

type contextFallbackModeOption ContextFallbackMode

func (o contextFallbackModeOption) apply(c contextFallbackConfig) contextFallbackConfig {
	c.mode = ContextFallbackMode(o)
	return c
}

// WithContextFallbackMode sets the incoming span context mode. The default is
// [ContextFallbackParent]. An invalid mode disables fallback and reports an error
// to the OpenTelemetry error handler when [WithContextFallback] is called.
func WithContextFallbackMode(mode ContextFallbackMode) ContextFallbackOption {
	return contextFallbackModeOption(mode)
}

type contextFallbackOption struct {
	sdktrace.TracerProviderOption
	spanContext trace.SpanContext
	baggage     baggage.Baggage
	mode        ContextFallbackMode
}

// Experimental identifies this option to SDK versions that support experimental
// options. Versions that do not support context fallback ignore this option.
func (contextFallbackOption) Experimental() {}

// ContextFallback returns the immutable snapshot consumed by the SDK.
func (o contextFallbackOption) ContextFallback() (trace.SpanContext, baggage.Baggage, int) {
	return o.spanContext, o.baggage, int(o.mode)
}

// WithContextFallback configures incoming span context and baggage defaults for a
// tracer provider. It extracts once from carrier using propagator when this
// function is called. Later carrier changes are ignored. Carrier key
// normalization and lookup behavior are defined by the supplied carrier.
//
// Fallback applies only when the starting context contains an empty span context
// and [trace.WithNewRoot] is absent. Existing span contexts, including partial
// ones, take precedence. The default [ContextFallbackParent] mode uses valid
// incoming span context as the parent; [ContextFallbackLink] adds it as a link
// before sampling. Both modes merge incoming baggage with application baggage,
// preserving application values. Descendants inherit baggage normally without
// reapplying defaults.
//
// This is a persistent, provider-wide fallback: independent parentless operations
// also receive it, including later requests without trace headers. Applications
// need to restrict the provider to work belonging to the incoming invocation,
// use [trace.WithNewRoot] for independent operations, or extract into an explicit
// context and pass that context instead of configuring fallback.
//
// Without this option, fallback is disabled. If multiple options are supplied to
// a provider, the last one wins. [ContextFallbackNone] skips extraction entirely.
// In enabled modes, a nil propagator or carrier disables fallback and reports an
// error to the OpenTelemetry error handler. An invalid mode also disables
// fallback and reports an error. Only extracted span context and
// baggage are retained; no active span is created or process-wide context set.
// Extracted baggage exceeding W3C limits is reduced to a valid partial result,
// and the error is reported to the OpenTelemetry error handler.
func WithContextFallback(
	propagator propagation.TextMapPropagator,
	carrier propagation.TextMapCarrier,
	opts ...ContextFallbackOption,
) sdktrace.TracerProviderOption {
	config := contextFallbackConfig{mode: ContextFallbackParent}
	for _, opt := range opts {
		config = opt.apply(config)
	}

	o := contextFallbackOption{mode: config.mode}
	switch config.mode {
	case ContextFallbackNone:
		return o
	case ContextFallbackParent, ContextFallbackLink:
	default:
		otel.Handle(fmt.Errorf("invalid context fallback mode: %d", config.mode))
		o.mode = ContextFallbackNone
		return o
	}
	if propagator == nil || carrier == nil {
		otel.Handle(errors.New("context fallback requires a non-nil propagator and carrier"))
		o.mode = ContextFallbackNone
		return o
	}

	extracted := propagator.Extract(context.Background(), carrier)
	if extracted == nil {
		otel.Handle(errors.New("context fallback propagator returned a nil context"))
		o.mode = ContextFallbackNone
		return o
	}
	o.spanContext = trace.SpanContextFromContext(extracted)
	var err error
	o.baggage, err = baggage.New(baggage.FromContext(extracted).Members()...)
	if err != nil {
		otel.Handle(fmt.Errorf("invalid context fallback baggage: %w", err))
	}
	return o
}
