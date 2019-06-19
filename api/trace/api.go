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
	"time"

	"google.golang.org/grpc/codes"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/log"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/api/stats"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
)

type (
	Tracer interface {
		Start(context.Context, string, ...SpanOption) (context.Context, Span)

		WithSpan(
			ctx context.Context,
			operation string,
			body func(ctx context.Context) error,
		) error

		WithService(name string) Tracer
		WithComponent(name string) Tracer
		WithResources(res ...core.KeyValue) Tracer

		// Note: see https://github.com/opentracing/opentracing-go/issues/127
		Inject(context.Context, Span, Injector)

		// ScopeID returns the resource scope of this tracer.
		scope.Scope
	}

	Span interface {
		scope.Mutable

		log.Interface

		stats.Interface

		SetError(bool)

		Tracer() Tracer

		Finish()

		// IsRecordingEvents returns true is the span is active and recording events is enabled.
		IsRecordingEvents() bool

		// SpancContext returns span context of the span. Return SpanContext is usable
		// even after the span is finished.
		SpanContext() core.SpanContext

		SetStatus(codes.Code)
	}

	Injector interface {
		Inject(core.SpanContext, tag.Map)
	}

	// SpanOption apply changes to SpanOptions.
	SpanOption func(*SpanOptions)

	SpanOptions struct {
		attributes  []core.KeyValue
		startTime   time.Time
		reference   Reference
		recordEvent bool
	}

	Reference struct {
		core.SpanContext
		RelationshipType
	}

	RelationshipType int
)

const (
	ChildOfRelationship RelationshipType = iota
	FollowsFromRelationship
)

func GlobalTracer() Tracer {
	if t := global.Load(); t != nil {
		return t.(Tracer)
	}
	return empty
}

func SetGlobalTracer(t Tracer) {
	global.Store(t)
}

func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	return GlobalTracer().Start(ctx, name, opts...)
}

func Active(ctx context.Context) Span {
	span, _ := scope.Active(ctx).(*span)
	return span
}

func WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return GlobalTracer().WithSpan(ctx, name, body)
}

func SetError(ctx context.Context, v bool) {
	Active(ctx).SetError(v)
}

func Inject(ctx context.Context, injector Injector) {
	span := Active(ctx)
	if span == nil {
		return
	}

	span.Tracer().Inject(ctx, span, injector)
}

func WithStartTime(t time.Time) SpanOption {
	return func(o *SpanOptions) {
		o.startTime = t
	}
}

func WithAttributes(attrs ...core.KeyValue) SpanOption {
	return func(o *SpanOptions) {
		o.attributes = attrs
	}
}

func WithRecordEvents() SpanOption {
	return func(o *SpanOptions) {
		o.recordEvent = true
	}
}

func ChildOf(sc core.SpanContext) SpanOption {
	return func(o *SpanOptions) {
		o.reference = Reference{
			SpanContext:      sc,
			RelationshipType: ChildOfRelationship,
		}
	}
}

func FollowsFrom(sc core.SpanContext) SpanOption {
	return func(o *SpanOptions) {
		o.reference = Reference{
			SpanContext:      sc,
			RelationshipType: FollowsFromRelationship,
		}
	}
}
