package correlation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type correlationsType struct{}

// SetHookFunc describes a type of a callback that is called when
// storing baggage in the context.
type SetHookFunc func(context.Context) context.Context

// GetHookFunc describes a type of a callback that is called when
// getting baggage from the context.
type GetHookFunc func(context.Context, Map) Map

// value under this key is either of type Map or correlationsData
var correlationsKey = &correlationsType{}

type correlationsData struct {
	m       Map
	setHook SetHookFunc
	getHook GetHookFunc
}

func (d correlationsData) isHookless() bool {
	return d.setHook == nil && d.getHook == nil
}

type hookKind int

const (
	hookKindSet hookKind = iota
	hookKindGet
)

func (d *correlationsData) overrideHook(kind hookKind, setHook SetHookFunc, getHook GetHookFunc) {
	switch kind {
	case hookKindSet:
		d.setHook = setHook
	case hookKindGet:
		d.getHook = getHook
	}
}

// WithSetHook installs a hook function that will be invoked every
// time WithMap is called. To avoid unnecessary callback invocations
// (recursive or not), the callback can temporarily clear the hooks
// from the context with the WithNoHooks function.
//
// Note that NewContext also calls WithMap, so the hook will be
// invoked.
//
// Passing nil HookFunc creates a context with no set hook to call.
//
// This function should not be used by applications or libraries. It
// is mostly for interoperation with other observability APIs.
func WithSetHook(ctx context.Context, hook SetHookFunc) context.Context {
	return withHook(ctx, hookKindSet, hook, nil)
}

// WithGetHook installs a hook function that will be invoked every
// time FromContext is called. To avoid unnecessary callback
// invocations (recursive or not), the callback can temporarily clear
// the hooks from the context with the WithNoHooks function.
//
// Note that NewContext also calls FromContext, so the hook will be
// invoked.
//
// Passing nil HookFunc creates a context with no get hook to call.
//
// This function should not be used by applications or libraries. It
// is mostly for interoperation with other observability APIs.
func WithGetHook(ctx context.Context, hook GetHookFunc) context.Context {
	return withHook(ctx, hookKindGet, nil, hook)
}

func withHook(ctx context.Context, kind hookKind, setHook SetHookFunc, getHook GetHookFunc) context.Context {
	switch v := ctx.Value(correlationsKey).(type) {
	case correlationsData:
		v.overrideHook(kind, setHook, getHook)
		if v.isHookless() {
			return context.WithValue(ctx, correlationsKey, v.m)
		}
		return context.WithValue(ctx, correlationsKey, v)
	case Map:
		return withOneHookAndMap(ctx, kind, setHook, getHook, v)
	default:
		m := NewEmptyMap()
		return withOneHookAndMap(ctx, kind, setHook, getHook, m)
	}
}

func withOneHookAndMap(ctx context.Context, kind hookKind, setHook SetHookFunc, getHook GetHookFunc, m Map) context.Context {
	d := correlationsData{m: m}
	d.overrideHook(kind, setHook, getHook)
	if d.isHookless() {
		return ctx
	}
	return context.WithValue(ctx, correlationsKey, d)
}

// WithNoHooks creates a context with all the hooks disabled. Also
// returns old set and get hooks. This function can be used to
// temporarily clear the context from hooks and then reinstate them by
// calling WithSetHook and WithGetHook functions passing the hooks
// returned by this function.
func WithNoHooks(ctx context.Context) (context.Context, SetHookFunc, GetHookFunc) {
	switch v := ctx.Value(correlationsKey).(type) {
	case correlationsData:
		return context.WithValue(ctx, correlationsKey, v.m), v.setHook, v.getHook
	default:
		return ctx, nil, nil
	}
}

// WithMap enters a Map into a new Context.
func WithMap(ctx context.Context, m Map) context.Context {
	switch v := ctx.Value(correlationsKey).(type) {
	case correlationsData:
		v.m = m
		ctx = context.WithValue(ctx, correlationsKey, v)
		if v.setHook != nil {
			ctx = v.setHook(ctx)
		}
		return ctx
	default:
		return context.WithValue(ctx, correlationsKey, m)
	}
}

// WithMap enters a key:value set into a new Context.
func NewContext(ctx context.Context, keyvalues ...core.KeyValue) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(MapUpdate{
		MultiKV: keyvalues,
	}))
}

// FromContext gets the current Map from a Context.
func FromContext(ctx context.Context) Map {
	switch v := ctx.Value(correlationsKey).(type) {
	case correlationsData:
		if v.getHook != nil {
			return v.getHook(ctx, v.m)
		}
		return v.m
	case Map:
		return v
	default:
		return NewEmptyMap()
	}
}
