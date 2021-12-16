package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
)

type (
	Accumulator struct {
		callbacksLock sync.Mutex
		callbacks     []*callback
	}

	stateEntry struct {
		value    number.Number
		modified uint64
	}

	State struct {
		lock     sync.Mutex
		started  uint64
		finished uint64
		tmpSort  attribute.Sortable
		store    map[*instrument]map[attribute.Set]stateEntry
	}

	instrument struct {
		descriptor sdkapi.Descriptor
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

func New() *Accumulator {
	return &Accumulator{}
}

func NewState() *State {
	return &State{
		store: map[*instrument]map[attribute.Set]stateEntry{},
	}
}

func (m *Accumulator) NewCallback(instruments []sdkapi.Instrument, function func(context.Context) error) (sdkapi.Callback, error) {
	cb := &callback{
		function:    function,
		instruments: instruments,
	}
	// TODO assign instruments, check errors

	m.callbacksLock.Lock()
	defer m.callbacksLock.Unlock()
	m.callbacks = append(m.callbacks, cb)
	return cb, nil
}

func (cb *callback) Instruments() []sdkapi.Instrument {
	return cb.instruments
}

// NewInstrument implements sdkapi.MetricImpl.
func (m *Accumulator) NewInstrument(descriptor sdkapi.Descriptor) (sdkapi.Instrument, error) {
	return &instrument{
		descriptor: descriptor,
	}, nil
}

func (a *Accumulator) getCallbacks() []*callback {
	a.callbacksLock.Lock()
	defer a.callbacksLock.Unlock()
	return a.callbacks
}

func (s *State) startCollect() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.started != s.finished {
		return false
	}
	s.started++
	return true
}

func (s *State) finishCollect() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.finished++
}

func (a *Accumulator) Collect(state *State) error {
	ctx := context.WithValue(
		context.Background(),
		contextKey{},
		state,
	)

	if !state.startCollect() {
		return fmt.Errorf("invalid state transition")
	}

	for _, cb := range a.getCallbacks() {
		cb.function(ctx)
	}

	state.finishCollect()
	return nil
}

func (i *instrument) Descriptor() sdkapi.Descriptor {
	return i.descriptor
}

func (inst *instrument) Capture(ctx context.Context, value number.Number, attrs []attribute.KeyValue) {
	valid := ctx.Value(contextKey{})
	if valid == nil {
		otel.Handle(fmt.Errorf("async instrument used outside of callback"))
		return
	}
	state := valid.(*State)
	state.lock.Lock()
	defer state.lock.Unlock()

	idata, ok := state.store[inst]

	if !ok {
		idata = map[attribute.Set]stateEntry{}
		state.store[inst] = idata
	}

	aset := attribute.NewSetWithSortable(attrs, &state.tmpSort)
	idata[aset] = stateEntry{
		value:    value,
		modified: state.started,
	}
}
