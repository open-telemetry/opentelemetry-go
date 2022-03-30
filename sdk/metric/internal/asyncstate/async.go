package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	apiInstrument "go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	readerState struct {
		lock  sync.Mutex
		store map[attribute.Set]viewstate.Accumulator
	}

	Instrument struct {
		apiInstrument.Asynchronous

		descriptor sdkinstrument.Descriptor
		compiled   viewstate.Instrument
		state      map[*reader.Reader]*readerState
	}

	Callback struct {
		function    func(context.Context)
		instruments map[*Instrument]struct{}
	}

	readerCallback struct {
		*reader.Reader
		*Callback
	}

	observer[N number.Any, Traits traits.Any[N]] struct {
		instrument.Asynchronous

		inst *Instrument
	}

	contextKey struct{}
)

func NewInstrument(desc sdkinstrument.Descriptor, compiled viewstate.Instrument) *Instrument {
	return &Instrument{
		descriptor: desc,
		compiled:   compiled,
		state:      map[*reader.Reader]*readerState{},
	}
}

func NewObserver[N number.Any, Traits traits.Any[N]](inst *Instrument) observer[N, Traits] {
	return observer[N, Traits]{inst: inst}
}

func (inst *Instrument) Descriptor() sdkinstrument.Descriptor {
	return inst.descriptor
}

func NewCallback(instruments []apiInstrument.Asynchronous, function func(context.Context)) (*Callback, error) {
	cb := &Callback{
		function:    function,
		instruments: map[*Instrument]struct{}{},
	}

	for _, inst := range instruments {
		ai, ok := inst.(*Instrument)
		if !ok {
			return nil, fmt.Errorf("asynchronous instrument does not belong to this provider")
		}
		cb.instruments[ai] = struct{}{}
	}

	return cb, nil
}

func (c *Callback) Run(ctx context.Context, r *reader.Reader) {
	c.function(context.WithValue(ctx, contextKey{}, readerCallback{
		Reader:   r,
		Callback: c,
	}))
}

func (rs *readerState) accumulate() {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	for _, capt := range rs.store {
		capt.Accumulate()
	}
}

func (inst *Instrument) Collect(r *reader.Reader, sequence reader.Sequence, output *[]reader.Instrument) {
	rs := inst.state[r]

	rs.accumulate()

	inst.compiled.Collect(r, sequence, output)

	// Reset the instruments used; the view state will remember what it needs.
	rs.store = map[attribute.Set]viewstate.Accumulator{}
}

func (o observer[N, Traits]) Observe(ctx context.Context, value N, attrs ...attribute.KeyValue) {
	if o.inst == nil {
		return
	}

	lookup := ctx.Value(contextKey{})
	if lookup == nil {
		otel.Handle(fmt.Errorf("async instrument used outside of callback"))
		return
	}

	rc := lookup.(readerCallback)
	if _, ok := rc.Callback.instruments[o.inst]; !ok {
		otel.Handle(fmt.Errorf("async instrument not declared for use in callback"))
	}

	se := o.inst.get(rc.Reader, attrs)
	se.(viewstate.AccumulatorUpdater[N]).Update(value)
}

func (inst *Instrument) get(r *reader.Reader, attrs []attribute.KeyValue) viewstate.Accumulator {
	rs := inst.state[r]
	rs.lock.Lock()
	defer rs.lock.Unlock()

	aset := attribute.NewSet(attrs...)
	se, has := rs.store[aset]
	if !has {
		se = inst.compiled.NewAccumulator(attrs, r)
		rs.store[aset] = se
	}
	return se
}
