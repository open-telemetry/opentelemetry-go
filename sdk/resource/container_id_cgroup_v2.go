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
	"regexp"
	"strings"
)

const cgroupV2Path = "/proc/self/mountinfo"

var cgroupV2ContainerIDRe = regexp.MustCompile(`^.*/.+/([\w+-.]{64})/.*$`)

func getContainerIDFromCGroupV2() (string, error) {
	return getContainerIDFromCGroupFile(cgroupV2Path, getContainerIDFromCgroupV2Line)
}

// getContainerIDFromCgroupV2Line returns the ID of the container from one string line.
func getContainerIDFromCgroupV2Line(line string) string {
	// Only match line contains "cpuset"
	if !strings.Contains(line, "hostname") {
		return ""
	}

	matches := cgroupV2ContainerIDRe.FindStringSubmatch(line)
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}
