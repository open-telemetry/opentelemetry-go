package viewstate

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
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

		lock           sync.Mutex
		namedFactories map[string]*viewCollectorFactory
	}

	// vCF is configured one per instrument with all
	// pre-calculated view behaviors.
	viewCollectorFactory struct {
		state         *State
		configuration []viewConfiguration
	}

	viewConfiguration struct {
		settings  aggregatorSettings
		behaviors []viewBehavior
	}

	viewBehavior struct {
		// copied out of the configuration struct
		// @@@ name, aggregation kind, temporality choice, etc.
	}

	// vC is returned by the factory New() method.  each label set
	// has one of these allocated with state for each matching view.
	viewCollector struct {
		states []viewConfigState
		// inputs []export.Aggregator
		// outputs []viewMatchState
	}

	viewConfigState interface {
		update()
		send()
	}

	aggregatorSettings struct {
		aggregation.Kind
	}
)

func aggregatorSettingsFor(desc sdkapi.Descriptor) aggregatorSettings {
	switch desc.InstrumentKind() {
	case sdkapi.HistogramInstrumentKind:
		return aggregatorSettings{
			Kind: aggregation.HistogramKind,
		}
	case sdkapi.GaugeObserverInstrumentKind:
		return aggregatorSettings{
			Kind: aggregation.LastValueKind,
		}
	default:
		return aggregatorSettings{
			Kind: aggregation.SumKind,
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
		definitions:    defs,
		library:        lib,
		hasDefault:     hasDefault,
		namedFactories: map[string]*viewCollectorFactory{},
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
	settingBehaviors := map[aggregatorSettings][]viewBehavior{}
	for _, def := range v.definitions {
		if !def.Matches(v.library, desc) {
			continue
		}

		// Note: aggregatorSettings is a stand-in for a
		// complete aggregator configuration, which currently
		// only includes the aggregation Kind.
		var as aggregatorSettings
		switch as.Kind {
		case aggregation.SumKind, aggregation.HistogramKind, aggregation.LastValueKind:
			as.Kind = def.Aggregation()
		default:
			as = aggregatorSettingsFor(desc)
		}

		settingBehaviors[as] = append(settingBehaviors[as], configViewBehavior(def))
	}
	// If there were no matching views, set the default aggregation.
	if len(settingBehaviors) == 0 {
		if !v.hasDefault {
			return nil, nil
		}

		as := aggregatorSettingsFor(desc)
		settingBehaviors[as] = append(settingBehaviors[as], defaultViewBehavior(desc))
	}

	// Form the configuration.
	var configuration []viewConfiguration
	for settings, behaviors := range settingBehaviors {
		configuration = append(configuration, viewConfiguration{
			settings:  settings,
			behaviors: behaviors,
		})
	}

	// Develop the list of basic aggregations.
	vcf := &viewCollectorFactory{
		state:         v,
		configuration: configuration,
	}

	// Reconcile.  Build a representation of the output names.
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, cfg := range configuration {
		for _, behavior := range cfg.behaviors {

			outputName := behavior.Name()
			if _, has := v.namedFactories[outputName]; has {
				return nil, fmt.Errorf("duplicate view name configured")
			}
			v.namedFactories[outputName] = vcf
		}
	}

	return vcf, nil
}

func (factory *viewCollectorFactory) New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector {
	var states []viewConfigState
	for _, config := range factory.configuration {
		// 1 is the input aggregator
		// 1 for each matching view
		cnt := 1 + len(config.behaviors)

		// 1 for each synchronous aggregator, as these need a temporary.
		if desc.InstrumentKind().Synchronous() {
			cnt++
		}

		switch config.settings.Kind {
		case aggregation.SumKind:
			// TODO: Add a type for []sum.Aggregator that will Merge()
			// from position 0 into position 1..N.
			// TODO: The two methods, Update() and Send(), are not
			// part of that interface.  Only Send() there ^^^.
			// ALSO: The Update() method here will be another aggregator
			// of the same type only for synchronous instruments.  For the
			// asynchronous instruments, the Update() method here will be
			// a dedicated asynchronous observer aggregator.
			//
			// ^^^ Thus we still need count + 1 for synchronous, and
			// only count number of aggregators to be allocated for
			// asynchronous.
			//
			// So all of this is about a memory optimization, still
			// worth it?  If yes, then we need one type per slice to
			// handle this Send() or we need some other kind of
			// reflection.  It will be much simpler if we allocate
			// aggregators independently.
			aggs := sum.New(cnt)

			states = append(states, &aggs[0])
			//outputs = append(outputs, viewMatchState{aggs})
		case aggregation.LastValueKind:
			aggs := lastvalue.New(cnt)
			inputs = append(inputs, &aggs[0])
			//outputs = append(outputs, viewMatchState{aggs})
		case aggregation.HistogramKind:
			aggs := histogram.New(cnt, desc)
			inputs = append(inputs, &aggs[0])
			//outputs = append(outputs, viewMatchState{aggs})
		}
	}
	return &viewCollector{
		inputs: inputs,
	}
}

func (v *viewCollector) Update(number number.Number, desc *sdkapi.Descriptor) {
	for _, input := range v.inputs {
		input.Update(number, desc)
	}
}

func (v *viewCollector) Send(desc *sdkapi.Descriptor) error {

	for _, output := range v.outputs {
		output(desc)
	}

	return nil
}
