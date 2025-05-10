// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package epoch stores time as nanoseconds since the Unix epoch as a more
// efficient alternative to time.Time.
package epoch // import "go.opentelemetry.io/otel/internal/epoch"

import (
	"sync/atomic"
	"time"
)

// Nanos is the nanoseconds since the Unix epoch.
//
//nolint:recvcheck
type Nanos int64

// NanosNowOrDefault returns nanoseconds from t. If t is the zero-value, returns
// the current, monotonic time in nanoseconds since the epoch.
func NanosNowOrDefault(t time.Time) Nanos {
	if t.IsZero() {
		return Nanos(time.Now().UnixNano())
	}
	return Nanos(t.UnixNano())
}

// ToTime converts the Nanos to a [time.Time].
func (u Nanos) ToTime() time.Time {
	if u == 0 {
		return time.Time{}
	}
	return time.Unix(0, int64(u))
}

// SwapIfZero atomically swaps the value of u with new if u is zero.
// Returns true if the value was swapped.
func (u *Nanos) SwapIfZero(ns Nanos) bool {
	p := (*int64)(u)
	return atomic.CompareAndSwapInt64(p, 0, int64(ns))
}

// Load atomically loads the value of u.
func (u *Nanos) Load() Nanos {
	return Nanos(atomic.LoadInt64((*int64)(u)))
}
