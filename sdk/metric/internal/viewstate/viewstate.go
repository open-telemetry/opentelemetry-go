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
		Send(CollectorFactory) error
	}

	Updater[N number.Any] interface {
		Update(number N)
	}

	CollectorFactory interface {
		// New returns a Collector that also implements Updater[N]
		New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector
	}

	State struct {
		library instrumentation.Library
		readers []viewReader
	}

	// vCF is configured one per instrument with all
	// pre-calculated view behaviors.
	viewCollectorFactory struct {
		state         *State
		configuration []viewConfiguration
	}

	viewConfiguration struct {
		newFunc   func() viewCollector
		behaviors []viewBehavior
	}

	viewCollector interface {
		Send(viewConfiguration) error
	}

	viewBehavior struct {
		// TODO: this is not an efficient way to represent the
		// calculated behavior, as this struct contains every
		// option.  Replace with pointer? Small struct?
		view views.View

		reader *viewReader
	}

	viewCollectors []viewCollector

	viewReader struct {
		config      reader.Reader
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

func New(lib instrumentation.Library, readerConfig []reader.Reader) *State {

	// TODO: error checking here, such as:
	// - empty (?)
	// - duplicate name
	// - invalid inst/number/aggregation kind
	// - both instrument name and regexp
	// - schemaURL or Version without library name
	// - empty attribute keys
	// - Name w/o SingleInst
	readers := make([]viewReader, len(readerConfig))
	for i, r := range readerConfig {
		readers[i].outputNames = map[string]struct{}{}
		readers[i].config = r
	}

	return &State{
		library: lib,
		readers: readers,
	}
}

func configViewBehavior(v views.View, r *viewReader) viewBehavior {
	return viewBehavior{
		reader: r,
		view:   v,
	}
}

func defaultViewBehavior(desc sdkapi.Descriptor, r *viewReader) viewBehavior {
	as := aggregatorSettingsFor(desc)
	return viewBehavior{
		view:   views.New(views.WithAggregation(as.kind)),
		reader: r,
	}
}

func (vb viewBehavior) Name() string {
	return vb.view.Name()
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
	addBehavior := func(readerIdx int, as aggregatorSettings, behavior viewBehavior) {
		ss := fmt.Sprint(as)
		allBehaviors[ss] = settingsBehaviors{
			settings:  as,
			behaviors: append(allBehaviors[ss].behaviors, behavior),
		}

	}

	for readerIdx := range v.readers {
		matchCount := 0
		reader := &v.readers[readerIdx]
		for _, def := range v.readers[readerIdx].config.Views() {
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

			addBehavior(readerIdx, as, configViewBehavior(def, reader))
		}

		// If there were no matching views, set the default aggregation.
		if matchCount == 0 {
			if !reader.config.HasDefaultView() {
				continue
			}

			addBehavior(readerIdx, aggregatorSettingsFor(desc), defaultViewBehavior(desc, reader))
		}
	}
	// When there are no matches for any reader, return a nil factory.
	if len(allBehaviors) == 0 {
		return nil, nil
	}

	vcf := &viewCollectorFactory{state: v}

	for _, reader := range v.readers {
		reader.lock.Lock()
		defer reader.lock.Unlock()
	}

	for _, sbs := range allBehaviors {
		valid := 0
		for _, behavior := range sbs.behaviors {
			if _, has := behavior.reader.outputNames[behavior.Name()]; !has {
				behavior.reader.outputNames[behavior.Name()] = struct{}{}
				valid++
			} else {
				otel.Handle(fmt.Errorf("duplicate view name registered"))
			}
		}
		if valid == 0 {
			continue
		}

		var cfg viewConfiguration
		switch desc.NumberKind() {
		case number.Int64Kind:
			cfg = buildView[int64, traits.Int64](desc, sbs.settings, sbs.behaviors)
		case number.Float64Kind:
			cfg = buildView[float64, traits.Float64](desc, sbs.settings, sbs.behaviors)
		}
		vcf.configuration = append(vcf.configuration, cfg)
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

func buildView[N number.Any, Traits traits.Any[N]](desc sdkapi.Descriptor, settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	// @@@ TODO in both code paths, we are dropping `behaviors`. Seems to re-enter via Send()
	// so now the CollectorFactory is there. The behaviors aren't worth storing in every record,
	// so we don't.
	if desc.InstrumentKind().Synchronous() {
		return buildSyncView[N, Traits](settings, behaviors)
	}
	return buildAsyncView[N, Traits](settings, behaviors)
}

func newSyncConfig[N number.Any, Traits traits.Any[N], Methods aggregator.Methods[N, Storage, Config], Storage, Config any](behaviors []viewBehavior, cfg *Config) viewConfiguration {
	return viewConfiguration{
		behaviors: behaviors,
		newFunc: func() viewCollector {
			aa := &syncCollector[N, Methods, Storage, Config]{}
			aa.Init(*cfg)
			return aa
		},
	}
}

func newAsyncConfig[N number.Any, Traits traits.Any[N], Methods aggregator.Methods[N, Storage, Config], Storage, Config any](behaviors []viewBehavior, cfg *Config) viewConfiguration {
	return viewConfiguration{
		behaviors: behaviors,
		newFunc: func() viewCollector {
			aa := &asyncCollector[N, Methods, Storage, Config]{}
			aa.Init(*cfg)
			return aa
		},
	}
}

func buildSyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	switch settings.kind {
	case aggregation.LastValueKind:
		return newSyncConfig[N, Traits, lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]], lastvalue.State[N, Traits], lastvalue.Config](behaviors, &settings.lvcfg)
	case aggregation.HistogramKind:
		return newSyncConfig[N, Traits, histogram.Methods[N, Traits, histogram.State[N, Traits]], histogram.State[N, Traits], histogram.Config](behaviors, &settings.hcfg)
	default:
		return newSyncConfig[N, Traits, sum.Methods[N, Traits, sum.State[N, Traits]], sum.State[N, Traits], sum.Config](behaviors, &settings.scfg)
	}
}

func buildAsyncView[N number.Any, Traits traits.Any[N]](settings aggregatorSettings, behaviors []viewBehavior) viewConfiguration {
	switch settings.kind {
	case aggregation.LastValueKind:
		return newAsyncConfig[N, Traits, lastvalue.Methods[N, Traits, lastvalue.State[N, Traits]], lastvalue.State[N, Traits], lastvalue.Config](behaviors, &settings.lvcfg)
	case aggregation.HistogramKind:
		return newAsyncConfig[N, Traits, histogram.Methods[N, Traits, histogram.State[N, Traits]], histogram.State[N, Traits], histogram.Config](behaviors, &settings.hcfg)
	default:
		return newAsyncConfig[N, Traits, sum.Methods[N, Traits, sum.State[N, Traits]], sum.State[N, Traits], sum.Config](behaviors, &settings.scfg)
	}
}

func (factory *viewCollectorFactory) New(kvs []attribute.KeyValue, desc *sdkapi.Descriptor) Collector {
	collectors := make(viewCollectors, 0, len(factory.configuration))
	for idx, vc := range factory.configuration {
		collectors[idx] = vc.newFunc()
	}
	return collectors
}

func (v viewCollectors) Send(cfactory CollectorFactory) error {
	vcf, ok := cfactory.(*viewCollectorFactory)
	if !ok {
		return fmt.Errorf("wrong factory")
	}
	for i, collector := range v {
		collector.Send(vcf.configuration[i])
	}

	return nil
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

func (sc *syncCollector[N, Methods, Storage, Config]) Send(viewConfiguration) error {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)

	// @@@ do something
	return nil
}

func (ac *asyncCollector[N, Methods, Storage, Config]) Init(cfg Config) {
	var methods Methods
	ac.current = 0
	methods.Init(&ac.snapshot, cfg)
}

func (ac *asyncCollector[N, Methods, Storage, Config]) Update(number N) {
	ac.current = number
}

func (ac *asyncCollector[N, Methods, Storage, Config]) Send(viewConfiguration) error {
	var methods Methods
	methods.SynchronizedMove(&ac.snapshot, nil)
	methods.Update(&ac.snapshot, ac.current)
	ac.current = 0
	// @@@ do something
	return nil
}
