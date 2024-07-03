// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrop(t *testing.T) {
	t.Run("Int64", testDropFiltered[int64])
	t.Run("Float64", testDropFiltered[float64])
}

func testDropFiltered[N int64 | float64](t *testing.T) {
	r := Drop[N]()

	var dest []Exemplar
	r.Collect(&dest)

	assert.Len(t, dest, 0, "non-sampled context should not be offered")
}
