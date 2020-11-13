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

package propagation_test

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/propagation"
)

type ctxKeyType uint

var (
	ctxKey ctxKeyType = 0
)

type carrier []string

func (c *carrier) Get(string) string { return "" }

func (c *carrier) Set(setter, _ string) {
	*c = append(*c, setter)
}

type propagator struct {
	Name string
}

func (p propagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	carrier.Set(p.Name, "")
}

func (p propagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	v := ctx.Value(ctxKey)
	if v == nil {
		ctx = context.WithValue(ctx, ctxKey, []string{p.Name})
	} else {
		orig := v.([]string)
		ctx = context.WithValue(ctx, ctxKey, append(orig, p.Name))
	}
	return ctx
}

func (p propagator) Fields() []string { return []string{p.Name} }

func TestCompositeTextMapPropagatorFields(t *testing.T) {
	a, b1, b2 := propagator{"a"}, propagator{"b"}, propagator{"b"}

	want := map[string]struct{}{
		"a": {},
		"b": {},
	}
	got := propagation.NewCompositeTextMapPropagator(a, b1, b2).Fields()
	if len(got) != len(want) {
		t.Fatalf("invalid fields from composite: %v (want %v)", got, want)
	}
	for _, v := range got {
		if _, ok := want[v]; !ok {
			t.Errorf("invalid field returned from composite: %q", v)
		}
	}
}

func TestCompositeTextMapPropagatorInject(t *testing.T) {
	a, b := propagator{"a"}, propagator{"b"}

	c := make(carrier, 0, 2)
	propagation.NewCompositeTextMapPropagator(a, b).Inject(context.Background(), &c)

	if got := strings.Join([]string(c), ","); got != "a,b" {
		t.Errorf("invalid inject order: %s", got)
	}
}

func TestCompositeTextMapPropagatorExtract(t *testing.T) {
	a, b := propagator{"a"}, propagator{"b"}

	ctx := context.Background()
	ctx = propagation.NewCompositeTextMapPropagator(a, b).Extract(ctx, nil)

	v := ctx.Value(ctxKey)
	if v == nil {
		t.Fatal("no composite extraction")
	}
	if got := strings.Join(v.([]string), ","); got != "a,b" {
		t.Errorf("invalid extract order: %s", got)
	}
}
