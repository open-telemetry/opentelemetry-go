// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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

func TestGetContainerIDFromLine(t *testing.T) {
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
			name:                "with colon",
			line:                " 13:name=systemd:/kuberuntime/containerd/kubepods-pod872d2066_00ef_48ea_a7d8_51b18b72d739:cri-containerd:e857a4bf05a69080a759574949d7a0e69572e27647800fa7faff6a05a8332aa1",
			expectedContainerID: "e857a4bf05a69080a759574949d7a0e69572e27647800fa7faff6a05a8332aa1",
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
			containerID := getContainerIDFromLine(tc.line)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
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
			containerID := getContainerIDFromReader(tc.reader)
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
		name                string
		cgroupFileNotExist  bool
		openFileError       error
		content             string
		expectedContainerID string
		expectedError       bool
	}{
		{
			name:               "the cgroup file does not exist",
			cgroupFileNotExist: true,
		},
		{
			name:          "error when opening cgroup file",
			openFileError: errors.New("test"),
			expectedError: true,
		},
		{
			name:                "cgroup file",
			content:             "1:name=systemd:/podruntime/docker/kubepods/docker-dc579f8a8319c8cf7d38e1adf263bc08d23",
			expectedContainerID: "dc579f8a8319c8cf7d38e1adf263bc08d23",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			osStat = func(name string) (os.FileInfo, error) {
				if tc.cgroupFileNotExist {
					return nil, os.ErrNotExist
				}
				return nil, nil
			}

			osOpen = func(name string) (io.ReadCloser, error) {
				if tc.openFileError != nil {
					return nil, tc.openFileError
				}
				return io.NopCloser(strings.NewReader(tc.content)), nil
			}

			containerID, err := getContainerIDFromCGroup()
			assert.Equal(t, tc.expectedError, err != nil)
			assert.Equal(t, tc.expectedContainerID, containerID)
		})
	}
}
