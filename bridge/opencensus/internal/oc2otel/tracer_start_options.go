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

package oc2otel // import "go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"

import (
	"fmt"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func StartOptions(optFuncs []octrace.StartOption) ([]trace.SpanStartOption, error) {
	var ocOpts octrace.StartOptions
	for _, fn := range optFuncs {
		fn(&ocOpts)
	}

	var otelOpts []trace.SpanStartOption
	switch ocOpts.SpanKind {
	case octrace.SpanKindClient:
		otelOpts = append(otelOpts, trace.WithSpanKind(trace.SpanKindClient))
	case octrace.SpanKindServer:
		otelOpts = append(otelOpts, trace.WithSpanKind(trace.SpanKindServer))
	case octrace.SpanKindUnspecified:
		otelOpts = append(otelOpts, trace.WithSpanKind(trace.SpanKindUnspecified))
	}

	var err error
	if ocOpts.Sampler != nil {
		err = fmt.Errorf("unsupported sampler: %v", ocOpts.Sampler)
	}
	return otelOpts, err
}
