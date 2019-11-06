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
)

type (
	ctxCorrelationsType struct{}
	ctxBaggageType      struct{}
)

var (
	ctxCorrelationsKey = &ctxCorrelationsType{}
	ctxBaggageKey      = &ctxBaggageType{}
)

func NewCorrelationsContextKV(ctx context.Context, correlations ...Correlation) context.Context {
	return withCorrelationsMapUpdate(ctx, CorrelationsUpdate{
		MultiKV: correlations,
	})
}

func NewCorrelationsContextMap(ctx context.Context, correlations Correlations) context.Context {
	return withCorrelationsMapUpdate(ctx, CorrelationsUpdate{
		Map: correlations,
	})
}

func CorrelationsFromContext(ctx context.Context) Correlations {
	if m, ok := ctx.Value(ctxCorrelationsKey).(Correlations); ok {
		return m
	}
	return NewEmptyCorrelations()
}

func NewBaggageContext(ctx context.Context, baggage Baggage) context.Context {
	return context.WithValue(ctx, ctxBaggageKey, baggage)
}

func BaggageFromContext(ctx context.Context) Baggage {
	if m, ok := ctx.Value(ctxBaggageKey).(Baggage); ok {
		return m
	}
	return NewEmptyBaggage()
}

func withCorrelationsMapUpdate(ctx context.Context, update CorrelationsUpdate) context.Context {
	return context.WithValue(ctx, ctxCorrelationsKey, CorrelationsFromContext(ctx).Apply(update))
}

/*
// TODO(krnowak): I don't know what's the point of this function…

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
*/
