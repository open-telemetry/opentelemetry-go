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

package testtrace

import (
	"encoding/binary"
	"sync"

	"go.opentelemetry.io/otel/api/trace"
)

type Generator interface {
	TraceID() trace.ID
	SpanID() trace.SpanID
}

var _ Generator = (*CountGenerator)(nil)

// CountGenerator is a simple Generator that can be used to create unique, albeit deterministic,
// trace and span IDs.
type CountGenerator struct {
	lock        sync.Mutex
	traceIDHigh uint64
	traceIDLow  uint64
	spanID      uint64
}

func NewCountGenerator() *CountGenerator {
	return &CountGenerator{}
}

func (g *CountGenerator) TraceID() trace.ID {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.traceIDHigh == g.traceIDLow {
		g.traceIDHigh++
	} else {
		g.traceIDLow++
	}

	var traceID trace.ID

	binary.BigEndian.PutUint64(traceID[0:8], g.traceIDLow)
	binary.BigEndian.PutUint64(traceID[8:], g.traceIDHigh)

	return traceID
}

func (g *CountGenerator) SpanID() trace.SpanID {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.spanID++

	var spanID trace.SpanID

	binary.BigEndian.PutUint64(spanID[:], g.spanID)

	return spanID
}
