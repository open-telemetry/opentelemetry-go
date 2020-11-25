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

package global_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/oteltest"
)

func TestTextMapPropagatorDelegation(t *testing.T) {
	global.ResetForTest()
	ctx := context.Background()
	carrier := oteltest.NewTextMapCarrier(nil)

	// The default should be a noop.
	initial := global.TextMapPropagator()
	initial.Inject(ctx, carrier)
	ctx = initial.Extract(ctx, carrier)
	if !carrier.GotN(t, 0) || !carrier.SetN(t, 0) {
		return
	}

	// Make sure the delegate woks as expected.
	delegate := oteltest.NewTextMapPropagator("test")
	delegate.Inject(ctx, carrier)
	ctx = delegate.Extract(ctx, carrier)
	if !delegate.InjectedN(t, carrier, 1) || !delegate.ExtractedN(t, ctx, 1) {
		return
	}

	// The initial propagator should use the delegate after it is set as the
	// global.
	global.SetTextMapPropagator(delegate)
	initial.Inject(ctx, carrier)
	ctx = initial.Extract(ctx, carrier)
	delegate.InjectedN(t, carrier, 2)
	delegate.ExtractedN(t, ctx, 2)
}

func TestTextMapPropagatorDelegationNil(t *testing.T) {
	global.ResetForTest()
	ctx := context.Background()
	carrier := oteltest.NewTextMapCarrier(nil)

	// The default should be a noop.
	initial := global.TextMapPropagator()
	initial.Inject(ctx, carrier)
	ctx = initial.Extract(ctx, carrier)
	if !carrier.GotN(t, 0) || !carrier.SetN(t, 0) {
		return
	}

	// Delegation to nil should not make a change.
	global.SetTextMapPropagator(nil)
	initial.Inject(ctx, carrier)
	initial.Extract(ctx, carrier)
	if !carrier.GotN(t, 0) || !carrier.SetN(t, 0) {
		return
	}
}

func TestTextMapPropagatorFields(t *testing.T) {
	global.ResetForTest()
	initial := global.TextMapPropagator()
	delegate := oteltest.NewTextMapPropagator("test")
	delegateFields := delegate.Fields()

	// Sanity check on the initial Fields.
	if got := initial.Fields(); fieldsEqual(got, delegateFields) {
		t.Fatalf("testing fields (%v) matched Noop fields (%v)", delegateFields, got)
	}
	global.SetTextMapPropagator(delegate)
	// Check previous returns from global not correctly delegate.
	if got := initial.Fields(); !fieldsEqual(got, delegateFields) {
		t.Errorf("global TextMapPropagator.Fields returned %v instead of delegating, want (%v)", got, delegateFields)
	}
	// Check new calls to global.
	if got := global.TextMapPropagator().Fields(); !fieldsEqual(got, delegateFields) {
		t.Errorf("global TextMapPropagator.Fields returned %v, want (%v)", got, delegateFields)
	}
}

func fieldsEqual(f1, f2 []string) bool {
	if len(f1) != len(f2) {
		return false
	}
	for i := range f1 {
		if f1[i] != f2[i] {
			return false
		}
	}
	return true
}
