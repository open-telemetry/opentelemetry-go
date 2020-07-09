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

package standard

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	otelkv "go.opentelemetry.io/otel/api/kv"
)

type tlsOption int

const (
	noTLS tlsOption = iota
	withTLS
)

func TestNetAttributesFromHTTPRequest(t *testing.T) {
	type testcase struct {
		name string

		network string

		method     string
		requestURI string
		proto      string
		remoteAddr string
		host       string
		url        *url.URL
		header     http.Header

		expected []otelkv.KeyValue
	}
	testcases := []testcase{
		{
			name:       "stripped, tcp",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
			},
		},
		{
			name:       "stripped, udp",
			network:    "udp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.UDP"),
			},
		},
		{
			name:       "stripped, ip",
			network:    "ip",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP"),
			},
		},
		{
			name:       "stripped, unix",
			network:    "unix",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "Unix"),
			},
		},
		{
			name:       "stripped, other",
			network:    "nih",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "other"),
			},
		},
		{
			name:       "with remote ip and port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
			},
		},
		{
			name:       "with remote name and port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "example.com:56",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.name", "example.com"),
				otelkv.Int("net.peer.port", 56),
			},
		},
		{
			name:       "with remote ip only",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
			},
		},
		{
			name:       "with remote name only",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "example.com",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.name", "example.com"),
			},
		},
		{
			name:       "with remote port only",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: ":56",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
			},
		},
		{
			name:       "with host name only",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "example.com",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.name", "example.com"),
			},
		},
		{
			name:       "with host ip only",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "4.3.2.1",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.ip", "4.3.2.1"),
			},
		},
		{
			name:       "with host name and port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "example.com:78",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.name", "example.com"),
				otelkv.Int("net.host.port", 78),
			},
		},
		{
			name:       "with host ip and port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "4.3.2.1:78",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.ip", "4.3.2.1"),
				otelkv.Int("net.host.port", 78),
			},
		},
		{
			name:       "with host name and bogus port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "example.com:qwerty",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.name", "example.com"),
			},
		},
		{
			name:       "with host ip and bogus port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "4.3.2.1:qwerty",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.ip", "4.3.2.1"),
			},
		},
		{
			name:       "with empty host and port",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       ":80",
			url: &url.URL{
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
			},
		},
		{
			name:       "with host ip and port in headers",
			network:    "tcp",
			method:     "GET",
			requestURI: "/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "",
			url: &url.URL{
				Path: "/user/123",
			},
			header: http.Header{
				"Host": []string{"4.3.2.1:78"},
			},
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.ip", "4.3.2.1"),
				otelkv.Int("net.host.port", 78),
			},
		},
		{
			name:       "with host ip and port in url",
			network:    "tcp",
			method:     "GET",
			requestURI: "http://4.3.2.1:78/user/123",
			proto:      "HTTP/1.0",
			remoteAddr: "1.2.3.4:56",
			host:       "",
			url: &url.URL{
				Host: "4.3.2.1:78",
				Path: "/user/123",
			},
			header: nil,
			expected: []otelkv.KeyValue{
				otelkv.String("net.transport", "IP.TCP"),
				otelkv.String("net.peer.ip", "1.2.3.4"),
				otelkv.Int("net.peer.port", 56),
				otelkv.String("net.host.ip", "4.3.2.1"),
				otelkv.Int("net.host.port", 78),
			},
		},
	}
	for idx, tc := range testcases {
		r := testRequest(tc.method, tc.requestURI, tc.proto, tc.remoteAddr, tc.host, tc.url, tc.header, noTLS)
		got := NetAttributesFromHTTPRequest(tc.network, r)
		assertElementsMatch(t, tc.expected, got, "testcase %d - %s", idx, tc.name)
	}
}

func TestEndUserAttributesFromHTTPRequest(t *testing.T) {
	r := testRequest("GET", "/user/123", "HTTP/1.1", "", "", nil, http.Header{}, withTLS)
	var expected []otelkv.KeyValue
	got := EndUserAttributesFromHTTPRequest(r)
	assert.ElementsMatch(t, expected, got)
	r.SetBasicAuth("admin", "password")
	expected = []otelkv.KeyValue{otelkv.String("enduser.id", "admin")}
	got = EndUserAttributesFromHTTPRequest(r)
	assert.ElementsMatch(t, expected, got)
}

func TestHTTPServerAttributesFromHTTPRequest(t *testing.T) {
	type testcase struct {
		name string

		serverName string
		route      string

		method        string
		requestURI    string
		proto         string
		remoteAddr    string
		host          string
		url           *url.URL
		header        http.Header
		tls           tlsOption
		contentLength int64

		expected []otelkv.KeyValue
	}
	testcases := []testcase{
		{
			name:       "stripped",
			serverName: "",
			route:      "",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "http"),
				otelkv.String("http.flavor", "1.0"),
			},
		},
		{
			name:       "with server name",
			serverName: "my-server-name",
			route:      "",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "http"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.server_name", "my-server-name"),
			},
		},
		{
			name:       "with tls",
			serverName: "my-server-name",
			route:      "",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.server_name", "my-server-name"),
			},
		},
		{
			name:       "with route",
			serverName: "my-server-name",
			route:      "/user/:id",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.server_name", "my-server-name"),
				otelkv.String("http.route", "/user/:id"),
			},
		},
		{
			name:       "with host",
			serverName: "my-server-name",
			route:      "/user/:id",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.server_name", "my-server-name"),
				otelkv.String("http.route", "/user/:id"),
				otelkv.String("http.host", "example.com"),
			},
		},
		{
			name:       "with user agent",
			serverName: "my-server-name",
			route:      "/user/:id",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.server_name", "my-server-name"),
				otelkv.String("http.route", "/user/:id"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
			},
		},
		{
			name:       "with proxy info",
			serverName: "my-server-name",
			route:      "/user/:id",
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
				"X-Forwarded-For": []string{"1.2.3.4"},
			},
			tls: withTLS,
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.server_name", "my-server-name"),
				otelkv.String("http.route", "/user/:id"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
				otelkv.String("http.client_ip", "1.2.3.4"),
			},
		},
		{
			name:       "with http 1.1",
			serverName: "my-server-name",
			route:      "/user/:id",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.1"),
				otelkv.String("http.server_name", "my-server-name"),
				otelkv.String("http.route", "/user/:id"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
				otelkv.String("http.client_ip", "1.2.3.4"),
			},
		},
		{
			name:       "with http 2",
			serverName: "my-server-name",
			route:      "/user/:id",
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "2"),
				otelkv.String("http.server_name", "my-server-name"),
				otelkv.String("http.route", "/user/:id"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
				otelkv.String("http.client_ip", "1.2.3.4"),
			},
		},
		{
			name:          "with content length",
			method:        "GET",
			requestURI:    "/user/123",
			contentLength: 100,
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.target", "/user/123"),
				otelkv.String("http.scheme", "http"),
				otelkv.Int64("http.request_content_length", 100),
			},
		},
	}
	for idx, tc := range testcases {
		r := testRequest(tc.method, tc.requestURI, tc.proto, tc.remoteAddr, tc.host, tc.url, tc.header, tc.tls)
		r.ContentLength = tc.contentLength
		got := HTTPServerAttributesFromHTTPRequest(tc.serverName, tc.route, r)
		assertElementsMatch(t, tc.expected, got, "testcase %d - %s", idx, tc.name)
	}
}

func TestHTTPAttributesFromHTTPStatusCode(t *testing.T) {
	expected := []otelkv.KeyValue{
		otelkv.Int("http.status_code", 404),
		otelkv.String("http.status_text", "Not Found"),
	}
	got := HTTPAttributesFromHTTPStatusCode(http.StatusNotFound)
	assertElementsMatch(t, expected, got, "with valid HTTP status code")
	assert.ElementsMatch(t, expected, got)
	expected = []otelkv.KeyValue{
		otelkv.Int("http.status_code", 499),
	}
	got = HTTPAttributesFromHTTPStatusCode(499)
	assertElementsMatch(t, expected, got, "with invalid HTTP status code")
}

func TestSpanStatusFromHTTPStatusCode(t *testing.T) {
	for code := 0; code < 1000; code++ {
		expected := getExpectedGRPCCodeForHTTPCode(code)
		got, _ := SpanStatusFromHTTPStatusCode(code)
		assert.Equalf(t, expected, got, "%s vs %s", expected, got)
	}
}

func getExpectedGRPCCodeForHTTPCode(code int) codes.Code {
	if http.StatusText(code) == "" {
		return codes.Unknown
	}
	switch code {
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	}
	category := code / 100
	if category < 4 {
		return codes.OK
	}
	if category < 5 {
		return codes.InvalidArgument
	}
	return codes.Internal
}

func assertElementsMatch(t *testing.T, expected, got []otelkv.KeyValue, format string, args ...interface{}) {
	if !assert.ElementsMatchf(t, expected, got, format, args...) {
		t.Log("expected:", kvStr(expected))
		t.Log("got:", kvStr(got))
	}
}

func testRequest(method, requestURI, proto, remoteAddr, host string, u *url.URL, header http.Header, tlsopt tlsOption) *http.Request {
	major, minor := protoToInts(proto)
	var tlsConn *tls.ConnectionState
	switch tlsopt {
	case noTLS:
	case withTLS:
		tlsConn = &tls.ConnectionState{}
	}
	return &http.Request{
		Method:     method,
		URL:        u,
		Proto:      proto,
		ProtoMajor: major,
		ProtoMinor: minor,
		Header:     header,
		Host:       host,
		RemoteAddr: remoteAddr,
		RequestURI: requestURI,
		TLS:        tlsConn,
	}
}

func protoToInts(proto string) (int, int) {
	switch proto {
	case "HTTP/1.0":
		return 1, 0
	case "HTTP/1.1":
		return 1, 1
	case "HTTP/2.0":
		return 2, 0
	}
	// invalid proto
	return 13, 42
}

func kvStr(kvs []otelkv.KeyValue) string {
	sb := strings.Builder{}
	sb.WriteRune('[')
	for idx, kv := range kvs {
		if idx > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString((string)(kv.Key))
		sb.WriteString(": ")
		sb.WriteString(kv.Value.Emit())
	}
	sb.WriteRune(']')
	return sb.String()
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

		expected []otelkv.KeyValue
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "http"),
				otelkv.String("http.flavor", "1.0"),
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.host", "example.com"),
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.0"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "1.1"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
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
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "https"),
				otelkv.String("http.flavor", "2"),
				otelkv.String("http.host", "example.com"),
				otelkv.String("http.user_agent", "foodownloader"),
			},
		},
		{
			name:   "with content length",
			method: "GET",
			url: &url.URL{
				Path: "/user/123",
			},
			contentLength: 100,
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "http"),
				otelkv.Int64("http.request_content_length", 100),
			},
		},
		{
			name:   "with empty method (fallback to GET)",
			method: "",
			url: &url.URL{
				Path: "/user/123",
			},
			expected: []otelkv.KeyValue{
				otelkv.String("http.method", "GET"),
				otelkv.String("http.url", "/user/123"),
				otelkv.String("http.scheme", "http"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := testRequest(tc.method, tc.requestURI, tc.proto, tc.remoteAddr, tc.host, tc.url, tc.header, tc.tls)
			r.ContentLength = tc.contentLength
			got := HTTPClientAttributesFromHTTPRequest(r)
			assert.ElementsMatch(t, tc.expected, got)
		})
	}
}
