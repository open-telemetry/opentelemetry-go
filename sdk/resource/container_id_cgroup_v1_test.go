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
			line:                "13:cpuset:/podruntime/docker/kubepods/ac679f8a8319c8cf7d38e1adf263bc08d23.aaaa",
			expectedContainerID: "ac679f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "with prefix and suffix",
			line:                "13:cpuset:/podruntime/docker/kubepods/crio-dc679f8a8319c8cf7d38e1adf263bc08d23.stuff",
			expectedContainerID: "dc679f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "no prefix and suffix",
			line:                "13:cpuset:/pod/d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
			expectedContainerID: "d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356",
		},
		{
			name:                "with space",
			line:                " 13:cpuset:/pod/d86d75589bf6cc254f3e2cc29debdf85dde404998aa128997a819ff991827356 ",
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
		{
			name:                "minikube containerd cgroup",
			line:                "11:cpuset:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236",
			expectedContainerID: "58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236",
		},
		{
			name:                "minikube docker cgroup",
			line:                "5:cpuset:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope",
			expectedContainerID: "3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b",
		},
		{
			name: "podman cgroup",
			line: "14:name=systemd:/user.slice/user-1000.slice/user@1000.service/app.slice/podman.service",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromCgroupV1Line(tc.line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
