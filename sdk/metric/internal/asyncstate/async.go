package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	apiInstrument "go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
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
		state      map[*reader.ReaderConfig]*readerState
	}

	Callback struct {
		function    func(context.Context)
		instruments map[*Instrument]struct{}
	}

	readerCallback struct {
		*reader.ReaderConfig
		*Callback
	}

	observer[N number.Any, Traits traits.Any[N]] struct {
		instrument.Asynchronous

		inst *Instrument
	}

	memberInstrument interface {
		instrument() *Instrument
	}

	contextKey struct{}
)

var _ memberInstrument = observer[int64, traits.Int64]{}
var _ memberInstrument = observer[float64, traits.Float64]{}

func NewInstrument(desc sdkinstrument.Descriptor, compiled viewstate.Instrument, readers []*reader.ReaderConfig) *Instrument {
	state := map[*reader.ReaderConfig]*readerState{}
	for _, r := range readers {
		state[r] = &readerState{
			store: map[attribute.Set]viewstate.Accumulator{},
		}
	}
	return &Instrument{
		descriptor: desc,
		compiled:   compiled,
		state:      state,
	}
}

func NewObserver[N number.Any, Traits traits.Any[N]](inst *Instrument) observer[N, Traits] {
	return observer[N, Traits]{inst: inst}
}

func NewCallback(instruments []apiInstrument.Asynchronous, function func(context.Context)) (*Callback, error) {
	cb := &Callback{
		function:    function,
		instruments: map[*Instrument]struct{}{},
	}

	for _, inst := range instruments {
		ai, ok := inst.(memberInstrument)
		if !ok {
			return nil, fmt.Errorf("asynchronous instrument does not belong to this provider: %T", inst)
		}
		cb.instruments[ai.instrument()] = struct{}{}
	}

	return cb, nil
}

func (c *Callback) Run(ctx context.Context, r *reader.ReaderConfig) {
	c.function(context.WithValue(ctx, contextKey{}, readerCallback{
		ReaderConfig: r,
		Callback:     c,
	}))
}

func (inst *Instrument) AccumulateFor(r *reader.ReaderConfig) {
	rs := inst.state[r]

	// This limits concurrent asynchronous collection, which is
	// only needed in stateful configurations (i.e.,
	// cumulative-to-delta). TODO: does it matter that this blocks
	// concurrent Prometheus scrapers concurrently? (I think not.)
	rs.lock.Lock()
	defer rs.lock.Unlock()

	for _, capt := range rs.store {
		capt.Accumulate()
	}

	// Reset the instruments used; the view state will remember
	// what it needs.
	rs.store = map[attribute.Set]viewstate.Accumulator{}
}

func (o observer[N, Traits]) instrument() *Instrument {
	return o.inst
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

	if err := aggregator.RangeTest[N, Traits](value, &o.inst.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	se := o.inst.get(rc.ReaderConfig, attrs)
	se.(viewstate.AccumulatorUpdater[N]).Update(value)
}

func (inst *Instrument) get(r *reader.ReaderConfig, attrs []attribute.KeyValue) viewstate.Accumulator {
	rs := inst.state[r]
	rs.lock.Lock()
	defer rs.lock.Unlock()

	aset := attribute.NewSet(attrs...)
	se, has := rs.store[aset]
	if !has {
		se = inst.compiled.NewAccumulator(aset, r)
		rs.store[aset] = se
	}
	return se
}
