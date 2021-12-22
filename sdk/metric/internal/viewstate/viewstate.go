package viewstate

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/views"
)

type (
	Collector interface {
		Update(number number.Number, descriptor *sdkapi.Descriptor)
		Send(final bool) error
	}

	CollectorFactory interface {
		// Note allowed to modify input, beacuse the sync
		// codepath has fingerprinted them as-is.
		New(kvs []attribute.KeyValue) Collector
	}

	State struct {
		// configuration

		hasDefault  bool
		library     instrumentation.Library
		definitions []views.View

		// state

		lock   sync.Mutex
		output map[string]*viewCollectorFactory
	}

	viewCollector struct {
		inputs []export.Aggregator
	}

	viewCollectorFactory struct {
		state   *State
		aggSet  map[aggregation.Kind]viewMatchState
		matches []views.View
	}

	viewMatchState struct {
		// @@@ something to encode the number of distinct pipelines
		// that receive aggregations of a particular kind.
	}
)

func aggregationKindFor(ik sdkapi.InstrumentKind) aggregation.Kind {
	switch ik {
	case sdkapi.HistogramInstrumentKind:
		return aggregation.HistogramKind
	case sdkapi.GaugeObserverInstrumentKind:
		return aggregation.LastValueKind
	default:
		return aggregation.SumKind
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
		output:      map[string]*viewCollectorFactory{},
	}
}

// NewFactory is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *State) NewFactory(desc sdkapi.Descriptor) (CollectorFactory, error) {
	// TODO: This should be conditioned on the configuration of
	// the aggregator/exporter too, but we're not configuring them yet.

	// Compute the set of matching views.
	var matches []views.View
	for _, def := range v.definitions {
		if def.Matches(v.library, desc) {
			matches = append(matches, def)
		}
	}

	// Compute the set of requested aggregations.
	aggSet := map[aggregation.Kind]int{}
	for _, match := range matches {
		va := match.Aggregation()
		if va == "" {
			aggSet[aggregationKindFor(desc.InstrumentKind())]++
		} else {
			aggSet[va]++
		}
	}
	// If there were no matching views, set the default aggregation.
	if matches == nil {
		if !v.hasDefault {
			return nil, nil
		}
		aggSet[aggregationKindFor(desc.InstrumentKind())] = struct{}{}
	}

	// Develop the list of basic aggregations.
	vcf := &viewCollectorFactory{
		state:   v,
		aggSet:  aggSet,
		matches: matches,
	}

	// Reconcile.  Build a representation of the output names.
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, match := range matches {
		var outputName string
		if match.HasName() {
			outputName = match.Name()
		} else {
			outputName = desc.Name()
		}

		if _, has := v.output[outputName]; has {
			return nil, fmt.Errorf("duplicate view name configured")
		}
		v.output[outputName] = vcf
	}

	return vcf, nil
}

func (v *viewCollectorFactory) New(kvs []attribute.KeyValue) Collector {
	var aggs []export.Aggregator
	for _, ak := range v.akinds {
		// Notes:
		//
		// Need to know how many aggregators are needed to construct
		// the output pipeline.
		//
		// 1 is the base case, for the collector itself
		// 1 for each matching view

		var agg export.Aggregator
		switch ak {
		case aggregation.SumKind:
		case aggregation.LastValueKind:
		case aggregation.HistogramKind:
		}
		aggs = append(aggs, agg)
	}
	return &viewCollector{}
}

func (v *viewCollector) Update(number number.Number, descriptor *sdkapi.Descriptor) {
	for _, input := range v.inputs {
		input.Update(number, descriptor)
	}
}

func (v *viewCollector) Send(final bool) error {
	// @@@ Call SynchronizedMove on each input

	return nil
}
