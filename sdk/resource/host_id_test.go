// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	expectedHostID = "f2c668b579780554f70f72a063dc0864"

	readFileNoError = func(filename string) (string, error) {
		return expectedHostID + "\n", nil
	}

	readFileError = func(filename string) (string, error) {
		return "", errors.New("not found")
	}

	execCommandNoError = func(string, ...string) (string, error) {
		return expectedHostID + "\n", nil
	}

	execCommandError = func(string, ...string) (string, error) {
		return "", errors.New("not found")
	}
)

func SetDefaultHostIDProvider() {
	SetHostIDProvider(defaultHostIDProvider)
}

func SetHostIDProvider(hostIDProvider hostIDProvider) {
	hostID = hostIDProvider
}

func TestHostIDReaderBSD(t *testing.T) {
	tt := []struct {
		name            string
		fileReader      fileReader
		commandExecutor commandExecutor
		expectedHostID  string
		expectError     bool
	}{
		{
			name:            "hostIDReaderBSD valid primary",
			fileReader:      readFileNoError,
			commandExecutor: execCommandError,
			expectedHostID:  expectedHostID,
			expectError:     false,
		},
		{
			name:            "hostIDReaderBSD invalid primary",
			fileReader:      readFileError,
			commandExecutor: execCommandNoError,
			expectedHostID:  expectedHostID,
			expectError:     false,
		},
		{
			name:            "hostIDReaderBSD invalid primary and secondary",
			fileReader:      readFileError,
			commandExecutor: execCommandError,
			expectedHostID:  "",
			expectError:     true,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			reader := hostIDReaderBSD{
				readFile:    tc.fileReader,
				execCommand: tc.commandExecutor,
			}
			hostID, err := reader.read()
			require.Equal(t, tc.expectError, err != nil)
			require.Equal(t, tc.expectedHostID, hostID)
		})
	}
}

func TestHostIDReaderLinux(t *testing.T) {
	readFilePrimaryError := func(filename string) (string, error) {
		if filename == "/var/lib/dbus/machine-id" {
			return readFileNoError(filename)
		}
		return readFileError(filename)
	}

	tt := []struct {
		name           string
		fileReader     fileReader
		expectedHostID string
		expectError    bool
	}{
		{
			name:           "hostIDReaderLinux valid primary",
			fileReader:     readFileNoError,
			expectedHostID: expectedHostID,
			expectError:    false,
		},
		{
			name:           "hostIDReaderLinux invalid primary",
			fileReader:     readFilePrimaryError,
			expectedHostID: expectedHostID,
			expectError:    false,
		},
		{
			name:           "hostIDReaderLinux invalid primary and secondary",
			fileReader:     readFileError,
			expectedHostID: "",
			expectError:    true,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			reader := hostIDReaderLinux{
				readFile: tc.fileReader,
			}
			hostID, err := reader.read()
			require.Equal(t, tc.expectError, err != nil)
			require.Equal(t, tc.expectedHostID, hostID)
		})
	}
}

func TestHostIDReaderDarwin(t *testing.T) {
	validOutput := `+-o J316sAP  <class IOPlatformExpertDevice, id 0x10000024d, registered, matched, active, busy 0 (132196 ms), retain 37>
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
	execCommandValid := func(string, ...string) (string, error) {
		return validOutput, nil
	}

	execCommandInvalid := func(string, ...string) (string, error) {
		return "wasn't expecting this", nil
	}

	tt := []struct {
		name            string
		fileReader      fileReader
		commandExecutor commandExecutor
		expectedHostID  string
		expectError     bool
	}{
		{
			name:            "hostIDReaderDarwin valid output",
			commandExecutor: execCommandValid,
			expectedHostID:  "81895B8D-9EF9-4EBB-B5DE-B00069CF53F0",
			expectError:     false,
		},
		{
			name:            "hostIDReaderDarwin invalid output",
			commandExecutor: execCommandInvalid,
			expectedHostID:  "",
			expectError:     true,
		},
		{
			name:            "hostIDReaderDarwin error",
			commandExecutor: execCommandError,
			expectedHostID:  "",
			expectError:     true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			reader := hostIDReaderDarwin{
				execCommand: tc.commandExecutor,
			}
			hostID, err := reader.read()
			require.Equal(t, tc.expectError, err != nil)
			require.Equal(t, tc.expectedHostID, hostID)
		})
	}
}
