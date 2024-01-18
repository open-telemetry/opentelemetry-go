// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

func TestSeverity(t *testing.T) {
	testCases := []struct {
		name     string
		severity log.Severity
		value    int
	}{
		{
			name:     "SeverityTrace",
			severity: log.SeverityTrace,
			value:    1,
		},
		{
			name:     "SeverityTrace1",
			severity: log.SeverityTrace1,
			value:    1,
		},
		{
			name:     "SeverityTrace2",
			severity: log.SeverityTrace2,
			value:    2,
		},
		{
			name:     "SeverityTrace3",
			severity: log.SeverityTrace3,
			value:    3,
		},
		{
			name:     "SeverityTrace4",
			severity: log.SeverityTrace4,
			value:    4,
		},
		{
			name:     "SeverityDebug",
			severity: log.SeverityDebug,
			value:    5,
		},
		{
			name:     "SeverityDebug1",
			severity: log.SeverityDebug1,
			value:    5,
		},
		{
			name:     "SeverityDebug2",
			severity: log.SeverityDebug2,
			value:    6,
		},
		{
			name:     "SeverityDebug3",
			severity: log.SeverityDebug3,
			value:    7,
		},
		{
			name:     "SeverityDebug4",
			severity: log.SeverityDebug4,
			value:    8,
		},
		{
			name:     "SeverityInfo",
			severity: log.SeverityInfo,
			value:    9,
		},
		{
			name:     "SeverityInfo1",
			severity: log.SeverityInfo1,
			value:    9,
		},
		{
			name:     "SeverityInfo2",
			severity: log.SeverityInfo2,
			value:    10,
		},
		{
			name:     "SeverityInfo3",
			severity: log.SeverityInfo3,
			value:    11,
		},
		{
			name:     "SeverityInfo4",
			severity: log.SeverityInfo4,
			value:    12,
		},
		{
			name:     "SeverityWarn",
			severity: log.SeverityWarn,
			value:    13,
		},
		{
			name:     "SeverityWarn1",
			severity: log.SeverityWarn1,
			value:    13,
		},
		{
			name:     "SeverityWarn2",
			severity: log.SeverityWarn2,
			value:    14,
		},
		{
			name:     "SeverityWarn3",
			severity: log.SeverityWarn3,
			value:    15,
		},
		{
			name:     "SeverityWarn4",
			severity: log.SeverityWarn4,
			value:    16,
		},
		{
			name:     "SeverityError",
			severity: log.SeverityError,
			value:    17,
		},
		{
			name:     "SeverityError1",
			severity: log.SeverityError1,
			value:    17,
		},
		{
			name:     "SeverityError2",
			severity: log.SeverityError2,
			value:    18,
		},
		{
			name:     "SeverityError3",
			severity: log.SeverityError3,
			value:    19,
		},
		{
			name:     "SeverityError4",
			severity: log.SeverityError4,
			value:    20,
		},
		{
			name:     "SeverityFatal",
			severity: log.SeverityFatal,
			value:    21,
		},
		{
			name:     "SeverityFatal1",
			severity: log.SeverityFatal1,
			value:    21,
		},
		{
			name:     "SeverityFatal2",
			severity: log.SeverityFatal2,
			value:    22,
		},
		{
			name:     "SeverityFatal3",
			severity: log.SeverityFatal3,
			value:    23,
		},
		{
			name:     "SeverityFatal4",
			severity: log.SeverityFatal4,
			value:    24,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.value, int(tc.severity))
		})
	}
}
