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

// Mapping is the interface of an exponential histogram mapper.
type Mapping interface {
	// MapToIndex maps positive floating point values to indexes
	// corresponding to Scale().  Implementations are not expected
	// to handle zeros, +Inf, NaN, or negative values.
	MapToIndex(value float64) int32

	// LowerBoundary returns the lower boundary of a given bucket
	// index.  The index is expected to map onto a range that is
	// at least partially inside the range of normalized floating
	// point values.  If the corresponding bucket's upper boundary
	// is less than or equal to 0x1p-1022, ErrUnderflow will be
	// returned.  If the corresponding bucket's lower boundary is
	// greater than math.MaxFloat64, ErrOverflow will be returned.
	LowerBoundary(index int32) (float64, error)

	// Scale returns the parameter that controls the resolution of
	// this mapping.  For details see:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/datamodel.md#exponential-scale
	Scale() int32
}

var (
	// ErrUnderflow is returned when computing the lower boundary
	// of an index that maps into a denormalized floating point value.
	ErrUnderflow = fmt.Errorf("underflow")
	// ErrOverflow is returned when computing the lower boundary
	// of an index that maps into +Inf.
	ErrOverflow = fmt.Errorf("overflow")
)
