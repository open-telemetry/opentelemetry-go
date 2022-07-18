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

package exporttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestResourceMetricsComparison(t *testing.T) {
	a := metricdata.ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("resource", "a")),
	}

	b := metricdata.ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("resource", "b")),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{Scope: instrumentation.Scope{Name: "b"}},
		},
	}

	AssertResourceMetricsEqual(t, a, a)
	AssertResourceMetricsEqual(t, b, b)

	equal, explination := CompareResourceMetrics(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Resource and ScopeMetrics do not match")
}
