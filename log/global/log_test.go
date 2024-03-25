// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/log/global"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
)

func TestMultipleGlobalLoggerProvider(t *testing.T) {
	type provider struct{ log.LoggerProvider }

	p1, p2 := provider{}, noop.NewLoggerProvider()

	SetLoggerProvider(&p1)
	SetLoggerProvider(p2)

	assert.Equal(t, p2, GetLoggerProvider())
}
