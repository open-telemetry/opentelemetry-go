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

	viewBehavior interface {
		terminal() *viewTerminal
		descriptor() *sdkapi.Descriptor
	}

	configureBehavior struct {
		desc sdkapi.Descriptor
		term *viewTerminal
		view *views.View
	}

	specificBehavior[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]] struct {
		desc sdkapi.Descriptor
		term *viewTerminal
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

var defaultViews = []views.View{
	views.New(views.WithAggregation("histogram")), // Histogram
	views.New(views.WithAggregation("gauge")),     // Gauge
	views.New(views.WithAggregation("sum")),       // Counter
	views.New(views.WithAggregation("sum")),       // UpDownCounter
	views.New(views.WithAggregation("sum")),       // AsyncCounter
	views.New(views.WithAggregation("sum")),       // AsyncUpDownCounter
}

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

func viewDescriptor(instrument sdkapi.Descriptor, v *views.View) sdkapi.Descriptor {
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

func (vb *specificBehavior[N, Storage, Config, Methods]) terminal() *viewTerminal {
	return vb.term
}

func (vb *specificBehavior[N, Storage, Config, Methods]) descriptor() *sdkapi.Descriptor {
	return &vb.desc
}

// NewFactory is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *State) NewFactory(instrument sdkapi.Descriptor) (*Factory, error) {
	// Compute the set of matching views.
	type (
		settingsBehaviors struct {
			settings  aggregatorSettings
			behaviors []configureBehavior
		}
	)

	allBehaviors := map[string]settingsBehaviors{}
	addBehavior := func(readerIdx int, as aggregatorSettings, behavior configureBehavior) {
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
			view := &v.terminals[terminalIdx].reader.Views()[viewIdx]
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
				as = aggregatorSettingsFor(instrument)
			}

			addBehavior(terminalIdx, as, configureBehavior{
				term: terminal,
				view: view,
				desc: viewDescriptor(instrument, view),
			})
		}

		// If there were no matching views, set the default aggregation.
		if matchCount == 0 {
			if !terminal.reader.HasDefaultView() {
				continue
			}

			addBehavior(terminalIdx, aggregatorSettingsFor(instrument), configureBehavior{
				term: terminal,
				view: &defaultViews[instrument.InstrumentKind()],
				desc: instrument,
			})
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
			if _, has := behavior.term.outputNames[behavior.desc.Name()]; !has {
				behavior.term.outputNames[behavior.desc.Name()] = struct{}{}
				valid++
			} else {
				otel.Handle(fmt.Errorf("duplicate view name registered"))
			}
		}
		if valid == 0 {
			continue
		}

		var compiled compiledView
		switch instrument.NumberKind() {
		case number.Int64Kind:
			compiled = buildView[int64, traits.Int64](instrument, sbs.settings, sbs.behaviors)
		case number.Float64Kind:
			compiled = buildView[float64, traits.Float64](instrument, sbs.settings, sbs.behaviors)
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

func buildView[N number.Any, Traits traits.Any[N]](instrument sdkapi.Descriptor, settings aggregatorSettings, behaviors []configureBehavior) compiledView {
	if instrument.InstrumentKind().Synchronous() {
		return buildSyncView[N, Traits](settings, behaviors)
	}
	return buildAsyncView[N, Traits](settings, behaviors)
}

func compileBehaviors[
	N number.Any,
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](behaviors []configureBehavior) []viewBehavior {
	output := make([]viewBehavior, len(behaviors))
	for i, behavior := range behaviors {
		output[i] = &specificBehavior[N, Storage, Config, Methods]{
			term: behavior.term,
			desc: behavior.desc,
		}
	}
	return output
}

func newSyncConfig[
	N number.Any,
	Traits traits.Any[N],
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](behaviors []configureBehavior, cfg *Config) compiledView {
	return compiledView{
		behaviors: compileBehaviors[N, Storage, Config, Methods](behaviors),
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
	Storage, Config any,
	Methods aggregator.Methods[N, Storage, Config],
](behaviors []configureBehavior, cfg *Config) compiledView {
	return compiledView{
		behaviors: compileBehaviors[N, Storage, Config, Methods](behaviors),
		newIntermediate: func() viewIntermediate {
			aa := &asyncCollector[N, Methods, Storage, Config]{}
			aa.Init(*cfg)
			return aa
		},
	}
}

func buildSyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []configureBehavior) compiledView {
	switch settings.kind {
	case aggregation.LastValueKind:
		return newSyncConfig[
			N,
			Traits,
			lastvalue.State[N, Traits],
			lastvalue.Config,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
		](behaviors, &settings.lvcfg)
	case aggregation.HistogramKind:
		return newSyncConfig[
			N,
			Traits,
			histogram.State[N, Traits],
			histogram.Config,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](behaviors, &settings.hcfg)
	default:
		return newSyncConfig[
			N,
			Traits,
			sum.State[N, Traits],
			sum.Config,
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](behaviors, &settings.scfg)
	}
}

func buildAsyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []configureBehavior) compiledView {
	switch settings.kind {
	case aggregation.LastValueKind:
		return newAsyncConfig[
			N,
			Traits,
			lastvalue.State[N, Traits],
			lastvalue.Config,
			lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]],
		](behaviors, &settings.lvcfg)
	case aggregation.HistogramKind:
		return newAsyncConfig[
			N,
			Traits,
			histogram.State[N, Traits],
			histogram.Config,
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](behaviors, &settings.hcfg)
	default:
		return newAsyncConfig[
			N,
			Traits,
			sum.State[N, Traits],
			sum.Config,
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](behaviors, &settings.scfg)
	}
}

// New returns a Collector that also implements Updater[N]
func (factory *Factory) New(kvs []attribute.KeyValue, instrument *sdkapi.Descriptor) Collector {
	if instrument.NumberKind() == number.Float64Kind {
		return newCollector[float64](factory, kvs, instrument)
	}
	return newCollector[int64](factory, kvs, instrument)
}

func newCollector[N number.Any](factory *Factory, kvs []attribute.KeyValue, instrument *sdkapi.Descriptor) CollectorUpdater[N] {
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
	for _, behavior := range vc.behaviors {
		behavior.terminal().reader.Process(behavior.descriptor(), agg)
	}

	return nil
}
