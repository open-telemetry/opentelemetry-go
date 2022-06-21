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

import "go.opentelemetry.io/otel/attribute"

// dropAgg drops all recorded data and returns an empty Aggregation.
type dropAgg[N int64 | float64] struct{}

// NewDrop returns an Aggregator that drops all recorded data and returns an
// empty Aggregation.
func NewDrop[N int64 | float64]() Aggregator[N] { return &dropAgg[N]{} }

func (s *dropAgg[N]) Record(N, *attribute.Set) {}

func (s *dropAgg[N]) Aggregate() []Aggregation { return nil }
