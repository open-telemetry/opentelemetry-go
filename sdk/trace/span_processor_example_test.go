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
	"time"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type exporter struct{}

func (e exporter) ExportSpans(context.Context, []*export.SpanData) error { return nil }
func (e exporter) Shutdown(context.Context) error                        { return nil }

// LowPassFilter is a SpanProcessor that drops short lived spans.
type LowPassFilter struct {
	// Next is the next SpanProcessor in the chain.
	Next SpanProcessor
	// Cutoff is the duration under which spans are dropped.
	Cutoff time.Duration
}

func (f LowPassFilter) OnStart(sd *export.SpanData) { f.Next.OnStart(sd) }
func (f LowPassFilter) Shutdown()                   { f.Next.Shutdown() }
func (f LowPassFilter) ForceFlush()                 { f.Next.ForceFlush() }
func (f LowPassFilter) OnEnd(sd *export.SpanData) {
	if sd.EndTime.Sub(sd.StartTime) < f.Cutoff {
		// Drop short lived spans.
		return
	}
	f.Next.OnEnd(sd)
}

// HighPassFilter is a SpanProcessor that drops long lived spans.
type HighPassFilter struct {
	// Next is the next SpanProcessor in the chain.
	Next SpanProcessor
	// Cutoff is the duration over which spans are dropped.
	Cutoff time.Duration
}

func (f HighPassFilter) OnStart(sd *export.SpanData) { f.Next.OnStart(sd) }
func (f HighPassFilter) Shutdown()                   { f.Next.Shutdown() }
func (f HighPassFilter) ForceFlush()                 { f.Next.ForceFlush() }
func (f HighPassFilter) OnEnd(sd *export.SpanData) {
	if sd.EndTime.Sub(sd.StartTime) > f.Cutoff {
		// Drop long lived spans.
		return
	}
	f.Next.OnEnd(sd)
}

func ExampleSpanProcessor() {
	exportSpanProcessor := NewSimpleSpanProcessor(exporter{})

	// Build a band-pass filter to only allow spans shorter than an minute and
	// longer than a second to be exported with the exportSpanProcessor.
	bandPassFilter := LowPassFilter{
		Next: HighPassFilter{
			Next:   exportSpanProcessor,
			Cutoff: time.Minute,
		},
		Cutoff: time.Second,
	}

	traceProvider := NewProvider()
	traceProvider.RegisterSpanProcessor(bandPassFilter)
	// ...
}
