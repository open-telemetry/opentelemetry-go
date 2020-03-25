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

package filters

import (
	"net/http"
	"net/url"
	"testing"

	"go.opentelemetry.io/otel/plugin/othttp"
)

type scenario struct {
	name   string
	filter othttp.Filter
	req    *http.Request
	exp    bool
}

func TestAny(t *testing.T) {
	for _, s := range []scenario{
		{
			name:   "no matching filters",
			filter: Any(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/boo", Host: "baz.bar:8080"}},
			exp:    false,
		},
		{
			name:   "one matching filter",
			filter: Any(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/foo", Host: "baz.bar:8080"}},
			exp:    true,
		},
		{
			name:   "all matching filters",
			filter: Any(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/foo", Host: "bar.baz:8080"}},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestAll(t *testing.T) {
	for _, s := range []scenario{
		{
			name:   "no matching filters",
			filter: All(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/boo", Host: "baz.bar:8080"}},
			exp:    false,
		},
		{
			name:   "one matching filter",
			filter: All(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/foo", Host: "baz.bar:8080"}},
			exp:    false,
		},
		{
			name:   "all matching filters",
			filter: All(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/foo", Host: "bar.baz:8080"}},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestNone(t *testing.T) {
	for _, s := range []scenario{
		{
			name:   "no matching filters",
			filter: None(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/boo", Host: "baz.bar:8080"}},
			exp:    true,
		},
		{
			name:   "one matching filter",
			filter: None(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/foo", Host: "baz.bar:8080"}},
			exp:    false,
		},
		{
			name:   "all matching filters",
			filter: None(Path("/foo"), Hostname("bar.baz")),
			req:    &http.Request{URL: &url.URL{Path: "/foo", Host: "bar.baz:8080"}},
			exp:    false,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestNot(t *testing.T) {
	req := &http.Request{URL: &url.URL{Path: "/foo", Host: "bar.baz:8080"}}
	filter := Path("/foo")
	if filter(req) == Not(filter)(req) {
		t.Error("Not filter should invert the result of the supplied filter")
	}
}

func TestPathPrefix(t *testing.T) {
	for _, s := range []scenario{
		{
			name:   "non-matching prefix",
			filter: PathPrefix("/foo"),
			req:    &http.Request{URL: &url.URL{Path: "/boo/far", Host: "baz.bar:8080"}},
			exp:    false,
		},
		{
			name:   "matching prefix",
			filter: PathPrefix("/foo"),
			req:    &http.Request{URL: &url.URL{Path: "/foo/bar", Host: "bar.baz:8080"}},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestMethod(t *testing.T) {
	for _, s := range []scenario{
		{
			name:   "non-matching method",
			filter: Method(http.MethodGet),
			req:    &http.Request{Method: http.MethodHead, URL: &url.URL{Path: "/boo/far", Host: "baz.bar:8080"}},
			exp:    false,
		},
		{
			name:   "matching method",
			filter: Method(http.MethodGet),
			req:    &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/boo/far", Host: "baz.bar:8080"}},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestQuery(t *testing.T) {
	matching, _ := url.Parse("http://bar.baz:8080/foo/bar?key=value")
	nonMatching, _ := url.Parse("http://bar.baz:8080/foo/bar?key=other")
	for _, s := range []scenario{
		{
			name:   "non-matching query parameter",
			filter: Query("key", "value"),
			req:    &http.Request{Method: http.MethodHead, URL: nonMatching},
			exp:    false,
		},
		{
			name:   "matching query parameter",
			filter: Query("key", "value"),
			req:    &http.Request{Method: http.MethodGet, URL: matching},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestQueryContains(t *testing.T) {
	matching, _ := url.Parse("http://bar.baz:8080/foo/bar?key=value")
	nonMatching, _ := url.Parse("http://bar.baz:8080/foo/bar?key=other")
	for _, s := range []scenario{
		{
			name:   "non-matching query parameter",
			filter: QueryContains("key", "alu"),
			req:    &http.Request{Method: http.MethodHead, URL: nonMatching},
			exp:    false,
		},
		{
			name:   "matching query parameter",
			filter: QueryContains("key", "alu"),
			req:    &http.Request{Method: http.MethodGet, URL: matching},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestHeader(t *testing.T) {
	matching := http.Header{}
	matching.Add("key", "value")
	nonMatching := http.Header{}
	nonMatching.Add("key", "other")
	for _, s := range []scenario{
		{
			name:   "non-matching query parameter",
			filter: Header("key", "value"),
			req:    &http.Request{Method: http.MethodHead, Header: nonMatching},
			exp:    false,
		},
		{
			name:   "matching query parameter",
			filter: Header("key", "value"),
			req:    &http.Request{Method: http.MethodGet, Header: matching},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}

func TestHeaderContains(t *testing.T) {
	matching := http.Header{}
	matching.Add("key", "value")
	nonMatching := http.Header{}
	nonMatching.Add("key", "other")
	for _, s := range []scenario{
		{
			name:   "non-matching query parameter",
			filter: HeaderContains("key", "alu"),
			req:    &http.Request{Method: http.MethodHead, Header: nonMatching},
			exp:    false,
		},
		{
			name:   "matching query parameter",
			filter: HeaderContains("key", "alu"),
			req:    &http.Request{Method: http.MethodGet, Header: matching},
			exp:    true,
		},
	} {
		res := s.filter(s.req)
		if s.exp != res {
			t.Errorf("Failed testing %q. Expected %t, got %t", s.name, s.exp, res)
		}
	}
}
