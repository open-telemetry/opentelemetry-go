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

package resource_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func mockRuntimeProviders() {
	resource.SetRuntimeProviders(
		fakeRuntimeNameProvider,
		fakeRuntimeVersionProvider,
		func() string { return "LINUX" },
		fakeRuntimeArchProvider,
	)

	resource.SetOSDescriptionProvider(
		func() (string, error) { return "Test", nil },
	)
}

func TestMapRuntimeOSToSemconvOSType(t *testing.T) {
	tt := []struct {
		Name   string
		Goos   string
		OSType attribute.KeyValue
	}{
		{"Apple Darwin", "darwin", semconv.OSTypeDarwin},
		{"DragonFly BSD", "dragonfly", semconv.OSTypeDragonflyBSD},
		{"FreeBSD", "freebsd", semconv.OSTypeFreeBSD},
		{"Linux", "linux", semconv.OSTypeLinux},
		{"NetBSD", "netbsd", semconv.OSTypeNetBSD},
		{"OpenBSD", "openbsd", semconv.OSTypeOpenBSD},
		{"Oracle Solaris", "solaris", semconv.OSTypeSolaris},
		{"Microsoft Windows", "windows", semconv.OSTypeWindows},
		{"Unknown", "unknown", semconv.OSTypeKey.String("unknown")},
		{"UNKNOWN", "UNKNOWN", semconv.OSTypeKey.String("unknown")},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			osTypeAttribute := resource.MapRuntimeOSToSemconvOSType(tc.Goos)
			require.EqualValues(t, osTypeAttribute, tc.OSType)
		})
	}
}
