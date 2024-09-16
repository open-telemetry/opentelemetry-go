// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/schema/internal"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFileFormatField(t *testing.T) {
	// Invalid file format version numbers.
	assert.Error(t, CheckFileFormatField("not a semver", 1, 0))
	assert.Error(t, CheckFileFormatField("2.0.0", 1, 0))
	assert.Error(t, CheckFileFormatField("1.1.0", 1, 0))
	assert.Error(t, CheckFileFormatField("1.1.0", -1, 0))
	assert.Error(t, CheckFileFormatField("1.1.0", 1, -2))

	assert.Error(t, CheckFileFormatField("1.2.0", 1, 1))

	// Valid cases.
	assert.NoError(t, CheckFileFormatField("1.0.0", 1, 0))
	assert.NoError(t, CheckFileFormatField("1.0.1", 1, 0))
	assert.NoError(t, CheckFileFormatField("1.0.10000-alpha+4857", 1, 0))

	assert.NoError(t, CheckFileFormatField("1.0.0", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.0.1", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.0.10000-alpha+4857", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.1.0", 1, 1))
	assert.NoError(t, CheckFileFormatField("1.1.1", 1, 1))
}
