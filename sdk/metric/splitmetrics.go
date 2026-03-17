// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// splitResourceMetrics splits a metricdata.ResourceMetrics into multiple ResourceMetrics, sequentially,
// ensuring no ResourceMetrics has more than `size` data points. It does not mutate the `src` object.
func splitResourceMetrics(size int, src *metricdata.ResourceMetrics) []*metricdata.ResourceMetrics {
	// TODO: implement
	return []*metricdata.ResourceMetrics{src}
}
