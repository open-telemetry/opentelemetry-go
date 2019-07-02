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

import (
	"testing"
)

// GlobalTracer return tracer registered with global registry.
// If no tracer is registered then an instance of noop Tracer is returned.
func TestGlobalTracer(t *testing.T) {
	tracer := GlobalTracer()
	if tracer == nil {
		t.Fatalf("Failed to get tracer\n")
	}
}
