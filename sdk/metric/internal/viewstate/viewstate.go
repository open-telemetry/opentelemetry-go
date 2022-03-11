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

	State struct {
		library   instrumentation.Library
		terminals []*viewTerminal
	}

	Factory struct {
		compiled []compiledView
	}

	compiledView interface {
		newInstrument(kvs []attribute.KeyValue) viewInstrument
	}

	configuredBehavior struct {
		instrument sdkapi.Descriptor
		term       *viewTerminal
		view       views.View
		settings   aggregatorSettings
		metric     *viewMetric
	}

	aggregatorSettings struct {
		kind  aggregation.Kind
		hcfg  histogram.Config
		scfg  sum.Config
		lvcfg lastvalue.Config
	}

	viewMultiInstrument[N number.Any] struct {
		instruments []viewInstrumentUpdater[N]
	}

	viewInstrument interface {
		collect()
	}

	viewInstrumentUpdater[N number.Any] interface {
		viewInstrument
		update(value N)
	}

	viewTerminal struct {
		reader      *reader.Reader
		lock        sync.Mutex
		outputNames map[string]struct{}
		metrics     []*viewMetric
	}

	viewMetric struct {
		lock    sync.Mutex
		desc    sdkapi.Descriptor
		streams map[attribute.Set]aggregation.Aggregation
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
		metric    *viewMetric
		aggConfig *Config
		viewKeys  attribute.Filter
	}

	compiledAsyncView[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		metric    *viewMetric
		aggConfig *Config
		viewKeys  attribute.Filter
	}
)

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

func New(lib instrumentation.Library, readerConfig []*reader.Reader) *State {

	// TODO: error checking here, such as:
	// - empty (?)
	// - duplicate name
	// - invalid inst/number/aggregation kind
	// - both instrument name and regexp
	// - schemaURL or Version without library name
	// - empty attribute keys
	// - Name w/o SingleInst
	terminals := make([]*viewTerminal, len(readerConfig))
	for i, r := range readerConfig {
		terminals[i] = &viewTerminal{
			outputNames: map[string]struct{}{},
			reader:      r,
		}
	}

	return &State{
		library:   lib,
		terminals: terminals,
	}
}

// NewFactory is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *State) NewFactory(instrument sdkapi.Descriptor) *Factory {
	var configs []configuredBehavior

	for _, terminal := range v.terminals {
		matchCount := 0
		for _, view := range terminal.reader.Views() {
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
				as = aggregatorSettingsFor(instrument, terminal.reader.Defaults())
			}

			if as.kind == aggregation.DropKind {
				continue
			}

			configs = append(configs, configuredBehavior{
				instrument: instrument,
				term:       terminal,
				view:       view,
				settings:   as,
			})
		}

		// If there were no matching views, set the default aggregation.
		// TODO: If the `default_enabled` variable is disabled, continue
		if matchCount == 0 {
			as := aggregatorSettingsFor(instrument, terminal.reader.Defaults())
			if as.kind == aggregation.DropKind {
				continue
			}

			configs = append(configs, configuredBehavior{
				instrument: instrument,
				term:       terminal,
				view:       views.New(views.WithAggregation(as.kind)),
				settings:   as,
			})
		}
	}
	// When there are no matches for any terminal, return a nil factory.
	if len(configs) == 0 {
		return nil
	}

	vcf := &Factory{}

	for _, terminal := range v.terminals {
		terminal.lock.Lock()
		defer terminal.lock.Unlock()
	}

	for _, config := range configs {
		viewDesc := viewDescriptor(config.instrument, config.view)

		if _, has := config.term.outputNames[viewDesc.Name()]; has {
			otel.Handle(fmt.Errorf("duplicate view name registered"))
			continue
		}
		config.metric = &viewMetric{
			desc:   viewDesc,
			streams: map[attribute.Set]aggregation.Aggregation{},
		}
		config.term.outputNames[viewDesc.Name()] = struct{}{}
		config.term.metrics = append(config.term.metrics, config.metric)

		var compiled compiledView
		switch viewDesc.NumberKind() {
		case number.Int64Kind:
			compiled = buildView[int64, traits.Int64](config)
		case number.Float64Kind:
			compiled = buildView[float64, traits.Float64](config)
		}
		vcf.compiled = append(vcf.compiled, compiled)
	}

	if len(vcf.compiled) == 0 {
		return nil
	}

	return vcf
}

func histogramDefaultsFor(kind number.Kind) histogram.Defaults {
	if kind == number.Int64Kind {
		return histogram.Int64Defaults{}
	}
	return histogram.Float64Defaults{}
}

func buildView[N number.Any, Traits traits.Any[N]](config configuredBehavior) compiledView {
	if config.metric.desc.InstrumentKind().Synchronous() {
		return buildSyncView[N, Traits](config)
	}
	return buildAsyncView[N, Traits](config)
}

func newSyncConfig[
	N number.Any,
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](config configuredBehavior, aggConfig *Config) compiledView {
	return &compiledSyncView[N, Storage, Config, Methods]{
		metric:    config.metric,
		aggConfig: aggConfig,
		viewKeys:  config.view.Keys(),
	}
}

func (csv *compiledSyncView[N, Storage, Config, Methods]) newInstrument(kvs []attribute.KeyValue) viewInstrument {
	sc := &syncCollector[N, Storage, Config, Methods]{
		output: findOutput[N, Storage, Config, Methods](csv.metric, csv.aggConfig, csv.viewKeys, kvs),
	}
	sc.init(*csv.aggConfig)
	return sc
}

func newAsyncConfig[
	N number.Any,
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](config configuredBehavior, aggConfig *Config) compiledView {
	return &compiledAsyncView[N, Storage, Config, Methods]{
		metric:    config.metric,
		aggConfig: aggConfig,
		viewKeys:  config.view.Keys(),
	}
}

func (cav *compiledAsyncView[N, Storage, Config, Methods]) newInstrument(kvs []attribute.KeyValue) viewInstrument {
	sc := &asyncCollector[N, Storage, Config, Methods]{
		output: findOutput[N, Storage, Config, Methods](cav.metric, cav.aggConfig, cav.viewKeys, kvs),
	}
	sc.init(*cav.aggConfig)
	return sc
}

func buildSyncView[N number.Any, Traits traits.Any[N]](config configuredBehavior) compiledView {
	switch config.settings.kind {
	case aggregation.LastValueKind:
		return newSyncConfig[
			N,
			lastvalue.State[N, Traits],
			lastvalue.Config,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
		](config, &config.settings.lvcfg)
	case aggregation.HistogramKind:
		return newSyncConfig[
			N,
			histogram.State[N, Traits],
			histogram.Config,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config, &config.settings.hcfg)
	default:
		return newSyncConfig[
			N,
			sum.State[N, Traits],
			sum.Config,
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config, &config.settings.scfg)
	}
}

func buildAsyncView[N number.Any, Traits traits.Any[N]](config configuredBehavior) compiledView {
	switch config.settings.kind {
	case aggregation.LastValueKind:
		return newAsyncConfig[
			N,
			lastvalue.State[N, Traits],
			lastvalue.Config,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
		](config, &config.settings.lvcfg)
	case aggregation.HistogramKind:
		return newAsyncConfig[
			N,
			histogram.State[N, Traits],
			histogram.Config,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config, &config.settings.hcfg)
	default:
		return newAsyncConfig[
			N,
			sum.State[N, Traits],
			sum.Config,
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config, &config.settings.scfg)
	}
}

// New returns a Collector that also implements Updater[N]
func (factory *Factory) New(kvs []attribute.KeyValue, instrument *sdkapi.Descriptor) Collector {
	if instrument.NumberKind() == number.Float64Kind {
		return newMultiInstrument[float64](factory, kvs, instrument)
	}
	return newMultiInstrument[int64](factory, kvs, instrument)
}

func newMultiInstrument[N number.Any](factory *Factory, kvs []attribute.KeyValue, instrument *sdkapi.Descriptor) CollectorUpdater[N] {
	instruments := make([]viewInstrumentUpdater[N], 0, len(factory.compiled))
	for idx, vc := range factory.compiled {
		instruments[idx] = vc.newInstrument(kvs).(viewInstrumentUpdater[N])
	}
	return &viewMultiInstrument[N]{
		instruments: instruments,
	}
}

func (c *viewMultiInstrument[N]) Collect() {
	for _, intermediate := range c.instruments {
		intermediate.collect()
	}
}

func (c *viewMultiInstrument[N]) Update(value N) {
	for _, intermediate := range c.instruments {
		intermediate.update(value)
	}
}

func (sc *syncCollector[N, Storage, Config, Methods]) init(cfg Config) {
	var methods Methods
	methods.Init(&sc.current, cfg)
	methods.Init(&sc.snapshot, cfg)
}

func (sc *syncCollector[N, Storage, Config, Methods]) update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncCollector[N, Storage, Config, Methods]) collect() {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)
	methods.Merge(&sc.snapshot, sc.output)
}

func (ac *asyncCollector[N, Storage, Config, Methods]) init(cfg Config) {
	var methods Methods
	ac.current = 0
	methods.Init(&ac.snapshot, cfg)
}

func (ac *asyncCollector[N, Storage, Config, Methods]) update(number N) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	ac.current = number
}

func (ac *asyncCollector[N, Storage, Config, Methods]) collect() {
	ac.lock.Lock()
	defer ac.lock.Unlock()

	var methods Methods
	methods.SynchronizedMove(&ac.snapshot, nil)
	methods.Update(&ac.snapshot, ac.current)
	ac.current = 0

	methods.Merge(&ac.snapshot, ac.output)
}

func findOutput[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]](metric *viewMetric, aggConfig *Config, viewKeys attribute.Filter, kvs []attribute.KeyValue) *Storage {
	set, _ := attribute.NewSetWithFiltered(kvs, viewKeys)
	return findOrCreate[N, Storage, Config, Methods](
		metric,
		set,
		aggConfig,
	)
}

func findOrCreate[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]](
	metric *viewMetric,
	attrs attribute.Set,
	aggConfig *Config,
) *Storage {
	var methods Methods
	metric.lock.Lock()
	defer metric.lock.Unlock()

	if aggr, has := metric.streams[attrs]; has {
		return methods.Storage(aggr)
	}

	ns := new(Storage)
	methods.Init(ns, *aggConfig)
	metric.streams[attrs] = methods.Aggregation(ns)
	return ns
}
