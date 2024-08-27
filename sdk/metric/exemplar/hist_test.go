// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import "testing"

func TestHist(t *testing.T) {
	bounds := []float64{0, 100}
	t.Run("Int64", ReservoirTest[int64](func(int) (Reservoir, int) {
		return Histogram(bounds), len(bounds)
	}))

	t.Run("Float64", ReservoirTest[float64](func(int) (Reservoir, int) {
		return Histogram(bounds), len(bounds)
	}))
}
