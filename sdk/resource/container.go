// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"

	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

type containerIDProvider func() (string, error)

var (
	containerID         containerIDProvider = getContainerIDFromCGroup
	cgroupContainerIDRe                     = regexp.MustCompile(`^.*/(?:.*[-:])?([0-9a-f]{32,})(?:\.|\s*$)`)
	mountInfoContainerIDRe                  = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type cgroupContainerIDDetector struct{}

const (
	cgroupPath    = "/proc/self/cgroup"
	mountInfoPath = "/proc/self/mountinfo"
)

// Detect returns a *Resource that describes the id of the container.
// If no container id found, an empty resource will be returned.
func (cgroupContainerIDDetector) Detect(context.Context) (*Resource, error) {
	containerID, err := containerID()
	if err != nil {
		return nil, err
	}

	if containerID == "" {
		return Empty(), nil
	}
	return NewWithAttributes(semconv.SchemaURL, semconv.ContainerID(containerID)), nil
}

var (
	defaultOSStat = os.Stat
	osStat        = defaultOSStat

	defaultOSOpen = func(name string) (io.ReadCloser, error) {
		return os.Open(name)
	}
	osOpen = defaultOSOpen
)

// getContainerIDFromCGroup returns the id of the container from the cgroup
// file. It tries cgroup v1 (/proc/self/cgroup) first, then falls back to
// cgroup v2 (/proc/self/mountinfo). If no container id found, an empty string
// will be returned.
func getContainerIDFromCGroup() (string, error) {
	if _, err := osStat(cgroupPath); !errors.Is(err, os.ErrNotExist) {
		file, err := osOpen(cgroupPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		if id := getContainerIDFromReader(file); id != "" {
			return id, nil
		}
	}

	// Fall back to cgroup v2: read /proc/self/mountinfo.
	return getContainerIDFromMountInfo()
}

// getContainerIDFromMountInfo returns the id of the container from the
// mountinfo file. If no container id found, an empty string will be returned.
func getContainerIDFromMountInfo() (string, error) {
	if _, err := osStat(mountInfoPath); errors.Is(err, os.ErrNotExist) {
		return "", nil
	}

	file, err := osOpen(mountInfoPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return getContainerIDFromMountInfoReader(file), nil
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
	matches := cgroupContainerIDRe.FindStringSubmatch(line)
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}

// getContainerIDFromMountInfoReader scans mountinfo lines for a container ID
// using three strategies in priority order across the entire file:
//
//  1. Runtime-specific prefixes (/crio-, cri-containerd:) - first match wins.
//  2. Generic /containers/ lines - last 64-hex path segment across all matching
//     lines (handles CRI-O infra vs workload container ordering).
//  3. Hostname-gated fallback - first 64-hex segment from a "hostname" line
//     (catches containerd-minikube /sandboxes/ paths).
func getContainerIDFromMountInfoReader(reader io.Reader) string {
	var genericMatch, hostnameMatch string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if id := getContainerIDFromMountInfoLine(line); id != "" {
			return id
		}

		if strings.Contains(line, "/containers/") {
			if id := getLastHexSegment(line); id != "" {
				genericMatch = id
			}
		}

		if hostnameMatch == "" && strings.Contains(line, "hostname") {
			if id := getLastHexSegment(line); id != "" {
				hostnameMatch = id
			}
		}
	}

	if genericMatch != "" {
		return genericMatch
	}
	return hostnameMatch
}

// getContainerIDFromMountInfoLine extracts a container ID from a mountinfo
// line using runtime-specific prefixes (/crio- and cri-containerd:).
func getContainerIDFromMountInfoLine(line string) string {
	for _, prefix := range []string{"/crio-", "cri-containerd:"} {
		idx := strings.LastIndex(line, prefix)
		if idx == -1 {
			continue
		}
		start := idx + len(prefix)
		if start+64 > len(line) {
			continue
		}
		candidate := line[start : start+64]
		if mountInfoContainerIDRe.MatchString(candidate) {
			return candidate
		}
	}
	return ""
}

// getLastHexSegment returns the last path segment in line that is a
// 64-character lowercase hex string, or "" if none is found.
func getLastHexSegment(line string) string {
	var result string
	for _, segment := range strings.Split(line, "/") {
		if mountInfoContainerIDRe.MatchString(segment) {
			result = segment
		}
	}
	return result
}
