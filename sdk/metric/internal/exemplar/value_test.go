// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	const iVal, fVal = int64(43), float64(0.3)
	i, f, bad := NewValue[int64](iVal), NewValue[float64](fVal), Value{}

	assert.Equal(t, Int64ValueType, i.Type())
	assert.Equal(t, iVal, i.Int64())
	assert.Equal(t, float64(0), i.Float64())

	assert.Equal(t, Float64ValueType, f.Type())
	assert.Equal(t, fVal, f.Float64())
	assert.Equal(t, int64(0), f.Int64())

	assert.Equal(t, UnknownValueType, bad.Type())
	assert.Equal(t, float64(0), bad.Float64())
	assert.Equal(t, int64(0), bad.Int64())
}
