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

// +build !go1.14

package filters

import (
	"net/http"
	"net/textproto"
	"strings"

	"go.opentelemetry.io/otel/plugin/othttp"
)

// Header returns a Filter that returns true if the request
// includes a header k with a value equal to v.
func Header(k, v string) othttp.Filter {
	return func(r *http.Request) bool {
		for _, hv := range r.Header[textproto.CanonicalMIMEHeaderKey(k)] {
			if v == hv {
				return true
			}
		}
		return false
	}
}

// HeaderContains returns a Filter that returns true if the request
// includes a header k with a value that contains v.
func HeaderContains(k, v string) othttp.Filter {
	return func(r *http.Request) bool {
		for _, hv := range r.Header[textproto.CanonicalMIMEHeaderKey(k)] {
			if strings.Contains(hv, v) {
				return true
			}
		}
		return false
	}
}
