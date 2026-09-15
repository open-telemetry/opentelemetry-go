// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import "strings"

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
