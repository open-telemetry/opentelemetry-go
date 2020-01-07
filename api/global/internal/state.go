package internal

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/internal"
)

type (
	scopeHolder struct {
		scope.Scope
	}
)

func init() {
	ResetForTest()
}

// Scope is the internal implementation for global.Scope().
func Scope() scope.Scope {
	return internal.GlobalScope.Load().(scopeHolder).Scope
}

// SetScope is the internal implementation for global.SetScope().
func SetScope(sc scope.Scope) {
	first := false
	internal.GlobalDelegateOnce.Do(func() {
		current := Scope()
		currentProvider := current.Provider()
		newProvider := sc.Provider()

		first = true

		if currentProvider.Meter() == newProvider.Meter() {
			// Setting the global scope to former default is nonsense, panic.
			// Panic is acceptable because we are likely still early in the
			// process lifetime.
			panic("invalid Provider, the global instance cannot be reinstalled")
		} else if deft, ok := currentProvider.Tracer().(*tracer); ok {
			deft.deferred.setDelegate(sc)
		} else {
			panic("impossible error")
		}
	})
	if !first {
		panic("global scope has already been initialized")
	}
	internal.GlobalScope.Store(scopeHolder{Scope: sc})
}

func defaultScopeValue() *atomic.Value {
	v := &atomic.Value{}
	d := newDeferred()
	v.Store(scopeHolder{
		Scope: scope.NewProvider(&d.tracer, &d.meter, &d.propagators).New(),
	})
	return v
}

// ResetForTest restores the initial global state, for testing purposes.
func ResetForTest() {
	internal.GlobalScope = defaultScopeValue()
	internal.GlobalDelegateOnce = sync.Once{}
}
