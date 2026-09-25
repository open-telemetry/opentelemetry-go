// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/internal/baggage"
)

func TestContext(t *testing.T) {
	ctx := t.Context()
	assert.Equal(t, Baggage{}, FromContext(ctx))

	b := Baggage{list: baggage.List{"key": baggage.Item{Value: "val"}}}
	ctx = ContextWithBaggage(ctx, b)
	assert.Equal(t, b, FromContext(ctx))

	ctx = ContextWithoutBaggage(ctx)
	assert.Equal(t, Baggage{}, FromContext(ctx))
}

func TestContextWithSetHook(t *testing.T) {
	var called bool
	f := func(ctx context.Context, _ Baggage) context.Context {
		called = true
		return ctx
	}

	ctx := t.Context()
	ctx = ContextWithSetHook(ctx, f)
	assert.False(t, called, "SetHookFunc called when setting hook")
	ctx = ContextWithBaggage(ctx, Baggage{})
	assert.True(t, called, "SetHookFunc not called when setting Baggage")

	// Ensure resetting the hook works.
	called = false
	ctx = ContextWithSetHook(ctx, f)
	assert.False(t, called, "SetHookFunc called when re-setting hook")
	ContextWithBaggage(ctx, Baggage{})
	assert.True(t, called, "SetHookFunc not called when re-setting Baggage")
}

func TestContextWithGetHook(t *testing.T) {
	var called bool
	f := func(_ context.Context, bag Baggage) Baggage {
		called = true
		return bag
	}

	ctx := t.Context()
	ctx = ContextWithGetHook(ctx, f)
	assert.False(t, called, "GetHookFunc called when setting hook")
	_ = FromContext(ctx)
	assert.True(t, called, "GetHookFunc not called when getting Baggage")

	// Ensure resetting the hook works.
	called = false
	ctx = ContextWithGetHook(ctx, f)
	assert.False(t, called, "GetHookFunc called when re-setting hook")
	_ = FromContext(ctx)
	assert.True(t, called, "GetHookFunc not called when re-getting Baggage")
}

func TestContextWithGetHookReplacesBaggage(t *testing.T) {
	replacement := Baggage{list: baggage.List{"replaced": baggage.Item{Value: "yes"}}}
	f := func(_ context.Context, _ Baggage) Baggage {
		return replacement
	}

	ctx := t.Context()
	ctx = ContextWithGetHook(ctx, f)
	ctx = ContextWithBaggage(ctx, Baggage{list: baggage.List{"key": baggage.Item{Value: "val"}}})

	assert.Equal(t, replacement, FromContext(ctx))
}
