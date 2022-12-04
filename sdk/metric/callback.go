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
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
)

type callback[N int64 | float64] struct {
	observe func(context.Context, N, ...attribute.KeyValue)
	newIter func(context.Context) (iterator[N], error)
}

func (c callback[N]) collect(ctx context.Context) error {
	iter, err := c.newIter(ctx)
	if err != nil {
		return err
	}

	for iter.Next() {
		val, attrs := iter.Yield()
		c.observe(ctx, val, attrs...)
	}
	return nil
}

type iterator[N int64 | float64] interface {
	Len() int
	Next() bool
	Yield() (N, []attribute.KeyValue)
}

func newInt64Iter(f asyncint64.Callback) func(context.Context) (iterator[int64], error) {
	return func(ctx context.Context) (iterator[int64], error) {
		o, err := f(ctx)
		return iterInt64{idx: -1, observ: o}, err
	}
}

type iterInt64 struct {
	idx    int
	observ []asyncint64.Observation
}

func (i iterInt64) Len() int {
	return len(i.observ)
}

func (i iterInt64) Next() bool {
	i.idx++
	return i.idx < i.Len()
}

func (i iterInt64) Yield() (int64, []attribute.KeyValue) {
	if i.observ == nil || i.idx < 0 || i.idx >= i.Len() {
		return 0, nil
	}
	o := i.observ[i.idx]
	return o.Value, o.Attributes
}

func newFloat64Iter(f asyncfloat64.Callback) func(context.Context) (iterator[float64], error) {
	return func(ctx context.Context) (iterator[float64], error) {
		o, err := f(ctx)
		return iterFloat64{idx: -1, observ: o}, err
	}
}

type iterFloat64 struct {
	idx    int
	observ []asyncfloat64.Observation
}

func (i iterFloat64) Len() int {
	return len(i.observ)
}

func (i iterFloat64) Next() bool {
	i.idx++
	return i.idx < i.Len()
}

func (i iterFloat64) Yield() (float64, []attribute.KeyValue) {
	if i.observ == nil || i.idx < 0 || i.idx >= i.Len() {
		return 0, nil
	}
	o := i.observ[i.idx]
	return o.Value, o.Attributes
}
