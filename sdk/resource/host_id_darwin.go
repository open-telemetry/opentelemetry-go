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
	"strings"
)

// hostIDReaderDarwin implements hostIDReader
type hostIDReaderDarwin struct {
	execCommand commandExecutor
}

// read executes `ioreg -rd1 -c "IOPlatformExpertDevice"` and parses host id
// from the IOPlatformUUID line. If the command fails or the uuid cannot be
// parsed an error will be returned.
func (r *hostIDReaderDarwin) read() (string, error) {
	result, err := r.execCommand("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	if err != nil {
		return "", err
	}

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "IOPlatformUUID") {
			parts := strings.Split(line, " = ")
			if len(parts) == 2 {
				return strings.Trim(parts[1], "\""), nil
			}
			break
		}
	}

	return "", errors.New("could not parse IOPlatformUUID")
}

var platformHostIDReader hostIDReader = &hostIDReaderDarwin{
	execCommand: execCommand,
}
