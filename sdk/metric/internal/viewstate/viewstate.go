package viewstate

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/views"
)

type (
	Compiler struct {
		library instrumentation.Library
		views   []views.View
		readers []*reader.Reader

		// names is the per-reader map of output names,
		// indexed by the reader's position in `readers`.
		namesLock sync.Mutex
		names     []map[string]struct{}
	}

	Instrument interface {
		// NewAccumulator returns a new Accumulator bound to
		// the attributes `kvs`.  If reader == nil the
		// accumulator applies to all readers, otherwise it
		// applies to the specific reader.
		NewAccumulator(kvs []attribute.KeyValue, reader *reader.Reader) Accumulator

		Collector
	}

	Collector interface {
		// Collect transfers aggregated data from the
		// Accumulators into the output struct.
		Collect(reader *reader.Reader, sequence reader.Sequence, output *[]reader.Instrument)
	}

	Accumulator interface {
		Accumulate()
	}

	Updater[N number.Any] interface {
		Update(value N)
	}

	AccumulatorUpdater[N number.Any] interface {
		Accumulator
		Updater[N]
	}

	multiInstrument[N number.Any] map[*reader.Reader][]Instrument

	multiAccumulator[N number.Any] []Accumulator

	configuredBehavior struct {
		desc     sdkapi.Descriptor
		view     views.View
		settings aggregatorSettings
		reader   *reader.Reader
	}

	aggregatorSettings struct {
		kind  aggregation.Kind
		hcfg  histogram.Config
		scfg  sum.Config
		lvcfg lastvalue.Config
	}

	baseMetric[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		lock     sync.Mutex
		desc     sdkapi.Descriptor
		acfg     *Config
		data     map[attribute.Set]*Storage
		viewKeys attribute.Filter
	}

	compiledSyncInstrument[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		baseMetric[N, Storage, Config, Methods]
	}

	compiledAsyncInstrument[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		baseMetric[N, Storage, Config, Methods]
	}

	statelessSyncProcess[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		compiledSyncInstrument[N, Storage, Config, Methods]
	}

	statefulAsyncProcess[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		compiledAsyncInstrument[N, Storage, Config, Methods]
		prior map[attribute.Set]*Storage
	}

	statefulSyncProcess[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		compiledSyncInstrument[N, Storage, Config, Methods]
	}

	statelessAsyncProcess[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		compiledAsyncInstrument[N, Storage, Config, Methods]
	}

	syncAccumulator[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		current  Storage
		snapshot Storage
		output   *Storage
	}

	asyncAccumulator[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		lock     sync.Mutex
		current  N
		snapshot Storage
		output   *Storage
	}
)

func New(lib instrumentation.Library, views []views.View, readers []*reader.Reader) *Compiler {

	// TODO: error checking here, such as:
	// - empty (?)
	// - duplicate name
	// - invalid inst/number/aggregation kind
	// - both instrument name and regexp
	// - schemaURL or Version without library name
	// - empty attribute keys
	// - Name w/o SingleInst

	names := make([]map[string]struct{}, len(readers))
	for i := range names {
		names[i] = map[string]struct{}{}
	}

	return &Compiler{
		library: lib,
		views:   views,
		readers: readers,
		names:   names,
	}
}

// Compile is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Accumulators throughout its lifetime.
func (v *Compiler) Compile(instrument sdkapi.Descriptor) Instrument {
	configs := make([][]configuredBehavior, len(v.readers))
	matches := make([]views.View, 0, len(v.views))

	for _, view := range v.views {
		if !view.Matches(v.library, instrument) {
			continue
		}
		matches = append(matches, view)
	}

	for readerIdx, r := range v.readers {
		for _, view := range matches {
			var as aggregatorSettings
			switch view.Aggregation() {
			case aggregation.SumKind, aggregation.LastValueKind:
				// These have no options
				as.kind = view.Aggregation()
			case aggregation.HistogramKind:
				as.kind = view.Aggregation()
				as.hcfg = histogram.NewConfig(
					// @@@ per-reader histogram defaults
					histogramDefaultsFor(instrument.NumberKind()),
					view.HistogramOptions()...,
				)
			default:
				as = aggregatorSettingsFor(instrument, r.Defaults())
			}

			if as.kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				desc:     instrument,
				view:     view,
				settings: as,
				reader:   r,
			})
		}

		// If there were no matching views, set the default aggregation.
		if len(matches) == 0 {
			as := aggregatorSettingsFor(instrument, r.Defaults())
			if as.kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				desc:     instrument,
				view:     views.New(views.WithAggregation(as.kind)),
				settings: as,
				reader:   r,
			})
		}
	}

	compiled := map[*reader.Reader][]Instrument{}

	v.namesLock.Lock()
	defer v.namesLock.Unlock()

	for readerIdx, readerList := range configs {
		r := v.readers[readerIdx]

		for _, config := range readerList {
			config.desc = viewDescriptor(config.desc, config.view)

			if _, has := v.names[readerIdx][config.desc.Name()]; has {
				otel.Handle(fmt.Errorf("duplicate view name registered"))
				continue
			}
			v.names[readerIdx][config.desc.Name()] = struct{}{}

			var one Instrument
			switch config.desc.NumberKind() {
			case number.Int64Kind:
				one = buildView[int64, traits.Int64](config)
			case number.Float64Kind:
				one = buildView[float64, traits.Float64](config)
			}

			compiled[r] = append(compiled[r], one)
		}
	}

	switch len(compiled) {
	case 0:
		return nil // TODO does this require a Noop?
	case 1:
		// As a special case, recognize the case where there
		// is only one reader and only one view to bypass the
		// map[...][]Instrument wrapper.
		for _, list := range compiled {
			if len(list) == 1 {
				return list[0]
			}
		}
	}
	if instrument.NumberKind() == number.Int64Kind {
		return multiInstrument[int64](compiled)
	}
	return multiInstrument[float64](compiled)
}

func aggregatorSettingsFor(desc sdkapi.Descriptor, defaults reader.DefaultsFunc) aggregatorSettings {
	aggr, _ := defaults(desc.InstrumentKind())
	return aggregatorSettings{
		kind: aggr,
	}
}

func viewDescriptor(instrument sdkapi.Descriptor, v views.View) sdkapi.Descriptor {
	ikind := instrument.InstrumentKind()
	nkind := instrument.NumberKind()
	name := instrument.Name()
	description := instrument.Description()
	unit := instrument.Unit()
	if v.HasName() {
		name = v.Name()
	}
	if v.Description() != "" {
		description = instrument.Description()
	}
	return sdkapi.NewDescriptor(name, ikind, nkind, description, unit)
}

func histogramDefaultsFor(kind number.Kind) histogram.Defaults {
	if kind == number.Int64Kind {
		return histogram.Int64Defaults{}
	}
	return histogram.Float64Defaults{}
}

func buildView[N number.Any, Traits traits.Any[N]](config configuredBehavior) Instrument {
	if config.desc.InstrumentKind().Synchronous() {
		return compileSync[N, Traits](config)
	}
	return compileAsync[N, Traits](config)
}

func newSyncView[
	N number.Any,
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](config configuredBehavior, aggConfig *Config) Instrument {
	_, tempo := config.reader.Defaults()(config.desc.InstrumentKind())
	metric := baseMetric[N, Storage, Config, Methods]{
		desc:     config.desc,
		acfg:     aggConfig,
		data:     map[attribute.Set]*Storage{},
		viewKeys: config.view.Keys(),
	}
	instrument := compiledSyncInstrument[N, Storage, Config, Methods]{
		baseMetric: metric,
	}
	if tempo == aggregation.DeltaTemporality {
		return &statelessSyncProcess[N, Storage, Config, Methods]{
			compiledSyncInstrument: instrument,
		}
	}

	return &statefulSyncProcess[N, Storage, Config, Methods]{
		compiledSyncInstrument: instrument,
	}
}

func newAsyncView[
	N number.Any,
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](config configuredBehavior, aggConfig *Config) Instrument {
	_, tempo := config.reader.Defaults()(config.desc.InstrumentKind())
	metric := baseMetric[N, Storage, Config, Methods]{
		desc:     config.desc,
		acfg:     aggConfig,
		data:     map[attribute.Set]*Storage{},
		viewKeys: config.view.Keys(),
	}
	instrument := compiledAsyncInstrument[N, Storage, Config, Methods]{
		baseMetric: metric,
	}
	if tempo == aggregation.DeltaTemporality {
		return &statefulAsyncProcess[N, Storage, Config, Methods]{
			compiledAsyncInstrument: instrument,

			// Note: this is extra storage relative to the
			// other three methods.
			prior: map[attribute.Set]*Storage{},
		}
	}

	return &statelessAsyncProcess[N, Storage, Config, Methods]{
		compiledAsyncInstrument: instrument,
	}
}

func compileSync[N number.Any, Traits traits.Any[N]](config configuredBehavior) Instrument {
	switch config.settings.kind {
	case aggregation.LastValueKind:
		return newSyncView[
			N,
			lastvalue.State[N, Traits],
			lastvalue.Config,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
		](config, &config.settings.lvcfg)
	case aggregation.HistogramKind:
		return newSyncView[
			N,
			histogram.State[N, Traits],
			histogram.Config,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config, &config.settings.hcfg)
	default:
		return newSyncView[
			N,
			sum.State[N, Traits],
			sum.Config,
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config, &config.settings.scfg)
	}
}

func compileAsync[N number.Any, Traits traits.Any[N]](config configuredBehavior) Instrument {
	switch config.settings.kind {
	case aggregation.LastValueKind:
		return newAsyncView[
			N,
			lastvalue.State[N, Traits],
			lastvalue.Config,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
		](config, &config.settings.lvcfg)
	case aggregation.HistogramKind:
		return newAsyncView[
			N,
			histogram.State[N, Traits],
			histogram.Config,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config, &config.settings.hcfg)
	default:
		return newAsyncView[
			N,
			sum.State[N, Traits],
			sum.Config,
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config, &config.settings.scfg)
	}
}

// NewAccumulator returns a Accumulator for multiple views of the same instrument.
func (mi multiInstrument[N]) NewAccumulator(kvs []attribute.KeyValue, reader *reader.Reader) Accumulator {
	var collectors []Accumulator
	// Note: This runtime switch happens because we're using the same API for
	// both async and sync instruments, whereas the APIs are not symmetrical.
	if reader == nil {
		for _, list := range mi {
			collectors = make([]Accumulator, 0, len(mi)*len(list))
		}
		for _, list := range mi {
			for _, inst := range list {
				collectors = append(collectors, inst.NewAccumulator(kvs, nil))
			}
		}
	} else {
		insts := mi[reader]

		collectors = make([]Accumulator, 0, len(insts))

		for _, inst := range insts {
			collectors = append(collectors, inst.NewAccumulator(kvs, reader))
		}
	}
	return multiAccumulator[N](collectors)
}

// multiAccumulator

func (c multiAccumulator[N]) Accumulate() {
	for _, coll := range c {
		coll.Accumulate()
	}
}

func (c multiAccumulator[N]) Update(value N) {
	for _, coll := range c {
		coll.(Updater[N]).Update(value)
	}
}

// syncAccumulator

func (sc *syncAccumulator[N, Storage, Config, Methods]) Update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncAccumulator[N, Storage, Config, Methods]) Accumulate() {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)
	methods.Merge(&sc.snapshot, sc.output)
}

// asyncAccumulator

func (ac *asyncAccumulator[N, Storage, Config, Methods]) Update(number N) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	ac.current = number
}

func (ac *asyncAccumulator[N, Storage, Config, Methods]) Accumulate() {
	ac.lock.Lock()
	defer ac.lock.Unlock()

	var methods Methods
	methods.Reset(&ac.snapshot)
	methods.Update(&ac.snapshot, ac.current)
	methods.Merge(&ac.snapshot, ac.output)
}

// baseMetric

func (metric *baseMetric[N, Storage, Config, Methods]) initStorage(s *Storage) {
	var methods Methods
	methods.Init(s, *metric.acfg)
}

func (metric *baseMetric[N, Storage, Config, Methods]) findOutput(
	kvs []attribute.KeyValue,
) *Storage {
	set, _ := attribute.NewSetWithFiltered(kvs, metric.viewKeys)

	metric.lock.Lock()
	defer metric.lock.Unlock()

	storage, has := metric.data[set]
	if has {
		return storage
	}

	ns := metric.newStorage()
	metric.data[set] = ns
	return ns
}

func (metric *baseMetric[N, Storage, Config, Methods]) newStorage() *Storage {
	ns := new(Storage)
	metric.initStorage(ns)
	return ns
}

// NewAccumulator

// NewAccumulator returns a Accumulator for a synchronous instrument view.
func (csv *compiledSyncInstrument[N, Storage, Config, Methods]) NewAccumulator(kvs []attribute.KeyValue, _ *reader.Reader) Accumulator {
	sc := &syncAccumulator[N, Storage, Config, Methods]{}
	csv.initStorage(&sc.current)
	csv.initStorage(&sc.snapshot)

	sc.output = csv.findOutput(kvs)

	return sc
}

// NewAccumulator returns a Accumulator for an asynchronous instrument view.
func (cav *compiledAsyncInstrument[N, Storage, Config, Methods]) NewAccumulator(kvs []attribute.KeyValue, _ *reader.Reader) Accumulator {
	ac := &asyncAccumulator[N, Storage, Config, Methods]{}

	cav.initStorage(&ac.snapshot)
	ac.current = 0
	ac.output = cav.findOutput(kvs)

	return ac
}

// Collect collects for multiple instruments
func (mi multiInstrument[N]) Collect(reader *reader.Reader, sequence reader.Sequence, output *[]reader.Instrument) {
	for _, inst := range mi[reader] {
		inst.Collect(reader, sequence, output)
	}
}

// Collect is for Cumulative, Stateless asynchronous instruments
func (p *statelessAsyncProcess[N, Storage, Config, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	*output = append(*output, reader.Instrument{
		Instrument:  p.desc,
		Temporality: aggregation.CumulativeTemporality,
	})
	ioutput := &(*output)[len(*output)-1]

	for set, storage := range p.data {
		// Copy the underlying storage.
		ioutput.Series = append(ioutput.Series, reader.Series{
			Attributes:  set,
			Aggregation: methods.Aggregation(storage),
			Start:       seq.Start,
			End:         seq.Now,
		})
	}
	// Reset the entire map.
	p.data = map[attribute.Set]*Storage{}
}

// Collect is for Delta, Stateful synchronous instruments
func (p *statefulSyncProcess[N, Storage, Config, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	*output = append(*output, reader.Instrument{
		Instrument:  p.desc,
		Temporality: aggregation.CumulativeTemporality,
	})
	ioutput := &(*output)[len(*output)-1]

	for set, storage := range p.data {
		// Note: No reset in this process.
		ioutput.Series = append(ioutput.Series, reader.Series{
			Attributes:  set,
			Aggregation: methods.Aggregation(storage), // This is a direct reference to the underlying storage.
			Start:       seq.Last,
			End:         seq.Now,
		})
	}
}

// Delta (Stateless)
func (p *statelessSyncProcess[N, Storage, Config, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	*output = append(*output, reader.Instrument{
		Instrument:  p.desc,
		Temporality: aggregation.DeltaTemporality,
	})
	ioutput := &(*output)[len(*output)-1]

	for set, storage := range p.data {
		if !methods.HasData(storage) {
			delete(p.data, set)
			continue
		}

		// Copy and reset the underlying storage.
		// Note: this could be re-used memory from the last collection.
		ns := p.newStorage()
		methods.Merge(ns, storage)
		methods.Reset(storage)

		ioutput.Series = append(ioutput.Series, reader.Series{
			Attributes:  set,
			Aggregation: methods.Aggregation(ns),
			Start:       seq.Last,
			End:         seq.Now,
		})
	}
}

func (p *statefulAsyncProcess[N, Storage, Config, Methods]) Collect(_ *reader.Reader, sequence reader.Sequence, output *[]reader.Instrument) {
	// @@@ HERE YOU ARE
}
