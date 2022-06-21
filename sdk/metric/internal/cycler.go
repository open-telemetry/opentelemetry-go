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

//go:build go1.18
// +build go1.18

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

// Cycler cycles aggregation periods. It will handle any state progression
// from one period to the next based on the temporality of the cycling.
type Cycler interface {
	// Cycle returns an []Aggregation for the current period. If the cycler
	// merges state from previous periods into the current, the []Aggregation
	// returned reflects this.
	Cycle() []Aggregation

	// TODO: Replace the return type with []export.Aggregation once #2961 is
	// merged.
}

// deltaCylcer cycles aggregation periods by returning the aggregation
// produces from that period only. No state is maintained from one period to
// the next.
type deltaCylcer[N int64 | float64] struct {
	aggregator Aggregator[N]
}

func NewDeltaCylcer[N int64 | float64](a Aggregator[N]) Cycler {
	return deltaCylcer[N]{aggregator: a}
}

func (c deltaCylcer[N]) Cycle() []Aggregation {
	return c.aggregator.flush()
}

// cumulativeCylcer cycles aggregation periods by returning the cumulative
// aggregation from its start time until the current period.
type cumulativeCylcer[N int64 | float64] struct {
	// TODO: implement a cumulative storing field.
	aggregator Aggregator[N]
}

func NewCumulativeCylcer[N int64 | float64](a Aggregator[N]) Cycler {
	c := cumulativeCylcer[N]{aggregator: a}

	// TODO: Initialize a new cumulative storage.

	return c
}

func (c cumulativeCylcer[N]) Cycle() []Aggregation {
	// TODO: Update cumulative storage of aggregations and return them.

	// FIXME: currently this returns a delta representation of the
	// aggregation. When the cumulative storage is complete it should return a
	// cumulative representation.
	return c.aggregator.flush()
}
