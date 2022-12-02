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
			name:                "resolv.conf",
			line:                "472 456 254:1 /docker/containers/dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183/resolv.conf /etc/resolv.conf rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183",
		},
		{
			name:                "hostname",
			line:                "473 456 254:1 /docker/containers/dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183/hostname /etc/hostname rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183",
		},
		{
			name:                "host",
			line:                "474 456 254:1 /docker/containers/dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183/hosts /etc/hosts rw,relatime - ext4 /dev/vda1 rw",
			expectedContainerID: "dc64b5743252dbaef6e30521c34d6bbd1620c8ce65bdb7bf9e7143b61bb5b183",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containerID := getContainerIDFromCgroupV2Line(tc.line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
