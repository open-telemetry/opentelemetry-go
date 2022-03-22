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

package envconfig // import "go.opentelemetry.io/otel/exporters/otlp/internal/envconfig"

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const WeakCertificate = `
-----BEGIN CERTIFICATE-----
MIIBhzCCASygAwIBAgIRANHpHgAWeTnLZpTSxCKs0ggwCgYIKoZIzj0EAwIwEjEQ
MA4GA1UEChMHb3RlbC1nbzAeFw0yMTA0MDExMzU5MDNaFw0yMTA0MDExNDU5MDNa
MBIxEDAOBgNVBAoTB290ZWwtZ28wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAS9
nWSkmPCxShxnp43F+PrOtbGV7sNfkbQ/kxzi9Ego0ZJdiXxkmv/C05QFddCW7Y0Z
sJCLHGogQsYnWJBXUZOVo2MwYTAOBgNVHQ8BAf8EBAMCB4AwEwYDVR0lBAwwCgYI
KwYBBQUHAwEwDAYDVR0TAQH/BAIwADAsBgNVHREEJTAjgglsb2NhbGhvc3SHEAAA
AAAAAAAAAAAAAAAAAAGHBH8AAAEwCgYIKoZIzj0EAwIDSQAwRgIhANwZVVKvfvQ/
1HXsTvgH+xTQswOwSSKYJ1cVHQhqK7ZbAiEAus8NxpTRnp5DiTMuyVmhVNPB+bVH
Lhnm4N/QDk5rek0=
-----END CERTIFICATE-----
`

type testOption struct {
	TestString   string
	TestDuration time.Duration
	TestHeaders  map[string]string
	TestURL      *url.URL
	TestTLS      *tls.Config
}

func TestEnvConfig(t *testing.T) {
	parsedURL, err := url.Parse("https://example.com")
	assert.NoError(t, err)

	options := []testOption{}
	for _, testcase := range []struct {
		name            string
		reader          EnvOptionsReader
		configs         []ConfigFn
		expectedOptions []testOption
	}{
		{
			name: "with no namespace and a matching key",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "world"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithString("HELLO", func(v string) {
					options = append(options, testOption{TestString: v})
				}),
			},
			expectedOptions: []testOption{
				{
					TestString: "world",
				},
			},
		},
		{
			name: "with no namespace and a non-matching key",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "world"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithString("HOLA", func(v string) {
					options = append(options, testOption{TestString: v})
				}),
			},
			expectedOptions: []testOption{},
		},
		{
			name: "with a namespace and a matching key",
			reader: EnvOptionsReader{
				Namespace: "MY_NAMESPACE",
				GetEnv: func(n string) string {
					if n == "MY_NAMESPACE_HELLO" {
						return "world"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithString("HELLO", func(v string) {
					options = append(options, testOption{TestString: v})
				}),
			},
			expectedOptions: []testOption{
				{
					TestString: "world",
				},
			},
		},
		{
			name: "with no namespace and a non-matching key",
			reader: EnvOptionsReader{
				Namespace: "MY_NAMESPACE",
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "world"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithString("HELLO", func(v string) {
					options = append(options, testOption{TestString: v})
				}),
			},
			expectedOptions: []testOption{},
		},
		{
			name: "with a duration config",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "60"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithDuration("HELLO", func(v time.Duration) {
					options = append(options, testOption{TestDuration: v})
				}),
			},
			expectedOptions: []testOption{
				{
					TestDuration: 60_000_000, // 60 milliseconds
				},
			},
		},
		{
			name: "with an invalid duration config",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "world"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithDuration("HELLO", func(v time.Duration) {
					options = append(options, testOption{TestDuration: v})
				}),
			},
			expectedOptions: []testOption{},
		},
		{
			name: "with headers",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "userId=42,userName=alice"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithHeaders("HELLO", func(v map[string]string) {
					options = append(options, testOption{TestHeaders: v})
				}),
			},
			expectedOptions: []testOption{
				{
					TestHeaders: map[string]string{
						"userId":   "42",
						"userName": "alice",
					},
				},
			},
		},
		{
			name: "with invalid headers",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "world"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithHeaders("HELLO", func(v map[string]string) {
					options = append(options, testOption{TestHeaders: v})
				}),
			},
			expectedOptions: []testOption{
				{
					TestHeaders: map[string]string{},
				},
			},
		},
		{
			name: "with URL",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "https://example.com"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithURL("HELLO", func(v *url.URL) {
					options = append(options, testOption{TestURL: v})
				}),
			},
			expectedOptions: []testOption{
				{
					TestURL: parsedURL,
				},
			},
		},
		{
			name: "with invalid URL",
			reader: EnvOptionsReader{
				GetEnv: func(n string) string {
					if n == "HELLO" {
						return "i nvalid://url"
					}
					return ""
				},
			},
			configs: []ConfigFn{
				WithURL("HELLO", func(v *url.URL) {
					options = append(options, testOption{TestURL: v})
				}),
			},
			expectedOptions: []testOption{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.reader.Apply(testcase.configs...)
			assert.Equal(t, testcase.expectedOptions, options)
			options = []testOption{}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	tlsCert, err := createTLSConfig([]byte(WeakCertificate))
	assert.NoError(t, err)

	reader := EnvOptionsReader{
		GetEnv: func(n string) string {
			if n == "CERTIFICATE" {
				return "/path/cert.pem"
			}
			return ""
		},
		ReadFile: func(p string) ([]byte, error) {
			if p == "/path/cert.pem" {
				return []byte(WeakCertificate), nil
			}
			return []byte{}, nil
		},
	}

	var option testOption
	reader.Apply(
		WithTLSConfig("CERTIFICATE", func(v *tls.Config) {
			option = testOption{TestTLS: v}
		}))

	// nolint:staticcheck // ignoring tlsCert.RootCAs.Subjects is deprecated ERR because cert does not come from SystemCertPool.
	assert.Equal(t, tlsCert.RootCAs.Subjects(), option.TestTLS.RootCAs.Subjects())
}

func TestStringToHeader(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  map[string]string
	}{
		{
			name:  "simple test",
			value: "userId=alice",
			want:  map[string]string{"userId": "alice"},
		},
		{
			name:  "simple test with spaces",
			value: " userId = alice  ",
			want:  map[string]string{"userId": "alice"},
		},
		{
			name:  "multiples headers encoded",
			value: "userId=alice,serverNode=DF%3A28,isProduction=false",
			want: map[string]string{
				"userId":       "alice",
				"serverNode":   "DF:28",
				"isProduction": "false",
			},
		},
		{
			name:  "invalid headers format",
			value: "userId:alice",
			want:  map[string]string{},
		},
		{
			name:  "invalid key",
			value: "%XX=missing,userId=alice",
			want: map[string]string{
				"userId": "alice",
			},
		},
		{
			name:  "invalid value",
			value: "missing=%XX,userId=alice",
			want: map[string]string{
				"userId": "alice",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, stringToHeader(tt.value))
		})
	}
}
