package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/pipeline"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	// State is the object used to maintain independent collection
	// state for each asynchronous meter.
	State struct {
		// pipe is the pipeline.Register number of this state.
		pipe int

		// lock protects against errant use of the instrument
		// w/ copied context after the callback returns.
		lock sync.Mutex

		// store is a map from instrument to set of values
		// observed during one collection.
		store map[*Instrument]map[attribute.Set]viewstate.Accumulator
	}

	// Instrument is the implementation object associated with one
	// asynchronous instrument.
	Instrument struct {
		// opaque is used to ensure that callbacks are
		// registered with instruments from the same provider.
		opaque interface{}

		// compiled is the per-pipeline compiled instrument.
		compiled pipeline.Register[viewstate.Instrument]

		// descriptor describes the API-level instrument.
		//
		// Note: Not clear why we need this. It's used for a
		// range test, but shouldn't the range test be
		// performed by the aggregator?  If a View is allowed
		// to reconfigure the aggregation in ways that change
		// semantics, should the range test be based on the
		// aggregation, not the original instrument?
		descriptor sdkinstrument.Descriptor
	}

	// contextKey is used with context.WithValue() to lookup
	// per-reader state from within an executing callback
	// function.
	contextKey struct{}
)

func NewState(pipe int) *State {
	return &State{
		pipe:  pipe,
		store: map[*Instrument]map[attribute.Set]viewstate.Accumulator{},
	}
}

// NewInstrument returns a new Instrument; this compiles individual
// instruments for each reader.
func NewInstrument(desc sdkinstrument.Descriptor, opaque interface{}, compiled pipeline.Register[viewstate.Instrument]) *Instrument {
	return &Instrument{
		opaque:     opaque,
		descriptor: desc,
		compiled:   compiled,
	}
}

// SnapshotAndProcess calls SnapshotAndProcess() on each of the pending
// aggregations for a given reader.
func (inst *Instrument) SnapshotAndProcess(state *State) {
	state.lock.Lock()
	defer state.lock.Unlock()

	for _, acc := range state.store[inst] {
		acc.SnapshotAndProcess()
	}
}

func (inst *Instrument) get(cs *callbackState, attrs []attribute.KeyValue) viewstate.Accumulator {
	comp := inst.compiled[cs.state.pipe]

	cs.state.lock.Lock()
	defer cs.state.lock.Unlock()

	aset := attribute.NewSet(attrs...)
	imap, has := cs.state.store[inst]

	if !has {
		imap = map[attribute.Set]viewstate.Accumulator{}
		cs.state.store[inst] = imap
	}

	se, has := imap[aset]
	if !has {
		se = comp.NewAccumulator(aset)
		imap[aset] = se
	}
	return se
}

func capture[N number.Any, Traits number.Traits[N]](ctx context.Context, inst *Instrument, value N, attrs []attribute.KeyValue) {
	lookup := ctx.Value(contextKey{})
	if lookup == nil {
		otel.Handle(fmt.Errorf("async instrument used outside of callback"))
		return
	}

	cs := lookup.(*callbackState)
	cb := cs.getCallback()
	if cb == nil {
		otel.Handle(fmt.Errorf("async instrument used after callback return"))
		return
	}
	if _, ok := cb.instruments[inst]; !ok {
		otel.Handle(fmt.Errorf("async instrument not declared for use in callback"))
		return
	}

	inst.get(cs, attrs).(viewstate.Updater[N]).Update(value)
}
