// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

type testError string

var _ error = testError("")

func newTestError(s string) error {
	return testError(s)
}

func (e testError) Error() string {
	return string(e)
}
