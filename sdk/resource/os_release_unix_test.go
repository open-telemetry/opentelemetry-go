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

//go:build aix || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix dragonfly freebsd linux netbsd openbsd solaris zos

package resource_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/resource"
)

func TestParseOSReleaseFile(t *testing.T) {
	osReleaseUbuntu := bytes.NewBufferString(`NAME="Ubuntu"
VERSION="20.04.2 LTS (Focal Fossa)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.2 LTS"
VERSION_ID="20.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=focal
UBUNTU_CODENAME=focal`)

	parsedUbuntu := map[string]string{
		"NAME":               "Ubuntu",
		"VERSION":            "20.04.2 LTS (Focal Fossa)",
		"ID":                 "ubuntu",
		"ID_LIKE":            "debian",
		"PRETTY_NAME":        "Ubuntu 20.04.2 LTS",
		"VERSION_ID":         "20.04",
		"HOME_URL":           "https://www.ubuntu.com/",
		"SUPPORT_URL":        "https://help.ubuntu.com/",
		"BUG_REPORT_URL":     "https://bugs.launchpad.net/ubuntu/",
		"PRIVACY_POLICY_URL": "https://www.ubuntu.com/legal/terms-and-policies/privacy-policy",
		"VERSION_CODENAME":   "focal",
		"UBUNTU_CODENAME":    "focal",
	}

	osReleaseDebian := bytes.NewBufferString(`PRETTY_NAME="Debian GNU/Linux 10 (buster)"
NAME="Debian GNU/Linux"
VERSION_ID="10"
VERSION="10 (buster)"
VERSION_CODENAME=buster
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"`)

	parsedDebian := map[string]string{
		"PRETTY_NAME":      "Debian GNU/Linux 10 (buster)",
		"NAME":             "Debian GNU/Linux",
		"VERSION_ID":       "10",
		"VERSION":          "10 (buster)",
		"VERSION_CODENAME": "buster",
		"ID":               "debian",
		"HOME_URL":         "https://www.debian.org/",
		"SUPPORT_URL":      "https://www.debian.org/support",
		"BUG_REPORT_URL":   "https://bugs.debian.org/",
	}

	osReleaseAlpine := bytes.NewBufferString(`NAME="Alpine Linux"
ID=alpine
VERSION_ID=3.13.4
PRETTY_NAME="Alpine Linux v3.13"
HOME_URL="https://alpinelinux.org/"
BUG_REPORT_URL="https://bugs.alpinelinux.org/"`)

	parsedAlpine := map[string]string{
		"NAME":           "Alpine Linux",
		"ID":             "alpine",
		"VERSION_ID":     "3.13.4",
		"PRETTY_NAME":    "Alpine Linux v3.13",
		"HOME_URL":       "https://alpinelinux.org/",
		"BUG_REPORT_URL": "https://bugs.alpinelinux.org/",
	}

	osReleaseMock := bytes.NewBufferString(`
# This line should be skipped

QUOTED1="Quoted value 1"
QUOTED2='Quoted value 2'
ESCAPED1="\$HOME"
ESCAPED2="\"release\""
ESCAPED3="rock\'n\'roll"
ESCAPED4="\\var"

=line with missing key should be skipped

PROP1=name=john
	PROP2  =  Value  
PROP3='This value will be overwritten by the next one'
PROP3='Final value'`)

	parsedMock := map[string]string{
		"QUOTED1":  "Quoted value 1",
		"QUOTED2":  "Quoted value 2",
		"ESCAPED1": "$HOME",
		"ESCAPED2": `"release"`,
		"ESCAPED3": "rock'n'roll",
		"ESCAPED4": `\var`,
		"PROP1":    "name=john",
		"PROP2":    "Value",
		"PROP3":    "Final value",
	}

	tt := []struct {
		Name      string
		OSRelease io.Reader
		Parsed    map[string]string
	}{
		{"Ubuntu", osReleaseUbuntu, parsedUbuntu},
		{"Debian", osReleaseDebian, parsedDebian},
		{"Alpine", osReleaseAlpine, parsedAlpine},
		{"Mock", osReleaseMock, parsedMock},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.ParseOSReleaseFile(tc.OSRelease)
			require.EqualValues(t, tc.Parsed, result)
		})
	}
}

func TestSkip(t *testing.T) {
	tt := []struct {
		Name     string
		Line     string
		Expected bool
	}{
		{"Empty string", "", true},
		{"Only whitespace", "   ", true},
		{"Hashtag prefix 1", "# Sample text", true},
		{"Hashtag prefix 2", "  # Sample text", true},
		{"Hashtag and whitespace 1", "#  ", true},
		{"Hashtag and whitespace 2", "  #", true},
		{"Hashtag and whitespace 3", "  #  ", true},
		{"Nonempty string", "Sample text", false},
		{"Nonempty string with whitespace around", " Sample text ", false},
		{"Nonempty string with middle hashtag", "Sample #text", false},
		{"Nonempty string with ending hashtag", "Sample text #", false},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.Skip(tc.Line)
			require.EqualValues(t, tc.Expected, result)
		})
	}
}

func TestParse(t *testing.T) {
	tt := []struct {
		Name          string
		Line          string
		ExpectedKey   string
		ExpectedValue string
		OK            bool
	}{
		{"Empty string", "", "", "", false},
		{"No separator", "wrong", "", "", false},
		{"Empty key", "=john", "", "", false},
		{"Empty key value", "=", "", "", false},
		{"Empty value", "name=", "name", "", true},
		{"Key value 1", "name=john", "name", "john", true},
		{"Key value 2", "name=john=dev", "name", "john=dev", true},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			key, value, ok := resource.Parse(tc.Line)
			require.EqualValues(t, tc.ExpectedKey, key)
			require.EqualValues(t, tc.ExpectedValue, value)
			require.EqualValues(t, tc.OK, ok)
		})
	}
}

func TestUnquote(t *testing.T) {
	tt := []struct {
		Name     string
		Text     string
		Expected string
	}{
		{"Empty string", ``, ``},
		{"Single double quote", `"`, `"`},
		{"Single single quote", `'`, `'`},
		{"Empty double quotes", `""`, ``},
		{"Empty single quotes", `''`, ``},
		{"Empty mixed quotes 1", `"'`, `"'`},
		{"Empty mixed quotes 2", `'"`, `'"`},
		{"Double quotes", `"Sample text"`, `Sample text`},
		{"Single quotes", `'Sample text'`, `Sample text`},
		{"Half-open starting double quote", `"Sample text`, `"Sample text`},
		{"Half-open ending double quote", `Sample text"`, `Sample text"`},
		{"Half-open starting single quote", `'Sample text`, `'Sample text`},
		{"Half-open ending single quote", `Sample text'`, `Sample text'`},
		{"Double double quotes", `""Sample text""`, `"Sample text"`},
		{"Double single quotes", `''Sample text''`, `'Sample text'`},
		{"Mismatch quotes 1", `"Sample text'`, `"Sample text'`},
		{"Mismatch quotes 2", `'Sample text"`, `'Sample text"`},
		{"No quotes", `Sample text`, `Sample text`},
		{"Internal double quote", `Sample "text"`, `Sample "text"`},
		{"Internal single quote", `Sample 'text'`, `Sample 'text'`},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.Unquote(tc.Text)
			require.EqualValues(t, tc.Expected, result)
		})
	}
}

func TestUnescape(t *testing.T) {
	tt := []struct {
		Name     string
		Text     string
		Expected string
	}{
		{"Empty string", ``, ``},
		{"Escaped dollar sign", `\$var`, `$var`},
		{"Escaped double quote", `\"var`, `"var`},
		{"Escaped single quote", `\'var`, `'var`},
		{"Escaped backslash", `\\var`, `\var`},
		{"Escaped backtick", "\\`var", "`var"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.Unescape(tc.Text)
			require.EqualValues(t, tc.Expected, result)
		})
	}
}

func TestBuildOSRelease(t *testing.T) {
	tt := []struct {
		Name     string
		Values   map[string]string
		Expected string
	}{
		{"Nil values", nil, ""},
		{"Empty values", map[string]string{}, ""},
		{"Name and version only", map[string]string{
			"NAME":    "Ubuntu",
			"VERSION": "20.04.2 LTS (Focal Fossa)",
		}, "Ubuntu 20.04.2 LTS (Focal Fossa)"},
		{"Name and version preferred", map[string]string{
			"NAME":        "Ubuntu",
			"VERSION":     "20.04.2 LTS (Focal Fossa)",
			"VERSION_ID":  "20.04",
			"PRETTY_NAME": "Ubuntu 20.04.2 LTS",
		}, "Ubuntu 20.04.2 LTS (Focal Fossa)"},
		{"Version ID fallback", map[string]string{
			"NAME":       "Ubuntu",
			"VERSION_ID": "20.04",
		}, "Ubuntu 20.04"},
		{"Pretty name fallback due to missing name", map[string]string{
			"VERSION":     "20.04.2 LTS (Focal Fossa)",
			"PRETTY_NAME": "Ubuntu 20.04.2 LTS",
		}, "Ubuntu 20.04.2 LTS"},
		{"Pretty name fallback due to missing version", map[string]string{
			"NAME":        "Ubuntu",
			"PRETTY_NAME": "Ubuntu 20.04.2 LTS",
		}, "Ubuntu 20.04.2 LTS"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.BuildOSRelease(tc.Values)
			require.EqualValues(t, tc.Expected, result)
		})
	}
}
