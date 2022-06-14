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

package aggtor // import "go.opentelemetry.io/otel/sdk/metric/internal/aggtor"

// Aggregator forms an aggregation from a collection of recorded measurements.
type Aggregator[N int64 | float64] interface {
	// Record includes value in the aggregation.
	Record(value N)

	// Aggregate returns an aggregation of all recorded values.
	Aggregate() Aggregation
}
