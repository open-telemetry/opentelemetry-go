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

package global

import (
	"sync/atomic"

	"go.opentelemetry.io/api/trace"
)

type globalProvider struct {
	p trace.Provider
}

var globalP atomic.Value

// TraceProvider returns the registered global trace provider.
// If none is registered then an instance of NoopTraceProvider is returned.
// Use the trace provider to create a named tracer. E.g.
//     tracer := global.TraceProvider().GetTracer("example.com/foo")
func TraceProvider() trace.Provider {
	if gp := globalP.Load(); gp != nil {
		return gp.(globalProvider).p
	}
	return trace.NoopTraceProvider{}
}

// SetTraceProvider registers p as the global trace provider.
func SetTraceProvider(p trace.Provider) {
	globalP.Store(globalProvider{p: p})
}
