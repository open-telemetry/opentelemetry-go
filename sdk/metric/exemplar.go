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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"os"

	"go.opentelemetry.io/otel/sdk/metric/internal/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/internal/x"
)

// reservoirFunc returns the appropriately configured exemplar reservoir
// creation func based on the passed InstrumentKind and user defined
// environment variables.
//
// Note: This will only return non-nil values when the experimental exemplar
// feature is enabled.
func reservoirFunc[N int64 | float64](agg Aggregation) func() exemplar.Reservoir[N] {
	if !x.Exemplars.Enabled() {
		return nil
	}

	// https://github.com/open-telemetry/opentelemetry-specification/blob/d4b241f451674e8f611bb589477680341006ad2b/specification/configuration/sdk-environment-variables.md#exemplar
	const filterEnvKey = "OTEL_METRICS_EXEMPLAR_FILTER"

	var fltr exemplar.Filter[N]
	switch os.Getenv(filterEnvKey) {
	case "always_on":
		fltr = exemplar.AlwaysSample[N]
	case "always_off":
		fltr = exemplar.NeverSample[N]
	case "trace_based":
		fallthrough
	default:
		fltr = exemplar.TraceBasedSample[N]
	}

	// TODO: This is not defined by the specification, nor is the mechanism to
	// configure it.
	const defaultFixedSize = 1

	// https://github.com/open-telemetry/opentelemetry-specification/blob/d4b241f451674e8f611bb589477680341006ad2b/specification/metrics/sdk.md#exemplar-defaults
	resF := func() func() exemplar.Reservoir[N] {
		a, ok := agg.(AggregationExplicitBucketHistogram)
		if ok && len(a.Boundaries) > 1 {
			cp := make([]float64, len(a.Boundaries))
			copy(cp, a.Boundaries)
			return func() exemplar.Reservoir[N] {
				bounds := cp
				return exemplar.Histogram[N](bounds)
			}
		}

		return func() exemplar.Reservoir[N] {
			return exemplar.FixedSize[N](defaultFixedSize)
		}
	}()

	return func() exemplar.Reservoir[N] {
		return exemplar.Filtered[N](resF(), fltr)
	}
}
