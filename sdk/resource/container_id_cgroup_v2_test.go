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

func TestGetContainerIDFromCgroupV2Line(t *testing.T) {
	testCases := []struct {
		name                string
		line                string
		expectedContainerID string
	}{
		{
			name: "empty - 1",
			line: "456 375 0:143 / / rw,relatime master:175 - overlay overlay rw,lowerdir=/var/lib/docker/overlay2/l/37L57D2IM7MEWLVE2Q2ECNDT67:/var/lib/docker/overlay2/l/46FCA2JFPCSNFGAR5TSYLLNHLK,upperdir=/var/lib/docker/overlay2/4e82c300793d703c19bdf887bfdad8b0354edda884ea27a8a2df89ab292719a4/diff,workdir=/var/lib/docker/overlay2/4e82c300793d703c19bdf887bfdad8b0354edda884ea27a8a2df89ab292719a4/work",
		},
		{
			name: "empty - 2",
			line: "457 456 0:146 / /proc rw,nosuid,nodev,noexec,relatime - proc proc rw",
		},
		{
			name: "empty - 3",
			line: "383 457 0:147 /null /proc/kcore rw,nosuid - tmpfs tmpfs rw,size=65536k,mode=755",
		},
		{
			name:                "hostname",
			line:                "473 456 254:1 /docker/containers/dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183",
		},
		{
			name:                "minikube containerd mountinfo",
			line:                "1537 1517 8:1 /var/lib/containerd/io.containerd.grpc.v1.cri/sandboxes/fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6/hostname /etc/hostname rw,relatime - ext4 /dev/sda1 rw",
			expectedContainerID: "fb5916a02feca96bdeecd8e062df9e5e51d6617c8214b5e1f3ff9320f4402ae6",
		},
		{
			name:                "minikube docker mountinfo",
			line:                "2327 2307 8:1 /var/lib/docker/containers/a1551a1d7e1881d6c18d2c9ec462cab6ad3666825f0adb2098e9d5b198fd7e19/hostname /etc/hostname rw,relatime - ext4 /dev/sda1 rw",
			expectedContainerID: "a1551a1d7e1881d6c18d2c9ec462cab6ad3666825f0adb2098e9d5b198fd7e19",
		},
		{
			name:                "minikube docker mountinfo 2",
			line:                "929 920 254:1 /docker/volumes/minikube/_data/lib/docker/containers/0eaa6718003210b6520f7e82d14b4c8d4743057a958a503626240f8d1900bc33/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "0eaa6718003210b6520f7e82d14b4c8d4743057a958a503626240f8d1900bc33",
		},
		{
			name:                "podman mountinfo",
			line:                "1096 1088 0:104 /containers/overlay-containers/1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809/userdata/hostname /etc/hostname rw,nosuid,nodev,relatime - tmpfs tmpfs rw,size=813800k,nr_inodes=203450,mode=700,uid=1000,gid=1000",
			expectedContainerID: "1a2de27e7157106568f7e081e42a8c14858c02bd9df30d6e352b298178b46809",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromCgroupV2Line(tc.line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
