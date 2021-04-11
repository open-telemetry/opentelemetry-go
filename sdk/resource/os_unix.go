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

// osDescription returns a human readable OS version information string. The final
// string combines the data of the os-release file (where available) and the result
// of the `uname` system call.
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

// uname issues a uname(2) system call (or equivalent on systems which doesn't
// have one) and formats the output in a single string, similar to the output
// of the `uname` commandline program. The final string resembles the one
// obtained with a call to `uname -snrvm`.
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

// osRelease builds a string describing the operating system release based on the
// properties of the os-release file. If no os-release file is found, or if the
// required properties to build de release description string are missing, an empty
// string is returned instead. For more information about the os-release file, see:
// https://www.freedesktop.org/software/systemd/man/os-release.html
func osRelease() string {
	file, err := getOSReleaseFile()
	if err != nil {
		return ""
	}

	defer file.Close()

	values := parseOSReleaseFile(file)

	return buildOSRelease(values)
}

// getOSReleaseFile returns a *os.File pointing to one of the well-known os-release
// files, according to their order of preference. If no file can be opened, it
// returns the last error.
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

// parseOSReleaseFile process the file pointed by `file` as an os-release file and
// returns a map with the key-values contained in it. Empty lines or lines starting
// with a '#' character are ignored, as well as lines with the missing key=value
// separator. Values are unquoted and unescaped.
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

// skip returns true if the line is blank or starts with a '#' character, and
// therefore should be skipped from processing.
func skip(line string) bool {
	line = strings.TrimSpace(line)

	return len(line) == 0 || strings.HasPrefix(line, "#")
}

// parse attempts to split the provided line on the first '=' character, and then
// sanitize each side of the split before returning them as a key-value pair.
func parse(line string) (string, string, bool) {
	parts := strings.SplitN(line, "=", 2)

	if len(parts) != 2 {
		return "", "", false
	}

	key := strings.TrimSpace(parts[0])
	value := unescape(unquote(strings.TrimSpace(parts[1])))

	return key, value, true
}

// unquote checks whether the string `s` is quoted with double or single quotes
// and, if so, returns a version of the string without them. Otherwise it returns
// the provided string unchanged.
func unquote(s string) string {
	if (s[0] == '"' || s[0] == '\'') && s[0] == s[len(s)-1] {
		return s[1 : len(s)-1]
	}

	return s
}

// unescape removes the `\` prefix from some characters that are expected
// to have it added in front of them for escaping purposes.
func unescape(s string) string {
	s = strings.ReplaceAll(s, `\$`, `$`)
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	s = strings.ReplaceAll(s, "\\`", "`")

	return s
}

// buildOSRelease builds a string describing the OS release based on the properties
// available on the provided map. It favors a combination of the `NAME` and `VERSION`
// properties as first option (falling back to `VERSION_ID` if `VERSION` isn't
// found), and using `PRETTY_NAME` alone if some of the previous are not present. If
// none of these properties are found, it returns an empty string.
//
// The rationale behind not using `PRETTY_NAME` as first choice was that, for some
// Linux distributions, it doesn't include the same detail that can be found on the
// individual `NAME` and `VERSION` properties, and combining `PRETTY_NAME` with
// other properties can produce "pretty" redundant strings in some cases.
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
