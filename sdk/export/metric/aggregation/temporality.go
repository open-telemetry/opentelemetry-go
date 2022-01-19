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

package aggregation // import "go.opentelemetry.io/otel/sdk/export/metric/aggregation"

import (
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Temporality = aggregation.Temporality

const (
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	CumulativeTemporality = aggregation.CumulativeTemporality

	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	DeltaTemporality = aggregation.DeltaTemporality
)

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
func ConstantTemporalitySelector(t Temporality) TemporalitySelector {
	return aggregation.ConstantTemporalitySelector(t)
}

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
func CumulativeTemporalitySelector() TemporalitySelector {
	return aggregation.CumulativeTemporalitySelector()
}

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
func DeltaTemporalitySelector() TemporalitySelector {
	return aggregation.DeltaTemporalitySelector()
}

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
func StatelessTemporalitySelector() TemporalitySelector {
	return aggregation.StatelessTemporalitySelector()
}

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type TemporalitySelector = aggregation.TemporalitySelector
