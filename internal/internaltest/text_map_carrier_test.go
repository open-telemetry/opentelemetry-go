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

package internaltest

import (
	"reflect"
	"testing"
)

var (
	key, value = "test", "true"
)

func TestTextMapCarrierKeys(t *testing.T) {
	tmc := NewTextMapCarrier(map[string]string{key: value})
	expected, actual := []string{key}, tmc.Keys()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected tmc.Keys() to be %v but it was %v", expected, actual)
	}
}

func TestTextMapCarrierGet(t *testing.T) {
	tmc := NewTextMapCarrier(map[string]string{key: value})
	tmc.GotN(t, 0)
	if got := tmc.Get("empty"); got != "" {
		t.Errorf("TextMapCarrier.Get returned %q for an empty key", got)
	}
	tmc.GotKey(t, "empty")
	tmc.GotN(t, 1)
	if got := tmc.Get(key); got != value {
		t.Errorf("TextMapCarrier.Get(%q) returned %q, want %q", key, got, value)
	}
	tmc.GotKey(t, key)
	tmc.GotN(t, 2)
}

func TestTextMapCarrierSet(t *testing.T) {
	tmc := NewTextMapCarrier(nil)
	tmc.SetN(t, 0)
	tmc.Set(key, value)
	if got, ok := tmc.data[key]; !ok {
		t.Errorf("TextMapCarrier.Set(%q,%q) failed to store pair", key, value)
	} else if got != value {
		t.Errorf("TextMapCarrier.Set(%q,%q) stored (%q,%q), not (%q,%q)", key, value, key, got, key, value)
	}
	tmc.SetKeyValue(t, key, value)
	tmc.SetN(t, 1)
}

func TestTextMapCarrierReset(t *testing.T) {
	tmc := NewTextMapCarrier(map[string]string{key: value})
	tmc.GotN(t, 0)
	tmc.SetN(t, 0)
	tmc.Reset(nil)
	tmc.GotN(t, 0)
	tmc.SetN(t, 0)
	if got := tmc.Get(key); got != "" {
		t.Error("TextMapCarrier.Reset() failed to clear initial data")
	}
	tmc.GotN(t, 1)
	tmc.GotKey(t, key)
	tmc.Set(key, value)
	tmc.SetKeyValue(t, key, value)
	tmc.SetN(t, 1)
	tmc.Reset(nil)
	tmc.GotN(t, 0)
	tmc.SetN(t, 0)
	if got := tmc.Get(key); got != "" {
		t.Error("TextMapCarrier.Reset() failed to clear data")
	}
}
