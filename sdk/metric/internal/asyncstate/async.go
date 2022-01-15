package asyncstate

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/asyncfloat64"
	"go.opentelemetry.io/otel/metric/asyncint64"
	apiInstrument "go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	Accumulator struct {
		callbacksLock sync.Mutex
		callbacks     []*callback

		instrumentsLock sync.Mutex
		instruments     []apiInstrument.Asynchronous
	}

	State struct {
		collectLock sync.Mutex

		storeLock sync.Mutex
		store     map[*instrument]map[attribute.Set]viewstate.Collector
		tmpSort   attribute.Sortable
	}

	instrument struct {
		descriptor sdkapi.Descriptor
		cfactory   viewstate.CollectorFactory
		callback   *callback
	}

	callback struct {
		function    func(context.Context) error
		instruments []apiInstrument.Asynchronous
	}

	common struct {
		accumulator *Accumulator
		registry    *registry.State
		views       *viewstate.State
	}

	Int64Instruments   struct{ common }
	Float64Instruments struct{ common }

	observer[N number.Any, Traits traits.Any[N]] struct {
		*instrument
	}

	contextKey struct{}
)

var (
	_ apiInstrument.Asynchronous = &instrument{}
)

func (inst *instrument) Asynchronous() {}

// implements registry.hasDescriptor
func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (cb *callback) Instruments() []apiInstrument.Asynchronous {
	return cb.instruments
}

func New() *Accumulator {
	return &Accumulator{}
}

func NewState() *State {
	return &State{
		store: map[*instrument]map[attribute.Set]viewstate.Collector{},
	}
}

func (m *Accumulator) NewCallback(instruments []apiInstrument.Asynchronous, function func(context.Context) error) (metric.Callback, error) {
	cb := &callback{
		function:    function,
		instruments: instruments,
	}

	m.callbacksLock.Lock()
	defer m.callbacksLock.Unlock()

	for _, inst := range instruments {
		ai, ok := inst.(*instrument)
		if !ok {
			return nil, fmt.Errorf("asynchronous instrument does not belong to this provider")
		}
		if ai.descriptor.InstrumentKind().Synchronous() {
			return nil, fmt.Errorf("synchronous instrument with asynchronous callback")
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
		if err := cb.function(ctx); err != nil {
			// TODO: handle errors.
		}
	}

	for inst, states := range state.store {
		for _, entry := range states {
			// Note: not deleting anything is a safe approach for
			// async view calculation.  This can be improved with
			// optional per-stream start timestamps and staleness
			// markers, but is complicated by views correctness
			// considerations, see text in the supplemental
			// guidelines.
			if err := entry.Send(inst.cfactory); err != nil {
				// TODO handle errors
			}
		}
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
	state.storeLock.Lock()
	defer state.storeLock.Unlock()

	idata, ok := state.store[inst]

	if !ok {
		idata = map[attribute.Set]viewstate.Collector{}
		state.store[inst] = idata
	}

	aset := attribute.NewSetWithSortable(attrs, &state.tmpSort)
	se, has := idata[aset]
	if !has {
		se = inst.cfactory.New(attrs, &inst.descriptor)
		idata[aset] = se
	}
	se.(viewstate.Updater[N]).Update(value)
}

func (a *Accumulator) Int64Instruments(reg *registry.State, views *viewstate.State) asyncint64.Instruments {
	return Int64Instruments{
		common: common{
			accumulator: a,
			registry:    reg,
			views:       views,
		},
	}
}

func (a *Accumulator) Float64Instruments(reg *registry.State, views *viewstate.State) asyncfloat64.Instruments {
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
		func(desc sdkapi.Descriptor) (*instrument, error) {
			cfactory, err := c.views.NewFactory(desc)
			if err != nil {
				return nil, err
			}
			inst := &instrument{
				descriptor: desc,
				cfactory:   cfactory,
			}

			c.accumulator.instrumentsLock.Lock()
			defer c.accumulator.instrumentsLock.Unlock()

			c.accumulator.instruments = append(c.accumulator.instruments, inst)
			return inst, nil
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
