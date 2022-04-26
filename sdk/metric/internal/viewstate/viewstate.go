package viewstate

import (
	"reflect"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

// Compiler implements Views for a single Meter.  A single Compiler
// controls the namespace of metric instruments output and reports
// conflicting definitions for the same name.
//
// Information flows through the Compiler as follows:
//
// When new instruments are created:
// - The Compiler.Compile() method returns an Instrument value and possible
//   duplicate or semantic conflict error.
//
// When instruments are used:
// - The Instrument.NewAccumulator() method returns an Accumulator value for each attribute.Set used
// - The Accumulator.Update() aggregates one value for each measurement.
//
// During collection:
// - The Accumulator.SnapshotAndProcess() method captures the current value
//   and conveys it to the output storage
// - The Compiler.Collectors() interface returns one Collector per output
//   Metric in the Meter (duplicate definitions included).
// - The Collector.Collect() method outputs one Point for each attribute.Set
//   in the result.
type Compiler struct {
	// views is the configuration of this compiler.
	views *view.Views

	// library is the value used fr matching
	// instrumentation library information.
	library instrumentation.Library

	// lock protects collectors and names.
	lock sync.Mutex

	// collectors is the de-duplicated list of metric outputs, which may
	// contain conflicting identities.
	collectors []data.Collector

	// names is the map of output names for metrics
	// produced by this compiler.
	names map[string][]leafInstrument
}

// Instrument is a compiled implementation of an instrument
// corresponding with one or more instrument-view behaviors.
type Instrument interface {
	// NewAccumulator returns an Accumulator and an Updater[N]
	// matching the number type of the API-level instrument.
	//
	// Callers are expected to type-assert Updater[int64] or
	// Updater[float64] before calling Update().
	//
	// The caller's primary responsibility is to maintain
	// the collection of Accumulators that had Update()
	// called since the last collection and to ensure that each
	// of them has SnapshotAndProcess() called.
	NewAccumulator(kvs attribute.Set) Accumulator
}

// Updater captures single measurements, for N an int64 or float64.
type Updater[N number.Any] interface {
	// Update captures a single measurement.  For synchronous
	// instruments, this passes directly through to the
	// aggregator.  For asynchronous instruments, the last value
	// is captured by the accumulator snapshot.
	Update(value N)
}

// Accumulator is an intermediate interface used for short-term
// aggregation.  Every Accumulator is also an Updater.  The owner of
// an Accumulator is responsible for maintaining the the current set
// of Accumulators, defined as those which have been Updated and not
// yet had SnapshotAndProcess() called.
type Accumulator interface {
	// SnapshotAndProcess() takes a snapshot of data aggregated
	// through Update() and simultaneously resets the current
	// aggregator.  The attribute.Set is possibly filtered, after
	// which the snapshot is merged into the output.
	//
	// There is no return value from this method; the caller can
	// safely forget an Accumulator after this method is called,
	// provided Update is not used again.
	SnapshotAndProcess()
}

// leafInstrument is one of the (synchronous or asynchronous),
// (cumulative or delta) instrument implementations.  This is used in
// duplicate conflict detection and resolution.
type leafInstrument interface {
	// Instrument is the form returned by Compile().
	Instrument
	// Collector is the form returned in Collectors().
	data.Collector
	// Duplicate is how other instruments this in a conflict.
	Duplicate

	// mergeDescription handles the special case allowing
	// descriptions to be merged instead of conflict.
	mergeDescription(string)
}

// singleBehavior is one instrument-view behavior, including the
// original instrument details, the aggregation kind and temporality,
// aggregator configuration, and optional keys to filter.
type singleBehavior struct {
	// fromName is the original instrument name
	fromName string

	// desc is the output of the view, including naming,
	// description and unit.  This includes the original
	// instrument's instrument kind and number kind.
	desc sdkinstrument.Descriptor

	// kind is the aggregation indicated by this view behavior.
	kind aggregation.Kind

	// tempo is the configured aggregation temporality.
	tempo aggregation.Temporality

	// acfg is the aggregator configuration.
	acfg aggregator.Config

	// keysSet (if non-nil) is an attribute set containing each
	// key being filtered with a zero value.  This is used to
	// compare against potential duplicates for having the
	// same/different filter.
	keysSet *attribute.Set // With Int(0)

	// keysFilter (if non-nil) is the constructed keys filter.
	keysFilter *attribute.Filter
}

// New returns a compiler for library given configured views.
func New(library instrumentation.Library, views *view.Views) *Compiler {
	return &Compiler{
		library: library,
		views:   views,
		names:   map[string][]leafInstrument{},
	}
}

func (v *Compiler) Collectors() []data.Collector {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.collectors
}

// Compile is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Accumulators throughout its lifetime.
func (v *Compiler) Compile(instrument sdkinstrument.Descriptor) (Instrument, ViewConflictsBuilder) {
	var behaviors []singleBehavior
	var matches []view.ClauseConfig

	for _, view := range v.views.Clauses {
		if !view.Matches(v.library, instrument) {
			continue
		}
		matches = append(matches, view)
	}

	for _, view := range matches {
		akind := view.Aggregation()
		if akind == aggregation.DropKind {
			continue
		}
		if akind == aggregation.UndefinedKind {
			akind = v.views.Defaults.Aggregation(instrument.Kind)
		}

		cf := singleBehavior{
			fromName: instrument.Name,
			desc:     viewDescriptor(instrument, view),
			kind:     akind,
			acfg:     viewAggConfig(&v.views.Defaults, akind, instrument.Kind, instrument.NumberKind, view.AggregatorConfig()),
			tempo:    v.views.Defaults.Temporality(instrument.Kind),
		}

		keys := view.Keys()
		if keys != nil {
			cf.keysSet = keysToSet(view.Keys())
			cf.keysFilter = keysToFilter(view.Keys())
		}
		behaviors = append(behaviors, cf)
	}

	// If there were no matching views, set the default aggregation.
	if len(matches) == 0 {
		akind := v.views.Defaults.Aggregation(instrument.Kind)
		if akind != aggregation.DropKind {
			behaviors = append(behaviors, singleBehavior{
				fromName: instrument.Name,
				desc:     instrument,
				kind:     akind,
				acfg:     viewAggConfig(&v.views.Defaults, akind, instrument.Kind, instrument.NumberKind, aggregator.Config{}),
				tempo:    v.views.Defaults.Temporality(instrument.Kind),
			})
		}
	}

	v.lock.Lock()
	defer v.lock.Unlock()

	var conflicts ViewConflictsBuilder
	var compiled []Instrument

	for _, behavior := range behaviors {
		// the following checks semantic compatibility
		// and if necessary fixes the aggregation kind
		// to the default, via in place update.
		semanticErr := checkSemanticCompatibility(instrument.Kind, &behavior.kind)

		existingInsts := v.names[behavior.desc.Name]
		var leaf leafInstrument

		// Scan the existing instruments for a match.
		for _, inst := range existingInsts {
			// Test for equivalence among the fields that we
			// cannot merge or will not convert, means the
			// testing everything except the description for
			// equality.

			if inst.Aggregation() != behavior.kind {
				continue
			}
			if inst.Descriptor().Kind.Synchronous() != behavior.desc.Kind.Synchronous() {
				continue
			}

			if inst.Descriptor().Unit != behavior.desc.Unit {
				continue
			}
			if inst.Descriptor().NumberKind != behavior.desc.NumberKind {
				continue
			}
			if !equalConfigs(inst.Config(), behavior.acfg) {
				continue
			}

			// For attribute keys, test for equal nil-ness or equal value.
			instKeys := inst.Keys()
			confKeys := behavior.keysSet
			if (instKeys == nil) != (confKeys == nil) {
				continue
			}
			if instKeys != nil && *instKeys != *confKeys {
				continue
			}
			// We can return the previously-compiled instrument,
			// we may have different descriptions and that is
			// specified to choose the longer one.
			inst.mergeDescription(behavior.desc.Description)
			leaf = inst
			break
		}
		if leaf == nil {
			switch behavior.desc.NumberKind {
			case number.Int64Kind:
				leaf = buildView[int64, number.Int64Traits](behavior)
			case number.Float64Kind:
				leaf = buildView[float64, number.Float64Traits](behavior)
			}

			v.collectors = append(v.collectors, leaf)
			existingInsts = append(existingInsts, leaf)
			v.names[behavior.desc.Name] = existingInsts

		}
		if len(existingInsts) > 1 || semanticErr != nil {
			c := Conflict{
				Semantic:   semanticErr,
				Duplicates: make([]Duplicate, len(existingInsts)),
			}
			for i := range existingInsts {
				c.Duplicates[i] = existingInsts[i]
			}
			conflicts.Add(v.views.Name, c)
		}
		compiled = append(compiled, leaf)
	}
	return Combine(instrument, compiled...), conflicts
}

// buildView compiles either a synchronous or asynchronous instrument
// given its behavior and generic number type/traits.
func buildView[N number.Any, Traits number.Traits[N]](behavior singleBehavior) leafInstrument {
	if behavior.desc.Kind.Synchronous() {
		return compileSync[N, Traits](behavior)
	}
	return compileAsync[N, Traits](behavior)
}

// newSyncView returns a compiled synchronous instrument.  If the view
// calls for delta temporality, a stateless instrument is returned,
// otherwise for cumulative temporality a stateful instrument will be
// used.  I.e., Delta->Stateless, Cumulative->Stateful.
func newSyncView[
	N number.Any,
	Storage any,
	Methods aggregator.Methods[N, Storage],
](behavior singleBehavior) leafInstrument {
	metric := instrumentBase[N, Storage, Methods]{
		fromName:   behavior.fromName,
		desc:       behavior.desc,
		acfg:       behavior.acfg,
		data:       map[attribute.Set]*Storage{},
		keysSet:    behavior.keysSet,
		keysFilter: behavior.keysFilter,
	}
	instrument := compiledSyncBase[N, Storage, Methods]{
		instrumentBase: metric,
	}
	if behavior.tempo == aggregation.DeltaTemporality {
		return &statelessSyncInstrument[N, Storage, Methods]{
			compiledSyncBase: instrument,
		}
	}

	return &statefulSyncInstrument[N, Storage, Methods]{
		compiledSyncBase: instrument,
	}
}

// compileSync calls newSyncView to compile a synchronous
// instrument with specific aggregator storage and methods.
func compileSync[N number.Any, Traits number.Traits[N]](behavior singleBehavior) leafInstrument {
	switch behavior.kind {
	case aggregation.HistogramKind:
		return newSyncView[
			N,
			histogram.State[N, Traits],
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](behavior)
	case aggregation.NonMonotonicSumKind:
		return newSyncView[
			N,
			sum.State[N, Traits, sum.NonMonotonic],
			sum.Methods[N, Traits, sum.NonMonotonic, sum.State[N, Traits, sum.NonMonotonic]],
		](behavior)
	default: // e.g., aggregation.MonotonicSumKind
		return newSyncView[
			N,
			sum.State[N, Traits, sum.Monotonic],
			sum.Methods[N, Traits, sum.Monotonic, sum.State[N, Traits, sum.Monotonic]],
		](behavior)
	}
}

// newAsyncView returns a compiled asynchronous instrument.  If the
// view calls for delta temporality, a stateful instrument is
// returned, otherwise for cumulative temporality a stateless
// instrument will be used.  I.e., Cumulative->Stateless,
// Delta->Stateful.
func newAsyncView[
	N number.Any,
	Storage any,
	Methods aggregator.Methods[N, Storage],
](behavior singleBehavior) leafInstrument {
	metric := instrumentBase[N, Storage, Methods]{
		fromName:   behavior.fromName,
		desc:       behavior.desc,
		acfg:       behavior.acfg,
		data:       map[attribute.Set]*Storage{},
		keysSet:    behavior.keysSet,
		keysFilter: behavior.keysFilter,
	}
	instrument := compiledAsyncBase[N, Storage, Methods]{
		instrumentBase: metric,
	}
	if behavior.tempo == aggregation.DeltaTemporality {
		return &statefulAsyncInstrument[N, Storage, Methods]{
			compiledAsyncBase: instrument,
		}
	}

	return &statelessAsyncInstrument[N, Storage, Methods]{
		compiledAsyncBase: instrument,
	}
}

// compileAsync calls newAsyncView to compile an asynchronous
// instrument with specific aggregator storage and methods.
func compileAsync[N number.Any, Traits number.Traits[N]](behavior singleBehavior) leafInstrument {
	switch behavior.kind {
	case aggregation.MonotonicSumKind:
		return newAsyncView[
			N,
			sum.State[N, Traits, sum.Monotonic],
			sum.Methods[N, Traits, sum.Monotonic, sum.State[N, Traits, sum.Monotonic]],
		](behavior)
	case aggregation.NonMonotonicSumKind:
		return newAsyncView[
			N,
			sum.State[N, Traits, sum.NonMonotonic],
			sum.Methods[N, Traits, sum.NonMonotonic, sum.State[N, Traits, sum.NonMonotonic]],
		](behavior)
	default: // e.g., aggregation.GaugeKind
		return newAsyncView[
			N,
			gauge.State[N, Traits],
			gauge.Methods[N, Traits, gauge.State[N, Traits]],
		](behavior)
	}
}

// Combine accepts a variable number of Instruments to combine.  If 0
// items, nil is returned. If 1 item, the item itself is return.
// otherwise, a multiInstrument of the appropriate number kind is returned.
func Combine(desc sdkinstrument.Descriptor, insts ...Instrument) Instrument {
	if len(insts) == 0 {
		return nil
	}
	if len(insts) == 1 {
		return insts[0]
	}
	if desc.NumberKind == number.Float64Kind {
		return multiInstrument[float64](insts)
	}
	return multiInstrument[int64](insts)
}

// multiInstrument is used by Combine() to combine the effects of
// multiple instrument-view behaviors.  These instruments produce
// multiAccumulators in NewAccumulator.
type multiInstrument[N number.Any] []Instrument

// NewAccumulator returns a Accumulator for multiple views of the same instrument.
func (mi multiInstrument[N]) NewAccumulator(kvs attribute.Set) Accumulator {
	accs := make([]Accumulator, 0, len(mi))

	for _, inst := range mi {
		accs = append(accs, inst.NewAccumulator(kvs))
	}
	return multiAccumulator[N](accs)
}

// Uses a int(0)-value attribute to identify distinct key sets.
func keysToSet(keys []attribute.Key) *attribute.Set {
	attrs := make([]attribute.KeyValue, len(keys))
	for i, key := range keys {
		attrs[i] = key.Int(0)
	}
	ns := attribute.NewSet(attrs...)
	return &ns
}

// keyFilter provides an attribute.Filter implementation based on a
// map[attribute.Key].
type keyFilter map[attribute.Key]struct{}

// filter is an attribute.Filter.
func (ks keyFilter) filter(kv attribute.KeyValue) bool {
	_, has := ks[kv.Key]
	return has
}

// keysToFilter constructs a keyFilter.
func keysToFilter(keys []attribute.Key) *attribute.Filter {
	kf := keyFilter{}
	for _, k := range keys {
		kf[k] = struct{}{}
	}
	var af attribute.Filter = kf.filter
	return &af
}

// equalConfigs compares two aggregator configurations.
func equalConfigs(a, b aggregator.Config) bool {
	return reflect.DeepEqual(a, b)
}

// histogramBoundariesFor returns the configured or
// number-kind-specific default histogram boundaries.
func histogramBoundariesFor(d *view.DefaultConfig, ik sdkinstrument.Kind, nk number.Kind, acfg aggregator.Config) []float64 {
	if len(acfg.Histogram.ExplicitBoundaries) != 0 {
		return acfg.Histogram.ExplicitBoundaries
	}

	cfg := d.AggregationConfig(ik, nk)
	if len(cfg.Histogram.ExplicitBoundaries) != 0 {
		return cfg.Histogram.ExplicitBoundaries
	}
	if nk == number.Int64Kind {
		return histogram.DefaultInt64Boundaries
	}
	return histogram.DefaultFloat64Boundaries
}

// viewAggConfig returns the aggregator configuration prescribed by a view clause.
func viewAggConfig(r *view.DefaultConfig, ak aggregation.Kind, ik sdkinstrument.Kind, nk number.Kind, vcfg aggregator.Config) aggregator.Config {
	if ak != aggregation.HistogramKind {
		return aggregator.Config{}
	}
	return aggregator.Config{
		Histogram: histogram.NewConfig(
			histogramBoundariesFor(r, ik, nk, vcfg),
		),
	}
}

// checkSemanticCompatibility checks whether an instrument /
// aggregator pairing is well defined.
//
// TODO(jmacd): There are a couple of specification questions about
// this worth raising.
func checkSemanticCompatibility(ik sdkinstrument.Kind, aggPtr *aggregation.Kind) error {
	agg := *aggPtr
	cat := agg.Category(ik)

	if agg == aggregation.AnySumKind {
		switch cat {
		case aggregation.MonotonicSumCategory, aggregation.HistogramCategory:
			agg = aggregation.MonotonicSumKind
		case aggregation.NonMonotonicSumCategory:
			agg = aggregation.NonMonotonicSumKind
		default:
			agg = aggregation.UndefinedKind
		}
		*aggPtr = agg
	}

	switch ik {
	case sdkinstrument.CounterKind, sdkinstrument.HistogramKind:
		switch cat {
		case aggregation.MonotonicSumCategory, aggregation.NonMonotonicSumCategory, aggregation.HistogramCategory:
			return nil
		}

	case sdkinstrument.UpDownCounterKind, sdkinstrument.UpDownCounterObserverKind:
		switch cat {
		case aggregation.NonMonotonicSumCategory:
			return nil
		}

	case sdkinstrument.CounterObserverKind:
		switch cat {
		case aggregation.NonMonotonicSumCategory, aggregation.MonotonicSumCategory:
			return nil
		}

	case sdkinstrument.GaugeObserverKind:
		switch cat {
		case aggregation.GaugeCategory:
			return nil
		}
	}

	*aggPtr = view.StandardAggregationKind(ik)
	return SemanticError{
		Instrument:  ik,
		Aggregation: agg,
	}
}

// viewDescriptor returns the modified sdkinstrument.Descriptor of a
// view.  It retains the original instrument kind, numebr kind, and
// unit, while allowing the name and description to change.
func viewDescriptor(instrument sdkinstrument.Descriptor, v view.ClauseConfig) sdkinstrument.Descriptor {
	ikind := instrument.Kind
	nkind := instrument.NumberKind
	name := instrument.Name
	description := instrument.Description
	unit := instrument.Unit
	if v.HasName() {
		name = v.Name()
	}
	if v.Description() != "" {
		description = v.Description()
	}
	return sdkinstrument.NewDescriptor(name, ikind, nkind, description, unit)
}
