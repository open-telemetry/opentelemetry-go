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

package oc2otel

import (
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func TestStartOptionsSpanKind(t *testing.T) {
	conv := map[int]trace.SpanKind{
		octrace.SpanKindClient:      trace.SpanKindClient,
		octrace.SpanKindServer:      trace.SpanKindServer,
		octrace.SpanKindUnspecified: trace.SpanKindUnspecified,
	}

	for oc, otel := range conv {
		ocOpts := []octrace.StartOption{octrace.WithSpanKind(oc)}
		otelOpts, err := StartOptions(ocOpts)
		if err != nil {
			t.Errorf("StartOptions errored: %v", err)
			continue
		}
		c := trace.NewSpanStartConfig(otelOpts...)
		if c.SpanKind() != otel {
			t.Errorf("conversion of SpanKind start option: got %v, want %v", c.SpanKind(), otel)
		}
	}
}

func TestStartOptionsSamplerErrors(t *testing.T) {
	ocOpts := []octrace.StartOption{octrace.WithSampler(octrace.AlwaysSample())}
	_, err := StartOptions(ocOpts)
	if err == nil {
		t.Error("StartOptions should error Sampler option")
	}
}
