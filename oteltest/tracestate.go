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

package oteltest

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TraceStateFromKeyValues is a convenience function to create a
// trace.TraceState from provided key/value pairs. There is no inverse to this
// function, returning attributes from a TraceState, because the TraceState,
// by definition from the W3C tracecontext specification, stores values as
// opaque strings.  Therefore, it is not possible to decode the original value
// type from TraceState. Be sure to not use this outside of testing purposes.
//
// Deprecated: use trace.ParseTraceState instead.
func TraceStateFromKeyValues(kvs ...attribute.KeyValue) (trace.TraceState, error) {
	if len(kvs) == 0 {
		return trace.TraceState{}, nil
	}

	members := make([]string, len(kvs))
	for i, kv := range kvs {
		members[i] = fmt.Sprintf("%s=%s", string(kv.Key), kv.Value.Emit())
	}
	return trace.ParseTraceState(strings.Join(members, ","))
}
