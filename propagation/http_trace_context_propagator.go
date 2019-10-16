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

package propagation

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.opentelemetry.io/api/trace"

	"go.opentelemetry.io/api/core"
	apipropagation "go.opentelemetry.io/api/propagation"
)

const (
	supportedVersion  = 0
	maxVersion        = 254
	TraceparentHeader = "Traceparent"
)

// HTTPTraceContextPropagator propagates SpanContext in W3C TraceContext format.
type HTTPTraceContextPropagator struct{}

var _ apipropagation.TextFormatPropagator = HTTPTraceContextPropagator{}
var traceCtxRegExp = regexp.MustCompile("^[0-9a-f]{2}-[a-f0-9]{32}-[a-f0-9]{16}-[a-f0-9]{2}-?")

func (hp HTTPTraceContextPropagator) Inject(ctx context.Context, supplier apipropagation.Supplier) {
	sc := trace.CurrentSpan(ctx).SpanContext()
	if sc.IsValid() {
		h := fmt.Sprintf("%.2x-%.16x%.16x-%.16x-%.2x",
			supportedVersion,
			sc.TraceID.High,
			sc.TraceID.Low,
			sc.SpanID,
			sc.TraceFlags&core.TraceFlagsSampled)
		supplier.Set(TraceparentHeader, h)
	}
}

func (hp HTTPTraceContextPropagator) Extract(ctx context.Context, supplier apipropagation.Supplier) core.SpanContext {
	h := supplier.Get(TraceparentHeader)
	if h == "" {
		return core.EmptySpanContext()
	}

	h = strings.Trim(h, "-")
	if !traceCtxRegExp.MatchString(h) {
		return core.EmptySpanContext()
	}

	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return core.EmptySpanContext()
	}

	if len(sections[0]) != 2 {
		return core.EmptySpanContext()
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return core.EmptySpanContext()
	}
	version := int(ver[0])
	if version > maxVersion {
		return core.EmptySpanContext()
	}

	if version == 0 && len(sections) != 4 {
		return core.EmptySpanContext()
	}

	if len(sections[1]) != 32 {
		return core.EmptySpanContext()
	}

	result, err := strconv.ParseUint(sections[1][0:16], 16, 64)
	if err != nil {
		return core.EmptySpanContext()
	}
	var sc core.SpanContext

	sc.TraceID.High = result

	result, err = strconv.ParseUint(sections[1][16:32], 16, 64)
	if err != nil {
		return core.EmptySpanContext()
	}
	sc.TraceID.Low = result

	if len(sections[2]) != 16 {
		return core.EmptySpanContext()
	}
	result, err = strconv.ParseUint(sections[2][0:], 16, 64)
	if err != nil {
		return core.EmptySpanContext()
	}
	sc.SpanID = result

	if len(sections[3]) != 2 {
		return core.EmptySpanContext()
	}
	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 || (version == 0 && opts[0] > 2) {
		return core.EmptySpanContext()
	}
	sc.TraceFlags = opts[0] &^ core.TraceFlagsUnused

	if !sc.IsValid() {
		return core.EmptySpanContext()
	}

	return sc
}

func (hp HTTPTraceContextPropagator) GetAllKeys() []string {
	return []string{TraceparentHeader}
}
