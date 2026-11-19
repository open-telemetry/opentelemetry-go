// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfilePartialSuccessError(t *testing.T) {
	err := ProfilePartialSuccessError(3, "some rejected")
	assert.ErrorIs(t, err, PartialSuccess{})
	assert.Equal(t, "OTLP partial success: some rejected (3 profiles rejected)", err.Error())
	assert.Equal(t, "rejected", PartialSuccess{}.ErrorType())

	emptyMsg := ProfilePartialSuccessError(1, "")
	assert.Equal(t, "OTLP partial success: empty message (1 profiles rejected)", emptyMsg.Error())
	assert.False(t, PartialSuccess{}.Is(errors.New("other")))
}
