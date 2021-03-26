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
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/export/trace/tracetest"
)

type AttrsFunc func(context.Context) []attribute.KeyValue

// annotator is a SpanProcessor that adds attributes to all started spans.
type annotator struct {
	// Next is the next SpanProcessor in the chain.
	Next SpanProcessor

	// AttrsFunc is called when a span is started and the returned attributes
	// are added to that span.
	AttrsFunc AttrsFunc
}

// NewAnnotator returns a SpanProcessor that adds attrs to all started spans.
func NewAnnotator(next SpanProcessor, attrsFunc AttrsFunc) SpanProcessor {
	return annotator{
		Next:      next,
		AttrsFunc: attrsFunc,
	}
}

func (a annotator) OnStart(parent context.Context, s ReadWriteSpan) {
	s.SetAttributes(a.AttrsFunc(parent)...)
	a.Next.OnStart(parent, s)
}
func (a annotator) Shutdown(ctx context.Context) error   { return a.Next.Shutdown(ctx) }
func (a annotator) ForceFlush(ctx context.Context) error { return a.Next.ForceFlush(ctx) }
func (a annotator) OnEnd(s ReadOnlySpan)                 { a.Next.OnEnd(s) }

func ExampleSpanProcessor_annotated() {
	// Sometimes information about a runtime environment can change
	// dynamically or be delayed from startup. Instead of continuously
	// recreating and distributing a TracerProvider with an immutable Resource
	// or delaying the startup of your application on a slow loading piece of
	// information annotate the created spans dynamically using a
	// SpanProcessor.

	// Use this chan to signal when an owner of the process is known.
	ownerCh := make(chan string, 1)
	ownerKey := attribute.Key("owner")

	exporter := tracetest.NewInMemoryExporter()
	a := NewAnnotator(
		// Chain the export pipeline downstream of this SpanProcessor.
		NewSimpleSpanProcessor(exporter),
		// Dynamically lookup the owner and annotate accordingly.
		func(ctx context.Context) []attribute.KeyValue {
			select {
			case name := <-ownerCh:
				return []attribute.KeyValue{ownerKey.String(name)}
			default:
				return []attribute.KeyValue{ownerKey.String("unknown")}
			}
		},
	)

	// Instead of waiting for the owner to be known before starting and
	// blocking here, start the tracing process and update when the
	// information becomes available.
	ctx := context.Background()
	tracer := NewTracerProvider(WithSpanProcessor(a)).Tracer("annotated")
	_, s0 := tracer.Start(ctx, "span0")

	// It was determined that Alice is the owner of this task, make sure all
	// subsequent spans are annotated appropriately.
	ownerCh <- "alice"

	_, s1 := tracer.Start(ctx, "span1")
	s0.End()
	s1.End()

	spans := exporter.GetSpans()
	a0, a1 := spans[0].Attributes[0], spans[1].Attributes[0]
	fmt.Printf("%s: %s\n", a0.Key, a0.Value.AsString())
	fmt.Printf("%s: %s\n", a1.Key, a1.Value.AsString())

	// Output:
	// owner: unknown
	// owner: alice
}
