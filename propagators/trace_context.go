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

package propagators

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"

	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

const (
	supportedVersion  = 0
	maxVersion        = 254
	traceparentHeader = "traceparent"
	tracestateHeader  = "tracestate"
)

type traceContextPropagatorKeyType uint

const (
	tracestateKey traceContextPropagatorKeyType = 0
)

// TraceContext propagates SpanContext in W3C TraceContext format.
//nolint:golint
type TraceContext struct{}

var _ propagation.HTTPPropagator = TraceContext{}
var traceCtxRegExp = regexp.MustCompile("^(?P<version>[0-9a-f]{2})-(?P<traceID>[a-f0-9]{32})-(?P<spanID>[a-f0-9]{16})-(?P<traceFlags>[a-f0-9]{2})(?:-.*)?$")

func (TraceContext) Inject(ctx context.Context, supplier propagation.HTTPSupplier) {
	tracestate := ctx.Value(tracestateKey)
	if state, ok := tracestate.(string); tracestate != nil && ok {
		supplier.Set(tracestateHeader, state)
	}

	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return
	}
	h := fmt.Sprintf("%.2x-%s-%s-%.2x",
		supportedVersion,
		sc.TraceID,
		sc.SpanID,
		sc.TraceFlags&trace.FlagsSampled)
	supplier.Set(traceparentHeader, h)
}

func (tc TraceContext) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	state := supplier.Get(tracestateHeader)
	if state != "" {
		ctx = context.WithValue(ctx, tracestateKey, state)
	}

	sc := tc.extract(supplier)
	if !sc.IsValid() {
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (TraceContext) extract(supplier propagation.HTTPSupplier) trace.SpanContext {
	h := supplier.Get(traceparentHeader)
	if h == "" {
		return trace.EmptySpanContext()
	}

	matches := traceCtxRegExp.FindStringSubmatch(h)

	if len(matches) == 0 {
		return trace.EmptySpanContext()
	}

	if len(matches) < 5 { // four subgroups plus the overall match
		return trace.EmptySpanContext()
	}

	if len(matches[1]) != 2 {
		return trace.EmptySpanContext()
	}
	ver, err := hex.DecodeString(matches[1])
	if err != nil {
		return trace.EmptySpanContext()
	}
	version := int(ver[0])
	if version > maxVersion {
		return trace.EmptySpanContext()
	}

	if version == 0 && len(matches) != 5 { // four subgroups plus the overall match
		return trace.EmptySpanContext()
	}

	if len(matches[2]) != 32 {
		return trace.EmptySpanContext()
	}

	var sc trace.SpanContext

	sc.TraceID, err = trace.IDFromHex(matches[2][:32])
	if err != nil {
		return trace.EmptySpanContext()
	}

	if len(matches[3]) != 16 {
		return trace.EmptySpanContext()
	}
	sc.SpanID, err = trace.SpanIDFromHex(matches[3])
	if err != nil {
		return trace.EmptySpanContext()
	}

	if len(matches[4]) != 2 {
		return trace.EmptySpanContext()
	}
	opts, err := hex.DecodeString(matches[4])
	if err != nil || len(opts) < 1 || (version == 0 && opts[0] > 2) {
		return trace.EmptySpanContext()
	}
	// Clear all flags other than the trace-context supported sampling bit.
	sc.TraceFlags = opts[0] & trace.FlagsSampled

	if !sc.IsValid() {
		return trace.EmptySpanContext()
	}

	return sc
}

func (TraceContext) GetAllKeys() []string {
	return []string{traceparentHeader, tracestateHeader}
}
