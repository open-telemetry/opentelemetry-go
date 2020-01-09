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

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

var _ apitrace.TracerSDK = &Tracer{}

func (tr *Tracer) Start(ctx context.Context, name core.Name, o ...apitrace.StartOption) (context.Context, apitrace.Span) {
	var opts apitrace.StartConfig
	var parent core.SpanContext
	var remoteParent bool

	//TODO [rghetia] : Add new option for parent. If parent is configured then use that parent.
	for _, op := range o {
		op(&opts)
	}

	if relation := opts.Relation; relation.SpanContext != core.EmptySpanContext() {
		switch relation.RelationshipType {
		case apitrace.ChildOfRelationship, apitrace.FollowsFromRelationship:
			parent = relation.SpanContext
			remoteParent = true
		default:
			// Future relationship types may have different behavior,
			// e.g., adding a `Link` instead of setting the `parent`
		}
	} else if p, ok := apitrace.SpanFromContext(ctx).(*span); ok {
		p.addChild()
		parent = p.spanContext
	}

	span := startSpanInternal(tr, name, parent, remoteParent, opts)
	for _, l := range opts.Links {
		span.addLink(l)
	}
	span.SetAttributes(opts.Attributes...)

	span.tracer = tr

	if span.IsRecording() {
		sps, _ := tr.spanProcessors.Load().(spanProcessorMap)
		for sp := range sps {
			sp.OnStart(span.data)
		}
	}

	ctx, end := startExecutionTracerTask(ctx, name.Base)
	span.executionTracerTaskEnd = end
	return apitrace.ContextWithSpan(ctx, span), span
}

func (tr *Tracer) WithSpan(ctx context.Context, name core.Name, body func(ctx context.Context) error) error {
	ctx, span := tr.Start(ctx, name)
	defer span.End()

	if err := body(ctx); err != nil {
		// TODO: set event with boolean attribute for error.
		return err
	}
	return nil
}
