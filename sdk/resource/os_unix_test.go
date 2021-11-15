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

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package resource_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	"go.opentelemetry.io/otel/sdk/resource"
)

func fakeUnameProvider(buf *unix.Utsname) error {
	copy(buf.Sysname[:], "Mock OS")
	copy(buf.Nodename[:], "DESKTOP-PC")
	copy(buf.Release[:], "5.0.0")
	copy(buf.Version[:], "#1 SMP Thu May 6 12:34:56 UTC 2021")
	copy(buf.Machine[:], "x86_64")

	return nil
}

func fakeUnameProviderWithError(buf *unix.Utsname) error {
	return fmt.Errorf("Error invoking uname(2)")
}

func TestUname(t *testing.T) {
	resource.SetUnameProvider(fakeUnameProvider)

	uname, err := resource.Uname()

	require.Equal(t, uname, "Mock OS DESKTOP-PC 5.0.0 #1 SMP Thu May 6 12:34:56 UTC 2021 x86_64")
	require.NoError(t, err)

	resource.SetDefaultUnameProvider()
}

func TestUnameError(t *testing.T) {
	resource.SetUnameProvider(fakeUnameProviderWithError)

	uname, err := resource.Uname()

	require.Empty(t, uname)
	require.Error(t, err)

	resource.SetDefaultUnameProvider()
}

func TestCharsToString(t *testing.T) {
	tt := []struct {
		Name     string
		Bytes    []byte
		Expected string
	}{
		{"Nil array", nil, ""},
		{"Empty array", []byte{}, ""},
		{"Empty string (null terminated)", []byte{0x00}, ""},
		{"Nonempty string (null terminated)", []byte{0x31, 0x32, 0x33, 0x00}, "123"},
		{"Nonempty string (non-null terminated)", []byte{0x31, 0x32, 0x33}, "123"},
		{"Nonempty string with values after null", []byte{0x31, 0x32, 0x33, 0x00, 0x34}, "123"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			result := resource.CharsToString(tc.Bytes)
			require.EqualValues(t, tc.Expected, result)
		})
	}
}

func TestGetFirstAvailableFile(t *testing.T) {
	tempDir := t.TempDir()

	file1, _ := ioutil.TempFile(tempDir, "candidate_")
	file2, _ := ioutil.TempFile(tempDir, "candidate_")

	filename1, filename2 := file1.Name(), file2.Name()

	tt := []struct {
		Name             string
		Candidates       []string
		ExpectedFileName string
		ExpectedErr      string
	}{
		{"Gets first, skip second candidate", []string{filename1, filename2}, filename1, ""},
		{"Skips first, gets second candidate", []string{"does_not_exists", filename2}, filename2, ""},
		{"Skips first, gets second, ignores third candidate", []string{"does_not_exists", filename2, filename1}, filename2, ""},
		{"No candidates (empty slice)", []string{}, "", "no candidate file available: []"},
		{"No candidates (nil slice)", nil, "", "no candidate file available: []"},
		{"Single nonexisting candidate", []string{"does_not_exists"}, "", "no candidate file available: [does_not_exists]"},
		{"Multiple nonexisting candidates", []string{"does_not_exists", "this_either"}, "", "no candidate file available: [does_not_exists this_either]"},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			file, err := resource.GetFirstAvailableFile(tc.Candidates)

			filename := ""
			if file != nil {
				filename = file.Name()
			}

			errString := ""
			if err != nil {
				errString = err.Error()
			}

			require.Equal(t, tc.ExpectedFileName, filename)
			require.Equal(t, tc.ExpectedErr, errString)
		})
	}
}
