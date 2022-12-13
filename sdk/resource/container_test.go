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

package resource

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setDefaultContainerProviders() {
	setContainerProviders(
		getContainerIDFromCGroup,
	)
}

func setContainerProviders(
	idProvider containerIDProvider,
) {
	containerID = idProvider
}

func TestGetContainerIDFromReader(t *testing.T) {
	testCases := []struct {
		name                string
		reader              io.Reader
		expectedContainerID string
	}{
		{
			name: "multiple lines",
			reader: strings.NewReader(`//
1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23
1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d24
`),
			expectedContainerID: "dc579f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name: "no container id",
			reader: strings.NewReader(`//
1:name=systemd:/podruntime/docker
`),
			expectedContainerID: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromReader(tc.reader, getContainerIDFromCgroupV1Line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}

func TestGetContainerIDFromCGroup(t *testing.T) {
	t.Cleanup(func() {
		osStat = defaultOSStat
		osOpen = defaultOSOpen
	})

	testCases := []struct {
		name                 string
		cgroupV1FileNotExist bool
		cgroupV2FileNotExist bool
		openFileError        error
		cgroupV1FileContent  string
		cgroupV2FileContent  string
		expectedContainerID  string
		expectedError        bool
	}{
		{
			name:                 "the cgroup file does not exist",
			cgroupV1FileNotExist: true,
			cgroupV2FileNotExist: true,
		},
		{
			name:          "error when opening cgroup file",
			openFileError: errors.New("test"),
			expectedError: true,
		},
		{
			name:                "cgroup file v1",
			cgroupV1FileContent: "1:cpuset:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23",
			expectedContainerID: "dc579f8a8319c8cf7d38e1adf263bc08d23",
		},
		{
			name:                "cgroup file v2",
			cgroupV2FileContent: "474 456 254:1 /docker/containers/dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183/hosts /etc/hosts rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183",
		},
		{
			name:                "both way fail",
			cgroupV1FileContent: " ",
			cgroupV2FileContent: " ",
		},
		{
			name: "minikube containerd cgroup",
			cgroupV1FileContent: `11:cpuset:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
10:hugetlb:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
9:pids:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
8:memory:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
7:net_cls,net_prio:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
6:perf_event:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
5:blkio:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
4:devices:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
3:freezer:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
2:cpu,cpuacct:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236
1:name=systemd:/kubepods/besteffort/pod28478e30-384f-41e5-9d85-eae249ae8506/58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236`,
			expectedContainerID: "58a77afcbf0b16959d526758f6696677c862517acc97a562dc5c5b09afbf5236",
		},
		{
			name: "minikube docker cgroup",
			cgroupV1FileContent: `11:blkio:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
10:perf_event:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
9:pids:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
8:memory:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
7:cpu,cpuacct:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
6:net_cls,net_prio:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
5:cpuset:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
4:devices:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
3:hugetlb:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
2:freezer:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope
1:name=systemd:/kubepods.slice/kubepods-besteffort.slice/kubepods-besteffort-pod350bff31_89d4_429e_b653_86d8167bc60e.slice/docker-3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b.scope`,
			expectedContainerID: "3a5881f09ab409d7ac174b59f20d003b28b76da368257eb1e3d23648920a742b",
		},
		{
			name: "podman cgroup",
			cgroupV1FileContent: `14:name=systemd:/user.slice/user-1000.slice/user@1000.service/app.slice/podman.service
13:rdma:/
12:pids:/user.slice/user-1000.slice/user@1000.service
11:hugetlb:/
10:net_prio:/
9:perf_event:/
8:net_cls:/
7:freezer:/
6:devices:/user.slice
5:blkio:/user.slice
4:cpuacct:/
3:cpu:/user.slice
2:cpuset:/
1:memory:/user.slice/user-1000.slice/user@1000.service
0::/user.slice/user-1000.slice/user@1000.service/app.slice/podman.service`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = func(name string) (os.FileInfo, error) {
				if tc.cgroupV1FileNotExist && name == cgroupV1Path {
					return nil, os.ErrNotExist
				}
				if tc.cgroupV2FileNotExist && name == cgroupV2Path {
					return nil, os.ErrNotExist
				}
				return nil, nil
			}

			osOpen = func(name string) (io.ReadCloser, error) {
				if tc.openFileError != nil {
					return nil, tc.openFileError
				}
				if name == cgroupV1Path {
					return io.NopCloser(strings.NewReader(tc.cgroupV1FileContent)), nil
				}
				return io.NopCloser(strings.NewReader(tc.cgroupV2FileContent)), nil
			}

			containerID, err := getContainerIDFromCGroup()
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
