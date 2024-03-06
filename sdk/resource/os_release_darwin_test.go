// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/resource"
)

func TestParsePlistFile(t *testing.T) {
	standardPlist := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>ProductBuildVersion</key>
	<string>20E232</string>
	<key>ProductCopyright</key>
	<string>1983-2021 Apple Inc.</string>
	<key>ProductName</key>
	<string>macOS</string>
	<key>ProductUserVisibleVersion</key>
	<string>11.3</string>
	<key>ProductVersion</key>
	<string>11.3</string>
	<key>iOSSupportVersion</key>
	<string>14.5</string>
</dict>
</plist>`)

	parsedPlist := map[string]string{
		"ProductBuildVersion":       "20E232",
		"ProductCopyright":          "1983-2021 Apple Inc.",
		"ProductName":               "macOS",
		"ProductUserVisibleVersion": "11.3",
		"ProductVersion":            "11.3",
		"iOSSupportVersion":         "14.5",
	}

	emptyPlist := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
</dict>
</plist>`)

	missingDictPlist := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
</plist>`)

	unknownElementsPlist := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<a>
		<b>123</b>
	</a>
	<key>ProductBuildVersion</key>
	<c>Value</c>
	<string>20E232</string>
	<d attr="1"></d>
</dict>
</plist>`)

	parsedUnknownElementsPlist := map[string]string{
		"ProductBuildVersion": "20E232",
	}

	tt := []struct {
		Name   string
		Plist  io.Reader
		Parsed map[string]string
	}{
		{"Standard", standardPlist, parsedPlist},
		{"Empty", emptyPlist, map[string]string{}},
		{"Missing dict", missingDictPlist, map[string]string{}},
		{"Unknown elements", unknownElementsPlist, parsedUnknownElementsPlist},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result, err := resource.ParsePlistFile(tc.Plist)

			require.Equal(t, tc.Parsed, result)
			require.NoError(t, err)
		})
	}
}

func TestParsePlistFileUnevenKeys(t *testing.T) {
	plist := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>ProductBuildVersion</key>
	<string>20E232</string>
	<key>ProductCopyright</key>
</dict>
</plist>`)

	result, err := resource.ParsePlistFile(plist)

	require.Nil(t, result)
	require.Error(t, err)
}

func TestParsePlistFileMalformed(t *testing.T) {
	plist := bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Product
</dict>
</plist>`)

	result, err := resource.ParsePlistFile(plist)

	require.Nil(t, result)
	require.Error(t, err)
}

func TestBuildOSRelease(t *testing.T) {
	tt := []struct {
		Name       string
		Properties map[string]string
		OSRelease  string
	}{
		{"Empty properties", map[string]string{}, ""},
		{"Empty properties (nil)", nil, ""},
		{"Missing product name", map[string]string{
			"ProductVersion":      "11.3",
			"ProductBuildVersion": "20E232",
		}, ""},
		{"Missing product version", map[string]string{
			"ProductName":         "macOS",
			"ProductBuildVersion": "20E232",
		}, ""},
		{"Missing product build version", map[string]string{
			"ProductName":    "macOS",
			"ProductVersion": "11.3",
		}, ""},
		{"All properties available", map[string]string{
			"ProductName":         "macOS",
			"ProductVersion":      "11.3",
			"ProductBuildVersion": "20E232",
		}, "macOS 11.3 (20E232)"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.BuildOSRelease(tc.Properties)
			require.Equal(t, tc.OSRelease, result)
		})
	}
}
