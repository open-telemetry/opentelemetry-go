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
			cgroupV1FileContent: "1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23",
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
