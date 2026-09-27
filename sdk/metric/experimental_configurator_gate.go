// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import "sync/atomic"

// versionedEnabled holds a configuratorMeter's enabled bool and the version it
// was last set under.
type versionedEnabled struct {
	// state contains a 63-bit version + 1-bit enabled.
	// These are contained together so a write can atomically compare
	// the stored version against its own and decide whether to overwrite.
	// Gated Add/Record/Observe calls still do a single atomic load, with no
	// extra allocation or synchronization added to the hot path.
	state atomic.Uint64
}

// Load reports the currently stored enabled bool.
func (ve *versionedEnabled) Load() bool {
	return ve.state.Load()&1 != 0
}

// StoreIfNewer sets enabled if version is at least as new as the version
// currently stored, and reports whether it did. Ties are allowed to
// overwrite: a given version identifies a single configuration decision, so
// two writes sharing a version can only ever be redoing the same decision,
// never disagreeing ones.
func (ve *versionedEnabled) StoreIfNewer( // nolint:revive  // enabled is not a control flag.
	version uint64,
	enabled bool,
) bool {
	nextVersionedState := version << 1
	if enabled {
		nextVersionedState |= 1
	}
	for {
		oldState := ve.state.Load()
		if version < (oldState >> 1) {
			return false
		}
		if ve.state.CompareAndSwap(oldState, nextVersionedState) {
			return true
		}
	}
}
