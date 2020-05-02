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

package trace

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/api/propagation"
)

const (
	supportedVersion  = 0
	maxVersion        = 254
	traceparentHeader = "traceparent"
)

// TraceContext propagates SpanContext in W3C TraceContext format.
//nolint:golint
type TraceContext struct{}

var _ propagation.HTTPPropagator = TraceContext{}
var traceCtxRegExp = regexp.MustCompile("^[0-9a-f]{2}-[a-f0-9]{32}-[a-f0-9]{16}-[a-f0-9]{2}-?")

// DefaultHTTPPropagator returns the default trace HTTP propagator.
func DefaultHTTPPropagator() propagation.HTTPPropagator {
	return TraceContext{}
}

func (TraceContext) Inject(ctx context.Context, supplier propagation.HTTPSupplier) {
	sc := SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return
	}
	h := fmt.Sprintf("%.2x-%s-%s-%.2x",
		supportedVersion,
		sc.TraceID,
		sc.SpanID,
		sc.TraceFlags&TraceFlagsSampled)
	supplier.Set(traceparentHeader, h)
}

func (tc TraceContext) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	sc := tc.extract(supplier)
	if !sc.IsValid() {
		return ctx
	}
	return ContextWithRemoteSpanContext(ctx, sc)
}

func (TraceContext) extract(supplier propagation.HTTPSupplier) SpanContext {
	h := supplier.Get(traceparentHeader)
	if h == "" {
		return EmptySpanContext()
	}

	h = strings.Trim(h, "-")
	if !traceCtxRegExp.MatchString(h) {
		return EmptySpanContext()
	}

	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return EmptySpanContext()
	}

	if len(sections[0]) != 2 {
		return EmptySpanContext()
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return EmptySpanContext()
	}
	version := int(ver[0])
	if version > maxVersion {
		return EmptySpanContext()
	}

	if version == 0 && len(sections) != 4 {
		return EmptySpanContext()
	}

	if len(sections[1]) != 32 {
		return EmptySpanContext()
	}

	var sc SpanContext

	sc.TraceID, err = TraceIDFromHex(sections[1][:32])
	if err != nil {
		return EmptySpanContext()
	}

	if len(sections[2]) != 16 {
		return EmptySpanContext()
	}
	sc.SpanID, err = SpanIDFromHex(sections[2])
	if err != nil {
		return EmptySpanContext()
	}

	if len(sections[3]) != 2 {
		return EmptySpanContext()
	}
	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 || (version == 0 && opts[0] > 2) {
		return EmptySpanContext()
	}
	sc.TraceFlags = opts[0] &^ TraceFlagsUnused

	if !sc.IsValid() {
		return EmptySpanContext()
	}

	return sc
}

func (TraceContext) GetAllKeys() []string {
	return []string{traceparentHeader}
}
