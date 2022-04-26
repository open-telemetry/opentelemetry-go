package viewstate

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

// compiledSyncBase is any synchronous instrument view.
type compiledSyncBase[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
	instrumentBase[N, Storage, Methods]
}

// NewAccumulator returns a Accumulator for a synchronous instrument view.
func (csv *compiledSyncBase[N, Storage, Methods]) NewAccumulator(kvs attribute.Set) Accumulator {
	sc := &syncAccumulator[N, Storage, Methods]{}
	csv.initStorage(&sc.current)
	csv.initStorage(&sc.snapshot)

	sc.findStorage = csv.storageFinder(kvs)

	return sc
}

// compiledSyncBase is any asynchronous instrument view.
type compiledAsyncBase[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
	instrumentBase[N, Storage, Methods]
}

// NewAccumulator returns a Accumulator for an asynchronous instrument view.
func (cav *compiledAsyncBase[N, Storage, Methods]) NewAccumulator(kvs attribute.Set) Accumulator {
	ac := &asyncAccumulator[N, Storage, Methods]{}

	cav.initStorage(&ac.snapshot)
	ac.current = 0
	ac.findStorage = cav.storageFinder(kvs)

	return ac
}

// statefulSyncInstrument is a synchronous instrument that maintains cumulative state.
type statefulSyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
	compiledSyncBase[N, Storage, Methods]
}

// Collect for synchronous cumulative temporality.
func (p *statefulSyncInstrument[N, Storage, Methods]) Collect(seq data.Sequence, output *[]data.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := p.appendInstrument(output)

	for set, storage := range p.data {
		p.appendPoint(ioutput, set, methods.ToAggregation(storage), seq.Start, seq.Now)
	}
}

// statelessSyncInstrument is a synchronous instrument that maintains no state.
type statelessSyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
	compiledSyncBase[N, Storage, Methods]
}

// Collect for synchronous delta temporality.
func (p *statelessSyncInstrument[N, Storage, Methods]) Collect(seq data.Sequence, output *[]data.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := p.appendInstrument(output)

	for set, storage := range p.data {
		if !methods.HasChange(storage) {
			delete(p.data, set)
			continue
		}

		// Possibly re-use the underlying storage.  For
		// synchronous instruments, where accumulation happens
		// between collection events (e.g., due to other
		// readers collecting), we must reset the storage now
		// or completely clear the map.
		point, exists := p.appendOrReusePoint(ioutput)
		if exists == nil {
			exists = p.newStorage()
		} else {
			methods.Reset(exists)
		}
		methods.Merge(exists, storage)

		point.Attributes = set
		point.Aggregation = methods.ToAggregation(exists)
		point.Start = seq.Last
		point.End = seq.Now

		methods.Reset(storage)
	}
}

// statelessAsyncInstrument is an asynchronous instrument that keeps
// maintains no state.
type statelessAsyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
	compiledAsyncBase[N, Storage, Methods]
}

// Collect for asychronous cumulative temporality.
func (p *statelessAsyncInstrument[N, Storage, Methods]) Collect(seq data.Sequence, output *[]data.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := p.appendInstrument(output)

	for set, storage := range p.data {
		// Copy the underlying storage.
		p.appendPoint(ioutput, set, methods.ToAggregation(storage), seq.Start, seq.Now)
	}

	// Reset the entire map.
	p.data = map[attribute.Set]*Storage{}
}

// statefulAsyncInstrument is an instrument that keeps asynchronous instrument state
// in order to perform cumulative to delta translation.
type statefulAsyncInstrument[N number.Any, Storage any, Methods aggregator.Methods[N, Storage]] struct {
	compiledAsyncBase[N, Storage, Methods]
	prior map[attribute.Set]*Storage
}

// Collect for asynchronous delta temporality.
func (p *statefulAsyncInstrument[N, Storage, Methods]) Collect(seq data.Sequence, output *[]data.Instrument) {
	var methods Methods

	p.lock.Lock()
	defer p.lock.Unlock()

	ioutput := p.appendInstrument(output)

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
		p.appendPoint(ioutput, set, methods.ToAggregation(storage), seq.Last, seq.Now)
	}
	// Copy the current to the prior and reset.
	p.prior = p.data
	p.data = map[attribute.Set]*Storage{}
}
