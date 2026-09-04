// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type lifecycleState uint8

const (
	lifecycleActive lifecycleState = iota
	lifecycleFinishPending
	lifecycleRetired
	lifecycleCollecting
)

const (
	// Lifecycle.state packs the lifecycle state into the high two bits and the
	// number of admitted measurements into the low 62 bits. Keeping both in one
	// atomic value prevents collection from closing admission between a
	// measurement's state check and its increment of the admitted count.
	stateShift = 62
	writerMask = uint64(1<<stateShift) - 1
	stateMask  = ^writerMask
)

// Lifecycle coordinates measurements, explicit finishes, collections, and
// retirement for one metric series. Its zero value is ready for use.
//
// Lifecycle methods are safe to call concurrently. A successful measurement
// acquisition and an emitted collection must each be completed exactly once.
// A Lifecycle must not be copied after first use.
//
// The valid transitions are active to finish-pending on Finish, finish-pending
// to active on measurement, active or finish-pending to collecting on
// collection, and collecting to active or retired when collection completes.
// Retired is terminal.
type Lifecycle struct {
	control sync.Mutex
	state   atomic.Uint64
	// finished is meaningful only while a pending finish is being collected.
	// Reactivation leaves the old value in place so measurements remain lock-free.
	finished time.Time
}

// Measurement represents admission to update one metric-series lifetime. A
// successful call to Lifecycle.AcquireMeasurement returns a Measurement, which
// must be released exactly once after the aggregate and exemplars are updated.
type Measurement struct {
	lifecycle *Lifecycle
}

// Release completes the admitted measurement.
func (m Measurement) Release() {
	m.lifecycle.releaseMeasurement()
}

// Collection describes an exclusive collection phase for a metric-series
// lifetime. If ShouldEmit reports true, Complete must be called exactly once
// after the aggregate and exemplars are read.
type Collection struct {
	lifecycle *Lifecycle
	time      time.Time
	emit      bool
	retire    bool
}

// Time returns the timestamp to use for the collected point.
func (c Collection) Time() time.Time {
	return c.time
}

// ShouldEmit reports whether the collection should emit a point.
func (c Collection) ShouldEmit() bool {
	return c.emit
}

// ShouldRetire reports whether Complete will retire the series lifetime.
func (c Collection) ShouldRetire() bool {
	return c.retire
}

// Complete ends the exclusive collection phase. It must be called exactly once
// for a Collection whose ShouldEmit method reports true.
func (c Collection) Complete() {
	state := lifecycleActive
	if c.retire {
		state = lifecycleRetired
	}
	c.lifecycle.finished = time.Time{}
	c.lifecycle.state.Store(withState(0, state))
	c.lifecycle.control.Unlock()
}

func stateOf(value uint64) lifecycleState {
	return lifecycleState(value >> stateShift)
}

func withState(value uint64, state lifecycleState) uint64 {
	return value&writerMask | uint64(state)<<stateShift
}

// AcquireMeasurement reserves this series lifetime for a measurement. It waits
// for an in-progress collection and returns false only if the lifetime is
// retired. A successful acquisition must be paired with exactly one call to
// Measurement.Release after the aggregate and exemplars are updated.
func (l *Lifecycle) AcquireMeasurement() (Measurement, bool) {
measurement:
	for {
		// Do not join the writer count after collection closes admission. This
		// lets collection drain the finite set of already-admitted measurements
		// even when new measurements continue to arrive.
		state := stateOf(l.state.Load())
		if state == lifecycleCollecting {
			l.waitForCollection()
			continue
		}
		if state == lifecycleRetired {
			return Measurement{}, false
		}

		previous := l.state.Add(1) - 1
		switch stateOf(previous) {
		case lifecycleActive:
			return Measurement{lifecycle: l}, true
		case lifecycleFinishPending:
			for {
				current := l.state.Load()
				switch stateOf(current) {
				case lifecycleActive:
					return Measurement{lifecycle: l}, true
				case lifecycleFinishPending:
					if l.state.CompareAndSwap(current, withState(current, lifecycleActive)) {
						return Measurement{lifecycle: l}, true
					}
				case lifecycleCollecting:
					l.releaseMeasurement()
					l.waitForCollection()
					continue measurement
				default:
					l.releaseMeasurement()
					return Measurement{}, false
				}
			}
		case lifecycleCollecting:
			l.releaseMeasurement()
			l.waitForCollection()
			continue measurement
		default:
			l.releaseMeasurement()
			return Measurement{}, false
		}
	}
}

func (l *Lifecycle) waitForCollection() {
	for stateOf(l.state.Load()) == lifecycleCollecting {
		runtime.Gosched()
	}
}

func (l *Lifecycle) releaseMeasurement() {
	// atomic.Uint64 has no subtraction operation. Adding MaxUint64 decrements
	// the packed writer count modulo 2^64. A matching acquisition guarantees
	// the count is nonzero, so the lifecycle bits are unaffected.
	l.state.Add(^uint64(0))
}

// Finish marks an active series lifetime to be retired by its next collection
// at the supplied time. It returns false if the lifetime is already pending or
// retired. A measurement admitted after Finish reactivates the lifetime.
func (l *Lifecycle) Finish(at time.Time) bool {
	l.control.Lock()
	defer l.control.Unlock()
	if stateOf(l.state.Load()) != lifecycleActive {
		return false
	}
	l.finished = at
	// The control lock keeps the lifecycle active until this atomic update.
	// Measurements may concurrently change only the writer count; Or preserves
	// those changes and is the linearization point for this finish.
	l.state.Or(uint64(lifecycleFinishPending) << stateShift)
	return true
}

// BeginCumulativeCollection closes measurement admission and waits for admitted
// measurements to complete. An active lifetime remains active when the returned
// Collection is completed; a finish-pending lifetime is retired.
func (l *Lifecycle) BeginCumulativeCollection(at time.Time) Collection {
	l.control.Lock()
	collection := l.beginCollection(at)
	if !collection.emit {
		l.control.Unlock()
	}
	return collection
}

// BeginDeltaCollection closes measurement admission and waits for admitted
// measurements to complete. Every emitted lifetime is retired when the returned
// Collection is completed.
func (l *Lifecycle) BeginDeltaCollection(at time.Time) Collection {
	l.control.Lock()
	collection := l.beginCollection(at)
	if !collection.emit {
		l.control.Unlock()
		return collection
	}
	collection.retire = true
	return collection
}

func (l *Lifecycle) beginCollection(at time.Time) Collection {
	state := stateOf(l.state.Load())
	if state == lifecycleRetired {
		return Collection{}
	}
	state = stateOf(l.state.Or(stateMask))
	for l.state.Load()&writerMask != 0 {
		runtime.Gosched()
	}

	if state == lifecycleFinishPending {
		return Collection{
			lifecycle: l,
			time:      l.finished,
			emit:      true,
			retire:    true,
		}
	}
	return Collection{lifecycle: l, time: at, emit: true}
}

// Retire permanently closes the series lifetime. It waits for admitted
// measurements to complete and has no effect if the lifetime is already
// retired.
func (l *Lifecycle) Retire() {
	l.control.Lock()
	defer l.control.Unlock()
	if stateOf(l.state.Load()) == lifecycleRetired {
		return
	}
	l.state.Or(stateMask)
	for l.state.Load()&writerMask != 0 {
		runtime.Gosched()
	}
	l.finished = time.Time{}
	l.state.Store(withState(0, lifecycleRetired))
}
