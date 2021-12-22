package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
)

type (
	Accumulator struct {
		callbacksLock sync.Mutex
		callbacks     []*callback
	}

	stateEntry struct {
		value     number.Number
		collector viewstate.Collector
	}

	State struct {
		collectLock sync.Mutex
		tmpSort     attribute.Sortable
		store       map[*instrument]map[attribute.Set]*stateEntry
	}

	instrument struct {
		descriptor sdkapi.Descriptor
		cfactory   viewstate.CollectorFactory
		callback   *callback
	}

	callback struct {
		function    func(context.Context) error
		instruments []sdkapi.Instrument
	}

	contextKey struct{}
)

var (
	_ sdkapi.Instrument = &instrument{}
)

func (cb *callback) Instruments() []sdkapi.Instrument {
	return cb.instruments
}

func New() *Accumulator {
	return &Accumulator{}
}

func NewState() *State {
	return &State{
		store: map[*instrument]map[attribute.Set]*stateEntry{},
	}
}

func (m *Accumulator) NewInstrument(descriptor sdkapi.Descriptor, cfactory viewstate.CollectorFactory) (sdkapi.Instrument, error) {
	return &instrument{
		cfactory:   cfactory,
		descriptor: descriptor,
	}, nil
}

func (m *Accumulator) NewCallback(instruments []sdkapi.Instrument, function func(context.Context) error) (sdkapi.Callback, error) {
	cb := &callback{
		function:    function,
		instruments: instruments,
	}

	m.callbacksLock.Lock()
	defer m.callbacksLock.Unlock()

	for _, inst := range instruments {
		if inst.Descriptor().InstrumentKind().Synchronous() {
			return nil, fmt.Errorf("synchronous instrument with asynchronous callback")
		}
		ai, ok := inst.Implementation().(*instrument)
		if !ok {
			return nil, fmt.Errorf("asynchronous instrument does not belong to this provider")
		}
		if ai.callback != nil {
			return nil, fmt.Errorf("asynchronous instrument already has a callback")
		}
		ai.callback = cb

	}

	m.callbacks = append(m.callbacks, cb)
	return cb, nil
}

func (a *Accumulator) getCallbacks() []*callback {
	a.callbacksLock.Lock()
	defer a.callbacksLock.Unlock()
	return a.callbacks
}

func (a *Accumulator) Collect(state *State) error {
	state.collectLock.Lock()
	defer state.collectLock.Unlock()

	ctx := context.WithValue(
		context.Background(),
		contextKey{},
		state,
	)

	// TODO: Add a timeout to the context.

	for _, cb := range a.getCallbacks() {
		// TODO: handle errors.
		_ = cb.function(ctx)
	}

	for inst, states := range state.store {
		for _, entry := range states {
			// Note: not deleting anything is a safe approach for
			// async view calculation.  This can be improved with
			// optional per-stream start timestamps and staleness
			// markers, but is complicated by views correctness
			// considerations, see text in the supplemental
			// guidelines.
			entry.collector.Update(entry.value, &inst.descriptor)
			entry.collector.Send(false)
		}
	}

	return nil
}

func (inst *instrument) Implementation() interface{} {
	return inst
}

func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (inst *instrument) Capture(ctx context.Context, value number.Number, attrs []attribute.KeyValue) {
	valid := ctx.Value(contextKey{})
	if valid == nil {
		otel.Handle(fmt.Errorf("async instrument used outside of callback"))
		return
	}
	state := valid.(*State)

	idata, ok := state.store[inst]

	if !ok {
		idata = map[attribute.Set]*stateEntry{}
		state.store[inst] = idata
	}

	aset := attribute.NewSetWithSortable(attrs, &state.tmpSort)
	se, has := idata[aset]
	if !has {
		se = &stateEntry{}
		se.collector = inst.cfactory.New(attrs)
		idata[aset] = se
	}
	se.value = value
}
