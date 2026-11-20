// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfilePartialSuccessError(t *testing.T) {
	err := ProfilePartialSuccessError(4, "rejected")
	assert.ErrorIs(t, err, PartialSuccess{})
	assert.Equal(t, "OTLP partial success: rejected (4 profiles rejected)", err.Error())
	assert.Equal(t, "rejected", PartialSuccess{}.ErrorType())
	assert.False(t, PartialSuccess{}.Is(errors.New("other")))

	emptyMsgErr := ProfilePartialSuccessError(2, "")
	assert.ErrorIs(t, emptyMsgErr, PartialSuccess{})
	assert.Equal(t, "OTLP partial success: empty message (2 profiles rejected)", emptyMsgErr.Error())

	zeroRejectedErr := ProfilePartialSuccessError(0, "warning only")
	assert.ErrorIs(t, zeroRejectedErr, PartialSuccess{})
	assert.Equal(t, "OTLP partial success: warning only (0 profiles rejected)", zeroRejectedErr.Error())
}
