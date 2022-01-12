package viewstate

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/views"
)

type (
	Collector interface {
		// @@@
		Update(number number.Number, desc *sdkapi.Descriptor)
		Send(desc *sdkapi.Descriptor) error
	}

	CollectorFactory interface {
		// New constructs a new collector.  This is not
		// allowed to modify input, beacuse the sync codepath
		// has fingerprinted them as-is.
		New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector
	}

	State struct {
		// configuration

		hasDefault  bool
		library     instrumentation.Library
		definitions []views.View

		// state

		lock        sync.Mutex
		outputNames map[string]struct{}
	}

	// vCF is configured one per instrument with all
	// pre-calculated view behaviors.
	viewCollectorFactory struct {
		state         *State
		configuration []viewConfiguration
	}

	viewConfiguration func() viewCollector

	viewBehavior struct {
		// copied out of the configuration struct
		// @@@ name, aggregation kind, temporality choice, etc.
	}

	viewCollector interface {
		// @@@ [N number.Any, Traits traits.Any[N]]
	}

	aggregatorSettings struct {
		kind  aggregation.Kind
		hcfg  histogram.Config
		scfg  sum.Config
		lvcfg lastvalue.Config
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

func New(lib instrumentation.Library, defs []views.View, hasDefault bool) *State {

	// TODO: error checking here, such as:
	// - empty (?)
	// - duplicate name
	// - invalid inst/number/aggregation kind
	// - both instrument name and regexp
	// - schemaURL or Version without library name
	// - empty attribute keys
	// - Name w/o SingleInst

	return &State{
		definitions: defs,
		library:     lib,
		hasDefault:  hasDefault,
		outputNames: map[string]struct{}{},
	}
}

func configViewBehavior(v views.View) viewBehavior {
	return viewBehavior{
		// @@@
	}
}

func defaultViewBehavior(desc sdkapi.Descriptor) viewBehavior {
	return viewBehavior{
		// @@@
	}
}

func (vb viewBehavior) Name() string {
	// @@@
	return ""
}

// NewFactory is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *State) NewFactory(desc sdkapi.Descriptor) (CollectorFactory, error) {
	// Compute the set of matching views.
	type settingsBehaviors struct {
		settings  aggregatorSettings
		behaviors []viewBehavior
	}
	allBehaviors := map[string]settingsBehaviors{}
	for _, def := range v.definitions {
		if !def.Matches(v.library, desc) {
			continue
		}

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

		ss := fmt.Sprint(as)
		allBehaviors[ss] = settingsBehaviors{
			settings:  as,
			behaviors: append(allBehaviors[ss].behaviors, configViewBehavior(def)),
		}
	}
	// If there were no matching views, set the default aggregation.
	if len(allBehaviors) == 0 {
		if !v.hasDefault {
			return nil, nil
		}

		as := aggregatorSettingsFor(desc)
		allBehaviors[fmt.Sprint(as)] = settingsBehaviors{
			settings:  as,
			behaviors: []viewBehavior{defaultViewBehavior(desc)},
		}
	}

	v.lock.Lock()
	defer v.lock.Unlock()

	addedHere := map[string]struct{}{}

	for _, sbs := range allBehaviors {
		for _, behavior := range sbs.behaviors {
			outputName := behavior.Name()
			_, has1 := v.outputNames[outputName]
			_, has2 := addedHere[outputName]
			if has1 || has2 {
				return nil, fmt.Errorf("duplicate view name configured: ", outputName)
			}
			addedHere[outputName] = struct{}{}
		}
	}

	vcf := &viewCollectorFactory{state: v}

	for _, sbs := range allBehaviors {
		var cfg viewConfiguration
		switch desc.NumberKind() {
		case number.Int64Kind:
			cfg = buildView[int64, traits.Int64](desc, sbs.settings, sbs.behaviors)
		case number.Float64Kind:
			cfg = buildView[float64, traits.Float64](desc, sbs.settings, sbs.behaviors)
		}
		vcf.configuration = append(vcf.configuration, cfg)
		for _, behavior := range sbs.behaviors {
			v.outputNames[behavior.Name()] = struct{}{}
		}
	}

	return vcf, nil
}

func histogramDefaultsFor(kind number.Kind) histogram.Defaults {
	if kind == number.Int64Kind {
		return histogram.Int64Defaults{}
	}
	return histogram.Float64Defaults{}
}

func buildView[N number.Any, Traits traits.Any[N]](desc sdkapi.Descriptor, settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	if desc.InstrumentKind().Synchronous() {
		return buildSyncView[N, Traits](settings, behaviors)
	}
	return buildAsyncView[N, Traits](settings, behaviors)
}

func buildSyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	switch settings.kind {
	case aggregation.LastValueKind:
		return func() viewCollector {
			aa := syncCollector[N, *lastvalue.Aggregator[N, Traits], lastvalue.Config]{}
			aa.Init(settings.lvcfg)
			return aa
		}
	case aggregation.HistogramKind:
		return func() viewCollector {
			aa := syncCollector[N, *histogram.Aggregator[N, Traits], histogram.Config]{}
			aa.Init(settings.hcfg)
			return aa
		}
	default:
		return func() viewCollector {
			aa := syncCollector[N, *sum.Aggregator[N, Traits], sum.Config]{}
			aa.Init(settings.scfg)
			return aa
		}
	}
}

func buildAsyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	switch settings.kind {
	case aggregation.LastValueKind:
		return func() viewCollector {
			aa := asyncCollector[N, *lastvalue.Aggregator[N, Traits], lastvalue.Config]{}
			aa.Init(settings.lvcfg)
			return aa
		}
	case aggregation.HistogramKind:
		return func() viewCollector {
			aa := asyncCollector[N, *lastvalue.Aggregator[N, Traits], histogram.Config]{}
			aa.Init(settings.hcfg)
			return aa
		}
	default:
		return func() viewCollector {
			aa := asyncCollector[N, *sum.Aggregator[N, Traits], sum.Config]{}
			aa.Init(settings.scfg)
			return aa
		}
	}
}

func (factory *viewCollectorFactory) New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector {
	states := make([]viewCollector, 0, len(factory.configuration))
	for _, vc := range factory.configuration {
		states = append(states, vc())
	}
	return &viewCollector{
		states: states,
	}
}

func (v *viewCollector) Update(number number.Number, desc *sdkapi.Descriptor) {
	for _, state := range v.states {
		state.update(number, desc)
	}
}

func (v *viewCollector) Send(desc *sdkapi.Descriptor) error {

	// for _, output := range v.outputs {
	// 	output(desc)
	// }

	return nil
}

type syncCollector[N number.Any, Agg aggregator.Any[N], Config any] struct {
	current  Agg
	snapshot Agg
}

func (sc *syncCollector[N, Agg, Config]) Init(cfg Config) {
	sc.current.Init(cfg)
	sc.snapshot.Init(cfg)
}

type asyncCollector[N number.Any, Agg aggregator.Any[N], Config any] struct {
	current  N
	snapshot Agg
}

func (ac *asyncCollector[N, Agg, Config]) Init(cfg Config) {
	ac.current = 0
	ac.snapshot.Init(cfg)
}
