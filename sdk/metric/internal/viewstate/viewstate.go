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

	viewConfiguration interface {
	}

	// viewConfigurationX struct {
	// 	settings  aggregatorSettings
	// 	behaviors []viewBehavior
	// }

	viewBehavior struct {
		// copied out of the configuration struct
		// @@@ name, aggregation kind, temporality choice, etc.
	}

	// vC is returned by the factory New() method.  each label set
	// has one of these allocated with state for each matching view.
	viewCollector struct {
		states []viewConfigState
	}

	viewSender interface {
		send(*sdkapi.Descriptor)
	}

	viewUpdater interface {
		Update(number.Number, *sdkapi.Descriptor)
	}

	viewConfigState interface {
		viewSender
		viewUpdater
	}

	aggregatorSettings struct {
		kind aggregation.Kind
		hcfg histogram.Config

	}
	
	syncCollector[N number.Any, Agg aggregator.Any[N]] struct {
		current  Agg
		snapshot Agg
	}

	asyncCollector[N number.Any, Agg aggregator.Any[N]] struct {
		current  N
		snapshot Agg
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
			settings: as,
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
			settings: as,
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

	for settings, behaviors := range allBehaviors {
		var cfg viewConfiguration
		switch desc.NumberKind() {
		case number.Int64Kind:
			cfg = buildViewConfiguration[int64, traits.Int64](desc, settings, behaviors)
		case number.Float64Kind:
			cfg = buildViewConfiguration[float64, traits.Float64](desc, settings, behaviors)
		}
		vcf.configuration = append(vcf.configuration, cfg)
		for _, behavior := range behaviors {
			v.outputNames[behavior.Name()] = struct{}{}
		}
	}

	return vcf, nil
}

func buildViewConfiguration[N number.Any, Traits traits.Any[N]](desc sdkapi.Descriptor, settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	if desc.InstrumentKind().Synchronous() {
		return buildSyncViewConfiguration[N, Traits](settings, behaviors)
	}
	return buildAsyncViewConfiguration[N, Traits](settings, behaviors)
}

func histogramDefaultsFor(kind number.Kind) histogram.Defaults {
	if kind == number.Int64Kind {
		return histogram.Int64Defaults{}
	}
	return histogram.Float64Defaults{}
}

func buildSyncViewConfiguration[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	switch settings.Kind {
	case aggregation.LastValueKind:
		aa := syncCollector[N, *lastvalue.Aggregator[N, Traits]]{}
		aa.current.Init(lastvalue.Config{})
		aa.snapshot.Init(lastvalue.Config{})
		return aa
	case aggregation.HistogramKind:
		aa := syncCollector[N, *histogram.Aggregator[N, Traits, histogram.Option]]{}
		aa.current.Init()  // ... options @@@
		aa.snapshot.Init() // ... options @@@
		return aa
	default:
		break
	}
	aa := syncCollector[N, sum.Aggregator[N, Traits]]{}
	aa.current.Init()
	aa.snapshot.Init()
	return aa
}

func buildAsyncViewConfiguration[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	switch settings.Kind {
	case aggregation.SumKind:
	case aggregation.LastValueKind:
	case aggregation.HistogramKind:
	}
}

func (factory *viewCollectorFactory) New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector {
	var updaters []viewUpdater
	for _, config := range factory.configuration {
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

// func (sa sumAggregators) send(desc *sdkapi.Descriptor) {
// 	for i := range sa[1:] {
// 		sa[i].Merge(&sa[0], desc)
// 	}
// }

// func (lva lastValueAggregators) send() {
// 	for i := range lva[1:] {
// 		lva[i].Merge(&lva[0], desc)
// 	}
// }

// func (ha histogramAggregators) send() {
// 	for i := range ha[1:] {
// 		ha[i].Merge(&ha[0], desc)
// 	}
// }
