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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// Performance note: there is still 1 obligatory allocation in the
// fast path of this code due to the sync.Map key.  Assuming Go will
// give us a generic form of sync.Map some time soon, the allocation
// cost of instrument.Current will be reduced to zero allocs in the
// fast path.  See also https://github.com/a8m/syncmap.

type (
	Instrument struct {
		instrument.Synchronous

		descriptor sdkinstrument.Descriptor
		compiled   viewstate.Instrument
		current    sync.Map // map[attribute.Set]*record
	}

	record struct {
		instrument.Synchronous

		refMapped  refcountMapped
		instrument *Instrument

		// updateCount is incremented on every Update.
		updateCount int64

		// collectedCount is set to updateCount on collection,
		// supports checking for no updates during a round.
		collectedCount int64

		accumulator viewstate.Accumulator
	}

	counter[N number.Any, Traits traits.Any[N]] struct {
		*Instrument
	}

	histogram[N number.Any, Traits traits.Any[N]] struct {
		*Instrument
	}
)

var (
	_ syncint64.Counter         = counter[int64, traits.Int64]{}
	_ syncint64.UpDownCounter   = counter[int64, traits.Int64]{}
	_ syncint64.Histogram       = histogram[int64, traits.Int64]{}
	_ syncfloat64.Counter       = counter[float64, traits.Float64]{}
	_ syncfloat64.UpDownCounter = counter[float64, traits.Float64]{}
	_ syncfloat64.Histogram     = histogram[float64, traits.Float64]{}
)

func NewInstrument(desc sdkinstrument.Descriptor, compiled viewstate.Instrument) *Instrument {
	return &Instrument{
		descriptor: desc,
		compiled:   compiled,
	}
}

func (inst *Instrument) Descriptor() sdkinstrument.Descriptor {
	return inst.descriptor
}

func NewCounter[N number.Any, Traits traits.Any[N]](inst *Instrument) counter[N, Traits] {
	return counter[N, Traits]{Instrument: inst}
}

func NewHistogram[N number.Any, Traits traits.Any[N]](inst *Instrument) histogram[N, Traits] {
	return histogram[N, Traits]{Instrument: inst}
}

func (c counter[N, Traits]) Add(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	if c.Instrument != nil {
		capture[N, Traits](ctx, c.Instrument, incr, attrs)
	}
}

func (h histogram[N, Traits]) Record(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	if h.Instrument != nil {
		capture[N, Traits](ctx, h.Instrument, incr, attrs)
	}
}

func (inst *Instrument) Collect(r *reader.Reader, sequence reader.Sequence, output *[]reader.Instrument) {
	inst.current.Range(func(key interface{}, value interface{}) bool {
		rec := value.(*record)
		any := inst.collectRecord(rec)

		if any != 0 {
			return true
		}
		// Having no updates since last collection, try to unmap:
		if unmapped := rec.refMapped.tryUnmap(); !unmapped {
			// The record is referenced by a binding, continue.
			return true
		}

		// If any other goroutines are now trying to re-insert this
		// entry in the map, they are busy calling Gosched() awaiting
		// this deletion:
		inst.current.Delete(key)

		// Last we'll see of this.
		_ = inst.collectRecord(rec)
		return true
	})
	inst.compiled.Collect(r, sequence, output)
}

func (inst *Instrument) collectRecord(rec *record) int {
	mods := atomic.LoadInt64(&rec.updateCount)
	coll := rec.collectedCount

	if mods == coll {
		return 0
	}
	// Updates happened in this interval,
	// collect and continue.
	rec.collectedCount = mods

	if rec.accumulator == nil {
		return 0
	}
	rec.accumulator.Accumulate()
	return 1
}

func capture[N number.Any, Traits traits.Any[N]](_ context.Context, inst *Instrument, num N, attrs []attribute.KeyValue) {
	// TODO: Here, this is the place to use context, extract baggage.

	rec, updater := acquireRecord[N](inst, attrs)
	defer rec.refMapped.unref()

	if err := aggregation.RangeTest[N, Traits](num, &rec.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	updater.(viewstate.AccumulatorUpdater[N]).Update(num)

	// Record was modified.
	atomic.AddInt64(&rec.updateCount, 1)
}

// acquireRecord gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
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
		// This group is no longer mapped, try
		// to add a new group below.
	}

	newRec := &record{
		refMapped:  refcountMapped{value: 2},
		instrument: inst,
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

	newRec.accumulator = inst.compiled.NewAccumulator(attrs, nil)

	return newRec, newRec.accumulator.(viewstate.Updater[N])
}
