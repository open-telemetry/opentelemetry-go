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

package othttp

import (
	"net/http"

	"go.opentelemetry.io/otel/api/kv"

	"go.opentelemetry.io/otel/api/trace"
)

// Attribute keys that can be added to a span.
const (
	HostKey       = kv.Key("http.host")        // the HTTP host (http.Request.Host)
	MethodKey     = kv.Key("http.method")      // the HTTP method (http.Request.Method)
	PathKey       = kv.Key("http.path")        // the HTTP path (http.Request.URL.Path)
	URLKey        = kv.Key("http.url")         // the HTTP URL (http.Request.URL.String())
	UserAgentKey  = kv.Key("http.user_agent")  // the HTTP user agent (http.Request.UserAgent())
	RouteKey      = kv.Key("http.route")       // the HTTP route (ex: /users/:id)
	RemoteAddrKey = kv.Key("http.remote_addr") // the network address of the client that sent the HTTP request (http.Request.RemoteAddr)
	StatusCodeKey = kv.Key("http.status_code") // if set, the HTTP status
	ReadBytesKey  = kv.Key("http.read_bytes")  // if anything was read from the request body, the total number of bytes read
	ReadErrorKey  = kv.Key("http.read_error")  // If an error occurred while reading a request, the string of the error (io.EOF is not recorded)
	WroteBytesKey = kv.Key("http.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
	WriteErrorKey = kv.Key("http.write_error") // if an error occurred while writing a reply, the string of the error (io.EOF is not recorded)
)

// Filter is a predicate used to determine whether a given http.request should
// be traced. A Filter must return true if the request should be traced.
type Filter func(*http.Request) bool

// Setup basic span attributes before so that they
// are available to be mutated if needed.
func setBasicAttributes(span trace.Span, r *http.Request) {
	span.SetAttributes(
		HostKey.String(r.Host),
		MethodKey.String(r.Method),
		PathKey.String(r.URL.Path),
		URLKey.String(r.URL.String()),
		UserAgentKey.String(r.UserAgent()),
	)
}
