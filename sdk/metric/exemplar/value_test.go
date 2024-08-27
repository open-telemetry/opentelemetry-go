// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	const iVal, fVal, nVal = int64(43), float64(0.3), int64(-42)
	i, f, n, bad := NewValue[int64](iVal), NewValue[float64](fVal), NewValue[int64](nVal), Value{}

	assert.Equal(t, Int64ValueType, i.Type())
	assert.Equal(t, iVal, i.Int64())
	assert.Equal(t, float64(0), i.Float64())

	assert.Equal(t, Float64ValueType, f.Type())
	assert.Equal(t, fVal, f.Float64())
	assert.Equal(t, int64(0), f.Int64())

	assert.Equal(t, Int64ValueType, n.Type())
	assert.Equal(t, nVal, n.Int64())
	assert.Equal(t, float64(0), i.Float64())

	assert.Equal(t, UnknownValueType, bad.Type())
	assert.Equal(t, float64(0), bad.Float64())
	assert.Equal(t, int64(0), bad.Int64())
}
