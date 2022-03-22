package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	apiInstrument "go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	Accumulator struct {
		callbacksLock sync.Mutex
		callbacks     []*callback

		instrumentsLock sync.Mutex
		instruments     []*instrument

		statesLock sync.Mutex
		states     map[*reader.Reader]*State
	}

	State struct {
		reader  *reader.Reader
		tmpSort attribute.Sortable
	}

	instrument struct {
		apiInstrument.Asynchronous

		descriptor sdkapi.Descriptor
		compiled   viewstate.Instrument
		storeLock  sync.Mutex
		store      map[attribute.Set]viewstate.Collector
		callback   *callback
	}

	callback struct {
		function    func(context.Context)
		instruments []apiInstrument.Asynchronous
	}

	common struct {
		accumulator *Accumulator
		registry    *registry.State
		views       *viewstate.Compiler
	}

	Int64Instruments   struct{ common }
	Float64Instruments struct{ common }

	observer[N number.Any, Traits traits.Any[N]] struct {
		*instrument
	}

	contextKey struct{}
)

// implements registry.hasDescriptor
func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (cb *callback) Instruments() []apiInstrument.Asynchronous {
	return cb.instruments
}

func New() *Accumulator {
	return &Accumulator{
		states: map[*reader.Reader]*State{},
	}
}

func (m *Accumulator) stateFor(reader *reader.Reader) *State {
	m.statesLock.Lock()
	defer m.statesLock.Unlock()
	if s, ok := m.states[reader]; ok {
		return s
	}
	s := &State{
		reader: reader,
	}
	m.states[reader] = s
	return s
}

func (m *Accumulator) RegisterCallback(instruments []apiInstrument.Asynchronous, function func(context.Context)) error {
	cb := &callback{
		function:    function,
		instruments: instruments,
	}

	m.callbacksLock.Lock()
	defer m.callbacksLock.Unlock()

	for _, inst := range instruments {
		ai, ok := inst.(*instrument)
		if !ok {
			return fmt.Errorf("asynchronous instrument does not belong to this provider")
		}
		if ai.descriptor.InstrumentKind().Synchronous() {
			return fmt.Errorf("synchronous instrument with asynchronous callback")
		}
		if ai.callback != nil {
			return fmt.Errorf("asynchronous instrument already has a callback")
		}
		ai.callback = cb

	}

	m.callbacks = append(m.callbacks, cb)
	return nil
}

func (a *Accumulator) getCallbacks() []*callback {
	a.callbacksLock.Lock()
	defer a.callbacksLock.Unlock()
	return a.callbacks
}

func (a *Accumulator) Collect(r *reader.Reader, output *[]reader.Instrument) error {
	state := a.stateFor(r)
	ctx := context.WithValue(
		context.Background(),
		contextKey{},
		state,
	)

	// TODO: Add a timeout to the context.

	for _, cb := range a.getCallbacks() {
		cb.function(ctx)
	}

	a.instrumentsLock.Lock()
	instruments := a.instruments
	a.instrumentsLock.Unlock()

	*output = make([]reader.Instrument, len(instruments))

	for instIdx, inst := range instruments {
		iout := &(*output)[instIdx]

		iout.Instrument = inst.descriptor
		iout.Temporality = 0 // @@@ Hey!!!
		
		inst.storeLock.Lock()
		// This iteration passes over each of the attribute
		// sets the user passed in.
		for _, coll := range inst.store {
			coll.Collect()
		}
		inst.storeLock.Unlock()

		inst.compiled.Write(r, &iout.Series)
	}

	return nil
}

func capture[N number.Any, Traits traits.Any[N]](ctx context.Context, inst *instrument, value N, attrs []attribute.KeyValue) {
	valid := ctx.Value(contextKey{})
	if valid == nil {
		otel.Handle(fmt.Errorf("async instrument used outside of callback"))
		return
	}
	state := valid.(*State)

	se := getStateEntry(state, inst, attrs)
	se.(viewstate.CollectorUpdater[N]).Update(value)
}

func getStateEntry(state *State, inst *instrument, attrs []attribute.KeyValue) viewstate.Collector {
	inst.storeLock.Lock()
	defer inst.storeLock.Unlock()

	aset := attribute.NewSetWithSortable(attrs, &state.tmpSort)
	se, has := inst.store[aset]
	if !has {
		se = inst.compiled.NewCollector(attrs, state.reader)
		inst.store[aset] = se
	}
	return se
}

func (a *Accumulator) Int64Instruments(reg *registry.State, views *viewstate.Compiler) asyncint64.InstrumentProvider {
	return Int64Instruments{
		common: common{
			accumulator: a,
			registry:    reg,
			views:       views,
		},
	}
}

func (a *Accumulator) Float64Instruments(reg *registry.State, views *viewstate.Compiler) asyncfloat64.InstrumentProvider {
	return Float64Instruments{
		common: common{
			accumulator: a,
			registry:    reg,
			views:       views,
		},
	}
}

func (o observer[N, Traits]) Observe(ctx context.Context, value N, attrs ...attribute.KeyValue) {
	if o.instrument != nil {
		capture[N, Traits](ctx, o.instrument, value, attrs)
	}
}

func (c common) newInstrument(name string, opts []apiInstrument.Option, nk number.Kind, ik sdkapi.InstrumentKind) (*instrument, error) {
	return registry.Lookup(
		c.registry,
		name, opts, nk, ik,
		func(desc sdkapi.Descriptor) *instrument {
			compiled := c.views.Compile(desc)
			inst := &instrument{
				descriptor: desc,
				compiled:   compiled,
				store:      map[attribute.Set]viewstate.Collector{},
			}

			c.accumulator.instrumentsLock.Lock()
			defer c.accumulator.instrumentsLock.Unlock()

			c.accumulator.instruments = append(c.accumulator.instruments, inst)
			return inst
		})
}

func (i Int64Instruments) Counter(name string, opts ...apiInstrument.Option) (asyncint64.Counter, error) {
	inst, err := i.newInstrument(name, opts, number.Int64Kind, sdkapi.CounterObserverInstrumentKind)
	return observer[int64, traits.Int64]{instrument: inst}, err
}

func (i Int64Instruments) UpDownCounter(name string, opts ...apiInstrument.Option) (asyncint64.UpDownCounter, error) {
	inst, err := i.newInstrument(name, opts, number.Int64Kind, sdkapi.UpDownCounterObserverInstrumentKind)
	return observer[int64, traits.Int64]{instrument: inst}, err
}

func (i Int64Instruments) Gauge(name string, opts ...apiInstrument.Option) (asyncint64.Gauge, error) {
	inst, err := i.newInstrument(name, opts, number.Int64Kind, sdkapi.GaugeObserverInstrumentKind)
	return observer[int64, traits.Int64]{instrument: inst}, err
}

func (f Float64Instruments) Counter(name string, opts ...apiInstrument.Option) (asyncfloat64.Counter, error) {
	inst, err := f.newInstrument(name, opts, number.Float64Kind, sdkapi.CounterObserverInstrumentKind)
	return observer[float64, traits.Float64]{instrument: inst}, err
}

func (f Float64Instruments) UpDownCounter(name string, opts ...apiInstrument.Option) (asyncfloat64.UpDownCounter, error) {
	inst, err := f.newInstrument(name, opts, number.Float64Kind, sdkapi.UpDownCounterObserverInstrumentKind)
	return observer[float64, traits.Float64]{instrument: inst}, err
}

func (f Float64Instruments) Gauge(name string, opts ...apiInstrument.Option) (asyncfloat64.Gauge, error) {
	inst, err := f.newInstrument(name, opts, number.Float64Kind, sdkapi.GaugeObserverInstrumentKind)
	return observer[float64, traits.Float64]{instrument: inst}, err
}
