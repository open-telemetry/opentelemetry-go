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

// The process global tracer could have process-wide resource
// tags applied directly, or we can have a SetGlobal tracer to
// install a default tracer w/ resources.
var global atomic.Value

var _ Tracer = noopTracer{}

// GlobalTracer return tracer registered with global registry.
// If no tracer is registered then an instance of noop Tracer is returned.
func GlobalTracer() Tracer {
	if t := global.Load(); t != nil {
		return t.(Tracer)
	}
	return noopTracer{}
}

// SetGlobalTracer sets provided tracer as a global tracer.
func SetGlobalTracer(t Tracer) {
	global.Store(t)
}
