// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
)

// providerProducer is the binding between the MeterProvider and the
// Reader.  This is the Producer instance that is passed to Register()
// for each Reader.
type providerProducer struct {
	lock        sync.Mutex
	provider    *MeterProvider
	pipe        int
	lastCollect time.Time
}

// producerFor returns the new Producer for calling Register.
func (mp *MeterProvider) producerFor(pipe int) Producer {
	return &providerProducer{
		provider:    mp,
		pipe:        pipe,
		lastCollect: mp.startTime,
	}
}

// Produce runs collection and produces a new metrics data object.
func (pp *providerProducer) Produce(inout *data.Metrics) data.Metrics {
	ordered := pp.provider.getOrdered()

	// Note: the Last time is only used in delta-temporality
	// scenarios.  This lock protects the only stateful change in
	// `pp` but does not prevent concurrent collection.  If a
	// delta-temporality reader were to call Produce
	// concurrently, the results would be be recorded with
	// non-overlapping timestamps but would have been collected in
	// an overlapping way.
	pp.lock.Lock()
	lastTime := pp.lastCollect
	nowTime := time.Now()
	pp.lastCollect = nowTime
	pp.lock.Unlock()

	var output data.Metrics
	if inout != nil {
		inout.Reset()
		output = *inout
	}

	output.Resource = pp.provider.cfg.res

	sequence := data.Sequence{
		Start: pp.provider.startTime,
		Last:  lastTime,
		Now:   nowTime,
	}

	// TODO: Add a timeout to the context.
	ctx := context.Background()

	for _, meter := range ordered {
		meter.collectFor(
			ctx,
			pp.pipe,
			sequence,
			&output,
		)
	}

	return output
}

// collectFor collects from a single meter.
func (m *meter) collectFor(ctx context.Context, pipe int, seq data.Sequence, output *data.Metrics) {
	// Use m.lock to briefly access the current lists: syncInsts, asyncInsts, callbacks
	m.lock.Lock()
	syncInsts := m.syncInsts
	asyncInsts := m.asyncInsts
	callbacks := m.callbacks
	m.lock.Unlock()

	asyncState := asyncstate.NewState(pipe)

	for _, cb := range callbacks {
		cb.Run(ctx, asyncState)
	}

	for _, inst := range syncInsts {
		inst.SnapshotAndProcess()
	}

	for _, inst := range asyncInsts {
		inst.SnapshotAndProcess(asyncState)
	}

	scope := data.ReallocateFrom(&output.Scopes)
	scope.Library = m.library

	for _, coll := range m.compilers[pipe].Collectors() {
		coll.Collect(seq, &scope.Instruments)
	}
}
