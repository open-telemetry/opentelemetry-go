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

package internal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/internal/sdk"
)

type SpanRecorder struct {
	Started []string
	Ended   []string
}

func (sr *SpanRecorder) OnStart(s *sdk.Span) {
	sr.Started = append(sr.Started, s.Name)
}

func (sr *SpanRecorder) OnEnd(s *sdk.Span) {
	sr.Ended = append(sr.Ended, s.Name)
}

func TestTraceWithSDK(t *testing.T) {
	internal.ResetForTest()

	ctx := context.Background()
	gtp := global.TraceProvider()
	tracer1 := gtp.Tracer("pre")
	_, span1 := tracer1.Start(ctx, "span1")

	// This should be dropped.
	if err := tracer1.WithSpan(ctx, "withSpan1", func(context.Context) error { return nil }); err != nil {
		t.Errorf("failed to wrap function with span prior to initialization: %v", err)
	}

	sr := SpanRecorder{}
	tp := sdk.NewTraceProvider(sdk.WithSpanRecorder(&sr))
	global.SetTraceProvider(tp)

	// This span was started before initialization, it is expected to be dropped.
	span1.End()

	// The existing Tracer should have been configured to now use the configured SDK.
	_, span2 := tracer1.Start(ctx, "span2")
	span2.End()
	if err := tracer1.WithSpan(ctx, "withSpan2", func(context.Context) error { return nil }); err != nil {
		t.Errorf("failed to wrap function with span post initialization: %v", err)
	}

	// The global trace Provider should now create Tracers that also use the newly configured SDK.
	tracer2 := gtp.Tracer("post")
	_, span3 := tracer2.Start(ctx, "span3")
	span3.End()
	if err := tracer2.WithSpan(ctx, "withSpan3", func(context.Context) error { return nil }); err != nil {
		t.Errorf("failed to wrap function with span post initialization with new tracer: %v", err)
	}

	expected := []string{"span2", "withSpan2", "span3", "withSpan3"}
	require.Equal(t, sr.Started, expected)
	require.Equal(t, sr.Ended, expected)
}
