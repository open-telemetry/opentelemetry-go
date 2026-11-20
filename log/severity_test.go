// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

func TestSeverity(t *testing.T) {
	// Test the Severity constants match the OTel values and short names.
	testCases := []struct {
		name     string
		severity log.Severity
		value    int
		str      string
	}{
		{
			name:     "SeverityUndefined",
			severity: log.SeverityUndefined,
			value:    0,
			str:      "UNDEFINED",
		},
		{
			name:     "SeverityTrace",
			severity: log.SeverityTrace,
			value:    1,
			str:      "TRACE",
		},
		{
			name:     "SeverityTrace1",
			severity: log.SeverityTrace1,
			value:    1,
			str:      "TRACE",
		},
		{
			name:     "SeverityTrace2",
			severity: log.SeverityTrace2,
			value:    2,
			str:      "TRACE2",
		},
		{
			name:     "SeverityTrace3",
			severity: log.SeverityTrace3,
			value:    3,
			str:      "TRACE3",
		},
		{
			name:     "SeverityTrace4",
			severity: log.SeverityTrace4,
			value:    4,
			str:      "TRACE4",
		},
		{
			name:     "SeverityDebug",
			severity: log.SeverityDebug,
			value:    5,
			str:      "DEBUG",
		},
		{
			name:     "SeverityDebug1",
			severity: log.SeverityDebug1,
			value:    5,
			str:      "DEBUG",
		},
		{
			name:     "SeverityDebug2",
			severity: log.SeverityDebug2,
			value:    6,
			str:      "DEBUG2",
		},
		{
			name:     "SeverityDebug3",
			severity: log.SeverityDebug3,
			value:    7,
			str:      "DEBUG3",
		},
		{
			name:     "SeverityDebug4",
			severity: log.SeverityDebug4,
			value:    8,
			str:      "DEBUG4",
		},
		{
			name:     "SeverityInfo",
			severity: log.SeverityInfo,
			value:    9,
			str:      "INFO",
		},
		{
			name:     "SeverityInfo1",
			severity: log.SeverityInfo1,
			value:    9,
			str:      "INFO",
		},
		{
			name:     "SeverityInfo2",
			severity: log.SeverityInfo2,
			value:    10,
			str:      "INFO2",
		},
		{
			name:     "SeverityInfo3",
			severity: log.SeverityInfo3,
			value:    11,
			str:      "INFO3",
		},
		{
			name:     "SeverityInfo4",
			severity: log.SeverityInfo4,
			value:    12,
			str:      "INFO4",
		},
		{
			name:     "SeverityWarn",
			severity: log.SeverityWarn,
			value:    13,
			str:      "WARN",
		},
		{
			name:     "SeverityWarn1",
			severity: log.SeverityWarn1,
			value:    13,
			str:      "WARN",
		},
		{
			name:     "SeverityWarn2",
			severity: log.SeverityWarn2,
			value:    14,
			str:      "WARN2",
		},
		{
			name:     "SeverityWarn3",
			severity: log.SeverityWarn3,
			value:    15,
			str:      "WARN3",
		},
		{
			name:     "SeverityWarn4",
			severity: log.SeverityWarn4,
			value:    16,
			str:      "WARN4",
		},
		{
			name:     "SeverityError",
			severity: log.SeverityError,
			value:    17,
			str:      "ERROR",
		},
		{
			name:     "SeverityError1",
			severity: log.SeverityError1,
			value:    17,
			str:      "ERROR",
		},
		{
			name:     "SeverityError2",
			severity: log.SeverityError2,
			value:    18,
			str:      "ERROR2",
		},
		{
			name:     "SeverityError3",
			severity: log.SeverityError3,
			value:    19,
			str:      "ERROR3",
		},
		{
			name:     "SeverityError4",
			severity: log.SeverityError4,
			value:    20,
			str:      "ERROR4",
		},
		{
			name:     "SeverityFatal",
			severity: log.SeverityFatal,
			value:    21,
			str:      "FATAL",
		},
		{
			name:     "SeverityFatal1",
			severity: log.SeverityFatal1,
			value:    21,
			str:      "FATAL",
		},
		{
			name:     "SeverityFatal2",
			severity: log.SeverityFatal2,
			value:    22,
			str:      "FATAL2",
		},
		{
			name:     "SeverityFatal3",
			severity: log.SeverityFatal3,
			value:    23,
			str:      "FATAL3",
		},
		{
			name:     "SeverityFatal4",
			severity: log.SeverityFatal4,
			value:    24,
			str:      "FATAL4",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.value, int(tc.severity), "value does not match OTel")
			assert.Equal(t, tc.str, tc.severity.String(), "string does not match OTel")
		})
	}
}
