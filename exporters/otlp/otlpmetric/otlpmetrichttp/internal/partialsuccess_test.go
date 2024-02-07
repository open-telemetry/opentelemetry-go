// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/partialsuccess_test.go

// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireErrorString(t *testing.T, expect string, err error) {
	t.Helper()
	require.NotNil(t, err)
	require.Error(t, err)
	require.True(t, errors.Is(err, PartialSuccess{}))

	const pfx = "OTLP partial success: "

	msg := err.Error()
	require.True(t, strings.HasPrefix(msg, pfx))
	require.Equal(t, expect, msg[len(pfx):])
}

func TestPartialSuccessFormat(t *testing.T) {
	requireErrorString(t, "empty message (0 metric data points rejected)", MetricPartialSuccessError(0, ""))
	requireErrorString(t, "help help (0 metric data points rejected)", MetricPartialSuccessError(0, "help help"))
	requireErrorString(t, "what happened (10 metric data points rejected)", MetricPartialSuccessError(10, "what happened"))
	requireErrorString(t, "what happened (15 spans rejected)", TracePartialSuccessError(15, "what happened"))
}
