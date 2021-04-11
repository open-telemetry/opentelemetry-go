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

// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

// osDescription returns a human readable OS version information string.
// It issues a uname(2) system call (or equivalent on systems which doesn't
// have one) and formats the output in a single string, similar to the output
// of the `uname` commandline program. The final string resembles the one
// obtained with a call to `uname -snrvm`.
func osDescription() (string, error) {
	uname, err := uname()
	if err != nil {
		return "", err
	}

	osRelease := osRelease()
	if osRelease != "" {
		return fmt.Sprintf("%s (%s)", osRelease, uname), nil
	}

	return uname, nil
}

func uname() (string, error) {
	var utsName unix.Utsname

	err := unix.Uname(&utsName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s %s %s %s",
		charsToString(utsName.Sysname[:]),
		charsToString(utsName.Nodename[:]),
		charsToString(utsName.Release[:]),
		charsToString(utsName.Version[:]),
		charsToString(utsName.Machine[:]),
	), nil
}

// charsToString converts a C-like null-terminated char array to a Go string.
func charsToString(charArray []byte) string {
	s := make([]byte, len(charArray))

	var i int
	for ; i < len(charArray) && charArray[i] != 0; i++ {
		s[i] = uint8(charArray[i])
	}

	return string(s[0:i])
}

func osRelease() string {
	file, err := getOSReleaseFile()
	if err != nil {
		return ""
	}

	defer file.Close()

	values := parseOSReleaseFile(file)

	return buildOSRelease(values)
}

func getOSReleaseFile() (*os.File, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		file, err = os.Open("/usr/lib/os-release")
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func parseOSReleaseFile(file *os.File) map[string]string {
	values := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if skip(line) {
			continue
		}

		key, value, ok := parse(line)
		if ok {
			values[key] = value
		}
	}

	return values
}

func skip(line string) bool {
	line = strings.TrimSpace(line)

	return len(line) == 0 || strings.HasPrefix(line, "#")
}

func parse(line string) (string, string, bool) {
	parts := strings.SplitN(line, "=", 2)

	if len(parts) != 2 {
		return "", "", false
	}

	key := strings.TrimSpace(parts[0])
	value := unescape(unquote(strings.TrimSpace(parts[1])))

	return key, value, true
}

func unquote(s string) string {
	if (s[0] == '"' || s[0] == '\'') && s[0] == s[len(s)-1] {
		return s[1 : len(s)-1]
	}

	return s
}

func unescape(s string) string {
	s = strings.ReplaceAll(s, `\$`, `$`)
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	s = strings.ReplaceAll(s, "\\`", "`")

	return s
}

func buildOSRelease(values map[string]string) string {
	var osRelease string

	name := values["NAME"]
	version := values["VERSION"]

	if version == "" {
		version = values["VERSION_ID"]
	}

	if name != "" && version != "" {
		osRelease = fmt.Sprintf("%s %s", name, version)
	} else {
		osRelease = values["PRETTY_NAME"]
	}

	return osRelease
}
