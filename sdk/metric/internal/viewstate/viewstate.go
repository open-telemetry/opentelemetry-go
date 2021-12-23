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

		lock   sync.Mutex
		output map[string]*viewCollectorFactory
	}

	viewCollector struct {
		inputs []export.Aggregator
	}

	viewCollectorFactory struct {
		state  *State
		aggSet map[aggregation.Kind]viewMatchState
	}

	viewMatchState struct {
		matches []viewMatch
	}

	viewMatch struct {
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

func configViewMatch(v views.View) viewMatch {
	return viewMatch{}
}

func defaultViewMatch(k sdkapi.InstrumentKind) viewMatch {
	return viewMatch{}
}

// NewFactory is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Collectors throughout its lifetime.
func (v *State) NewFactory(desc sdkapi.Descriptor) (CollectorFactory, error) {
	// TODO: This should be conditioned on the configuration of
	// the aggregator/exporter too, but we're not configuring them yet.

	// Compute the set of matching views.
	aggSet := map[aggregation.Kind]viewMatchState{}
	for _, def := range v.definitions {
		if !def.Matches(v.library, desc) {
			continue
		}

		var kind aggregation.Kind
		switch def.Aggregation() {
		case aggregation.SumKind, aggregation.HistogramKind, aggregation.LastValueKind:
			kind = def.Aggregation()
		default:
			kind = aggregationKindFor(desc.InstrumentKind())
		}
		current := aggSet[kind]
		current.matches = append(current.matches, configViewMatch(def))
		aggSet[kind] = current
	}
	// If there were no matching views, set the default aggregation.
	if len(aggSet) == 0 {
		if !v.hasDefault {
			return nil, nil
		}
		k := desc.InstrumentKind()
		aggSet[aggregationKindFor(k)] = viewMatchState{
			matches: []viewMatch{defaultViewMatch(k)},
		}
	}

	// Develop the list of basic aggregations.
	vcf := &viewCollectorFactory{
		state:  v,
		aggSet: aggSet,
	}

	// Reconcile.  Build a representation of the output names.
	v.lock.Lock()
	defer v.lock.Unlock()

	for _, vms := range vcf.aggSet {
		for _, vm := range vms.matches {
			var outputName string
			if vm.HasName() {
				outputName = vm.Name()
			} else {
				outputName = desc.Name()
			}

			if _, has := v.output[outputName]; has {
				return nil, fmt.Errorf("duplicate view name configured")
			}
			v.output[outputName] = vcf
		}
	}

	return vcf, nil
}

func (v *viewCollectorFactory) New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector {
	var inputs []export.Aggregator
	for ak, vms := range v.aggSet {
		// Notes:
		//
		// Need to know how many aggregators are needed to construct
		// the output pipeline.
		//
		// 1 is the base case, for the collector itself
		// 1 for each matching view
		cnt := 1 + len(vms.matches)

		switch ak {
		case aggregation.SumKind:
			aggs := sum.New(cnt)
			inputs = append(inputs, &aggs[0])
		case aggregation.LastValueKind:
			aggs := lastvalue.New(cnt)
			inputs = append(inputs, &aggs[0])
		case aggregation.HistogramKind:
			aggs := histogram.New(cnt, desc)
			inputs = append(inputs, &aggs[0])
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

func (v *viewCollector) Send(final bool) error {
	// @@@ Call SynchronizedMove on each input
	// @@@ Does _final_ matter?

	return nil
}
