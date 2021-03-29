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

package trace // import "go.opentelemetry.io/otel/trace"

import "context"

type traceContextKeyType int

const currentSpanKey traceContextKeyType = iota

// ContextWithSpan returns a copy of parent with span set as the current Span.
func ContextWithSpan(parent context.Context, span Span) context.Context {
	return context.WithValue(parent, currentSpanKey, span)
}

// ContextWithRemoteSpanContext returns a copy of parent with rsc set explicly
// as a remote SpanContext and as the current Span. The Span implementation
// that wraps rsc is non-recording and performs no operations other than to
// return rsc as the SpanContext from the SpanContext method.
func ContextWithRemoteSpanContext(parent context.Context, rsc SpanContext) context.Context {
	return ContextWithSpan(parent, nonRecordingSpan{sc: rsc.WithRemote(true)})
}

// SpanFromContext returns the current Span from ctx.
//
// If no Span is currently set in ctx an implementation of a Span that
// performs no operations is returned.
func SpanFromContext(ctx context.Context) Span {
	if span, ok := ctx.Value(currentSpanKey).(Span); ok {
		return span
	}
	return noopSpan{}
}

// SpanContextFromContext returns the current Span's SpanContext.
func SpanContextFromContext(ctx context.Context) SpanContext {
	return SpanFromContext(ctx).SpanContext()
}
