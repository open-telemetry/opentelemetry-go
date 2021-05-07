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

package v1_4_0

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv"
)

func NetAttributesFromHTTPRequest(network string, request *http.Request) []attribute.KeyValue {
	return semconv.NetAttributesFromHTTPRequest(network, request)
}

func EndUserAttributesFromHTTPRequest(request *http.Request) []attribute.KeyValue {
	return semconv.EndUserAttributesFromHTTPRequest(request)
}

func HTTPClientAttributesFromHTTPRequest(request *http.Request) []attribute.KeyValue {
	return semconv.HTTPClientAttributesFromHTTPRequest(request)
}

func HTTPServerMetricAttributesFromHTTPRequest(serverName string, request *http.Request) []attribute.KeyValue {
	return semconv.HTTPServerMetricAttributesFromHTTPRequest(serverName, request)
}

func HTTPServerAttributesFromHTTPRequest(serverName, route string, request *http.Request) []attribute.KeyValue {
	return semconv.HTTPServerAttributesFromHTTPRequest(serverName, route, request)
}

func HTTPAttributesFromHTTPStatusCode(code int) []attribute.KeyValue {
	return semconv.HTTPAttributesFromHTTPStatusCode(code)
}

func SpanStatusFromHTTPStatusCode(code int) (codes.Code, string) {
	return semconv.SpanStatusFromHTTPStatusCode(code)
}
