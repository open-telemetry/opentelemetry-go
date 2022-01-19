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
		Send(*Factory) error
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
		terminals []viewTerminal
	}

	Factory struct {
		state         *State
		configuration []compiledView
	}

	compiledView struct {
		behaviors       []viewBehavior
		newIntermediate func() viewIntermediate
	}

	viewBehavior struct {
		descriptor  sdkapi.Descriptor
		newReceiver func()
		terminal    *viewTerminal
	}

	viewIntermediate interface {
		Send(*compiledView) error
	}

	viewIntermediateUpdater[N number.Any] interface {
		viewIntermediate
		Updater[N]
	}
	
	viewCollector[N number.Any] struct {
		intermediates []viewIntermediateUpdater[N]
	}

	viewTerminal struct {
		reader      *reader.Reader
		lock        sync.Mutex
		outputNames map[string]struct{}
	}

	aggregatorSettings struct {
		kind  aggregation.Kind
		hcfg  histogram.Config
		scfg  sum.Config
		lvcfg lastvalue.Config
	}

	syncCollector[N number.Any, Methods aggregator.Methods[N, Storage, Config], Storage, Config any] struct {
		current  Storage
		snapshot Storage
	}

	asyncCollector[N number.Any, Methods aggregator.Methods[N, Storage, Config], Storage, Config any] struct {
		current  N
		snapshot Storage
	}
)

func aggregatorSettingsFor(desc sdkapi.Descriptor) aggregatorSettings {
	switch desc.InstrumentKind() {
	case sdkapi.HistogramInstrumentKind:
		return aggregatorSettings{
			kind: aggregation.HistogramKind,
		}
	case sdkapi.GaugeObserverInstrumentKind:
		return aggregatorSettings{
			kind: aggregation.LastValueKind,
		}
	default:
		return aggregatorSettings{
			kind: aggregation.SumKind,
		}
	}
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
	terminals := make([]viewTerminal, len(readerConfig))
	for i, r := range readerConfig {
		terminals[i].outputNames = map[string]struct{}{}
		terminals[i].reader = r
	}

	return &State{
		library:   lib,
		terminals: terminals,
	}
}

func configViewBehavior(instrument sdkapi.Descriptor, v *views.View, r *viewTerminal) viewBehavior {
	ikind := instrument.InstrumentKind()
	nkind := instrument.NumberKind()
	name := instrument.Name()
	desc := instrument.Description()
	unit := instrument.Unit()
	if v.HasName() {
		name = v.Name()
	}
	if v.Description() != "" {
		desc = instrument.Description()
	}
	return viewBehavior{
		terminal:   r,
		descriptor: sdkapi.NewDescriptor(name, ikind, nkind, desc, unit),
	}
}

func defaultViewBehavior(desc sdkapi.Descriptor, r *viewTerminal) viewBehavior {
	return viewBehavior{
		descriptor: desc,
		terminal:   r,
	}
}

func (vb *viewBehavior) Name() string {
	return vb.descriptor.Name()
}

// NewFactory is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *State) NewFactory(desc sdkapi.Descriptor) (*Factory, error) {
	// Compute the set of matching views.
	type settingsBehaviors struct {
		settings  aggregatorSettings
		behaviors []viewBehavior
	}

	allBehaviors := map[string]settingsBehaviors{}
	addBehavior := func(readerIdx int, as aggregatorSettings, behavior viewBehavior) {
		ss := fmt.Sprint(as)
		allBehaviors[ss] = settingsBehaviors{
			settings:  as,
			behaviors: append(allBehaviors[ss].behaviors, behavior),
		}

	}

	for terminalIdx := range v.terminals {
		matchCount := 0
		terminal := &v.terminals[terminalIdx]
		for viewIdx := range v.terminals[terminalIdx].reader.Views() {
			def := &v.terminals[terminalIdx].reader.Views()[viewIdx]
			if !def.Matches(v.library, desc) {
				continue
			}
			matchCount++
			var as aggregatorSettings
			switch def.Aggregation() {
			case aggregation.SumKind, aggregation.LastValueKind:
				// These have no options
				as.kind = def.Aggregation()
			case aggregation.HistogramKind:
				as.kind = def.Aggregation()
				as.hcfg = histogram.NewConfig(
					histogramDefaultsFor(desc.NumberKind()),
					def.HistogramOptions()...,
				)
			default:
				as = aggregatorSettingsFor(desc)
			}

			addBehavior(terminalIdx, as, configViewBehavior(desc, def, terminal))
		}

		// If there were no matching views, set the default aggregation.
		if matchCount == 0 {
			if !terminal.reader.HasDefaultView() {
				continue
			}

			addBehavior(terminalIdx, aggregatorSettingsFor(desc), defaultViewBehavior(desc, terminal))
		}
	}
	// When there are no matches for any terminal, return a nil factory.
	if len(allBehaviors) == 0 {
		return nil, nil
	}

	vcf := &Factory{state: v}

	for _, terminal := range v.terminals {
		terminal.lock.Lock()
		defer terminal.lock.Unlock()
	}

	for _, sbs := range allBehaviors {
		valid := 0
		for _, behavior := range sbs.behaviors {
			if _, has := behavior.terminal.outputNames[behavior.Name()]; !has {
				behavior.terminal.outputNames[behavior.Name()] = struct{}{}
				valid++
			} else {
				otel.Handle(fmt.Errorf("duplicate view name registered"))
			}
		}
		if valid == 0 {
			continue
		}

		var compiled compiledView
		switch desc.NumberKind() {
		case number.Int64Kind:
			compiled = buildView[int64, traits.Int64](desc, sbs.settings, sbs.behaviors)
		case number.Float64Kind:
			compiled = buildView[float64, traits.Float64](desc, sbs.settings, sbs.behaviors)
		}
		vcf.configuration = append(vcf.configuration, compiled)
	}

	if len(vcf.configuration) == 0 {
		return nil, nil
	}

	return vcf, nil
}

func histogramDefaultsFor(kind number.Kind) histogram.Defaults {
	if kind == number.Int64Kind {
		return histogram.Int64Defaults{}
	}
	return histogram.Float64Defaults{}
}

func buildView[N number.Any, Traits traits.Any[N]](desc sdkapi.Descriptor, settings aggregatorSettings, behaviors []viewBehavior) compiledView {
	if desc.InstrumentKind().Synchronous() {
		return buildSyncView[N, Traits](settings, behaviors)
	}
	return buildAsyncView[N, Traits](settings, behaviors)
}

func newSyncConfig[
	N number.Any,
	Traits traits.Any[N],
	Methods aggregator.Methods[N, Storage, Config],
	Storage, Config any,
](behaviors []viewBehavior, cfg *Config) compiledView {

	// behavior is only partly specified. 
	// for i := range behaviors {
	// 	behaviors[i]. 
	// }@ @@@
	return compiledView{
		behaviors: behaviors,
		newIntermediate: func() viewIntermediate {
			aa := &syncCollector[N, Methods, Storage, Config]{}
			aa.Init(*cfg)
			return aa
		},
	}
}

func newAsyncConfig[
	N number.Any,
	Traits traits.Any[N],
	Methods aggregator.Methods[N, Storage, Config],
	Storage, Config any,
](behaviors []viewBehavior, cfg *Config) compiledView {
	return compiledView{
		behaviors: behaviors,
		newIntermediate: func() viewIntermediate {
			aa := &asyncCollector[N, Methods, Storage, Config]{}
			aa.Init(*cfg)
			return aa
		},
	}
}

func buildSyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) compiledView {
	switch settings.kind {
	case aggregation.LastValueKind:
		return newSyncConfig[
			N,
			Traits,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
			lastvalue.State[N, Traits],
			lastvalue.Config,
		](behaviors, &settings.lvcfg)
	case aggregation.HistogramKind:
		return newSyncConfig[
			N,
			Traits,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
			histogram.State[N, Traits],
			histogram.Config,
		](behaviors, &settings.hcfg)
	default:
		return newSyncConfig[
			N,
			Traits,
			sum.Methods[N, Traits, sum.State[N, Traits]],
			sum.State[N, Traits],
			sum.Config,
		](behaviors, &settings.scfg)
	}
}

func buildAsyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) compiledView {
	switch settings.kind {
	case aggregation.LastValueKind:
		return newAsyncConfig[
			N,
			Traits,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
			lastvalue.State[N, Traits],
			lastvalue.Config,
		](behaviors, &settings.lvcfg)
	case aggregation.HistogramKind:
		return newAsyncConfig[
			N,
			Traits,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
			histogram.State[N, Traits],
			histogram.Config,
		](behaviors, &settings.hcfg)
	default:
		return newAsyncConfig[
			N,
			Traits,
			sum.Methods[N, Traits, sum.State[N, Traits]],
			sum.State[N, Traits],
			sum.Config,
		](behaviors, &settings.scfg)
	}
}

// New returns a Collector that also implements Updater[N]
func (factory *Factory) New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector {
	if desc.NumberKind() == number.Float64Kind {
		return newCollector[float64](factory, kvs, desc)
	}
	return newCollector[int64](factory, kvs, desc)
}

func newCollector[N number.Any](factory *Factory, kvs []attribute.KeyValue, desc *sdkapi.Descriptor) CollectorUpdater[N] {
	intermediates := make([]viewIntermediateUpdater[N], 0, len(factory.configuration))
	for idx, vc := range factory.configuration {
		intermediates[idx] = vc.newIntermediate().(viewIntermediateUpdater[N])
	}
	return &viewCollector[N]{
		intermediates: intermediates,
	}
}

func (c *viewCollector[N]) Send(factory *Factory) error {
	for i, intermediate := range c.intermediates {
		intermediate.Send(&factory.configuration[i])
	}
	return nil
}

func (c *viewCollector[N]) Update(value N) {
	for _, intermediate := range c.intermediates {
		intermediate.Update(value)
	}
}

func (sc *syncCollector[N, Methods, Storage, Config]) Init(cfg Config) {
	var methods Methods
	methods.Init(&sc.current, cfg)
	methods.Init(&sc.snapshot, cfg)
}

func (sc *syncCollector[N, Methods, Storage, Config]) Update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncCollector[N, Methods, Storage, Config]) Send(vc *compiledView) error {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)

	return vc.send(methods.Aggregation(&sc.snapshot))
}

func (ac *asyncCollector[N, Methods, Storage, Config]) Init(cfg Config) {
	var methods Methods
	ac.current = 0
	methods.Init(&ac.snapshot, cfg)
}

func (ac *asyncCollector[N, Methods, Storage, Config]) Update(number N) {
	ac.current = number
}

func (ac *asyncCollector[N, Methods, Storage, Config]) Send(vc *compiledView) error {
	var methods Methods
	methods.SynchronizedMove(&ac.snapshot, nil)
	methods.Update(&ac.snapshot, ac.current)
	ac.current = 0

	return vc.send(methods.Aggregation(&ac.snapshot))
}

func (vc *compiledView) send(agg aggregation.Aggregation) error {
	for i := range vc.behaviors {
		vc.behaviors[i].terminal.reader.Process(
			&vc.behaviors[i].descriptor,
			agg,
		)
	}

	return nil
}
