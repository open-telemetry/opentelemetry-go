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

package testtrace

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/distributedcontext"
	"go.opentelemetry.io/api/trace"
)

var _ trace.Span = (*Span)(nil)

type Span struct{}

func (s *Span) Tracer() trace.Tracer {
	return nil
}

func (s *Span) End(opts ...trace.EndOption) {}

func (s *Span) AddEvent(ctx context.Context, msg string, attrs ...core.KeyValue) {}

func (s *Span) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...core.KeyValue) {
}

func (s *Span) IsRecording() bool {
	return false
}

func (s *Span) AddLink(link trace.Link) {}

func (s *Span) Link(sc core.SpanContext, attrs ...core.KeyValue) {}

func (s *Span) SpanContext() core.SpanContext {
	return core.SpanContext{}
}

func (s *Span) SetStatus(code codes.Code) {}

func (s *Span) SetName(name string) {}

func (s *Span) SetAttribute(attr core.KeyValue) {}

func (s *Span) SetAttributes(attrs ...core.KeyValue) {}

func (s *Span) ModifyAttribute(mutator distributedcontext.Mutator) {}

func (s *Span) ModifyAttributes(mutators ...distributedcontext.Mutator) {}
