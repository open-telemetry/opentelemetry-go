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

package oteltest

import (
	"context"
	"testing"
)

var (
	key, value = "test", "true"
)

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
	tmc.Reset()
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
	tmc.Reset()
	tmc.GotN(t, 0)
	tmc.SetN(t, 0)
	if got := tmc.Get(key); got != "" {
		t.Error("TextMapCarrier.Reset() failed to clear data")
	}
}

func TestTextMapPropagatorInjectExtract(t *testing.T) {
	name := "testing"
	ctx := context.Background()
	carrier := NewTextMapCarrier(map[string]string{name: value})
	propagator := NewTextMapPropagator(name)

	propagator.Inject(ctx, carrier)
	// Carrier value overridden with state.
	if carrier.SetKeyValue(t, name, "1,0") {
		// Ensure nothing has been extracted yet.
		propagator.ExtractedN(t, ctx, 0)
		// Test the injection was counted.
		propagator.InjectedN(t, carrier, 1)
	}

	ctx = propagator.Extract(ctx, carrier)
	v := ctx.Value(ctxKeyType(name))
	if v == nil {
		t.Error("TextMapPropagator.Extract failed to extract state")
	}
	if s, ok := v.(state); !ok {
		t.Error("TextMapPropagator.Extract did not extract proper state")
	} else if s.Extractions != 1 {
		t.Error("TextMapPropagator.Extract did not increment state.Extractions")
	}
	if carrier.GotKey(t, name) {
		// Test the extraction was counted.
		propagator.ExtractedN(t, ctx, 1)
		// Ensure no additional injection was recorded.
		propagator.InjectedN(t, carrier, 1)
	}
}

func TestTextMapPropagatorFields(t *testing.T) {
	name := "testing"
	propagator := NewTextMapPropagator(name)
	if got := propagator.Fields(); len(got) != 1 {
		t.Errorf("TextMapPropagator.Fields returned %d fields, want 1", len(got))
	} else if got[0] != name {
		t.Errorf("TextMapPropagator.Fields returned %q, want %q", got[0], name)
	}
}

func TestNewStateEmpty(t *testing.T) {
	if want, got := (state{}), newState(""); got != want {
		t.Errorf("newState(\"\") returned %v, want %v", got, want)
	}
}
