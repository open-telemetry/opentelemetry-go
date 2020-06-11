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

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type testAgg struct{}

var _ export.Aggregator = (*testAgg)(nil)

func (ta *testAgg) Update(context.Context, metric.Number, *metric.Descriptor) error {
	return nil
}

func (ta *testAgg) Checkpoint(export.Aggregator, *metric.Descriptor) error {
	return nil
}

func (ta *testAgg) Merge(export.Aggregator, *metric.Descriptor) error {
	return nil
}

func TestUnslice(t *testing.T) {
	in := make([]testAgg, 2)

	a, b := Unslice2(in)

	require.Equal(t, a.(*testAgg), &in[0])
	require.Equal(t, b.(*testAgg), &in[1])
}
