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
	"context"
	"math/rand"
	"sync"

	"google.golang.org/grpc/codes"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	apitrace "github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/sdk/event"
)

type tracer struct {
	resources core.EventID
}

var (
	ServiceKey   = tag.New("service")
	ComponentKey = tag.New("component")
	ErrorKey     = tag.New("error")
	SpanIDKey    = tag.New("span_id")
	TraceIDKey   = tag.New("trace_id")
	MessageKey   = tag.New("message",
		tag.WithDescription("message text: info, error, etc"),
	)
)

var t = &tracer{}

// Register registers tracer to global registry and returns the registered tracer.
func Register() apitrace.Tracer {
	apitrace.SetGlobalTracer(t)
	return t
}

func (t *tracer) ScopeID() core.ScopeID {
	return t.resources.Scope()
}

func (t *tracer) WithResources(attributes ...core.KeyValue) apitrace.Tracer {
	s := scope.New(t.resources.Scope(), attributes...)
	return &tracer{
		resources: s.ScopeID().EventID,
	}
}

func (t *tracer) WithComponent(name string) apitrace.Tracer {
	return t.WithResources(ComponentKey.String(name))
}

func (t *tracer) WithService(name string) apitrace.Tracer {
	return t.WithResources(ServiceKey.String(name))
}

func (t *tracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	// TODO: use runtime/trace.WithRegion for execution tracer support
	// TODO: use runtime/pprof.Do for profile tags support
	ctx, span := t.Start(ctx, name)
	defer span.Finish()

	if err := body(ctx); err != nil {
		span.SetAttribute(ErrorKey.Bool(true))
		span.AddEvent(ctx, event.WithAttr("span error", MessageKey.String(err.Error())))
		return err
	}
	return nil
}

// Start starts a new span with provided name and span options.
// If parent span reference is provided in the span option then it is used as as parent.
// Otherwise, parent span reference is retrieved from current context.
// The new span uses the same TraceID as parent.
// If no parent is found then a root span is created and started with random TraceID.
// TODO: Add sampling logic.
func (t *tracer) Start(ctx context.Context, name string, opts ...apitrace.SpanOption) (context.Context, apitrace.Span) {
	var child core.SpanContext

	child.SpanID = rand.Uint64()

	o := &apitrace.SpanOptions{}

	for _, opt := range opts {
		opt(o)
	}

	var parentScope core.ScopeID

	if o.Reference.HasTraceID() {
		parentScope = o.Reference.Scope()
	} else {
		parentSpan, _ := apitrace.Active(ctx).(*span)
		parentScope = parentSpan.ScopeID()
	}

	if parentScope.HasTraceID() {
		parent := parentScope.SpanContext
		child.TraceIDHigh = parent.TraceIDHigh
		child.TraceIDLow = parent.TraceIDLow
	} else {
		child.TraceIDHigh = rand.Uint64()
		child.TraceIDLow = rand.Uint64()
	}

	childScope := core.ScopeID{
		SpanContext: child,
		EventID:     t.resources,
	}

	span := &span{
		spanContext: child,
		tracer:      t,
		recordEvent: o.RecordEvent,
		eventID: observer.Record(observer.Event{
			Time:    o.StartTime,
			Type:    observer.START_SPAN,
			Scope:   scope.New(childScope, o.Attributes...).ScopeID(),
			Context: ctx,
			Parent:  parentScope,
			String:  name,
		}),
	}
	return scope.SetActive(ctx, span), span
}

func (t *tracer) Inject(ctx context.Context, span apitrace.Span, injector apitrace.Injector) {
	injector.Inject(span.ScopeID().SpanContext, tag.FromContext(ctx))
}
