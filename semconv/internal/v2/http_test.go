// Copyright The OpenTelemetry Authors
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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

type tlsOption int

const (
	noTLS tlsOption = iota
	withTLS
)

var hc = &HTTPConv{
	NetConv: nc,

	EnduserIDKey:                 attribute.Key("enduser.id"),
	HTTPClientIPKey:              attribute.Key("http.client_ip"),
	HTTPFlavorKey:                attribute.Key("http.flavor"),
	HTTPMethodKey:                attribute.Key("http.method"),
	HTTPRequestContentLengthKey:  attribute.Key("http.request_content_length"),
	HTTPResponseContentLengthKey: attribute.Key("http.response_content_length"),
	HTTPRouteKey:                 attribute.Key("http.route"),
	HTTPSchemeHTTP:               attribute.String("http.scheme", "http"),
	HTTPSchemeHTTPS:              attribute.String("http.scheme", "https"),
	HTTPStatusCodeKey:            attribute.Key("http.status_code"),
	HTTPTargetKey:                attribute.Key("http.target"),
	HTTPURLKey:                   attribute.Key("http.url"),
	HTTPUserAgentKey:             attribute.Key("http.user_agent"),
}

func TestHTTPClientResponse(t *testing.T) {
	const stat, n = 201, 397
	resp := http.Response{
		StatusCode:    stat,
		ContentLength: n,
	}
	got := hc.ClientResponse(resp)
	assert.Equal(t, 2, cap(got), "slice capacity")
	assert.ElementsMatch(t, []attribute.KeyValue{
		attribute.Key("http.status_code").Int(stat),
		attribute.Key("http.response_content_length").Int(n),
	}, got)
}

func TestHTTPClientRequest(t *testing.T) {
	const (
		user  = "alice"
		n     = 128
		agent = "Go-http-client/1.1"
	)
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:8080",
			Path:   "/resource",
		},
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"User-Agent": []string{agent},
		},
		ContentLength: n,
	}
	req.SetBasicAuth(user, "pswrd")

	assert.Equal(
		t,
		[]attribute.KeyValue{
			attribute.String("http.method", "GET"),
			attribute.String("http.flavor", "1.0"),
			attribute.String("http.url", "http://127.0.0.1:8080/resource"),
			attribute.String("net.peer.name", "127.0.0.1"),
			attribute.Int("net.peer.port", 8080),
			attribute.String("http.user_agent", agent),
			attribute.Int("http.request_content_length", n),
			attribute.String("enduser.id", user),
		},
		hc.ClientRequest(req),
	)
}

func TestHTTPClientRequestRequired(t *testing.T) {
	req := new(http.Request)
	var got []attribute.KeyValue
	assert.NotPanics(t, func() { got = hc.ClientRequest(req) })
	want := []attribute.KeyValue{
		attribute.String("http.method", "GET"),
		attribute.String("http.flavor", ""),
		attribute.String("http.url", ""),
		attribute.String("net.peer.name", ""),
	}
	assert.Equal(t, want, got)
}

func TestHTTPServerRequest(t *testing.T) {
	got := make(chan *http.Request, 1)
	handler := func(w http.ResponseWriter, r *http.Request) {
		got <- r
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	srvURL, err := url.Parse(srv.URL)
	require.NoError(t, err)
	srvPort, err := strconv.ParseInt(srvURL.Port(), 10, 32)
	require.NoError(t, err)

	resp, err := srv.Client().Get(srv.URL)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	req := <-got
	peer, peerPort := splitHostPort(req.RemoteAddr)

	const user = "alice"
	req.SetBasicAuth(user, "pswrd")

	const clientIP = "127.0.0.5"
	req.Header.Add("X-Forwarded-For", clientIP)

	assert.ElementsMatch(t,
		[]attribute.KeyValue{
			attribute.String("http.method", "GET"),
			attribute.String("http.target", "/"),
			attribute.String("http.scheme", "http"),
			attribute.String("http.flavor", "1.1"),
			attribute.String("net.host.name", srvURL.Hostname()),
			attribute.Int("net.host.port", int(srvPort)),
			attribute.String("net.sock.peer.addr", peer),
			attribute.Int("net.sock.peer.port", peerPort),
			attribute.String("http.user_agent", "Go-http-client/1.1"),
			attribute.String("enduser.id", user),
			attribute.String("http.client_ip", clientIP),
		},
		hc.ServerRequest(req))
}

func TestHTTPServerRequestFailsGracefully(t *testing.T) {
	req := new(http.Request)
	var got []attribute.KeyValue
	assert.NotPanics(t, func() { got = hc.ServerRequest(req) })
	want := []attribute.KeyValue{
		attribute.String("http.method", "GET"),
		attribute.String("http.target", ""),
		attribute.String("http.scheme", "http"),
		attribute.String("http.flavor", ""),
		attribute.String("net.host.name", ""),
	}
	assert.ElementsMatch(t, want, got)
}

func TestMethod(t *testing.T) {
	assert.Equal(t, attribute.String("http.method", "POST"), hc.method("POST"))
	assert.Equal(t, attribute.String("http.method", "GET"), hc.method(""))
	assert.Equal(t, attribute.String("http.method", "garbage"), hc.method("garbage"))
}

func TestScheme(t *testing.T) {
	assert.Equal(t, attribute.String("http.scheme", "http"), hc.scheme(false))
	assert.Equal(t, attribute.String("http.scheme", "https"), hc.scheme(true))
}

func TestProto(t *testing.T) {
	tests := map[string]string{
		"HTTP/1.0": "1.0",
		"HTTP/1.1": "1.1",
		"HTTP/2":   "2.0",
		"HTTP/3":   "3.0",
		"SPDY":     "SPDY",
		"QUIC":     "QUIC",
		"other":    "other",
	}

	for proto, want := range tests {
		assert.Equal(t, attribute.String("http.flavor", want), hc.proto(proto))
	}
}

func TestServerClientIP(t *testing.T) {
	tests := []struct {
		xForwardedFor string
		want          string
	}{
		{"", ""},
		{"127.0.0.1", "127.0.0.1"},
		{"127.0.0.1,127.0.0.5", "127.0.0.1"},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, serverClientIP(test.xForwardedFor))
	}
}

func TestRequiredHTTPPort(t *testing.T) {
	tests := []struct {
		https bool
		port  int
		want  int
	}{
		{true, 443, -1},
		{true, 80, 80},
		{true, 8081, 8081},
		{false, 443, 443},
		{false, 80, -1},
		{false, 8080, 8080},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, requiredHTTPPort(test.https, test.port))
	}
}

/*
func TestHTTPAttributesFromHTTPStatusCode(t *testing.T) {
	expected := []attribute.KeyValue{
		attribute.Int("http.status_code", 404),
	}
	got := hc.HTTPAttributesFromHTTPStatusCode(http.StatusNotFound)
	assertElementsMatch(t, expected, got, "with valid HTTP status code")
	assert.ElementsMatch(t, expected, got)
	expected = []attribute.KeyValue{
		attribute.Int("http.status_code", 499),
	}
	got = hc.HTTPAttributesFromHTTPStatusCode(499)
	assertElementsMatch(t, expected, got, "with invalid HTTP status code")
}

func TestSpanStatusFromHTTPStatusCode(t *testing.T) {
	for code := 0; code < 1000; code++ {
		expected := getExpectedCodeForHTTPCode(code, trace.SpanKindClient)
		got, msg := SpanStatusFromHTTPStatusCode(code)
		assert.Equalf(t, expected, got, "%s vs %s", expected, got)

		_, valid := validateHTTPStatusCode(code)
		if !valid {
			assert.NotEmpty(t, msg, "message should be set if error cannot be inferred from code")
		} else {
			assert.Empty(t, msg, "message should not be set if error can be inferred from code")
		}
	}
}

func TestSpanStatusFromHTTPStatusCodeAndSpanKind(t *testing.T) {
	for code := 0; code < 1000; code++ {
		expected := getExpectedCodeForHTTPCode(code, trace.SpanKindClient)
		got, msg := SpanStatusFromHTTPStatusCodeAndSpanKind(code, trace.SpanKindClient)
		assert.Equalf(t, expected, got, "%s vs %s", expected, got)

		_, valid := validateHTTPStatusCode(code)
		if !valid {
			assert.NotEmpty(t, msg, "message should be set if error cannot be inferred from code")
		} else {
			assert.Empty(t, msg, "message should not be set if error can be inferred from code")
		}
	}
	code, _ := SpanStatusFromHTTPStatusCodeAndSpanKind(400, trace.SpanKindServer)
	assert.Equalf(t, codes.Unset, code, "message should be set if error cannot be inferred from code")
}

func getExpectedCodeForHTTPCode(code int, spanKind trace.SpanKind) codes.Code {
	if http.StatusText(code) == "" {
		return codes.Error
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
		return codes.Error
	}
	category := code / 100
	if category > 0 && category < 4 {
		return codes.Unset
	}
	if spanKind == trace.SpanKindServer && category == 4 {
		return codes.Unset
	}
	return codes.Error
}

func TestHTTPClientAttributesFromHTTPRequest(t *testing.T) {
	testCases := []struct {
		name string

		method        string
		requestURI    string
		proto         string
		remoteAddr    string
		host          string
		url           *url.URL
		header        http.Header
		tls           tlsOption
		contentLength int64

		expected []attribute.KeyValue
	}{
		{
			name:       "stripped",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    noTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "http"),
				attribute.String("http.flavor", "1.0"),
			},
		},
		{
			name:       "with tls",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
			},
		},
		{
			name:       "with host",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with host fallback",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "https://example.com/user/123"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with user agent",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent": []string{"foodownloader"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.host", "example.com"),
				attribute.String("http.user_agent", "foodownloader"),
			},
		},
		{
			name:       "with http 1.1",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.1",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent": []string{"foodownloader"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.1"),
				attribute.String("http.host", "example.com"),
				attribute.String("http.user_agent", "foodownloader"),
			},
		},
		{
			name:       "with http 2",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/2.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent": []string{"foodownloader"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "2"),
				attribute.String("http.host", "example.com"),
				attribute.String("http.user_agent", "foodownloader"),
			},
		},
		{
			name:   "with content length",
			method: "GET",
			url: &url.URL{
				Path: "/user/123",
			},
			contentLength: 100,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "http"),
				attribute.Int64("http.request_content_length", 100),
			},
		},
		{
			name:   "with empty method (fallback to GET)",
			method: "",
			url: &url.URL{
				Path: "/user/123",
			},
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "http"),
			},
		},
		{
			name:   "authentication information is stripped",
			method: "",
			url: &url.URL{
				Path: "/user/123",
				User: url.UserPassword("foo", "bar"),
			},
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.url", "/user/123"),
				attribute.String("http.scheme", "http"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := testRequest(tc.method, tc.requestURI, tc.proto, tc.remoteAddr, tc.host, tc.url, tc.header, tc.tls)
			r.ContentLength = tc.contentLength
			got := hc.HTTPClientAttributesFromHTTPRequest(r)
			assert.ElementsMatch(t, tc.expected, got)
		})
	}
}

func TestHTTPServerMetricAttributesFromHTTPRequest(t *testing.T) {
	type testcase struct {
		name          string
		serverName    string
		method        string
		requestURI    string
		proto         string
		remoteAddr    string
		host          string
		url           *url.URL
		header        http.Header
		tls           tlsOption
		contentLength int64
		expected      []attribute.KeyValue
	}
	testcases := []testcase{
		{
			name:       "stripped",
			serverName: "",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    noTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "http"),
				attribute.String("http.flavor", "1.0"),
			},
		},
		{
			name:       "with server name",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    noTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "http"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
			},
		},
		{
			name:       "with tls",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
			},
		},
		{
			name:       "with route",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
			},
		},
		{
			name:       "with host",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with host fallback",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Host: "example.com",
				Path: "/user/123",
			},
			header: nil,
			tls:    withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with user agent",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent": []string{"foodownloader"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with proxy info",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent":      []string{"foodownloader"},
				"X-Forwarded-For": []string{"203.0.113.195, 70.41.3.18, 150.172.238.178"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.server_name", "my-server-name"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with http 1.1",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.1",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent":      []string{"foodownloader"},
				"X-Forwarded-For": []string{"1.2.3.4"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "1.1"),
				attribute.String("http.server_name", "my-server-name"),
				attribute.String("http.host", "example.com"),
			},
		},
		{
			name:       "with http 2",
			serverName: "my-server-name",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/2.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"User-Agent":      []string{"foodownloader"},
				"X-Forwarded-For": []string{"1.2.3.4"},
			},
			tls: withTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "https"),
				attribute.String("http.flavor", "2"),
				attribute.String("http.server_name", "my-server-name"),
				attribute.String("http.host", "example.com"),
			},
		},
	}
	for idx, tc := range testcases {
		r := testRequest(tc.method, tc.requestURI, tc.proto, tc.remoteAddr, tc.host, tc.url, tc.header, tc.tls)
		r.ContentLength = tc.contentLength
		got := hc.HTTPServerMetricAttributesFromHTTPRequest(tc.serverName, r)
		assertElementsMatch(t, tc.expected, got, "testcase %d - %s", idx, tc.name)
	}
}

func TestHttpBasicAttributesFromHTTPRequest(t *testing.T) {
	type testcase struct {
		name          string
		method        string
		requestURI    string
		proto         string
		remoteAddr    string
		host          string
		url           *url.URL
		header        http.Header
		tls           tlsOption
		contentLength int64
		expected      []attribute.KeyValue
	}
	testcases := []testcase{
		{
			name:       "stripped",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			tls:    noTLS,
			expected: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.String("http.scheme", "http"),
				attribute.String("http.flavor", "1.0"),
				attribute.String("http.host", "example.com"),
			},
		},
	}
	for idx, tc := range testcases {
		r := testRequest(tc.method, tc.requestURI, tc.proto, tc.remoteAddr, tc.host, tc.url, tc.header, tc.tls)
		r.ContentLength = tc.contentLength
		got := hc.httpBasicAttributesFromHTTPRequest(r)
		assertElementsMatch(t, tc.expected, got, "testcase %d - %s", idx, tc.name)
	}
}
*/
