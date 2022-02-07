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
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Aggregation = aggregation.Aggregation

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Sum = aggregation.Sum

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Count = aggregation.Count

// Deprecated: Will be removed soon.
type Min interface {
	Aggregation
	Min() (number.Number, error)
}

// Deprecated: Will be removed soon.
type Max interface {
	Aggregation
	Max() (number.Number, error)
}

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type LastValue = aggregation.LastValue

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Buckets = aggregation.Buckets

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Histogram = aggregation.Histogram

// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
type Kind = aggregation.Kind

const (
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	SumKind = aggregation.SumKind
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	HistogramKind = aggregation.HistogramKind
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	LastValueKind = aggregation.LastValueKind
)

var (
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	ErrNegativeInput = aggregation.ErrNegativeInput
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	ErrNaNInput = aggregation.ErrNaNInput
	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	ErrInconsistentType = aggregation.ErrInconsistentType

	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	ErrNoCumulativeToDelta = aggregation.ErrNoCumulativeToDelta

	// Deprecated: use module "go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	ErrNoData = aggregation.ErrNoData
)
