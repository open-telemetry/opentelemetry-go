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

package trace

import (
	"reflect"
	"testing"
)

func init() {
}

func TestAdd(t *testing.T) {
	q := newEvictedQueue(3)
	q.Add("value1")
	q.Add("value2")
	if wantLen, gotLen := 2, len(q.Queue()); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
}

func (eq *evictedQueue) queueToArray() []string {
	arr := make([]string, 0)
	for _, value := range eq.Queue() {
		arr = append(arr, value.(string))
	}
	return arr
}

func TestDropCount(t *testing.T) {
	q := newEvictedQueue(3)
	q.Add("value1")
	q.Add("value2")
	q.Add("value3")
	q.Add("value1")
	q.Add("value4")
	if wantLen, gotLen := 3, len(q.Queue()); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
	if wantDropCount, gotDropCount := 2, q.DroppedCount(); wantDropCount != gotDropCount {
		t.Errorf("got drop count %d want %d", gotDropCount, wantDropCount)
	}
	wantArr := []string{"value3", "value1", "value4"}
	gotArr := q.queueToArray()

	if wantLen, gotLen := len(wantArr), len(gotArr); gotLen != wantLen {
		t.Errorf("got array len %d want %d", gotLen, wantLen)
	}

	if !reflect.DeepEqual(gotArr, wantArr) {
		t.Errorf("got array = %#v; want %#v", gotArr, wantArr)
	}
}
