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

package view // import "go.opentelemetry.io/otel/sdk/metric/view"

import (
	"fmt"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

func Example() {
	// The "active-users" instrument created by the
	// "github.com/super/noisy/instrumentation/package" your project includes
	// has a bug, it records a measurment any time a user has any activity.
	// This is causing a lot of strain on your program without providing any
	// value to you. The next version of
	// "github.com/super/noisy/instrumentation/package" corrects the
	// instrumentation to only record a value when a user logs in, but it
	// isn't out yet.
	//
	// Use a View to drop these measurments while you wait for the fix to come
	// from upstream.

	v, err := New(
		MatchInstrumentName("active-users"),
		MatchInstrumentationScope(instrumentation.Scope{
			Name:    "github.com/super/noisy/instrumentation/package",
			Version: "v0.22.0", // Only match the problematic instrumentation version.
		}),
		WithSetAggregation(aggregation.Drop{}),
	)
	if err != nil {
		panic(err)
	}

	// The SDK this view is registered with calls TransformInstrument when an
	// instrument is created. Test that our fix will work as intended.
	i, _ := v.TransformInstrument(Instrument{
		Name: "active-users",
		Scope: instrumentation.Scope{
			Name:    "github.com/super/noisy/instrumentation/package",
			Version: "v0.22.0",
		},
		Aggregation: aggregation.LastValue{},
	})
	fmt.Printf("Instrument{%q: %s}: %#v\n", i.Name, i.Scope.Version, i.Aggregation)

	// Also, ensure the next version will not be transformed.
	_, ok := v.TransformInstrument(Instrument{
		Name: "active-users",
		Scope: instrumentation.Scope{
			Name:    "github.com/super/noisy/instrumentation/package",
			Version: "v0.23.0",
		},
		Aggregation: aggregation.LastValue{},
	})
	fmt.Printf("Instrument{\"active-users\": v0.23.0} matched: %t\n", ok)
	// Output:
	//
	// Instrument{"active-users": v0.22.0}: aggregation.Drop{}
	// Instrument{"active-users": v0.23.0} matched: false
}

func ExampleMatchInstrumentName() {
	v, err := New(MatchInstrumentName("request-*")) // Wildcard match.
	if err != nil {
		panic(err)
	}

	for _, i := range []Instrument{
		{Name: "request-count"},
		{Name: "request-rate"},
		{Name: "latency"},
	} {
		// The SDK calls TransformInstrument when an instrument is created.
		_, ok := v.TransformInstrument(i)
		fmt.Printf("Instrument{%q} matched: %t\n", i.Name, ok)
	}
	// Output:
	// Instrument{"request-count"} matched: true
	// Instrument{"request-rate"} matched: true
	// Instrument{"latency"} matched: false
}

func ExampleMatchInstrumentKind() {
	v, err := New(MatchInstrumentKind(SyncCounter))
	if err != nil {
		panic(err)
	}

	for _, i := range []Instrument{
		{Name: "synchronous counter", Kind: SyncCounter},
		{Name: "synchronous histogram", Kind: SyncHistogram},
		{Name: "asynchronous counter", Kind: AsyncCounter},
	} {
		// The SDK calls TransformInstrument when an instrument is created.
		_, ok := v.TransformInstrument(i)
		fmt.Printf("Instrument{%q} matched: %t\n", i.Name, ok)
	}
	// Output:
	// Instrument{"synchronous counter"} matched: true
	// Instrument{"synchronous histogram"} matched: false
	// Instrument{"asynchronous counter"} matched: false
}

func ExampleMatchInstrumentationScope() {
	v, err := New(MatchInstrumentationScope(instrumentation.Scope{
		Name:    "custom/instrumentation/package",
		Version: "v0.22.0", // Only match this version of instrumentation.
	}))
	if err != nil {
		panic(err)
	}

	for _, i := range []Instrument{
		{Name: "v1.0.0 instrumentation", Scope: instrumentation.Scope{
			Name:    "custom/instrumentation/package",
			Version: "v1.0.0",
		}},
		{Name: "v0.22.0 instrumentation", Scope: instrumentation.Scope{
			Name:    "custom/instrumentation/package",
			Version: "v0.22.0",
		}},
	} {
		// The SDK calls TransformInstrument when an instrument is created.
		_, ok := v.TransformInstrument(i)
		fmt.Printf("Instrument{%q} matched: %t\n", i.Name, ok)
	}
	// Output:
	// Instrument{"v1.0.0 instrumentation"} matched: false
	// Instrument{"v0.22.0 instrumentation"} matched: true
}

func ExampleWithRename() {
	v, err := New(MatchInstrumentName("bad-name"), WithRename("good-name"))
	if err != nil {
		panic(err)
	}

	// The SDK calls TransformInstrument when an instrument is created.
	i, _ := v.TransformInstrument(Instrument{Name: "bad-name"})
	fmt.Printf("Instrument{%q}\n", i.Name)
	// Output: Instrument{"good-name"}
}

func ExampleWithSetDescription() {
	v, err := New(
		MatchInstrumentName("requests"),
		WithSetDescription("Number of requests received"),
	)
	if err != nil {
		panic(err)
	}

	// The SDK calls TransformInstrument when an instrument is created.
	i, _ := v.TransformInstrument(Instrument{
		Name:        "requests",
		Description: "incorrect description",
	})
	fmt.Printf("Instrument{%q: %s}\n", i.Name, i.Description)
	// Output: Instrument{"requests": Number of requests received}
}

func ExampleWithSetAggregation() {
	v, err := New(MatchInstrumentationScope(instrumentation.Scope{
		Name: "super/noisy/instrumentation/package",
	}), WithSetAggregation(aggregation.Drop{}))
	if err != nil {
		panic(err)
	}

	// The SDK calls TransformInstrument when an instrument is created.
	i, _ := v.TransformInstrument(Instrument{
		Name: "active-users",
		Scope: instrumentation.Scope{
			Name:    "super/noisy/instrumentation/package",
			Version: "v0.5.0",
		},
		Aggregation: aggregation.LastValue{},
	})
	fmt.Printf("Instrument{%q}: %#v\n", i.Name, i.Aggregation)
	// Output: Instrument{"active-users"}: aggregation.Drop{}
}
