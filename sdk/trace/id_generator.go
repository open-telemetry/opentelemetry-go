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

type randomIDGenerator struct {
	// pool is a pool of *rand.Rand objects. It provides a way to allow multiple
	// goroutines to concurrently generate random IDs without having to acquire
	// a mutex.
	pool *sync.Pool // of *rand.Rand
}

var _ IDGenerator = &randomIDGenerator{}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	r := gen.pool.Get().(*rand.Rand)
	defer gen.pool.Put(r)
	sid := trace.SpanID{}
	_, _ = r.Read(sid[:])
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *randomIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	r := gen.pool.Get().(*rand.Rand)
	defer gen.pool.Put(r)
	tid := trace.TraceID{}
	_, _ = r.Read(tid[:])
	sid := trace.SpanID{}
	_, _ = r.Read(sid[:])
	return tid, sid
}

func defaultIDGenerator() IDGenerator {
	return &randomIDGenerator{
		pool: &sync.Pool{
			New: func() interface{} {
				var rngSeed int64
				_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
				return rand.New(rand.NewSource(rngSeed))
			},
		},
	}
}
