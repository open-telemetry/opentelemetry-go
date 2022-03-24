package viewstate

import (
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
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

	Sequence struct {
		Number int64
		Start  time.Time
		Last   time.Time
		Now    time.Time
	}

	Instrument interface {
		// NewAccumulator returns a new Accumulator bound to
		// the attributes `kvs`.  If reader == nil the
		// accumulator applies to all readers, otherwise it
		// applies to the specific reader.
		NewAccumulator(kvs []attribute.KeyValue, reader *reader.Reader) Accumulator

		// Collect transfers aggregated data from the
		// Accumulators into the output struct.
		Collect(reader *reader.Reader, sequence Sequence, output *[]reader.Series)
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
	}

	aggregatorSettings struct {
		kind  aggregation.Kind
		hcfg  histogram.Config
		scfg  sum.Config
		lvcfg lastvalue.Config
	}

	viewMetric[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		lock sync.Mutex
		desc sdkapi.Descriptor
		data map[attribute.Set]*Storage
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

	compiledSyncView[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		*viewMetric[N, Storage, Config, Methods]

		aggConfig *Config
		viewKeys  attribute.Filter
	}

	compiledAsyncView[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		*viewMetric[N, Storage, Config, Methods]

		aggConfig *Config
		viewKeys  attribute.Filter
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

	for readerIdx, reader := range v.readers {
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
				as = aggregatorSettingsFor(instrument, reader.Defaults())
			}

			if as.kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				desc:     instrument,
				view:     view,
				settings: as,
			})
		}

		// If there were no matching views, set the default aggregation.
		if len(matches) == 0 {
			as := aggregatorSettingsFor(instrument, reader.Defaults())
			if as.kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				desc:     instrument,
				view:     views.New(views.WithAggregation(as.kind)),
				settings: as,
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
	return &compiledSyncView[N, Storage, Config, Methods]{
		viewMetric: &viewMetric[N, Storage, Config, Methods]{
			desc: config.desc,
			data: map[attribute.Set]*Storage{},
		},
		aggConfig: aggConfig,
		viewKeys:  config.view.Keys(),
	}
}

func newAsyncView[
	N number.Any,
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](config configuredBehavior, aggConfig *Config) Instrument {
	return &compiledAsyncView[N, Storage, Config, Methods]{
		viewMetric: &viewMetric[N, Storage, Config, Methods]{
			desc: config.desc,
			data: map[attribute.Set]*Storage{},
		},
		aggConfig: aggConfig,
		viewKeys:  config.view.Keys(),
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

// NewAccumulator returns a Accumulator for a synchronous instrument view.
func (csv *compiledSyncView[N, Storage, Config, Methods]) NewAccumulator(kvs []attribute.KeyValue, _ *reader.Reader) Accumulator {
	sc := &syncAccumulator[N, Storage, Config, Methods]{}
	sc.init(csv.viewMetric, *csv.aggConfig, csv.viewKeys, kvs)
	return sc
}

// NewAccumulator returns a Accumulator for an asynchronous instrument view.
func (cav *compiledAsyncView[N, Storage, Config, Methods]) NewAccumulator(kvs []attribute.KeyValue, _ *reader.Reader) Accumulator {
	sc := &asyncAccumulator[N, Storage, Config, Methods]{}
	sc.init(cav.viewMetric, *cav.aggConfig, cav.viewKeys, kvs)
	return sc
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

func (mi multiInstrument[N]) Collect(reader *reader.Reader, sequence Sequence, output *[]reader.Series) {
	for _, inst := range mi[reader] {
		inst.Collect(reader, sequence, output)
	}
}

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

func (sc *syncAccumulator[N, Storage, Config, Methods]) init(metric *viewMetric[N, Storage, Config, Methods], cfg Config, keys attribute.Filter, kvs []attribute.KeyValue) {
	var methods Methods
	methods.Init(&sc.current, cfg)
	methods.Init(&sc.snapshot, cfg)

	sc.output = metric.findOutput(cfg, keys, kvs)
}

func (sc *syncAccumulator[N, Storage, Config, Methods]) Update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncAccumulator[N, Storage, Config, Methods]) Accumulate() {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)
	methods.Merge(&sc.snapshot, sc.output)
}

func (ac *asyncAccumulator[N, Storage, Config, Methods]) init(metric *viewMetric[N, Storage, Config, Methods], cfg Config, keys attribute.Filter, kvs []attribute.KeyValue) {
	var methods Methods
	methods.Init(&ac.snapshot, cfg)
	ac.current = 0
	ac.output = metric.findOutput(cfg, keys, kvs)
}

func (ac *asyncAccumulator[N, Storage, Config, Methods]) Update(number N) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	ac.current = number
}

func (ac *asyncAccumulator[N, Storage, Config, Methods]) Accumulate() {
	ac.lock.Lock()
	defer ac.lock.Unlock()

	var methods Methods
	methods.SynchronizedMove(&ac.snapshot, nil)
	methods.Update(&ac.snapshot, ac.current)
	ac.current = 0

	methods.Merge(&ac.snapshot, ac.output)
}

func (metric *viewMetric[N, Storage, Config, Methods]) findOutput(
	cfg Config,
	viewKeys attribute.Filter,
	kvs []attribute.KeyValue,
) *Storage {
	set, _ := attribute.NewSetWithFiltered(kvs, viewKeys)

	metric.lock.Lock()
	defer metric.lock.Unlock()

	storage, has := metric.data[set]
	if has {
		return storage
	}

	ns := new(Storage)
	var methods Methods
	methods.Init(ns, cfg)
	return ns
}

func (metric *viewMetric[N, Storage, Config, Methods]) Descriptor() sdkapi.Descriptor {
	return metric.desc
}

func (metric *viewMetric[N, Storage, Config, Methods]) Collect(_ *reader.Reader, sequence Sequence, output *[]reader.Series) {
	var methods Methods
	metric.lock.Lock()
	defer metric.lock.Unlock()

	for set, storage := range metric.data {
		*output = append(*output, reader.Series{
			Attributes:  set,
			Aggregation: methods.Aggregation(storage),
			Start:       sequence.Start,
			End:         sequence.Now, // @@@
		})
	}
}
