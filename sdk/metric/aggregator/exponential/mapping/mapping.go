// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mapping // import "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"

import "fmt"

// Mapping is the interface of a mapper.
type Mapping interface {
	// MapToIndex maps positive floating point values to indexes
	// corresponding to Scale().  Implementations are not expected
	// to handle zeros, +Inf, NaN, or negative values.
	MapToIndex(value float64) int32
	LowerBoundary(index int32) (float64, error)
	Scale() int32
}

var (
	ErrUnderflow = fmt.Errorf("underflow")
	ErrOverflow  = fmt.Errorf("overflow")
)
