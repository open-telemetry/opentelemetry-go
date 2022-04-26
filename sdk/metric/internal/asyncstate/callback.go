package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric/instrument"
)

// Callback is the implementation object associated with one
// asynchronous callback.
type Callback struct {
	// function is the user-provided callback function.
	function func(context.Context)

	// instruments are the instruments permitted to be
	// used inside this callback.
	instruments map[*Instrument]struct{}
}

// NewCallback returns a new Callback; this checks that each of the
// provided instruments belongs to the same meter provider.
func NewCallback(instruments []instrument.Asynchronous, opaque interface{}, function func(context.Context)) (*Callback, error) {
	if len(instruments) == 0 {
		return nil, fmt.Errorf("asynchronous callback without instruments")
	}
	if function == nil {
		return nil, fmt.Errorf("asynchronous callback with nil function")
	}

	cb := &Callback{
		function:    function,
		instruments: map[*Instrument]struct{}{},
	}

	for _, inst := range instruments {
		ai, ok := inst.(memberInstrument)
		if !ok {
			return nil, fmt.Errorf("asynchronous instrument does not belong to this SDK: %T", inst)
		}
		thisInst := ai.instrument()
		if thisInst.opaque != opaque {
			return nil, fmt.Errorf("asynchronous instrument belongs to a different meter")
		}

		cb.instruments[thisInst] = struct{}{}
	}

	return cb, nil
}

// Run executes the callback after setting up the appropriate context
// for a specific reader.
func (c *Callback) Run(ctx context.Context, state *State) {
	cp := &callbackState{
		callback: c,
		state:    state,
	}
	c.function(context.WithValue(ctx, contextKey{}, cp))
	cp.invalidate()
}

// callbackState is used to lookup the current callback and
// pipeline from within an executing callback function.
type callbackState struct {
	// lock protects callback, see invalidate() and getCallback()
	lock sync.Mutex

	// callback is the currently running callback; this is set to nil
	// after the associated callback function returns.
	callback *Callback

	// state is a single collection of data.
	state *State
}

func (cp *callbackState) invalidate() {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	cp.callback = nil
}

func (cp *callbackState) getCallback() *Callback {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	return cp.callback
}
