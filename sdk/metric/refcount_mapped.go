// Copyright 2020, OpenTelemetry Authors
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

package metric

import (
	"sync/atomic"
)

// refcountMapped atomically counts the number of references (usages) of an entry
// while also keeping a state of mapped/unmapped into a different datastructure.
//
// ref() will fail if the entry is unmapped.
// unmap() will fail if the entry is in use.
//
// refcountMapped uses an atomic value where the least significant bit is used to
// keep the state of mapping and the rest of the bits are used for recounting.
type refcountMapped struct {
	// refcount has to be aligned for 64-bit atomic operations.
	value int64
}

// ref returns true if the entry is still mapped and increases the
// reference usages, if unmapped returns false.
func (rm *refcountMapped) ref() bool {
	if atomic.AddInt64(&rm.value, 2)&1 != 0 {
		// This entry was removed from the map between the moment
		// we got a reference to it (or will be removed very soon)
		// and here.
		return false
	}
	// At this moment it is guaranteed that the entry is in
	// the map and referenced (so it will not be unmapped).
	return true
}

func (rm *refcountMapped) unref() {
	atomic.AddInt64(&rm.value, -2)
}

// inUse returns true if there is a reference to the entry and it is not unmapped.
func (rm *refcountMapped) inUse() bool {
	val := atomic.LoadInt64(&rm.value)
	return val >= 2 && val&1 == 0
}

// unmap returns true if no references are active, and the
func (rm *refcountMapped) tryUnmap() bool {
	if atomic.LoadInt64(&rm.value) != 0 {
		return false
	}
	return atomic.CompareAndSwapInt64(
		&rm.value,
		0,
		1,
	)
}
