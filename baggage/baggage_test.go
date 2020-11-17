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

package baggage

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/internal/baggage"
	"go.opentelemetry.io/otel/label"
)

func TestBaggage(t *testing.T) {
	ctx := context.Background()
	ctx = baggage.ContextWithMap(ctx, baggage.NewEmptyMap())

	b := Set(ctx)
	if b.Len() != 0 {
		t.Fatalf("empty baggage returned a set with %d elements", b.Len())
	}

	first, second, third := label.Key("first"), label.Key("second"), label.Key("third")
	ctx = ContextWithValues(ctx, first.Bool(true), second.String("2"))
	m := baggage.MapFromContext(ctx)
	v, ok := m.Value(first)
	if !ok {
		t.Fatal("WithValues failed to set first value")
	}
	if !v.AsBool() {
		t.Fatal("WithValues failed to set first correct value")
	}
	v, ok = m.Value(second)
	if !ok {
		t.Fatal("WithValues failed to set second value")
	}
	if v.AsString() != "2" {
		t.Fatal("WithValues failed to set second correct value")
	}
	_, ok = m.Value(third)
	if ok {
		t.Fatal("WithValues set an unexpected third value")
	}

	b = Set(ctx)
	if b.Len() != 2 {
		t.Fatalf("Baggage returned a set with %d elements, want 2", b.Len())
	}

	v = Value(ctx, first)
	if v.Type() != label.BOOL || !v.AsBool() {
		t.Fatal("Value failed to get correct first value")
	}
	v = Value(ctx, second)
	if v.Type() != label.STRING || v.AsString() != "2" {
		t.Fatal("Value failed to get correct second value")
	}

	ctx = ContextWithoutValues(ctx, first)
	m = baggage.MapFromContext(ctx)
	_, ok = m.Value(first)
	if ok {
		t.Fatal("WithoutValues failed to remove a baggage value")
	}
	_, ok = m.Value(second)
	if !ok {
		t.Fatal("WithoutValues removed incorrect value")
	}

	ctx = ContextWithEmpty(ctx)
	m = baggage.MapFromContext(ctx)
	if m.Len() != 0 {
		t.Fatal("WithoutBaggage failed to clear baggage")
	}
}
