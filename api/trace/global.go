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

import "sync/atomic"

type globalProvider struct {
	p Provider
}

var globalP atomic.Value

// GlobalProvider returns trace provider registered with global registry.
// If no trace provider is registered then an instance of NoopTraceProvider is returned.
// Use the trace provider to create a named tracer. E.g.
//     tracer := trace.GlobalProvider().GetTracer("example.com/foo")
func GlobalProvider() Provider {
	if gp := globalP.Load(); gp != nil {
		return gp.(globalProvider).p
	}
	return NoopTraceProvider{}
}

// SetGlobalProvider sets the provider as a global trace provider.
func SetGlobalProvider(m Provider) {
	globalP.Store(globalProvider{p: m})
}
