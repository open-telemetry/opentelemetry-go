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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainerIDFromCgroupV1Line(t *testing.T) {
	testCases := []struct {
		name                string
		line                string
		expectedContainerID string
	}{
		{
			name:                "with suffix",
			line:                "13:name=systemd:/podruntime/docker/kubepods/ac679f8a8319c8cf7d38e1adf263bc08d23.aaaa",
			expectedContainerID: "ac679f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "with prefix and suffix",
			line:                "13:name=systemd:/podruntime/docker/kubepods/crio-dc679f8a8319c8cf7d38e1adf263bc08d23.stuff",
			expectedContainerID: "dc679f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "no prefix and suffix",
			line:                "13:name=systemd:/pod/d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
			expectedContainerID: "d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
		},
		{
			name:                "with space",
			line:                " 13:name=systemd:/pod/d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356 ",
			expectedContainerID: "d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
		},
		{
			name: "invalid hex string",
			line: "13:name=systemd:/podruntime/docker/kubepods/ac679f8a8319c8cf7d38e1adf263bc08d23zzzz",
		},
		{
			name: "no container id - 1",
			line: "pids: /",
		},
		{
			name: "no container id - 2",
			line: "pids: ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromCgroupV1Line(tc.line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
