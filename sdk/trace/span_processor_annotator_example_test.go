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
)

/*
Sometimes information about a runtime environment can change dynamically or be
delayed from startup. Instead of continuously recreating and distributing a
TracerProvider with an immutable Resource or delaying the startup of your
application on a slow-loading piece of information, annotate the created spans
dynamically using a SpanProcessor.
*/

var (
	// owner represents the owner of the application. In this example it is
	// stored as a simple string, but in real-world use this may be the
	// response to an asynchronous request.
	owner    = "unknown"
	ownerKey = attribute.Key("owner")
)

// Annotator is a SpanProcessor that adds attributes to all started spans.
type Annotator struct {
	// AttrsFunc is called when a span is started. The attributes it returns
	// are set on the Span being started.
	AttrsFunc func() []attribute.KeyValue
}

func (a Annotator) OnStart(_ context.Context, s ReadWriteSpan) { s.SetAttributes(a.AttrsFunc()...) }
func (a Annotator) Shutdown(context.Context) error             { return nil }
func (a Annotator) ForceFlush(context.Context) error           { return nil }
func (a Annotator) OnEnd(s ReadOnlySpan) {
	attr := s.Attributes()[0]
	fmt.Printf("%s: %s\n", attr.Key, attr.Value.AsString())
}

func ExampleSpanProcessor_annotated() {
	a := Annotator{
		AttrsFunc: func() []attribute.KeyValue {
			return []attribute.KeyValue{ownerKey.String(owner)}
		},
	}
	tracer := NewTracerProvider(WithSpanProcessor(a)).Tracer("annotated")

	// Simulate the situation where we want to annotate spans with an owner,
	// but at startup we do not now this information. Instead of waiting for
	// the owner to be known before starting and blocking here, start doing
	// work and update when the information becomes available.
	ctx := context.Background()
	_, s0 := tracer.Start(ctx, "span0")

	// Simulate an asynchronous call to determine the owner succeeding. We now
	// know that the owner of this application has been determined to be
	// Alice. Make sure all subsequent spans are annotated appropriately.
	owner = "alice"

	_, s1 := tracer.Start(ctx, "span1")
	s0.End()
	s1.End()

	// Output:
	// owner: unknown
	// owner: alice
}
