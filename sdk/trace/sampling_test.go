// Copyright 2020, OpenTelemetry Authors
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
	"testing"

	"go.opentelemetry.io/otel/api/core"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestAlwaysParentSampleWithParentSampled(t *testing.T) {
	sampler := sdktrace.AlwaysParentSample()
	traceID, _ := core.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := core.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := core.SpanContext{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: core.TraceFlagsSampled,
	}
	if sampler.ShouldSample(sdktrace.SamplingParameters{ParentContext: parentCtx}).Decision != sdktrace.RecordAndSampled {
		t.Error("Sampling decision should be RecordAndSampled")
	}
}

func TestAlwaysParentSampleWithParentNotSampled(t *testing.T) {
	sampler := sdktrace.AlwaysParentSample()
	traceID, _ := core.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := core.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := core.SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	}
	if sampler.ShouldSample(sdktrace.SamplingParameters{ParentContext: parentCtx}).Decision != sdktrace.NotRecord {
		t.Error("Sampling decision should be NotRecord")
	}
}
