// Copyright 2019, OpenTelemetry Authors
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

package distributedcontext

import (
	"context"
	"runtime/pprof"

	"go.opentelemetry.io/otel/api/core"
)

type ctxEntriesType struct{}

var (
	ctxEntriesKey = &ctxEntriesType{}
)

func WithMap(ctx context.Context, m Map) context.Context {
	return context.WithValue(ctx, ctxEntriesKey, m)
}

func NewContext(ctx context.Context, keyvalues ...core.KeyValue) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(MapUpdate{
		MultiKV: keyvalues,
	}))
}

func FromContext(ctx context.Context) Map {
	if m, ok := ctx.Value(ctxEntriesKey).(Map); ok {
		return m
	}
	return NewEmptyMap()
}

// Note: the golang pprof.Do API forces this memory allocation, we
// should file an issue about that.  (There's a TODO in the source.)
func Do(ctx context.Context, f func(ctx context.Context)) {
	m := FromContext(ctx)
	keyvals := make([]string, 0, 2*len(m.m))
	for k, v := range m.m {
		keyvals = append(keyvals, string(k), v.value.Emit())
	}
	pprof.Do(ctx, pprof.Labels(keyvals...), f)
}
