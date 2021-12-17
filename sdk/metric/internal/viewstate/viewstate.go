package viewstate

import (
	"context"
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
		Update(ctx context.Context, number number.Number, descriptor *sdkapi.Descriptor)
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
		akinds  []aggregation.Kind
		matches []views.View
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

// func aggregatorFor(ak aggregation.Kind) export.Aggregator {
// 	switch ak {
// 	case aggregation.SumKind:
// 	case aggregation.LastValueKind:
// 	case aggregation.HistogramKind:
// 	}
// }

func New(lib instrumentation.Library, defs []views.View, hasDefault bool) *State {

	// TODO: error checking here, such as:
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

func (v *State) NewFor(desc sdkapi.Descriptor) (CollectorFactory, error) {
	var matches []views.View

	// TODO: This should be conditioned on the configuraiton of
	// the aggregator too, but we're not configuring them yet.
	aggSet := map[aggregation.Kind]struct{}{}

	for _, def := range v.definitions {
		if def.Matches(v.library, desc) {
			matches = append(matches, def)
		}
	}
	for _, match := range matches {
		va := match.Aggregation()
		if va == "" {
			aggSet[aggregationKindFor(desc.InstrumentKind())] = struct{}{}
		} else {
			aggSet[va] = struct{}{}
		}
	}
	if matches == nil {
		if !v.hasDefault {
			return nil, nil
		}
		aggSet[aggregationKindFor(desc.InstrumentKind())] = struct{}{}
	}

	// Develop an output plan
	var aks []aggregation.Kind
	for ak := range aggSet {
		aks = append(aks, ak)
	}
	vcf := &viewCollectorFactory{
		state:   v,
		akinds:  aks,
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
	return &viewCollector{}
}

func (v *viewCollector) Update(_ context.Context, number number.Number, descriptor *sdkapi.Descriptor) {
	// TODO: baggage from context here

	for _, input := range v.inputs {
		input.Update(number, descriptor)
	}
}

func (v *viewCollector) Send(final bool) error {
	// @@@
	return nil
}
