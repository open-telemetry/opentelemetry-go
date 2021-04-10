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
	"fmt"
	"syscall"
)

// osDescription issues a uname(2) system call and formats the output in a single
// string, similar to the output of the `uname` commandline program. The final string
// resembles the one obtained with a call to `uname -snrvm`.
func osDescription() (string, error) {
	var utsName syscall.Utsname

	err := syscall.Uname(&utsName)
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
func charsToString(charArray []int8) string {
	s := make([]byte, len(charArray))

	var i int
	for ; i < len(charArray) && charArray[i] != 0; i++ {
		s[i] = uint8(charArray[i])
	}

	return string(s[0:i])
}
