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

package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/transform"

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

// Metrics returns a slice of OTLP Metric generated from ms. If ms contains
// invalid metric values, an error will be returned along with a slice that
// contains partial OTLP Metrics.
func Metrics(ms []metricdata.Metrics) ([]*mpb.Metric, error) {
	var errs []string
	out := make([]*mpb.Metric, 0, len(ms))
	for _, m := range ms {
		o, err := metric(m)
		if err != nil {
			errs = append(errs, err.Error())
		}
		out = append(out, o)
	}
	if len(errs) > 0 {
		return out, fmt.Errorf("transform Metrics: %s", strings.Join(errs, ", "))
	}
	return out, nil
}

func metric(m metricdata.Metrics) (*mpb.Metric, error) {
	var err error
	out := &mpb.Metric{
		Name:        m.Name,
		Description: m.Description,
		Unit:        string(m.Unit),
	}
	switch a := m.Data.(type) {
	case metricdata.Gauge[int64]:
		out.Data = Gauge[int64](a)
	case metricdata.Gauge[float64]:
		out.Data = Gauge[float64](a)
	case metricdata.Sum[int64]:
		out.Data, err = Sum[int64](a)
	case metricdata.Sum[float64]:
		out.Data, err = Sum[float64](a)
	case metricdata.Histogram:
		out.Data, err = Histogram(a)
	default:
		return out, fmt.Errorf("unknown aggregation: %T", a)
	}
	return out, err
}
