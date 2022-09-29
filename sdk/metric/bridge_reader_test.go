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

package metric

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testBridge struct {
	collectFunc func(context.Context) (metricdata.ScopeMetrics, error)
}

func (t *testBridge) Collect(ctx context.Context) (metricdata.ScopeMetrics, error) {
	return t.collectFunc(ctx)
}

func TestBridgeReader(t *testing.T) {
	resource1 := resource.NewWithAttributes("test1", attribute.String("name", "resource1"))
	sm1 := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{Name: "sm1"},
		Metrics: []metricdata.Metrics{
			{Name: "metrics1"},
		},
	}
	sm2 := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{Name: "sm2"},
		Metrics: []metricdata.Metrics{
			{Name: "metrics2"},
		},
	}
	rdr1 := &reader{
		collectFunc: func(context.Context) (metricdata.ResourceMetrics, error) {
			return metricdata.ResourceMetrics{
				Resource:     resource1,
				ScopeMetrics: []metricdata.ScopeMetrics{sm1},
			}, nil
		},
	}
	bridge := &testBridge{
		collectFunc: func(context.Context) (metricdata.ScopeMetrics, error) {
			return sm2, nil
		},
	}

	wantMetrics := metricdata.ResourceMetrics{
		Resource: resource1,
		ScopeMetrics: []metricdata.ScopeMetrics{
			sm1,
			sm2,
		},
	}

	rdr := NewBridgedReader(rdr1, bridge)

	rm, err := rdr.Collect(context.Background())
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantMetrics, rm, metricdatatest.IgnoreTimestamp())
}

func TestBridgeReaderError(t *testing.T) {
	resource1 := resource.NewWithAttributes("test1", attribute.String("name", "resource1"))
	sm1 := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{Name: "sm1"},
		Metrics: []metricdata.Metrics{
			{Name: "metrics1"},
		},
	}

	rdr1 := &reader{
		collectFunc: func(context.Context) (metricdata.ResourceMetrics, error) {
			return metricdata.ResourceMetrics{
				Resource:     resource1,
				ScopeMetrics: []metricdata.ScopeMetrics{sm1},
			}, nil
		},
	}
	bridge := &testBridge{
		collectFunc: func(context.Context) (metricdata.ScopeMetrics, error) {
			return metricdata.ScopeMetrics{}, fmt.Errorf("Test Error")
		},
	}

	wantMetrics := metricdata.ResourceMetrics{
		Resource: resource1,
		ScopeMetrics: []metricdata.ScopeMetrics{
			sm1,
		},
	}

	rdr := NewBridgedReader(rdr1, bridge)

	rm, err := rdr.Collect(context.Background())
	assert.Error(t, err)
	metricdatatest.AssertEqual(t, wantMetrics, rm, metricdatatest.IgnoreTimestamp())
}
