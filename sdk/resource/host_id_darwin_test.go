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
	"testing"

	"github.com/stretchr/testify/require"
)

const cmdOutput = `+-o J316sAP  <class IOPlatformExpertDevice, id 0x10000024d, registered, matched, active, busy 0 (132196 ms), retain 37>
{
  "IOPolledInterface" = "AppleARMWatchdogTimerHibernateHandler is not serializable"
  "#address-cells" = <02000000>
  "AAPL,phandle" = <01000000>
  "serial-number" = <94e1c79ec04cd3f153f600000000000000000000000000000000000000000000>
  "IOBusyInterest" = "IOCommand is not serializable"
  "target-type" = <"J316s">
  "platform-name" = <7436303030000000000000000000000000000000000000000000000000000000>
  "secure-root-prefix" = <"md">
  "name" = <"device-tree">
  "region-info" = <4c4c2f4100000000000000000000000000000000000000000000000000000000>
  "manufacturer" = <"Apple Inc.">
  "compatible" = <"J316sAP","MacBookPro18,1","AppleARM">
  "config-number" = <00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000>
  "IOPlatformSerialNumber" = "HDWLIF2LM7"
  "regulatory-model-number" = <4132343835000000000000000000000000000000000000000000000000000000>
  "time-stamp" = <"Fri Aug 5 20:25:38 PDT 2022">
  "clock-frequency" = <00366e01>
  "model" = <"MacBookPro18,1">
  "mlb-serial-number" = <5c92d268d6cd789e475ffafc0d363fc950000000000000000000000000000000>
  "model-number" = <5a31345930303136430000000000000000000000000000000000000000000000>
  "IONWInterrupts" = "IONWInterrupts"
  "model-config" = <"ICT;MoPED=0x03D053A605C84ED11C455A18D6C643140B41A239">
  "device_type" = <"bootrom">
  "#size-cells" = <02000000>
  "IOPlatformUUID" = "81895B8D-9EF9-4EBB-B5DE-B00069CF53F0"
}
`

func TestReaderValidOutput(t *testing.T) {
	expectedHostID := "81895B8D-9EF9-4EBB-B5DE-B00069CF53F0"
	reader := &hostIDReaderDarwin{
		execCommand: func(string, ...string) (string, error) {
			return cmdOutput, nil
		},
	}

	result, err := reader.read()
	require.NoError(t, err)
	require.Equal(t, expectedHostID, result)
}

func TestReaderInvalidOutput(t *testing.T) {
	reader := &hostIDReaderDarwin{
		execCommand: func(string, ...string) (string, error) {
			return "not expecting this", nil
		},
	}

	result, err := reader.read()
	require.Error(t, err)
	require.Empty(t, result)
}

func TestReaderError(t *testing.T) {
	reader := &hostIDReaderDarwin{
		execCommand: func(string, ...string) (string, error) {
			return "", errors.New("could not parse host id")
		},
	}

	result, err := reader.read()
	require.Error(t, err)
	require.Empty(t, result)
}
