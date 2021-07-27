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
	"fmt"

	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartOptions(optFns []octrace.StartOption, name string) []trace.SpanStartOption {
	var ocOpts octrace.StartOptions
	for _, fn := range optFns {
		fn(&ocOpts)
	}
	otOpts := []trace.SpanStartOption{}
	switch ocOpts.SpanKind {
	case octrace.SpanKindClient:
		otOpts = append(otOpts, trace.WithSpanKind(trace.SpanKindClient))
	case octrace.SpanKindServer:
		otOpts = append(otOpts, trace.WithSpanKind(trace.SpanKindServer))
	case octrace.SpanKindUnspecified:
		otOpts = append(otOpts, trace.WithSpanKind(trace.SpanKindUnspecified))
	}

	if ocOpts.Sampler != nil {
		otel.Handle(fmt.Errorf("ignoring custom sampler for span %q created by OpenCensus because OpenTelemetry does not support creating a span with a custom sampler", name))
	}
	return otOpts
}
