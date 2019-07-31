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

package trace

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"sync/atomic"

	apitrace "go.opentelemetry.io/api/trace"
)

var config atomic.Value // access atomically

func init() {
	gen := &defaultIDGenerator{}
	// initialize traceID and spanID generators.
	var rngSeed int64
	for _, p := range []interface{}{
		&rngSeed, &gen.traceIDAdd, &gen.nextSpanID, &gen.spanIDInc,
	} {
		_ = binary.Read(crand.Reader, binary.LittleEndian, p)
	}
	gen.traceIDRand = rand.New(rand.NewSource(rngSeed))
	gen.spanIDInc |= 1

	config.Store(&Config{
		DefaultSampler:          ProbabilitySampler(defaultSamplingProbability),
		IDGenerator:             gen,
		MaxAttributesPerSpan:    DefaultMaxAttributesPerSpan,
		MaxMessageEventsPerSpan: DefaultMaxMessageEventsPerSpan,
		MaxLinksPerSpan:         DefaultMaxLinksPerSpan,
	})
}

var tr *tracer
var registerOnce sync.Once

func Register() apitrace.Tracer {
	registerOnce.Do(func() {
		tr = &tracer{}
		apitrace.SetGlobalTracer(tr)
	})
	return tr
}
