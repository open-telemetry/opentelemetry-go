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
	"errors"
	"strings"

	"go.opentelemetry.io/otel/api/propagation"
)

const (
	// Default B3 Header names.
	B3SingleHeader       = "b3"
	B3DebugFlagHeader    = "x-b3-flags"
	B3TraceIDHeader      = "x-b3-traceid"
	B3SpanIDHeader       = "x-b3-spanid"
	B3SampledHeader      = "x-b3-sampled"
	B3ParentSpanIDHeader = "x-b3-parentspanid"

	b3TraceIDPadding = "0000000000000000"
)

var (
	empty = EmptySpanContext()

	errInvalidSampledByte        = errors.New("invalid B3 Sampled found")
	errInvalidSampledHeader      = errors.New("invalid B3 Sampled header found")
	errInvalidTraceIDHeader      = errors.New("invalid B3 TraceID header found")
	errInvalidSpanIDHeader       = errors.New("invalid B3 SpanID header found")
	errInvalidParentSpanIDHeader = errors.New("invalid B3 ParentSpanID header found")
	errInvalidScope              = errors.New("require either both TraceID and SpanID or none")
	errInvalidScopeParent        = errors.New("ParentSpanID requires both TraceID and SpanID to be available")
	errInvalidScopeParentSingle  = errors.New("ParentSpanID requires TraceID, SpanID and Sampled to be available")
	errEmptyContext              = errors.New("empty request context")
	errInvalidTraceIDValue       = errors.New("invalid B3 TraceID value found")
	errInvalidSpanIDValue        = errors.New("invalid B3 SpanID value found")
	errInvalidParentSpanIDValue  = errors.New("invalid B3 ParentSpanID value found")
)

// B3Encoding is a bitmask representation of the B3 encoding type.
type B3Encoding uint8

const (
	// MultipleHeader is a B3 encoding that uses multiple headers to
	// transmit tracing information all prefixed with `x-b3-`.
	MultipleHeader B3Encoding = 1
	// SingleHeader is a B3 encoding that uses a single header named `b3 to
	// transmit tracing information.
	SingleHeader B3Encoding = 2
)

// B3 propagator serializes SpanContext to/from B3 Headers.
// This propagator supports both versions of B3 headers,
//  1. Single Header:
//    b3: {TraceId}-{SpanId}-{SamplingState}-{ParentSpanId}
//  2. Multiple Headers:
//    x-b3-traceid: {TraceId}
//    x-b3-parentspanid: {ParentSpanId}
//    x-b3-spanid: {SpanId}
//    x-b3-sampled: {SamplingState}
//    x-b3-flags: {DebugFlag}
type B3 struct {
	// InjectEncoding are the B3 encodings used when injecting trace
	// information. If no encoding is specific it defaults to
	// `MultipleHeader`.
	InjectEncoding B3Encoding
}

func (b3 B3) supports(e B3Encoding) bool {
	return b3.InjectEncoding&e != 0
}

var _ propagation.HTTPPropagator = B3{}

// Inject injects a context into the supplier as B3 headers.
func (b3 B3) Inject(ctx context.Context, supplier propagation.HTTPSupplier) {
	sc := SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return
	}

	if b3.supports(SingleHeader) {
		header := []string{}
		if sc.TraceID.IsValid() && sc.SpanID.IsValid() {
			header = append(header, sc.TraceID.String(), sc.SpanID.String())
		}

		if sc.TraceFlags&FlagsUnset != FlagsUnset {
			if sc.IsSampled() {
				header = append(header, "1")
			} else {
				header = append(header, "0")
			}
		}

		supplier.Set(B3SingleHeader, strings.Join(header, "-"))
	}

	if b3.supports(MultipleHeader) || b3.InjectEncoding == 0 {
		if sc.TraceID.IsValid() && sc.SpanID.IsValid() {
			supplier.Set(B3TraceIDHeader, sc.TraceID.String())
			supplier.Set(B3SpanIDHeader, sc.SpanID.String())
		}

		if sc.TraceFlags&FlagsUnset != FlagsUnset {
			if sc.IsSampled() {
				supplier.Set(B3SampledHeader, "1")
			} else {
				supplier.Set(B3SampledHeader, "0")
			}
		}
	}
}

// Extract extracts a context from the supplier if it contains B3 headers.
func (b3 B3) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	var (
		sc  SpanContext
		err error
	)

	if h := supplier.Get(B3SingleHeader); h != "" {
		sc, err = extractSingle(h)
		if err == nil && sc.IsValid() {
			return ContextWithRemoteSpanContext(ctx, sc)
		}
	}

	var (
		traceID      = supplier.Get(B3TraceIDHeader)
		spanID       = supplier.Get(B3SpanIDHeader)
		parentSpanID = supplier.Get(B3ParentSpanIDHeader)
		sampled      = supplier.Get(B3SampledHeader)
		debugFlag    = supplier.Get(B3DebugFlagHeader)
	)
	sc, err = extractMultiple(traceID, spanID, parentSpanID, sampled, debugFlag)
	if err != nil || !sc.IsValid() {
		return ctx
	}
	return ContextWithRemoteSpanContext(ctx, sc)
}

func (b3 B3) GetAllKeys() []string {
	if b3.supports(SingleHeader) {
		return []string{B3SingleHeader}
	}
	return []string{B3TraceIDHeader, B3SpanIDHeader, B3SampledHeader}
}

// extractMultiple reconstructs a SpanContext from header values based on B3
// Multiple header. It is based on the implementation found here:
// https://github.com/openzipkin/zipkin-go/blob/v0.2.2/propagation/b3/spancontext.go
// and adapted to support a SpanContext.
func extractMultiple(traceID, spanID, parentSpanID, sampled, flags string) (SpanContext, error) {
	var (
		err           error
		requiredCount int
		sc            = SpanContext{}
	)

	// correct values for an existing sampled header are "0" and "1".
	// For legacy support and  being lenient to other tracing implementations we
	// allow "true" and "false" as inputs for interop purposes.
	switch strings.ToLower(sampled) {
	case "0", "false":
		sc.TraceFlags = FlagsNotSampled
	case "1", "true":
		sc.TraceFlags = FlagsSampled
	case "":
		sc.TraceFlags = FlagsUnset
	default:
		return empty, errInvalidSampledHeader
	}

	// The only accepted value for Flags is "1". This will set Debug to true. All
	// other values and omission of header will be ignored.
	if flags == "1" {
		// We do not track debug status, but the sampling needs to be unset.
		sc.TraceFlags = FlagsUnset
	}

	if traceID != "" {
		requiredCount++
		id := traceID
		if len(traceID) == 16 {
			// Pad 64-bit trace IDs.
			id = b3TraceIDPadding + traceID
		}
		if sc.TraceID, err = IDFromHex(id); err != nil {
			return empty, errInvalidTraceIDHeader
		}
	}

	if spanID != "" {
		requiredCount++
		if sc.SpanID, err = SpanIDFromHex(spanID); err != nil {
			return empty, errInvalidSpanIDHeader
		}
	}

	if requiredCount != 0 && requiredCount != 2 {
		return empty, errInvalidScope
	}

	if parentSpanID != "" {
		if requiredCount == 0 {
			return empty, errInvalidScopeParent
		}
		// Validate parent span ID but we do not use it so do not save it.
		if _, err = SpanIDFromHex(parentSpanID); err != nil {
			return empty, errInvalidParentSpanIDHeader
		}
	}

	return sc, nil
}

// extractSingle reconstructs a SpanContext from contextHeader based on a B3
// Single header. It is based on the implementation found here:
// https://github.com/openzipkin/zipkin-go/blob/v0.2.2/propagation/b3/spancontext.go
// and adapted to support a SpanContext.
func extractSingle(contextHeader string) (SpanContext, error) {
	if contextHeader == "" {
		return empty, errEmptyContext
	}

	var (
		sc       = SpanContext{}
		sampling string
	)

	headerLen := len(contextHeader)

	if headerLen == 1 {
		sampling = contextHeader
	} else if headerLen == 16 || headerLen == 32 {
		return empty, errInvalidScope
	} else if headerLen >= 16+16+1 {
		pos := 0
		var traceID string
		if string(contextHeader[16]) == "-" {
			// traceID must be 64 bits
			pos += 16 + 1 // {traceID}-
			traceID = b3TraceIDPadding + string(contextHeader[0:16])
		} else if string(contextHeader[32]) == "-" {
			// traceID must be 128 bits
			pos += 32 + 1 // {traceID}-
			traceID = string(contextHeader[0:32])
		} else {
			return empty, errInvalidTraceIDValue
		}
		var err error
		sc.TraceID, err = IDFromHex(traceID)
		if err != nil {
			return empty, errInvalidTraceIDValue
		}

		sc.SpanID, err = SpanIDFromHex(contextHeader[pos : pos+16])
		if err != nil {
			return empty, errInvalidSpanIDValue
		}
		pos += 16 // {traceID}-{spanID}

		if headerLen > pos {
			if headerLen == pos+1 {
				return empty, errInvalidSampledByte
			}
			pos++ // {traceID}-{spanID}-

			if headerLen == pos+1 {
				sampling = string(contextHeader[pos])
			} else if headerLen == pos+16 {
				return empty, errInvalidScopeParentSingle
			} else if headerLen == pos+1+16+1 {
				sampling = string(contextHeader[pos])
				pos += 1 + 1 // {traceID}-{spanID}-{sampling}-

				// Validate parent span ID but we do not use it so do not
				// save it.
				_, err = SpanIDFromHex(contextHeader[pos:])
				if err != nil {
					return empty, errInvalidParentSpanIDValue
				}
			} else {
				return empty, errInvalidParentSpanIDValue
			}
		}
	} else {
		return empty, errInvalidTraceIDValue
	}
	switch sampling {
	case "", "d":
		sc.TraceFlags = FlagsUnset
	case "1":
		sc.TraceFlags = FlagsSampled
	case "0":
		sc.TraceFlags = FlagsNotSampled
	default:
		return empty, errInvalidSampledByte
	}

	return sc, nil
}
