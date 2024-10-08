// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

// IDGenerator allows custom generators for TraceID and SpanID.
type IDGenerator interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// NewIDs returns a new trace and span ID.
	NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID)
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// NewSpanID returns a ID for a new span in the trace with traceID.
	NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}

// IDGeneratorRandom allows custom generators for TraceID and SpanID that comply
// with W3C Trace Context Level 2 randomness requirements.
type W3CTraceContextIDGenerator interface {
	// W3CTraceContextLevel2Random, when implemented by a
	// generator, indicates that this generator meets the
	// requirements
	W3CTraceContextLevel2Random()
}

type randomIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

var _ IDGenerator = &randomIDGenerator{}
var _ W3CTraceContextIDGenerator = &randomIDGenerator{}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := trace.SpanID{}
	for {
		_, _ = gen.randSource.Read(sid[:])
		if sid.IsValid() {
			break
		}
	}
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *randomIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	gen.Lock()
	defer gen.Unlock()
	tid := trace.TraceID{}
	sid := trace.SpanID{}
	for {
		_, _ = gen.randSource.Read(tid[:])
		if tid.IsValid() {
			break
		}
	}
	for {
		_, _ = gen.randSource.Read(sid[:])
		if sid.IsValid() {
			break
		}
	}
	return tid, sid
}

// W3CTraceContextLevel2Random declares meeting the W3C trace context
// level 2 randomness requirement.
func (gen *randomIDGenerator) W3CTraceContextLevel2Random() {}

func defaultIDGenerator() IDGenerator {
	gen := &randomIDGenerator{}
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}
