// Copyright 2019, OpenTelemetry Authors
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
	"unsafe"
)

func (l *sortedLabels) Len() int {
	return len(*l)
}

func (l *sortedLabels) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}

func (l *sortedLabels) Less(i, j int) bool {
	return (*l)[i].Key < (*l)[j].Key
}

func (m *SDK) addPrimary(rec *record) {
	for {
		rec.next.primary.store(m.records.primary.load())
		if atomic.CompareAndSwapPointer(
			&m.records.primary.ptr,
			rec.next.primary.ptr,
			unsafe.Pointer(rec),
		) {
			return
		}
	}
}

func (m *SDK) addReclaim(rec *record) {
	for {
		rec.next.reclaim.store(m.records.reclaim.load())
		if atomic.CompareAndSwapPointer(
			&m.records.reclaim.ptr,
			rec.next.reclaim.ptr,
			unsafe.Pointer(rec),
		) {
			return
		}
	}
}

func (s *singlePtr) swapNil() *record {
	for {
		newValue := unsafe.Pointer(nil)
		swapped := atomic.LoadPointer(&s.ptr)
		if atomic.CompareAndSwapPointer(&s.ptr, swapped, newValue) {
			return (*record)(swapped)
		}
	}
}

func (s *singlePtr) load() *record {
	return (*record)(atomic.LoadPointer(&s.ptr))
}

func (s *singlePtr) store(r *record) {
	atomic.StorePointer(&s.ptr, unsafe.Pointer(r))
}

func (s *singlePtr) clear() {
	atomic.StorePointer(&s.ptr, unsafe.Pointer(nil))
}
