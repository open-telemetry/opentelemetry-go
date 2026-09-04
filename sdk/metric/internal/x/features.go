// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"strconv"
	"strings"
)

// MetricExportBatchSize is an experimental feature flag that controls the
// max export batch size for metric data.
//
// To enable this feature set the OTEL_GO_X_METRIC_EXPORT_BATCH_SIZE environment
// variable to a positive integer value.
var MetricExportBatchSize = newFeature(
	[]string{"METRIC_EXPORT_BATCH_SIZE"},
	func(v string) (int, bool) {
		val, err := strconv.Atoi(v)
		if err == nil && val > 0 {
			return val, true
		}
		return 0, false
	},
)

// ParallelCallbacks is an experimental feature flag that enables running
// observable-instrument callbacks concurrently during a collection.
//
// To enable this feature set the OTEL_GO_X_PARALLEL_CALLBACKS environment
// variable to the case-insensitive string value of "true".
var ParallelCallbacks = newFeature(
	[]string{"PARALLEL_CALLBACKS"},
	func(v string) (bool, bool) {
		if strings.EqualFold(v, "true") {
			return true, true
		}
		return false, false
	},
)
