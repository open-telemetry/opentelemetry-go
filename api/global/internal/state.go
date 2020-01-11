package internal

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/internal"
)

func init() {
	ResetForTest()
}

// Scope is the internal implementation for global.Scope().
func Scope() scope.Scope {
	if sc, ok := (*atomic.Value)(atomic.LoadPointer(&internal.GlobalScope)).Load().(scope.Scope); ok {
		return sc
	}
	return scope.Scope{}
}

// SetScope is the internal implementation for global.SetScope().
func SetScope(sc scope.Scope) {
	first := false
	(*sync.Once)(atomic.LoadPointer(&internal.GlobalDelegateOnce)).Do(func() {
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
	(*atomic.Value)(atomic.LoadPointer(&internal.GlobalScope)).Store(sc)
}

func defaultScopeValue() *atomic.Value {
	v := &atomic.Value{}
	d := newDeferred()
	v.Store(scope.NewProvider(&d.tracer, &d.meter).New())
	return v
}

// ResetForTest restores the initial global state, for testing purposes.
func ResetForTest() {
	atomic.StorePointer((*unsafe.Pointer)(&internal.GlobalScope), unsafe.Pointer(defaultScopeValue()))
	atomic.StorePointer((*unsafe.Pointer)(&internal.GlobalDelegateOnce), unsafe.Pointer(&sync.Once{}))
}
