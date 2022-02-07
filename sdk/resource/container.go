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
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"

	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type containerIDProvider func() (string, error)

var (
	containerID containerIDProvider = getContainerIDFromCGroup
)

type cgroupContainerIDDetector struct{}

const cgroupPath = "/proc/self/cgroup"

// Detect returns a *Resource that describes the id of the container.
// If no container id found, an empty resource will be returned.
func (cgroupContainerIDDetector) Detect(ctx context.Context) (*Resource, error) {
	containerID, err := containerID()
	if err != nil {
		return nil, err
	}

	if containerID == "" {
		return Empty(), nil
	}
	return NewWithAttributes(semconv.SchemaURL, semconv.ContainerIDKey.String(containerID)), nil
}

// getContainerIDFromCGroup returns the id of the container from the cgroup file.
// If no container id found, an empty string will be returned.
func getContainerIDFromCGroup() (string, error) {
	if _, err := os.Stat(cgroupPath); errors.Is(err, os.ErrNotExist) {
		// File does not exists, skip
		return "", nil
	}

	file, err := os.Open(cgroupPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	containerID := getContainerIDFromReader(file)
	if containerID == "" {
		// Container ID not found
		return "", nil
	}
	return containerID, nil
}

// getContainerIDFromReader returns the id of the container from reader.
func getContainerIDFromReader(reader io.Reader) string {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if id := getContainerIDFromLine(line); id != "" {
			return id
		}
	}
	return ""
}

// getContainerIDFromLine returns the id of the container from one string line.
func getContainerIDFromLine(line string) string {
	line = strings.TrimSpace(line)

	lastSlashIndexOfLine := strings.LastIndexByte(line, '/')
	if lastSlashIndexOfLine == -1 {
		return ""
	}

	lastSection := line[lastSlashIndexOfLine+1:]
	dashIndex := strings.IndexByte(lastSection, '-')
	lastDotIndex := strings.LastIndexByte(lastSection, '.')

	startIndex := 0
	if dashIndex != -1 {
		startIndex = dashIndex + 1
	}

	endIndex := len(lastSection)
	if lastDotIndex != -1 {
		endIndex = lastDotIndex
	}

	containerID := lastSection[startIndex:endIndex]
	if !isHex(containerID) {
		return ""
	}
	return containerID
}

// isHex returns true when input is a hex string.
func isHex(h string) bool {
	for _, r := range h {
		switch {
		case 'a' <= r && r <= 'f':
			continue
		case '0' <= r && r <= '9':
			continue
		default:
			return false
		}
	}
	return true
}
