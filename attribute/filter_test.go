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

package attribute

import "testing"

func TestNewAllowKeysFilter(t *testing.T) {
	keys := []string{"zero", "one", "two"}
	attrs := []KeyValue{Int(keys[0], 0), Int(keys[1], 1), Int(keys[2], 2)}

	t.Run("Empty", func(t *testing.T) {
		empty := NewAllowKeysFilter()
		for _, kv := range attrs {
			if empty(kv) {
				t.Errorf("empty NewAllowKeysFilter filter accepted %v", kv)
			}
		}
	})

	t.Run("Partial", func(t *testing.T) {
		partial := NewAllowKeysFilter(Key(keys[0]), Key(keys[1]))
		for _, kv := range attrs[:2] {
			if !partial(kv) {
				t.Errorf("partial NewAllowKeysFilter filter denied %v", kv)
			}
		}
		if partial(attrs[2]) {
			t.Errorf("partial NewAllowKeysFilter filter accepted %v", attrs[2])
		}
	})

	t.Run("Full", func(t *testing.T) {
		full := NewAllowKeysFilter(Key(keys[0]), Key(keys[1]), Key(keys[2]))
		for _, kv := range attrs {
			if !full(kv) {
				t.Errorf("full NewAllowKeysFilter filter denied %v", kv)
			}
		}
	})
}

func TestNewDenyKeysFilter(t *testing.T) {
	keys := []string{"zero", "one", "two"}
	attrs := []KeyValue{Int(keys[0], 0), Int(keys[1], 1), Int(keys[2], 2)}

	t.Run("Empty", func(t *testing.T) {
		empty := NewDenyKeysFilter()
		for _, kv := range attrs {
			if !empty(kv) {
				t.Errorf("empty NewDenyKeysFilter filter denied %v", kv)
			}
		}
	})

	t.Run("Partial", func(t *testing.T) {
		partial := NewDenyKeysFilter(Key(keys[0]), Key(keys[1]))
		for _, kv := range attrs[:2] {
			if partial(kv) {
				t.Errorf("partial NewDenyKeysFilter filter accepted %v", kv)
			}
		}
		if !partial(attrs[2]) {
			t.Errorf("partial NewDenyKeysFilter filter denied %v", attrs[2])
		}
	})

	t.Run("Full", func(t *testing.T) {
		full := NewDenyKeysFilter(Key(keys[0]), Key(keys[1]), Key(keys[2]))
		for _, kv := range attrs {
			if full(kv) {
				t.Errorf("full NewDenyKeysFilter filter accepted %v", kv)
			}
		}
	})
}
