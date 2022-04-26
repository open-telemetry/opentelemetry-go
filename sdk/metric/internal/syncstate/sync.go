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

package syncstate

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/pipeline"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// Instrument maintains a mapping from attribute.Set to an internal
// record type for a single API-level instrument.  This type is
// organized so that a single attribute.Set lookup is performed
// regardless of the number of reader and instrument-view behaviors.
// Entries in the map have their accumulator's SnapshotAndProcess()
// method called whenever they are removed from the map, which can
// happen when any reader collects the instrument.
type Instrument struct {
	// descriptor is the API-provided descriptor for the
	// instrument, unmodified by views.
	descriptor sdkinstrument.Descriptor

	// compiled will be a single compiled instrument or a
	// multi-instrument in case of multiple view behaviors
	// and/or readers; these distinctions do not matter
	// for synchronous aggregation.
	compiled viewstate.Instrument

	// current is a synchronous form of map[attribute.Set]*record.
	current sync.Map
}

// NewInstruments builds a new synchronous instrument given the
// per-pipeline instrument-views compiled.  Note that the unused
// second parameter is an opaque value used in the asyncstate package,
// passed here to make these two packages generalize.
func NewInstrument(desc sdkinstrument.Descriptor, _ interface{}, compiled pipeline.Register[viewstate.Instrument]) *Instrument {
	return &Instrument{
		descriptor: desc,

		// Note that viewstate.Combine is used to eliminate
		// the per-pipeline distinction that is useful in the
		// asyncstate package.  Here, in the common case there
		// will be one pipeline and one view, such that
		// viewstate.Combine produces a single concrete
		// viewstate.Instrument.  Only when there are multiple
		// views or multiple pipelines will the combination
		// produce a viewstate.multiInstrment here.
		compiled: viewstate.Combine(desc, compiled...),
	}
}

// SnapshotAndProcess calls SnapshotAndProcess() for all live
// accumulators of this instrument.  Inactive accumulators will be
// subsequently removed from the map.
func (inst *Instrument) SnapshotAndProcess() {
	inst.current.Range(func(key interface{}, value interface{}) bool {
		rec := value.(*record)
		if rec.snapshotAndProcess() {
			return true
		}
		// Having no updates since last collection, try to unmap:
		if unmapped := rec.refMapped.tryUnmap(); !unmapped {
			// The record is still referenced, continue.
			return true
		}

		// If any other goroutines are now trying to re-insert this
		// entry in the map, they are busy calling Gosched() awaiting
		// this deletion:
		inst.current.Delete(key)

		// Last we'll see of this.
		_ = rec.snapshotAndProcess()
		return true
	})
}

// record consists of an accumulator, a reference count, the number of
// updates, and the number of collected updates.
type record struct {
	refMapped refcountMapped

	// updateCount is incremented on every Update.
	updateCount int64

	// collectedCount is set to updateCount on collection,
	// supports checking for no updates during a round.
	collectedCount int64

	// accumulator is can be a multi-accumulator if there
	// are multiple behaviors or multiple readers, but
	// these distinctions are not relevant for synchronous
	// instruments.
	accumulator viewstate.Accumulator
}

// snapshotAndProcessRecord checks whether the accumulator has been
// modified since the last collection (by any reader), returns a
// boolean indicating whether the record is active.  If active, calls
// SnapshotAndProcess on the associated accumulator and returns true.
// If updates happened since the last collection (by any reader),
// returns false.
func (rec *record) snapshotAndProcess() bool {
	mods := atomic.LoadInt64(&rec.updateCount)
	coll := rec.collectedCount

	if mods == coll {
		return false
	}
	// Updates happened in this interval, collect and continue.
	rec.collectedCount = mods

	rec.accumulator.SnapshotAndProcess()
	return true
}

// capture performs a single update for any synchronous instrument.
func capture[N number.Any, Traits number.Traits[N]](_ context.Context, inst *Instrument, num N, attrs []attribute.KeyValue) {
	if inst.compiled == nil {
		return
	}

	// Note: Here, this is the place to use context, e.g., extract baggage.

	rec, updater := acquireRecord[N](inst, attrs)
	defer rec.refMapped.unref()

	updater.Update(num)

	// Record was modified.
	atomic.AddInt64(&rec.updateCount, 1)
}

// acquireRecord gets or creates a `*record` corresponding to `attrs`,
// the input attributes.
func acquireRecord[N number.Any](inst *Instrument, attrs []attribute.KeyValue) (*record, viewstate.Updater[N]) {
	aset := attribute.NewSet(attrs...)
	if lookup, ok := inst.current.Load(aset); ok {
		// Existing record case.
		rec := lookup.(*record)

		if rec.refMapped.ref() {
			// At this moment it is guaranteed that the
			// record is in the map and will not be removed.
			return rec, rec.accumulator.(viewstate.Updater[N])
		}
		// This record is no longer mapped, try to add a new
		// record below.
	}

	newRec := &record{
		refMapped: refcountMapped{value: 2},
	}

	for {
		if found, loaded := inst.current.LoadOrStore(aset, newRec); loaded {
			oldRec := found.(*record)
			if oldRec.refMapped.ref() {
				return oldRec, oldRec.accumulator.(viewstate.Updater[N])
			}
			runtime.Gosched()
			continue
		}
		break
	}

	newRec.accumulator = inst.compiled.NewAccumulator(aset)
	return newRec, newRec.accumulator.(viewstate.Updater[N])
}
