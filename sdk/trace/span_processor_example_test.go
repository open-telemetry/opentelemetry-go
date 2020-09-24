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
	"time"

	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/export/trace/tracetest"
)

// DurationFilter is a SpanProcessor that filters spans that have lifetimes
// outside of a defined range.
type DurationFilter struct {
	// Next is the next SpanProcessor in the chain.
	Next SpanProcessor

	// Min is the duration under which spans are dropped.
	Min time.Duration
	// Max is the duration over which spans are dropped.
	Max time.Duration
}

func (f DurationFilter) OnStart(sd *export.SpanData) { f.Next.OnStart(sd) }
func (f DurationFilter) Shutdown()                   { f.Next.Shutdown() }
func (f DurationFilter) ForceFlush()                 { f.Next.ForceFlush() }
func (f DurationFilter) OnEnd(sd *export.SpanData) {
	if f.Min > 0 && sd.EndTime.Sub(sd.StartTime) < f.Min {
		// Drop short lived spans.
		return
	}
	if f.Max > 0 && sd.EndTime.Sub(sd.StartTime) > f.Max {
		// Drop long lived spans.
		return
	}
	f.Next.OnEnd(sd)
}

// InstrumentationBlacklist is a SpanProcessor that drops all spans from
// certain instrumentation.
type InstrumentationBlacklist struct {
	// Next is the next SpanProcessor in the chain.
	Next SpanProcessor

	// Blacklist is the set of instrumentation names for which spans will be
	// dropped.
	Blacklist map[string]bool
}

func (f InstrumentationBlacklist) OnStart(sd *export.SpanData) { f.Next.OnStart(sd) }
func (f InstrumentationBlacklist) Shutdown()                   { f.Next.Shutdown() }
func (f InstrumentationBlacklist) ForceFlush()                 { f.Next.ForceFlush() }
func (f InstrumentationBlacklist) OnEnd(sd *export.SpanData) {
	if f.Blacklist != nil && f.Blacklist[sd.InstrumentationLibrary.Name] {
		// Drop spans from this instrumentation
		return
	}
	f.Next.OnEnd(sd)
}

func ExampleSpanProcessor() {
	exportSP := NewSimpleSpanProcessor(tracetest.NewNoopExporter())

	// Build a SpanProcessor chain to filter out all spans from the pernicious
	// "naughty-instrumentation" dependency and only allow spans shorter than
	// an minute and longer than a second to be exported with the exportSP.
	filter := DurationFilter{
		Next: InstrumentationBlacklist{
			Next: exportSP,
			Blacklist: map[string]bool{
				"naughty-instrumentation": true,
			},
		},
		Min: time.Second,
		Max: time.Minute,
	}

	_ = NewTracerProvider(WithSpanProcessor(filter))
	// ...
}
