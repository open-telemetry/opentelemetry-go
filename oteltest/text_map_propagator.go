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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/propagation"
)

type ctxKeyType string

// TextMapCarrier provides a testing storage medium to for a
// TextMapPropagator. It records all the operations it performs.
type TextMapCarrier struct {
	mtx sync.Mutex

	gets []string
	sets [][2]string
	data map[string]string
}

// NewTextMapCarrier returns a new *TextMapCarrier populated with data.
func NewTextMapCarrier(data map[string]string) *TextMapCarrier {
	copied := make(map[string]string, len(data))
	for k, v := range data {
		copied[k] = v
	}
	return &TextMapCarrier{data: copied}
}

// Get returns the value associated with the passed key.
func (c *TextMapCarrier) Get(key string) string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.gets = append(c.gets, key)
	return c.data[key]
}

// GotKey tests if c.Get has been called for key.
func (c *TextMapCarrier) GotKey(t *testing.T, key string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, k := range c.gets {
		if k == key {
			return true
		}
	}
	t.Errorf("TextMapCarrier.Get(%q) has not been called", key)
	return false
}

// GotN tests if n calls to c.Get have been made.
func (c *TextMapCarrier) GotN(t *testing.T, n int) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if len(c.gets) != n {
		t.Errorf("TextMapCarrier.Get was called %d times, not %d", len(c.gets), n)
		return false
	}
	return true
}

// Set stores the key-value pair.
func (c *TextMapCarrier) Set(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.sets = append(c.sets, [2]string{key, value})
	c.data[key] = value
}

// SetKeyValue tests if c.Set has been called for the key-value pair.
func (c *TextMapCarrier) SetKeyValue(t *testing.T, key, value string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var vals []string
	for _, pair := range c.sets {
		if key == pair[0] {
			if value == pair[1] {
				return true
			}
			vals = append(vals, pair[1])
		}
	}
	if len(vals) > 0 {
		t.Errorf("TextMapCarrier.Set called with %q and %v values, but not %s", key, vals, value)
	}
	t.Errorf("TextMapCarrier.Set(%q,%q) has not been called", key, value)
	return false
}

// SetN tests if n calls to c.Set have been made.
func (c *TextMapCarrier) SetN(t *testing.T, n int) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if len(c.sets) != n {
		t.Errorf("TextMapCarrier.Set was called %d times, not %d", len(c.sets), n)
		return false
	}
	return true
}

// Reset zeros out the internal state recording of c.
func (c *TextMapCarrier) Reset() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.gets = nil
	c.sets = nil
	c.data = make(map[string]string)
}

type state struct {
	Injections  uint64
	Extractions uint64
}

func newState(encoded string) state {
	if encoded == "" {
		return state{}
	}
	split := strings.SplitN(encoded, ",", 2)
	injects, _ := strconv.ParseUint(split[0], 10, 64)
	extracts, _ := strconv.ParseUint(split[1], 10, 64)
	return state{
		Injections:  injects,
		Extractions: extracts,
	}
}

func (s state) String() string {
	return fmt.Sprintf("%d,%d", s.Injections, s.Extractions)
}

type TextMapPropagator struct {
	Name   string
	ctxKey ctxKeyType
}

func NewTextMapPropagator(name string) *TextMapPropagator {
	return &TextMapPropagator{Name: name, ctxKey: ctxKeyType(name)}
}

func (p *TextMapPropagator) stateFromContext(ctx context.Context) state {
	if v := ctx.Value(p.ctxKey); v != nil {
		if s, ok := v.(state); ok {
			return s
		}
	}
	return state{}
}

func (p *TextMapPropagator) stateFromCarrier(carrier propagation.TextMapCarrier) state {
	return newState(carrier.Get(p.Name))
}

// Inject set cross-cutting concerns for p from the Context into the carrier.
func (p *TextMapPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	s := p.stateFromContext(ctx)
	s.Injections++
	carrier.Set(p.Name, s.String())
}

// InjectedN tests if p has made n injections to carrier.
func (p *TextMapPropagator) InjectedN(t *testing.T, carrier *TextMapCarrier, n int) bool {
	if actual := p.stateFromCarrier(carrier).Injections; actual != uint64(n) {
		t.Errorf("TextMapPropagator{%q} injected %d times, not %d", p.Name, actual, n)
		return false
	}
	return true
}

// Extract reads cross-cutting concerns for p from the carrier into a Context.
func (p *TextMapPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	s := p.stateFromCarrier(carrier)
	s.Extractions++
	return context.WithValue(ctx, p.ctxKey, s)
}

// ExtractedN tests if p has made n extractions from the lineage of ctx.
// nolint (context is not first arg)
func (p *TextMapPropagator) ExtractedN(t *testing.T, ctx context.Context, n int) bool {
	if actual := p.stateFromContext(ctx).Extractions; actual != uint64(n) {
		t.Errorf("TextMapPropagator{%q} extracted %d time, not %d", p.Name, actual, n)
		return false
	}
	return true
}

// Fields returns p.Name as the key who's value is set with Inject.
func (p *TextMapPropagator) Fields() []string { return []string{p.Name} }
