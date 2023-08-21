// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/internaltest/text_map_carrier.go.tmpl

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

package internaltest // import "go.opentelemetry.io/otel/sdk/internal/internaltest"

import (
	"sync"
	"testing"

	"go.opentelemetry.io/otel/propagation"
)

// TextMapCarrier is a storage medium for a TextMapPropagator used in testing.
// The methods of a TextMapCarrier are concurrent safe.
type TextMapCarrier struct {
	mtx sync.Mutex

	gets []string
	sets [][2]string
	data map[string]string
}

var _ propagation.TextMapCarrier = (*TextMapCarrier)(nil)

// NewTextMapCarrier returns a new *TextMapCarrier populated with data.
func NewTextMapCarrier(data map[string]string) *TextMapCarrier {
	copied := make(map[string]string, len(data))
	for k, v := range data {
		copied[k] = v
	}
	return &TextMapCarrier{data: copied}
}

// Keys returns the keys for which this carrier has a value.
func (c *TextMapCarrier) Keys() []string {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	result := make([]string, 0, len(c.data))
	for k := range c.data {
		result = append(result, k)
	}
	return result
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

// Reset zeros out the recording state and sets the carried values to data.
func (c *TextMapCarrier) Reset(data map[string]string) {
	copied := make(map[string]string, len(data))
	for k, v := range data {
		copied[k] = v
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.gets = nil
	c.sets = nil
	c.data = copied
}
