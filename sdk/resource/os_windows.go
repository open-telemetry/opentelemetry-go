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

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// platformOSDescription returns a human readable OS version information string.
// It does so by querying registry values under the
// `SOFTWARE\Microsoft\Windows NT\CurrentVersion` key. The final string
// resembles the one displayed by the Version Reporter Applet (winver.exe).
func platformOSDescription() (string, error) {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)

	if err != nil {
		return "", err
	}

	defer k.Close()

	productName, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return "", err
	}

	displayVersion, _, err := k.GetStringValue("DisplayVersion")
	if err != nil {
		return "", err
	}

	releaseID, _, err := k.GetStringValue("ReleaseID")
	if err != nil {
		return "", err
	}

	currentMajorVersionNumber, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return "", err
	}

	currentMinorVersionNumber, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return "", err
	}

	currentBuildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		return "", err
	}

	ubr, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s (%s) [Version %d.%d.%s.%d]",
		productName,
		displayVersion,
		releaseID,
		currentMajorVersionNumber,
		currentMinorVersionNumber,
		currentBuildNumber,
		ubr,
	), nil
}
