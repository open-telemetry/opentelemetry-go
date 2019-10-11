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

	"go.opentelemetry.io/api/core"
	apitrace "go.opentelemetry.io/api/trace"
)

type tracer struct {
	name      string
	component string
	resources []core.KeyValue
}

var _ apitrace.Tracer = &tracer{}

func (tr *tracer) Start(ctx context.Context, name string, o ...apitrace.SpanOption) (context.Context, apitrace.Span) {
	var opts apitrace.SpanOptions
	var parent core.SpanContext
	var remoteParent bool

	//TODO [rghetia] : Add new option for parent. If parent is configured then use that parent.
	for _, op := range o {
		op(&opts)
	}

	// TODO: [rghetia] ChildOfRelationship is used to indicate that the parent is remote
	// and its context is received as part of a request. There are two possibilities
	// 1. Remote is trusted. So continue using same trace.
	//      tracer.Start(ctx, "some name", ChildOf(remote_span_context))
	// 2. Remote is not trusted. In this case create a root span and then add the remote as link
	//      span := tracer.Start(ctx, "some name")
	//      span.Link(remote_span_context, ChildOfRelationship)
	if opts.Reference.SpanContext != core.EmptySpanContext() &&
		opts.Reference.RelationshipType == apitrace.ChildOfRelationship {
		parent = opts.Reference.SpanContext
		remoteParent = true
	} else {
		if p := apitrace.CurrentSpan(ctx); p != nil {
			if sdkSpan, ok := p.(*span); ok {
				sdkSpan.addChild()
				parent = sdkSpan.spanContext
			}
		}
	}

	span := startSpanInternal(name, parent, remoteParent, opts)
	span.tracer = tr

	if span.IsRecording() {
		sps, _ := spanProcessors.Load().(spanProcessorMap)
		for sp := range sps {
			sp.OnStart(span.data)
		}
	}

	ctx, end := startExecutionTracerTask(ctx, name)
	span.executionTracerTaskEnd = end
	return apitrace.SetCurrentSpan(ctx, span), span
}

func (tr *tracer) WithSpan(ctx context.Context, name string, body func(ctx context.Context) error) error {
	ctx, span := tr.Start(ctx, name)
	defer span.End()

	if err := body(ctx); err != nil {
		// TODO: set event with boolean attribute for error.
		return err
	}
	return nil
}

func (tr *tracer) WithService(name string) apitrace.Tracer {
	tr.name = name
	return tr
}

// WithResources does nothing and returns noop implementation of apitrace.Tracer.
func (tr *tracer) WithResources(res ...core.KeyValue) apitrace.Tracer {
	tr.resources = res
	return tr
}

// WithComponent does nothing and returns noop implementation of apitrace.Tracer.
func (tr *tracer) WithComponent(component string) apitrace.Tracer {
	tr.component = component
	return tr
}
