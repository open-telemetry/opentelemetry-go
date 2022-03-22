package viewstate

import (
	"fmt"
	"sync"

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
	}

	Instrument interface {
		// reader == nil for all readers, else 1 reader
		NewCollector(kvs []attribute.KeyValue, reader *reader.Reader) Collector
	}

	Collector interface {
		Collect()
	}

	Updater[N number.Any] interface {
		Update(value N)
	}

	CollectorUpdater[N number.Any] interface {
		Collector
		Updater[N]
	}

	multiInstrument[N number.Any] map[*reader.Reader][]Instrument

	multiCollector[N number.Any] []Collector

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

	syncCollector[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		current  Storage
		snapshot Storage
		output   *Storage
	}

	asyncCollector[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		lock     sync.Mutex
		current  N
		snapshot Storage
		output   *Storage
	}

	compiledSyncView[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		metric    *viewMetric[N, Storage, Config, Methods]
		aggConfig *Config
		viewKeys  attribute.Filter
	}

	compiledAsyncView[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		metric    *viewMetric[N, Storage, Config, Methods]
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
	return &Compiler{
		library: lib,
		views:   views,
		readers: readers,
	}
}

// Compile is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *Compiler) Compile(instrument sdkapi.Descriptor) Instrument {
	configs := map[*reader.Reader][]configuredBehavior{}

	for _, reader := range v.readers {
		matchCount := 0
		for _, view := range v.views {
			if !view.Matches(v.library, instrument) {
				continue
			}
			matchCount++
			var as aggregatorSettings
			switch view.Aggregation() {
			case aggregation.SumKind, aggregation.LastValueKind:
				// These have no options
				as.kind = view.Aggregation()
			case aggregation.HistogramKind:
				as.kind = view.Aggregation()
				as.hcfg = histogram.NewConfig(
					histogramDefaultsFor(instrument.NumberKind()),
					view.HistogramOptions()...,
				)
			default:
				as = aggregatorSettingsFor(instrument, reader.Defaults())
			}

			if as.kind == aggregation.DropKind {
				continue
			}

			configs[reader] = append(configs[reader], configuredBehavior{
				desc:     instrument,
				view:     view,
				settings: as,
			})
		}

		// If there were no matching views, set the default aggregation.
		if matchCount == 0 {
			as := aggregatorSettingsFor(instrument, reader.Defaults())
			if as.kind == aggregation.DropKind {
				continue
			}

			configs[reader] = append(configs[reader], configuredBehavior{
				desc:     instrument,
				view:     views.New(views.WithAggregation(as.kind)),
				settings: as,
			})
		}
	}

	compiled := map[*reader.Reader][]Instrument{}

	for reader, list := range configs {
		for _, config := range list {
			config.desc = viewDescriptor(config.desc, config.view)

			var one Instrument
			switch config.desc.NumberKind() {
			case number.Int64Kind:
				one = buildView[int64, traits.Int64](config)
			case number.Float64Kind:
				one = buildView[float64, traits.Float64](config)
			}

			if available := reader.AcquireOutput(config.desc); !available {
				otel.Handle(fmt.Errorf("duplicate view name registered"))
				continue
			}

			compiled[reader] = append(compiled[reader], one)
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
		metric: &viewMetric[N, Storage, Config, Methods]{
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
](config configuredBehavior, aggConfig *Config) Instrument{
	return &compiledAsyncView[N, Storage, Config, Methods]{
		metric: &viewMetric[N, Storage, Config, Methods]{
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

// NewCollector returns a Collector for a synchronous instrument view.
func (csv *compiledSyncView[N, Storage, Config, Methods]) NewCollector(kvs []attribute.KeyValue, _ *reader.Reader) Collector {
	sc := &syncCollector[N, Storage, Config, Methods]{}
	sc.init(csv.metric, *csv.aggConfig, csv.viewKeys, kvs)
	return sc
}

// NewCollector returns a Collector for an asynchronous instrument view.
func (cav *compiledAsyncView[N, Storage, Config, Methods]) NewCollector(kvs []attribute.KeyValue, _ *reader.Reader) Collector {
	sc := &asyncCollector[N, Storage, Config, Methods]{}
	sc.init(cav.metric, *cav.aggConfig, cav.viewKeys, kvs)
	return sc
}

// NewCollector returns a Collector for multiple views of the same instrument.
func (mi multiInstrument[N]) NewCollector(kvs []attribute.KeyValue, reader *reader.Reader) Collector {
	var collectors []Collector
	// Note: This runtime switch happens because we're using the same API for
	// both async and sync instruments, whereas the APIs are not symmetrical.
	if reader == nil {
		for _, list := range mi {
			collectors = make([]Collector, 0, len(mi) * len(list))
		}
		for _, list := range mi {
			for _, inst := range list {
				collectors = append(collectors, inst.NewCollector(kvs, nil))
			}
		}
	} else {
		insts := mi[reader]

		collectors = make([]Collector, 0, len(insts))

		for _, inst := range insts{
			collectors = append(collectors, inst.NewCollector(kvs, reader))
		}
	}
	return multiCollector[N](collectors)
}

func (c multiCollector[N]) Collect() {
	for _, coll := range c {
		coll.Collect()
	}
}

func (c multiCollector[N]) Update(value N) {
	for _, coll := range c {
		coll.(Updater[N]).Update(value)
	}
}

func (sc *syncCollector[N, Storage, Config, Methods]) init(metric *viewMetric[N, Storage, Config, Methods], cfg Config, keys attribute.Filter, kvs []attribute.KeyValue) {
	var methods Methods
	methods.Init(&sc.current, cfg)
	methods.Init(&sc.snapshot, cfg)

	sc.output = metric.findOutput(cfg, keys, kvs)
}

func (sc *syncCollector[N, Storage, Config, Methods]) Update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncCollector[N, Storage, Config, Methods]) Collect() {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)
	methods.Merge(&sc.snapshot, sc.output)
}

func (ac *asyncCollector[N, Storage, Config, Methods]) init(metric *viewMetric[N, Storage, Config, Methods], cfg Config, keys attribute.Filter, kvs []attribute.KeyValue) {
	var methods Methods
	methods.Init(&ac.snapshot, cfg)
	ac.current = 0
	ac.output = metric.findOutput(cfg, keys, kvs)
}

func (ac *asyncCollector[N, Storage, Config, Methods]) Update(number N) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	ac.current = number
}

func (ac *asyncCollector[N, Storage, Config, Methods]) Collect() {
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
