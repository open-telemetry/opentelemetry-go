package viewstate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/views"
)

type (
	Compiler struct {
		library instrumentation.Library
		views   []views.View
		readers []*reader.Reader

		lock sync.Mutex

		// names is the per-reader map of output names,
		// indexed by the reader's position in `readers`.
		names []map[string][]leafInstrument

		// collectors is the de-duplicated list of metric outputs, which may
		// contain conflicting identities.
		collectors map[*reader.Reader][]Collector
	}

	Instrument interface {
		// NewAccumulator returns a new Accumulator bound to
		// the attributes `kvs`.  If reader == nil the
		// accumulator applies to all readers, otherwise it
		// applies to the specific reader.
		NewAccumulator(kvs attribute.Set, reader *reader.Reader) Accumulator
	}

	leafInstrument interface {
		Instrument
		Collector
		Duplicate

		mergeDescription(string)
	}

	Collector interface {
		// Collect transfers aggregated data from the
		// Accumulators into the output struct.
		Collect(reader *reader.Reader, sequence reader.Sequence, output *[]reader.Instrument)
	}

	Accumulator interface {
		Accumulate()
	}

	Updater[N number.Any] interface {
		Update(value N)
	}

	AccumulatorUpdater[N number.Any] interface {
		Accumulator
		Updater[N]
	}

	multiInstrument[N number.Any] map[*reader.Reader][]Instrument

	multiAccumulator[N number.Any] []Accumulator

	configuredBehavior struct {
		fromName string
		desc     sdkinstrument.Descriptor
		kind     aggregation.Kind
		acfg     aggregator.Config
		reader   *reader.Reader

		keysSet    *attribute.Set // With Int(0)
		keysFilter *attribute.Filter
	}

	baseMetric[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		lock     sync.Mutex
		fromName string
		desc     sdkinstrument.Descriptor
		acfg     aggregator.Config
		data     map[attribute.Set]*Storage

		keysSet    *attribute.Set
		keysFilter *attribute.Filter
	}

	compiledSyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		baseMetric[N, Storage, Methods]
	}

	compiledAsyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		baseMetric[N, Storage, Methods]
	}

	statelessSyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledSyncInstrument[N, Storage, Methods]
	}

	statefulSyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledSyncInstrument[N, Storage, Methods]
	}

	statelessAsyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledAsyncInstrument[N, Storage, Methods]
	}

	statefulAsyncProcess[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		compiledAsyncInstrument[N, Storage, Methods]
		prior map[attribute.Set]*Storage
	}

	syncAccumulator[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		current  Storage
		snapshot Storage
		output   *Storage
	}

	asyncAccumulator[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
		lock     sync.Mutex
		current  N
		snapshot Storage
		output   *Storage
	}
)

const noConflictsString = "no conflicts"

func New(lib instrumentation.Library, views []views.View, readers []*reader.Reader) *Compiler {

	// TODO: error checking here, such as:
	// - empty (?)
	// - duplicate name
	// - invalid inst/number/aggregation kind
	// - both instrument name and regexp
	// - schemaURL or Version without library name
	// - empty attribute keys
	// - Name w/o SingleInst

	names := make([]map[string][]leafInstrument, len(readers))
	for i := range names {
		names[i] = map[string][]leafInstrument{}
	}

	return &Compiler{
		library:    lib,
		views:      views,
		readers:    readers,
		names:      names,
		collectors: map[*reader.Reader][]Collector{},
	}
}

func (v *Compiler) Collectors(r *reader.Reader) []Collector {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.collectors[r]
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

func keysToFilter(keys []attribute.Key) *attribute.Filter {
	kf := keyFilter{}
	for _, k := range keys {
		kf[k] = struct{}{}
	}
	var af attribute.Filter = kf.filter
	return &af
}

type DuplicateConflicts map[*reader.Reader][]DuplicateConflict

type DuplicateConflict []Duplicate

type Duplicate interface {
	Aggregation() aggregation.Kind
	Descriptor() sdkinstrument.Descriptor
	Keys() *attribute.Set
	Config() aggregator.Config
	OriginalName() string
}

type duplicateInstrument leafInstrument

func (dc DuplicateConflicts) Error() string {
	total := 0
	for _, l := range dc {
		total += len(l)
	}
	// These are almost always duplicative, so we print only examples for one Reader.
	for _, byReader := range dc {
		if len(byReader) == 0 {
			break
		}
		if len(dc) == 1 {
			if len(byReader) == 1 {
				return byReader[0].Error()
			}
			return fmt.Sprintf("%d conflicts, e.g. %v", total, byReader[0])
		}
		return fmt.Sprintf("%d conflicts in %d readers, e.g. %v", total, len(dc), byReader[0])
	}
	return noConflictsString
}

func (DuplicateConflicts) Is(err error) bool {
	_, ok := err.(DuplicateConflicts)
	return ok
}

func (dc DuplicateConflict) Error() string {
	// Note: choose the first and last element of the current conflicts
	// list because they are ordered, and if the conflict in question is
	// new it will be the last item.
	if len(dc) < 2 {
		return noConflictsString
	}
	inst1 := dc[0]
	inst2 := dc[len(dc)-1]
	name1 := fullNameString(inst1)
	name2 := renameString(inst2)
	conf1 := shortString(inst1)
	conf2 := shortString(inst2)

	var s strings.Builder
	s.WriteString(name1)

	if conf1 != conf2 {
		s.WriteString(fmt.Sprintf(" conflicts %v, %v%v", conf1, conf2, name2))
	} else if !equalConfigs(inst1.Config(), inst2.Config()) {
		s.WriteString(" has conflicts: different aggregator configuration")
	} else {
		s.WriteString(" has conflicts: different attribute filters")
	}

	if len(dc) > 2 {
		s.WriteString(fmt.Sprintf(" and %d more", len(dc)-2))
	}
	return s.String()
}

func equalConfigs(a, b aggregator.Config) bool {
	return reflect.DeepEqual(a, b)
}

func histogramBoundariesFor(r *reader.Reader, ik sdkinstrument.Kind, nk number.Kind, acfg aggregator.Config) []float64 {
	if len(acfg.Histogram.ExplicitBoundaries) != 0 {
		return acfg.Histogram.ExplicitBoundaries
	}

	cfg := r.DefaultAggregationConfig(ik, nk)
	if len(cfg.Histogram.ExplicitBoundaries) != 0 {
		return cfg.Histogram.ExplicitBoundaries
	}
	if nk == number.Int64Kind {
		return histogram.DefaultInt64Boundaries
	}
	return histogram.DefaultFloat64Boundaries
}

func viewAggConfig(r *reader.Reader, ak aggregation.Kind, ik sdkinstrument.Kind, nk number.Kind, vcfg aggregator.Config) aggregator.Config {
	if ak != aggregation.HistogramKind {
		return aggregator.Config{}
	}
	return aggregator.Config{
		Histogram: histogram.NewConfig(
			histogramBoundariesFor(r, ik, nk, vcfg),
		),
	}
}

// Compile is called during NewInstrument by the Meter
// implementation, the result saved in the instrument and used to
// construct new Accumulators throughout its lifetime.
func (v *Compiler) Compile(instrument sdkinstrument.Descriptor) (Instrument, error) {
	configs := make([][]configuredBehavior, len(v.readers))
	matches := make([]views.View, 0, len(v.views))

	for _, view := range v.views {
		if !view.Matches(v.library, instrument) {
			continue
		}
		matches = append(matches, view)
	}

	for readerIdx, r := range v.readers {
		for _, view := range matches {
			kind := view.Aggregation()
			switch kind {
			case aggregation.SumKind, aggregation.GaugeKind, aggregation.HistogramKind, aggregation.DropKind:
			default:
				kind = aggregationConfigFor(instrument, r)
			}

			if kind == aggregation.DropKind {
				continue
			}

			cf := configuredBehavior{
				fromName: instrument.Name,
				desc:     viewDescriptor(instrument, view),
				kind:     kind,
				acfg:     viewAggConfig(r, kind, instrument.Kind, instrument.NumberKind, view.AggregatorConfig()),
				reader:   r,
			}

			// TODO: Add checks for semantic compatibility somehwere
			// around here.  E.g., no Gauges applied to Counters.

			keys := view.Keys()
			if keys != nil {
				cf.keysSet = keysToSet(view.Keys())
				cf.keysFilter = keysToFilter(view.Keys())
			}
			configs[readerIdx] = append(configs[readerIdx], cf)
		}

		// If there were no matching views, set the default aggregation.
		if len(matches) == 0 {
			kind := aggregationConfigFor(instrument, r)
			if kind == aggregation.DropKind {
				continue
			}

			configs[readerIdx] = append(configs[readerIdx], configuredBehavior{
				fromName: instrument.Name,
				desc:     instrument,
				kind:     kind,
				acfg:     viewAggConfig(r, kind, instrument.Kind, instrument.NumberKind, aggregator.Config{}),
				reader:   r,
			})
		}
	}

	compiled := map[*reader.Reader][]Instrument{}

	v.lock.Lock()
	defer v.lock.Unlock()

	var conflicts DuplicateConflicts
	readerConflicts := 0

	for readerIdx, behaviors := range configs {
		r := v.readers[readerIdx]
		names := v.names[readerIdx]
		var conflictsThisReader []DuplicateConflict

		for _, config := range behaviors {

			existingInsts := names[config.desc.Name]
			var leaf leafInstrument

			// Scan the existing instruments for a match.
			for _, inst := range existingInsts {
				// Test for equivalence among the fields that we
				// cannot merge or will not convert, means the
				// testing everything except the description for
				// equality.

				if inst.Aggregation() != config.kind {
					continue
				}
				// If the aggregation is a Sum and monotonicity is
				// different, a conflict.
				if config.kind == aggregation.SumKind &&
					config.desc.Kind.Monotonic() != inst.Descriptor().Kind.Monotonic() {
					continue
				}
				if inst.Descriptor().Kind.Synchronous() != config.desc.Kind.Synchronous() {
					continue
				}

				if inst.Descriptor().Unit != config.desc.Unit {
					continue
				}
				if inst.Descriptor().NumberKind != config.desc.NumberKind {
					continue
				}
				if !equalConfigs(inst.Config(), config.acfg) {
					continue
				}

				// For attribute keys, test for equal nil-ness or equal value.
				instKeys := inst.Keys()
				confKeys := config.keysSet
				if (instKeys == nil) != (confKeys == nil) {
					continue
				}
				if instKeys != nil && *instKeys != *confKeys {
					continue
				}
				// We can return the previously-compiled instrument,
				// we may have different descriptions and that is
				// specified to choose the longer one.
				inst.mergeDescription(config.desc.Description)
				leaf = inst
				break
			}
			if leaf == nil {
				switch config.desc.NumberKind {
				case number.Int64Kind:
					leaf = buildView[int64, traits.Int64](config)
				case number.Float64Kind:
					leaf = buildView[float64, traits.Float64](config)
				}

				// @@@ need to check for semantic conflicts, somehow

				v.collectors[r] = append(v.collectors[r], leaf)
				existingInsts = append(existingInsts, leaf)
				names[config.desc.Name] = existingInsts
			}
			if len(existingInsts) > 1 {
				dups := make([]Duplicate, len(existingInsts))
				for i := range dups {
					dups[i] = existingInsts[i]
				}
				conflictsThisReader = append(conflictsThisReader, dups)
			}
			compiled[r] = append(compiled[r], leaf)
		}
		if len(conflictsThisReader) != 0 {
			readerConflicts++
			if conflicts == nil {
				conflicts = DuplicateConflicts{}
			}
			conflicts[r] = append(conflicts[r], conflictsThisReader...)
		}
	}

	var err error
	if len(conflicts) != 0 {
		err = conflicts
	}

	switch len(compiled) {
	case 0:
		return nil, nil
	case 1:
		// As a special case, recognize the case where there
		// is only one reader and only one view to bypass the
		// map[...][]Instrument wrapper.
		for _, list := range compiled {
			if len(list) == 1 {
				return list[0], err
			}
		}
	}
	if instrument.NumberKind == number.Int64Kind {
		return multiInstrument[int64](compiled), err
	}
	return multiInstrument[float64](compiled), err
}

func aggregationConfigFor(desc sdkinstrument.Descriptor, r *reader.Reader) aggregation.Kind {
	return r.DefaultAggregation(desc.Kind)
}

func viewDescriptor(instrument sdkinstrument.Descriptor, v views.View) sdkinstrument.Descriptor {
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

func buildView[N number.Any, Traits traits.Any[N]](config configuredBehavior) leafInstrument {
	if config.desc.Kind.Synchronous() {
		return compileSync[N, Traits](config)
	}
	return compileAsync[N, Traits](config)
}

func newSyncView[
	N number.Any,
	Storage any,
	Methods aggregator.Methods[N, Storage],
](config configuredBehavior) leafInstrument {
	tempo := config.reader.DefaultTemporality(config.desc.Kind)
	metric := baseMetric[N, Storage, Methods]{
		fromName: config.fromName,
		desc:     config.desc,
		acfg:     config.acfg,
		data:     map[attribute.Set]*Storage{},

		keysSet:    config.keysSet,
		keysFilter: config.keysFilter,
	}
	instrument := compiledSyncInstrument[N, Storage, Methods]{
		baseMetric: metric,
	}
	if tempo == aggregation.DeltaTemporality {
		return &statelessSyncProcess[N, Storage, Methods]{
			compiledSyncInstrument: instrument,
		}
	}

	return &statefulSyncProcess[N, Storage, Methods]{
		compiledSyncInstrument: instrument,
	}
}

func newAsyncView[
	N number.Any,
	Storage any,
	Methods aggregator.Methods[N, Storage],
](config configuredBehavior) leafInstrument {
	tempo := config.reader.DefaultTemporality(config.desc.Kind)
	metric := baseMetric[N, Storage, Methods]{
		fromName: config.fromName,
		desc:     config.desc,
		acfg:     config.acfg,
		data:     map[attribute.Set]*Storage{},
		keysSet:  config.keysSet,
	}
	instrument := compiledAsyncInstrument[N, Storage, Methods]{
		baseMetric: metric,
	}
	if tempo == aggregation.DeltaTemporality {
		return &statefulAsyncProcess[N, Storage, Methods]{
			compiledAsyncInstrument: instrument,
		}
	}

	return &statelessAsyncProcess[N, Storage, Methods]{
		compiledAsyncInstrument: instrument,
	}
}

func compileSync[N number.Any, Traits traits.Any[N]](config configuredBehavior) leafInstrument {
	switch config.kind {
	case aggregation.HistogramKind:
		return newSyncView[
			N,
			histogram.State[N, Traits],
			histogram.Methods[N, Traits, histogram.State[N, Traits]],
		](config)
	default:
		// Note: this includes Drop and Gauge, which are prevented above.
		return newSyncView[
			N,
			sum.State[N, Traits],
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config)
	}
}

func compileAsync[N number.Any, Traits traits.Any[N]](config configuredBehavior) leafInstrument {
	switch config.kind {
	case aggregation.SumKind:
		return newAsyncView[
			N,
			sum.State[N, Traits],
			sum.Methods[N, Traits, sum.State[N, Traits]],
		](config)
	default:
		// Note: this includes Drop and Histogram, which are prevented above.
		return newAsyncView[
			N,
			gauge.State[N, Traits],
			gauge.Methods[N, Traits, gauge.State[N, Traits]],
		](config)
	}
}

// NewAccumulator returns a Accumulator for multiple views of the same instrument.
func (mi multiInstrument[N]) NewAccumulator(kvs attribute.Set, reader *reader.Reader) Accumulator {
	var collectors []Accumulator
	// Note: This runtime switch happens because we're using the same API for
	// both async and sync instruments, whereas the APIs are not symmetrical.
	if reader == nil {
		for _, list := range mi {
			collectors = make([]Accumulator, 0, len(mi)*len(list))
		}
		for _, list := range mi {
			for _, inst := range list {
				collectors = append(collectors, inst.NewAccumulator(kvs, nil))
			}
		}
	} else {
		insts := mi[reader]

		collectors = make([]Accumulator, 0, len(insts))

		for _, inst := range insts {
			collectors = append(collectors, inst.NewAccumulator(kvs, reader))
		}
	}
	return multiAccumulator[N](collectors)
}

// multiAccumulator

func (c multiAccumulator[N]) Accumulate() {
	for _, coll := range c {
		coll.Accumulate()
	}
}

func (c multiAccumulator[N]) Update(value N) {
	for _, coll := range c {
		coll.(Updater[N]).Update(value)
	}
}

// syncAccumulator

func (sc *syncAccumulator[N, Storage, Methods]) Update(number N) {
	var methods Methods
	methods.Update(&sc.current, number)
}

func (sc *syncAccumulator[N, Storage, Methods]) Accumulate() {
	var methods Methods
	methods.SynchronizedMove(&sc.current, &sc.snapshot)
	methods.Merge(sc.output, &sc.snapshot)
}

// asyncAccumulator

func (ac *asyncAccumulator[N, Storage, Methods]) Update(number N) {
	ac.lock.Lock()
	defer ac.lock.Unlock()
	ac.current = number
}

func (ac *asyncAccumulator[N, Storage, Methods]) Accumulate() {
	ac.lock.Lock()
	defer ac.lock.Unlock()

	var methods Methods
	methods.Reset(&ac.snapshot)
	methods.Update(&ac.snapshot, ac.current)
	methods.Merge(ac.output, &ac.snapshot)
}

// baseMetric

func (metric *baseMetric[N, Storage, Methods]) Aggregation() aggregation.Kind {
	var methods Methods
	return methods.Kind()
}

func (metric *baseMetric[N, Storage, Methods]) OriginalName() string {
	return metric.fromName
}

func (metric *baseMetric[N, Storage, Methods]) Descriptor() sdkinstrument.Descriptor {
	return metric.desc
}

func (metric *baseMetric[N, Storage, Methods]) Keys() *attribute.Set {
	return metric.keysSet
}

func (metric *baseMetric[N, Storage, Methods]) Config() aggregator.Config {
	return metric.acfg
}

func (metric *baseMetric[N, Storage, Methods]) initStorage(s *Storage) {
	var methods Methods
	methods.Init(s, metric.acfg)
}

func (metric *baseMetric[N, Storage, Methods]) mergeDescription(d string) {
	metric.lock.Lock()
	defer metric.lock.Unlock()
	if len(d) > len(metric.desc.Description) {
		metric.desc.Description = d
	}
}

func (metric *baseMetric[N, Storage, Methods]) findOutput(
	kvs attribute.Set,
) *Storage {
	if metric.keysFilter != nil {
		kvs, _ = attribute.NewSetWithFiltered(kvs.ToSlice(), *metric.keysFilter)
	}

	metric.lock.Lock()
	defer metric.lock.Unlock()

	storage, has := metric.data[kvs]
	if has {
		return storage
	}

	ns := metric.newStorage()
	metric.data[kvs] = ns
	return ns
}

func (metric *baseMetric[N, Storage, Methods]) newStorage() *Storage {
	ns := new(Storage)
	metric.initStorage(ns)
	return ns
}

func fullNameString(d Duplicate) string {
	return fmt.Sprintf("name %q%v", d.Descriptor().Name, renameString(d))
}

func renameString(d Duplicate) string {
	if d.OriginalName() == d.Descriptor().Name {
		return ""
	}
	return fmt.Sprintf(" (original %q)", d.OriginalName())
}

func shortString(d Duplicate) string {
	s := fmt.Sprintf("%v-%v-%v",
		strings.TrimSuffix(d.Descriptor().Kind.String(), "Kind"),
		strings.TrimSuffix(d.Descriptor().NumberKind.String(), "Kind"),
		d.Aggregation(),
	)
	if d.Descriptor().Unit != "" {
		s = fmt.Sprintf("%v-%v", s, d.Descriptor().Unit)
	}
	return s
}

// NewAccumulator returns a Accumulator for a synchronous instrument view.
func (csv *compiledSyncInstrument[N, Storage, Methods]) NewAccumulator(kvs attribute.Set, _ *reader.Reader) Accumulator {
	sc := &syncAccumulator[N, Storage, Methods]{}
	csv.initStorage(&sc.current)
	csv.initStorage(&sc.snapshot)

	sc.output = csv.findOutput(kvs)

	return sc
}

// NewAccumulator returns a Accumulator for an asynchronous instrument view.
func (cav *compiledAsyncInstrument[N, Storage, Methods]) NewAccumulator(kvs attribute.Set, _ *reader.Reader) Accumulator {
	ac := &asyncAccumulator[N, Storage, Methods]{}

	cav.initStorage(&ac.snapshot)
	ac.current = 0
	ac.output = cav.findOutput(kvs)

	return ac
}

func appendInstrument(insts *[]reader.Instrument, desc sdkinstrument.Descriptor, tempo aggregation.Temporality) *reader.Instrument {
	ioutput := reader.Reallocate(insts)
	ioutput.Descriptor = desc
	ioutput.Temporality = tempo
	return ioutput
}

func appendPoint(points *[]reader.Point, set attribute.Set, agg aggregation.Aggregation, start, end time.Time) {
	soutput := reader.Reallocate(points)
	soutput.Attributes = set
	soutput.Aggregation = agg
	soutput.Start = start
	soutput.End = end
}

// Collect for Synchronous Delta->Cumulative
func (p *statefulSyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.CumulativeTemporality)

	for set, storage := range p.data {
		// Note: No reset in this process.
		// This takes a direct reference to the underlying storage.
		appendPoint(&ioutput.Points, set, methods.Aggregation(storage), seq.Start, seq.Now)
	}
}

// Collect for Synchronous Delta->Delta
func (p *statelessSyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.DeltaTemporality)

	for set, storage := range p.data {
		// Note: this can't be a Gauge b/c synchronous instrument.
		if !methods.HasChange(storage) {
			delete(p.data, set)
			continue
		}

		// Copy and reset the underlying storage.
		// Note: this could be re-used memory from the last collection.
		ns := p.newStorage()
		methods.Merge(ns, storage)
		methods.Reset(storage)

		appendPoint(&ioutput.Points, set, methods.Aggregation(ns), seq.Last, seq.Now)
	}
}

// Collect for Asychronous Cumulative->Cumulative
func (p *statelessAsyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.CumulativeTemporality)

	for set, storage := range p.data {
		// Copy the underlying storage.
		appendPoint(&ioutput.Points, set, methods.Aggregation(storage), seq.Start, seq.Now)
	}

	// Reset the entire map.
	p.data = map[attribute.Set]*Storage{}
}

// Collect for Asynchronous Cumulative->Delta
func (p *statefulAsyncProcess[N, Storage, Methods]) Collect(_ *reader.Reader, seq reader.Sequence, output *[]reader.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := appendInstrument(output, p.desc, aggregation.DeltaTemporality)

	for set, storage := range p.data {
		pval, has := p.prior[set]
		if has {
			// This does `*pval := *storage - *pval`
			methods.SubtractSwap(storage, pval)

			// Skip the series if it has not changed.
			if !methods.HasChange(pval) {
				continue
			}
			// Output the difference except for Gauge, in
			// which case output the new value.
			if p.desc.Kind.HasTemporality() {
				storage = pval
			}
		}

		appendPoint(&ioutput.Points, set, methods.Aggregation(storage), seq.Last, seq.Now)
	}
	// Copy the current to the prior and reset.
	p.prior = p.data
	p.data = map[attribute.Set]*Storage{}
}
