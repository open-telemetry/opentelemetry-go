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
	"go.opentelemetry.io/sdk/trace/internal"
	"math/rand"
	"sync"
	"sync/atomic"

	apitrace "go.opentelemetry.io/api/trace"
)

var config atomic.Value // access atomically

func init() {
	config.Store(&Config{
		DefaultSampler:       ProbabilitySampler(defaultSamplingProbability),
		IDGenerator:          defIDGenerator(),
		MaxAttributesPerSpan: DefaultMaxAttributesPerSpan,
		MaxEventsPerSpan:     DefaultMaxEventsPerSpan,
		MaxLinksPerSpan:      DefaultMaxLinksPerSpan,
	})
}

var p *TraceProvider
var registerProviderOnce sync.Once

// RegisterProvider registers trace provider implementation as default Trace Provider.
// It creates single instance of trace provider and registers it once.
// Recommended use is to call RegisterProvider in main() of an
// application before calling any tracing api.
func RegisterProvider() apitrace.Provider {
	registerProviderOnce.Do(func() {
		p = &TraceProvider{namedTracer: map[string]*tracer{}}
		apitrace.SetGlobalProvider(p)
	})
	return p
}

var tr *tracer
var registerOnce sync.Once

// Register registers tracer implementation as default Tracer.
// It creates single instance of tracer and registers it once.
// Recommended use is to call Register in main() of an
// application before calling any tracing api.
func Register() apitrace.Tracer {
	registerOnce.Do(func() {
		tr = &tracer{}
		apitrace.SetGlobalTracer(tr)
	})
	return tr
}

func defIDGenerator() internal.IDGenerator {
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
	return gen
}