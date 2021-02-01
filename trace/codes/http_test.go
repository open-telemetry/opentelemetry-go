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

package codes

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpanStatusFromHTTPStatusCode(t *testing.T) {
	for code := 0; code < 1000; code++ {
		expected := getExpectedCodeForHTTPCode(code)
		got, _ := SpanStatusFromHTTPStatusCode(code)
		assert.Equalf(t, expected, got, "%s vs %s", expected, got)
	}
}

func getExpectedCodeForHTTPCode(code int) Code {
	if http.StatusText(code) == "" {
		return Error
	}
	switch code {
	case
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusTooManyRequests,
		http.StatusNotImplemented,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return Error
	}
	category := code / 100
	if category > 0 && category < 4 {
		return Unset
	}
	return Error
}
