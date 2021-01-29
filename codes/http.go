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

package codes // import "go.opentelemetry.io/otel/codes"

import (
	"fmt"
	"net/http"
)

type codeRange struct {
	fromInclusive int
	toInclusive   int
}

func (r codeRange) contains(code int) bool {
	return r.fromInclusive <= code && code <= r.toInclusive
}

var validRangesPerCategory = map[int][]codeRange{
	1: {
		{http.StatusContinue, http.StatusEarlyHints},
	},
	2: {
		{http.StatusOK, http.StatusAlreadyReported},
		{http.StatusIMUsed, http.StatusIMUsed},
	},
	3: {
		{http.StatusMultipleChoices, http.StatusUseProxy},
		{http.StatusTemporaryRedirect, http.StatusPermanentRedirect},
	},
	4: {
		{http.StatusBadRequest, http.StatusTeapot}, // yes, teapot is so useful…
		{http.StatusMisdirectedRequest, http.StatusUpgradeRequired},
		{http.StatusPreconditionRequired, http.StatusTooManyRequests},
		{http.StatusRequestHeaderFieldsTooLarge, http.StatusRequestHeaderFieldsTooLarge},
		{http.StatusUnavailableForLegalReasons, http.StatusUnavailableForLegalReasons},
	},
	5: {
		{http.StatusInternalServerError, http.StatusLoopDetected},
		{http.StatusNotExtended, http.StatusNetworkAuthenticationRequired},
	},
}

// SpanStatusFromHTTPStatusCode generates a status code and a message
// as specified by the OpenTelemetry specification for a span.
func SpanStatusFromHTTPStatusCode(code int) (Code, string) {
	spanCode, valid := func() (Code, bool) {
		category := code / 100
		ranges, ok := validRangesPerCategory[category]
		if !ok {
			return Error, false
		}
		ok = false
		for _, crange := range ranges {
			ok = crange.contains(code)
			if ok {
				break
			}
		}
		if !ok {
			return Error, false
		}
		if category > 0 && category < 4 {
			return Unset, true
		}
		return Error, true
	}()
	if !valid {
		return spanCode, fmt.Sprintf("Invalid HTTP status code %d", code)
	}
	return spanCode, fmt.Sprintf("HTTP status code: %d", code)
}
