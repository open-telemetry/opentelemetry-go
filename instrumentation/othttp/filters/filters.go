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

// Package filters provides a set of filters useful with the
// othttp.WithFilter() option to control which inbound requests are traced.
package filters

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/plugin/othttp"
)

// Any takes a list of Filters and returns a Filter that
// returns true if any Filter in the list returns true.
func Any(fs ...othttp.Filter) othttp.Filter {
	return func(r *http.Request) bool {
		for _, f := range fs {
			if f(r) {
				return true
			}
		}
		return false
	}
}

// All takes a list of Filters and returns a Filter that
// returns true only if all Filters in the list return true.
func All(fs ...othttp.Filter) othttp.Filter {
	return func(r *http.Request) bool {
		for _, f := range fs {
			if !f(r) {
				return false
			}
		}
		return true
	}
}

// None takes a list of Filters and returns a Filter that returns
// true only if none of the Filters in the list return true.
func None(fs ...othttp.Filter) othttp.Filter {
	return func(r *http.Request) bool {
		for _, f := range fs {
			if f(r) {
				return false
			}
		}
		return true
	}
}

// Not provides a convenience mechanism for inverting a Filter
func Not(f othttp.Filter) othttp.Filter {
	return func(r *http.Request) bool {
		return !f(r)
	}
}

// Hostname returns a Filter that returns true if the request's
// hostname matches the provided string.
func Hostname(h string) othttp.Filter {
	return func(r *http.Request) bool {
		return r.URL.Hostname() == h
	}
}

// Path returns a Filter that returns true if the request's
// path matches the provided string.
func Path(p string) othttp.Filter {
	return func(r *http.Request) bool {
		return r.URL.Path == p
	}
}

// PathPrefix returns a Filter that returns true if the request's
// path starts with the provided string.
func PathPrefix(p string) othttp.Filter {
	return func(r *http.Request) bool {
		return strings.HasPrefix(r.URL.Path, p)
	}
}

// Query returns a Filter that returns true if the request
// includes a query parameter k with a value equal to v.
func Query(k, v string) othttp.Filter {
	return func(r *http.Request) bool {
		for _, qv := range r.URL.Query()[k] {
			if v == qv {
				return true
			}
		}
		return false
	}
}

// QueryContains returns a Filter that returns true if the request
// includes a query parameter k with a value that contains v.
func QueryContains(k, v string) othttp.Filter {
	return func(r *http.Request) bool {
		for _, qv := range r.URL.Query()[k] {
			if strings.Contains(qv, v) {
				return true
			}
		}
		return false
	}
}

// Method returns a Filter that returns true if the request
// method is equal to the provided value.
func Method(m string) othttp.Filter {
	return func(r *http.Request) bool {
		return m == r.Method
	}
}
