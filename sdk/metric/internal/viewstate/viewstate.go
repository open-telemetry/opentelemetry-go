package viewstate

import (
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
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
		desc   sdkinstrument.Descriptor
		kind   aggregation.Kind
		keys   attribute.Filter
		acfg   aggregator.Config
		reader *reader.Reader
	}

	baseMetric[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		lock sync.Mutex
		desc sdkinstrument.Descriptor
		acfg aggregator.Config
		data map[attribute.Set]*Storage
		keys attribute.Filter
	}

	compiledSyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		baseMetric[N, Storage, Methods]
	}

	compiledAsyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		baseMetric[N, Storage, Methods]
	}

	statelessSyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledSyncInstrument[N, Storage, Methods]
	}

	statefulSyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledSyncInstrument[N, Storage, Methods]
	}

	statelessAsyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledAsyncInstrument[N, Storage, Methods]
	}

	statefulAsyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledAsyncInstrument[N, Storage, Methods]
		prior map[attribute.Set]*Storage
	}

	syncAccumulator[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		current  Storage
		snapshot Storage
		output   *Storage
	}

	asyncAccumulator[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
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
func (v *Compiler) Compile(instrument sdkinstrument.Descriptor) Instrument {
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
			var acfg aggregator.Config

			kind := view.Aggregation()
			switch kind {
			case aggregation.SumKind, aggregation.GaugeKind:
				// These have no options
			case aggregation.HistogramKind:
				acfg.Histogram = histogram.NewConfig(
					// @@@ per-reader histogram defaults
					histogramDefaultsFor(instrument.NumberKind),
					view.HistogramOptions()...,
				)
			default:
				kind = aggregatorConfigFor(instrument, r.Defaults())
			}

			if kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				desc:   viewDescriptor(instrument, view),
				keys:   view.Keys(),
				kind:   kind,
				acfg:   acfg,
				reader: r,
			})
		}

		// If there were no matching views, set the default aggregation.
		if len(matches) == 0 {
			kind := aggregatorConfigFor(instrument, r.Defaults())
			if kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				desc:   instrument,
				kind:   kind,
				reader: r,
			})
		}
	}

	compiled := map[*reader.Reader][]Instrument{}

	v.namesLock.Lock()
	defer v.namesLock.Unlock()

	for readerIdx, readerList := range configs {
		r := v.readers[readerIdx]

		for _, config := range readerList {
			if _, has := v.names[readerIdx][config.desc.Name]; has {
				otel.Handle(fmt.Errorf("duplicate view name registered"))
				continue
			}
			v.names[readerIdx][config.desc.Name] = struct{}{}

			var one Instrument
			switch config.desc.NumberKind {
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
		return nil
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
	if instrument.NumberKind == number.Int64Kind {
		return multiInstrument[int64](compiled)
	}
	return multiInstrument[float64](compiled)
}

func aggregatorConfigFor(desc sdkinstrument.Descriptor, defaults reader.DefaultsFunc) aggregation.Kind {
	aggr, _ := defaults(desc.Kind)
	return aggr
}

func viewDescriptor(instrument sdkinstrument.Descriptor, v views.View) sdkinstrument.Descriptor {
	ikind := instrument.Kind
	nkind := instrument.NumberKind
	name := instrument.Name
	description := instrument.Description
	unit := instrument.Unit
	if v.HasName() {
		name = v.Name()
	}
	if v.Description() != "" {
		description = instrument.Description
	}
	return sdkinstrument.NewDescriptor(name, ikind, nkind, description, unit)
}

func histogramDefaultsFor(kind number.Kind) histogram.Defaults {
	if kind == number.Int64Kind {
		return histogram.Int64Defaults{}
	}
	return histogram.Float64Defaults{}
}

func buildView[N number.Any, Traits traits.Any[N]](config configuredBehavior) Instrument {
	if config.desc.Kind.Synchronous() {
		return compileSync[N, Traits](config)
	}
	return compileAsync[N, Traits](config)
}

func newSyncView[
	N number.Any,
	Storage any,
	Methods aggregator.Methods[N, Storage],
](config configuredBehavior) Instrument {
	_, tempo := config.reader.Defaults()(config.desc.Kind)
	metric := baseMetric[N, Storage, Methods]{
		desc: config.desc,
		acfg: config.acfg,
		data: map[attribute.Set]*Storage{},
		keys: config.keys,
	}
	instrument := compiledSyncInstrument[N, Storage, Methods]{
		baseMetric: metric,
	}
	if tempo == aggregation.DeltaTemporality {
		return &statelessSyncProcess[N, Storage, Methods]{
			compiledSyncInstrument: instrument,
		}
	}

	return &statefulSyncProcess[N, Storage, Methods]{
		compiledSyncInstrument: instrument,
	}
}

func newAsyncView[
	N number.Any,
	Storage any,
	Methods aggregator.Methods[N, Storage],
](config configuredBehavior) Instrument {
	_, tempo := config.reader.Defaults()(config.desc.Kind)
	metric := baseMetric[N, Storage, Methods]{
		desc: config.desc,
		acfg: config.acfg,
		data: map[attribute.Set]*Storage{},
		keys: config.keys,
	}
	instrument := compiledAsyncInstrument[N, Storage, Methods]{
		baseMetric: metric,
	}
	if tempo == aggregation.DeltaTemporality {
		return &statefulAsyncProcess[N, Storage, Methods]{
			compiledAsyncInstrument: instrument,
		}
	}

	return &statelessAsyncProcess[N, Storage, Methods]{
		compiledAsyncInstrument: instrument,
	}
}

func compileSync[N number.Any, Traits traits.Any[N]](config configuredBehavior) Instrument {
	switch config.kind {
	case aggregation.GaugeKind:
		return newSyncView[
			N,
			gauge.State[N, Traits],
			gauge.Methods[N, Traits, gauge.State[N, Traits]],
		](config)
	case aggregation.HistogramKind:
		return newSyncView[
			N,
			histogram.State[N, Traits],
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config)
	default:
		return newSyncView[
			N,
			sum.State[N, Traits],
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config)
	}
}

func compileAsync[N number.Any, Traits traits.Any[N]](config configuredBehavior) Instrument {
	switch config.kind {
	case aggregation.GaugeKind:
		return newAsyncView[
			N,
			gauge.State[N, Traits],
			gauge.Methods[N, Traits, gauge.State[N, Traits]],
		](config)
	case aggregation.HistogramKind:
		return newAsyncView[
			N,
			histogram.State[N, Traits],
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config)
	default:
		return newAsyncView[
			N,
			sum.State[N, Traits],
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config)
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

func (sc *syncAccumulator[N, Storage, Methods]) Update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncAccumulator[N, Storage, Methods]) Accumulate() {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)
	methods.Merge(sc.output, &sc.snapshot)
}

// asyncAccumulator

func (ac *asyncAccumulator[N, Storage, Methods]) Update(number N) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	ac.current = number
}

func (ac *asyncAccumulator[N, Storage, Methods]) Accumulate() {
	ac.lock.Lock()
	defer ac.lock.Unlock()

	var methods Methods
	methods.Reset(&ac.snapshot)
	methods.Update(&ac.snapshot, ac.current)
	methods.Merge(&ac.snapshot, ac.output)
}

// baseMetric

func (metric *baseMetric[N, Storage, Methods]) initStorage(s *Storage) {
	var methods Methods
	methods.Init(s, metric.acfg)
}

func (metric *baseMetric[N, Storage, Methods]) findOutput(
	kvs []attribute.KeyValue,
) *Storage {
	set, _ := attribute.NewSetWithFiltered(kvs, metric.keys)

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

func (metric *baseMetric[N, Storage, Methods]) newStorage() *Storage {
	ns := new(Storage)
	metric.initStorage(ns)
	return ns
}

// NewAccumulator returns a Accumulator for a synchronous instrument view.
func (csv *compiledSyncInstrument[N, Storage, Methods]) NewAccumulator(kvs []attribute.KeyValue, _ *reader.Reader) Accumulator {
	sc := &syncAccumulator[N, Storage, Methods]{}
	csv.initStorage(&sc.current)
	csv.initStorage(&sc.snapshot)

	sc.output = csv.findOutput(kvs)

	return sc
}

// NewAccumulator returns a Accumulator for an asynchronous instrument view.
func (cav *compiledAsyncInstrument[N, Storage, Methods]) NewAccumulator(kvs []attribute.KeyValue, _ *reader.Reader) Accumulator {
	ac := &asyncAccumulator[N, Storage, Methods]{}

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

func reuseLast[T any](p *[]T) *T {
	if len(*p) < cap(*p) {
		(*p) = (*p)[0 : len(*p)+1 : cap(*p)]
	} else {
		var empty T
		(*p) = append(*p, empty)
	}
	return &(*p)[len(*p)-1]
}

func appendInstrument(insts *[]reader.Instrument, desc sdkinstrument.Descriptor, tempo aggregation.Temporality) *reader.Instrument {
	ioutput := reuseLast(insts)
	ioutput.Instrument = desc
	ioutput.Temporality = tempo
	return ioutput
}

func appendSeries(series *[]reader.Series, set attribute.Set, agg aggregation.Aggregation, start, end time.Time) {
	soutput := reuseLast(series)
	soutput.Attributes = set
	soutput.Aggregation = agg
	soutput.Start = start
	soutput.End = end
}

// Collect for Synchronous Delta->Cumulative
func (p *statefulSyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.CumulativeTemporality)

	for set, storage := range p.data {
		// Note: No reset in this process.
		// This takes a direct reference to the underlying storage.
		appendSeries(&ioutput.Series, set, methods.Aggregation(storage), seq.Start, seq.Now)
	}
}

// Collect for Synchronous Delta->Delta
func (p *statelessSyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.DeltaTemporality)

	for set, storage := range p.data {
		// Note: this can't be a Gauge b/c synchronous instrument.
		if !methods.HasChange(storage) {
			delete(p.data, set)
			continue
		}

		// Copy and reset the underlying storage.
		// Note: this could be re-used memory from the last collection.
		ns := p.newStorage()
		methods.Merge(ns, storage)
		methods.Reset(storage)

		appendSeries(&ioutput.Series, set, methods.Aggregation(ns), seq.Last, seq.Now)
	}
}

// Collect for Asychronous Cumulative->Cumulative
func (p *statelessAsyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.CumulativeTemporality)

	for set, storage := range p.data {
		// Copy the underlying storage.
		appendSeries(&ioutput.Series, set, methods.Aggregation(storage), seq.Start, seq.Now)
	}
	// Reset the entire map.
	p.data = map[attribute.Set]*Storage{}
}

// Collect for Asynchronous Cumulative->Delta
func (p *statefulAsyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.DeltaTemporality)

	for set, storage := range p.data {
		pval, has := p.prior[set]
		if has {
			// This does `*pval := *storage - *pval`
			methods.SubtractSwap(storage, pval)

			// Skip the series if it has not changed.
			if !methods.HasChange(pval) {
				continue
			}
			// Output the difference except for Gauge, in
			// which case output the new value.
			if p.desc.Kind.HasTemporality() {
				storage = pval
			}
		}

		appendSeries(&ioutput.Series, set, methods.Aggregation(storage), seq.Last, seq.Now)
	}
	// Copy the current to the prior and reset.
	p.prior = p.data
	p.data = map[attribute.Set]*Storage{}
}
