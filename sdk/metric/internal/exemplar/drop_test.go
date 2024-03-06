// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"testing"
)

func TestDrop(t *testing.T) {
	t.Run("Int64", ReservoirTest[int64](func(int) (Reservoir[int64], int) {
		return Drop[int64](), 0
	}))

	t.Run("Float64", ReservoirTest[float64](func(int) (Reservoir[float64], int) {
		return Drop[float64](), 0
	}))
}
