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

package trace_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/trace"
)

func TestStartWithRemoteSpanContext(t *testing.T) {
	trace.SetGlobalTracer(trace.PassThroughTracer{})
	sc := core.SpanContext{
		TraceID: core.TraceID{
			High: 0x0102030405060708,
			Low:  0x090a0b0c0d0e0f10,
		},
		SpanID:       0x0102030405060708,
		TraceOptions: 0x00,
	}
	ctx, span := trace.GlobalTracer().Start(context.Background(), "foo", trace.CopyOfRemote(sc))
	if _, ok := span.(*trace.PassThroughSpan); !ok {
		t.Errorf("Want PassThroughSpan but got %T\n", span)
	}
	currSpanCtx := trace.CurrentSpan(ctx).SpanContext()
	if diff := cmp.Diff(currSpanCtx, sc); diff != "" {
		t.Errorf("Want copy of span context but got %v\n", currSpanCtx)
	}
}

func TestStartWithoutRemoteSpanContext(t *testing.T) {
	trace.SetGlobalTracer(trace.PassThroughTracer{})
	ctx, span := trace.GlobalTracer().Start(context.Background(), "foo")
	if _, ok := span.(trace.NoopSpan); !ok {
		t.Errorf("Want NoopSpan but got %T\n", span)
	}
	currSpanCtx := trace.CurrentSpan(ctx).SpanContext()
	if diff := cmp.Diff(currSpanCtx, core.EmptySpanContext()); diff != "" {
		t.Errorf("Want Invalid span context but got %v\n", currSpanCtx)
	}
}
