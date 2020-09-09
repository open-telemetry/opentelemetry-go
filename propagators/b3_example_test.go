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

package propagators_test

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/propagators"
)

func ExampleB3() {
	b3 := propagators.B3{}
	// Register the B3 propagator globally.
	global.SetPropagators(propagation.New(
		propagation.WithExtractors(b3),
		propagation.WithInjectors(b3),
	))
}

func ExampleB3_injectEncoding() {
	// Create a B3 propagator configured to inject context with both multiple
	// and single header B3 HTTP encoding.
	b3 := propagators.B3{
		InjectEncoding: propagators.B3MultipleHeader | propagators.B3SingleHeader,
	}
	global.SetPropagators(propagation.New(
		propagation.WithExtractors(b3),
		propagation.WithInjectors(b3),
	))
}
