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

const cgroupV1Path = "/proc/self/cgroup"

var cgroupV1ContainerIDRe = regexp.MustCompile(`^.*/(?:.*-)?([\w+-]+)(?:\.|\s*$)`)

func getContainerIDFromCGroupV1() (string, error) {
	return getContainerIDFromCGroupFile(cgroupV1Path, getContainerIDFromCgroupV1Line)
}

// getContainerIDFromCgroupV1Line returns the ID of the container from one string line.
func getContainerIDFromCgroupV1Line(line string) string {
	// Only match line contains "cpuset"
	if !strings.Contains(line, "cpuset") {
		return ""
	}

	matches := cgroupV1ContainerIDRe.FindStringSubmatch(line)
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}
