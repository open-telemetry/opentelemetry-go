// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage

import "context"

type baggageContextKeyType int

const baggageKey baggageContextKeyType = iota

// SetHookFunc is a callback called when storing baggage in a context.
type SetHookFunc func(context.Context, Baggage) context.Context

// GetHookFunc is a callback called when getting baggage from a context.
type GetHookFunc func(context.Context, Baggage) Baggage

type baggageState struct {
	bag Baggage

	setHook SetHookFunc
	getHook GetHookFunc
}

// ContextWithSetHook returns a copy of parent with hook configured to be
// invoked every time ContextWithBaggage is called.
//
// This is meant for use by bridges between an OpenTelemetry Baggage and
// another baggage-like representation, such as the OpenTracing bridge, that
// need to be notified whenever baggage stored in a context changes so they
// can keep their own representation synchronized.
//
// Passing nil SetHookFunc creates a context with no set hook to call.
func ContextWithSetHook(parent context.Context, hook SetHookFunc) context.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.setHook = hook
	return context.WithValue(parent, baggageKey, s)
}

// ContextWithGetHook returns a copy of parent with hook configured to be
// invoked every time FromContext is called.
//
// This is meant for use by bridges between an OpenTelemetry Baggage and
// another baggage-like representation, such as the OpenTracing bridge, that
// need to contribute additional baggage members whenever baggage is read
// from a context.
//
// Passing nil GetHookFunc creates a context with no get hook to call.
func ContextWithGetHook(parent context.Context, hook GetHookFunc) context.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.getHook = hook
	return context.WithValue(parent, baggageKey, s)
}

// ContextWithBaggage returns a copy of parent with baggage.
func ContextWithBaggage(parent context.Context, b Baggage) context.Context {
	var s baggageState
	if v, ok := parent.Value(baggageKey).(baggageState); ok {
		s = v
	}

	s.bag = b
	ctx := context.WithValue(parent, baggageKey, s)
	if s.setHook != nil {
		ctx = s.setHook(ctx, b)
	}

	return ctx
}

// ContextWithoutBaggage returns a copy of parent with no baggage.
func ContextWithoutBaggage(parent context.Context) context.Context {
	return ContextWithBaggage(parent, Baggage{})
}

// FromContext returns the baggage contained in ctx.
func FromContext(ctx context.Context) Baggage {
	switch v := ctx.Value(baggageKey).(type) {
	case baggageState:
		if v.getHook != nil {
			return v.getHook(ctx, v.bag)
		}
		return v.bag
	default:
		return Baggage{}
	}
}
