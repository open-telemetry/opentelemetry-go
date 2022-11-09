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

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// Float64Counter is an instrument that records increasing incremental float64
// values.
//
// Warning: methods may be added to this interface in minor releases.
type Float64Counter interface {
	// Add records a positive increment measurement for attrs.
	//
	// The behavior of passing a non-positive increment is undefined. Specific
	// implementations may document their own behavior.
	Add(ctx context.Context, increment float64, attrs ...attribute.KeyValue)
}

// Int64Counter is an instrument that records increasing incremental int64
// values.
//
// Warning: methods may be added to this interface in minor releases.
type Int64Counter interface {
	// Add records a positive increment measurement for attrs.
	//
	// The behavior of passing a non-positive increment is undefined. Specific
	// implementations may document their own behavior.
	Add(ctx context.Context, increment int64, attrs ...attribute.KeyValue)
}

// Float64Counter is an instrument that records increasing or decreasing
// incremental float64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Float64UpDownCounter interface {
	// Add records an increment measurement for attrs.
	Add(ctx context.Context, increment float64, attrs ...attribute.KeyValue)
}

// Int64UpDownCounter is an instrument that records increasing or decreasing
// incremental int64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Int64UpDownCounter interface {
	// Add records an increment measurement for attrs.
	Add(ctx context.Context, increment int64, attrs ...attribute.KeyValue)
}

// Float64Histogram is an instrument that records a distribution of float64
// values.
//
// Warning: methods may be added to this interface in minor releases.
type Float64Histogram interface {
	// Record records value as part of a distribution for attrs.
	Record(ctx context.Context, value float64, attrs ...attribute.KeyValue)
}

// Int64Histogram is an instrument that records a distribution of int64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Int64Histogram interface {
	// Record records value as part of a distribution for attrs.
	Record(ctx context.Context, value int64, attrs ...attribute.KeyValue)
}
