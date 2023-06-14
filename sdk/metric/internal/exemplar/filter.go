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

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/internal/exemplar"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Filter determines filters measurements passed to a [Reservoir]. If true is
// return, the measurement will be considered for sampling.
//
// See [Filtered] for how to create a [Reservoir] that uses a Filter.
type Filter[N int64 | float64] func(context.Context, N, attribute.Set) bool

// AlwaysSample is a Filter that always signals measurements should be
// considered for sampling by a [Reservoir].
func AlwaysSample[N int64 | float64](context.Context, N, attribute.Set) bool {
	return true
}

// NeverSample is a Filter that always signals measurements should not be
// considered for sampling by a [Reservoir].
func NeverSample[N int64 | float64](context.Context, N, attribute.Set) bool {
	return false
}

// TraceBasedSample is a Filter that signals measurements should be considered
// for sampling by a [Reservoir] if the ctx contains a
// [go.opentelemetry.io/otel/trace.SpanContext] that is sampled.
func TraceBasedSample[N int64 | float64](ctx context.Context, _ N, _ attribute.Set) bool {
	return trace.SpanContextFromContext(ctx).IsSampled()
}

// Filtered returns a [Reservoir] wrapping r that will only offer measurements
// to r if f returns true.
func Filtered[N int64 | float64](r Reservoir[N], f Filter[N]) Reservoir[N] {
	return filtered[N]{Reservoir: r, Filter: f}
}

type filtered[N int64 | float64] struct {
	Reservoir[N]

	Filter Filter[N]
}

func (f filtered[N]) Offer(ctx context.Context, t time.Time, n N, a attribute.Set) {
	if f.Filter(ctx, n, a) {
		f.Reservoir.Offer(ctx, t, n, a)
	}
}
