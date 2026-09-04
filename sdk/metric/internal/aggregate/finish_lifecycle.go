// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"runtime"
	"sync/atomic"
	"time"
)

type finishLifecycleState uint8

const (
	lifecycleActive finishLifecycleState = iota
	lifecycleFinishPending
	lifecycleRetired
	lifecycleCollecting
)

const (
	// finishLifecycle.state packs the lifecycle state into the high two bits and
	// the number of admitted measurements into the low 62 bits. Keeping both in
	// one atomic value prevents collection from closing admission between a
	// measurement's state check and its increment of the admitted count.
	finishStateShift = 62
	finishWriterMask = uint64(1<<finishStateShift) - 1
	finishStateMask  = ^finishWriterMask
)

type finishCollectionMode uint8

const (
	finishCollectionKeepActive finishCollectionMode = iota
	finishCollectionRetireActive
)

type finishCollection struct {
	time   time.Time
	emit   bool
	retire bool
}

// finishLifecycle tracks the state of a finish-aware metric series. Its caller
// must serialize finish, collection, and retirement. Measurements are admitted
// atomically and may run concurrently.
//
// The valid transitions are active to pending on finish, pending to active on
// measurement, active or pending to collecting on collection, and collecting
// to active or retired when collection completes. Retired is terminal.
type finishLifecycle struct {
	state atomic.Uint64
	// finished is meaningful only while a pending finish is being collected.
	// Reactivation leaves the old value in place so measurements remain lock-free.
	finished time.Time
}

func finishStateOf(value uint64) finishLifecycleState {
	return finishLifecycleState(value >> finishStateShift)
}

func finishWithState(value uint64, state finishLifecycleState) uint64 {
	return value&finishWriterMask | uint64(state)<<finishStateShift
}

func (l *finishLifecycle) startMeasurement() bool {
measurement:
	for {
		// Do not join the writer count after collection closes admission. This
		// lets collection drain the finite set of already-admitted measurements
		// even when new measurements continue to arrive.
		state := finishStateOf(l.state.Load())
		if state == lifecycleCollecting {
			for finishStateOf(l.state.Load()) == lifecycleCollecting {
				runtime.Gosched()
			}
			continue
		}
		if state == lifecycleRetired {
			return false
		}

		previous := l.state.Add(1) - 1
		switch finishStateOf(previous) {
		case lifecycleActive:
			return true
		case lifecycleFinishPending:
			for {
				current := l.state.Load()
				switch finishStateOf(current) {
				case lifecycleActive:
					return true
				case lifecycleFinishPending:
					if l.state.CompareAndSwap(current, finishWithState(current, lifecycleActive)) {
						return true
					}
				case lifecycleCollecting:
					l.finishMeasurement()
					for finishStateOf(l.state.Load()) == lifecycleCollecting {
						runtime.Gosched()
					}
					continue measurement
				default:
					l.finishMeasurement()
					return false
				}
			}
		case lifecycleCollecting:
			l.finishMeasurement()
			for finishStateOf(l.state.Load()) == lifecycleCollecting {
				runtime.Gosched()
			}
			continue measurement
		default:
			l.finishMeasurement()
			return false
		}
	}
}

func (l *finishLifecycle) finishMeasurement() {
	// Adding the maximum uint64 value is equivalent to subtracting one.
	l.state.Add(^uint64(0))
}

func (l *finishLifecycle) finish(at time.Time) bool {
	if finishStateOf(l.state.Load()) != lifecycleActive {
		return false
	}
	l.finished = at
	l.state.Or(uint64(lifecycleFinishPending) << finishStateShift)
	return true
}

func (l *finishLifecycle) startCollection(
	at time.Time,
	mode finishCollectionMode,
) finishCollection {
	state := finishStateOf(l.state.Load())
	if state == lifecycleRetired {
		return finishCollection{}
	}
	state = finishStateOf(l.state.Or(finishStateMask))
	for l.state.Load()&finishWriterMask != 0 {
		runtime.Gosched()
	}

	if state == lifecycleFinishPending {
		return finishCollection{time: l.finished, emit: true, retire: true}
	}
	return finishCollection{
		time:   at,
		emit:   true,
		retire: mode == finishCollectionRetireActive,
	}
}

func (l *finishLifecycle) completeCollection(collection finishCollection) {
	state := lifecycleActive
	if collection.retire {
		state = lifecycleRetired
	}
	l.finished = time.Time{}
	l.state.Store(finishWithState(0, state))
}

func (l *finishLifecycle) retire() {
	if finishStateOf(l.state.Load()) == lifecycleRetired {
		return
	}
	l.state.Or(finishStateMask)
	for l.state.Load()&finishWriterMask != 0 {
		runtime.Gosched()
	}
	l.finished = time.Time{}
	l.state.Store(finishWithState(0, lifecycleRetired))
}
