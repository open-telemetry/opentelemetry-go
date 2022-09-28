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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"errors"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric/internal"
)

// cache is a locking storage used to quickly return already computed values. A
// registry type should be used with a cache for get and set operations of
// certain types.
//
// The zero value of a cache is empty and ready to use.
//
// A cache must not be copied after first use.
//
// All methods of a cache are safe to call concurrently.
type cache[K comparable, V any] struct {
	sync.Mutex
	data map[K]V
}

// GetOrSet returns the value stored in the cache for key if it exists.
// Otherwise, f is called and the returned value is set in the cache for key
// and returned.
//
// GetOrSet is safe to call concurrently. It will hold the cache lock, so f
// should not block excessively.
func (c *cache[K, V]) GetOrSet(key K, f func() V) V {
	c.Lock()
	defer c.Unlock()

	if c.data == nil {
		val := f()
		c.data = map[K]V{key: val}
		return val
	}
	if v, ok := c.data[key]; ok {
		return v
	}
	val := f()
	c.data[key] = val
	return val
}

// resolvedAggregators is the result of resolving aggregators for an instrument.
type resolvedAggregators[N int64 | float64] struct {
	aggregators []internal.Aggregator[N]
	err         error
}

type instrumentRegistry[N int64 | float64] struct {
	c *cache[instrumentID, any]
}

func newInstrumentRegistry[N int64 | float64](c *cache[instrumentID, any]) instrumentRegistry[N] {
	if c == nil {
		c = &cache[instrumentID, any]{}
	}
	return instrumentRegistry[N]{c: c}
}

var errExists = errors.New("instrument already exists for different number type")

func (q instrumentRegistry[N]) GetOrSet(key instrumentID, f func() ([]internal.Aggregator[N], error)) (aggs []internal.Aggregator[N], err error) {
	vAny := q.c.GetOrSet(key, func() any {
		a, err := f()
		return &resolvedAggregators[N]{
			aggregators: a,
			err:         err,
		}
	})

	switch v := vAny.(type) {
	case *resolvedAggregators[N]:
		aggs = v.aggregators
	default:
		err = errExists
	}
	return aggs, err
}
