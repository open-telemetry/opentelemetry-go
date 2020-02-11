// Copyright 2020, OpenTelemetry Authors
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

package correlation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type correlationsType struct{}

var correlationsKey = &correlationsType{}

// WithMap returns a context with the Map entered into it.
func WithMap(ctx context.Context, m Map) context.Context {
	return context.WithValue(ctx, correlationsKey, m)
}

// NewContext returns a context with the map from passed context
// updated with the passed key-value pairs.
func NewContext(ctx context.Context, keyvalues ...core.KeyValue) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(MapUpdate{
		MultiKV: keyvalues,
	}))
}

// FromContext gets the current Map from a Context.
func FromContext(ctx context.Context) Map {
	if m, ok := ctx.Value(correlationsKey).(Map); ok {
		return m
	}
	return NewEmptyMap()
}
