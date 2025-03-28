// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package resource_test

import (
	"fmt"
	"os"
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
	return fmt.Errorf("error invoking uname(2)")
}

func TestUname(t *testing.T) {
	resource.SetUnameProvider(fakeUnameProvider)

	uname, err := resource.Uname()

	require.Equal(t, "Mock OS DESKTOP-PC 5.0.0 #1 SMP Thu May 6 12:34:56 UTC 2021 x86_64", uname)
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

func TestGetFirstAvailableFile(t *testing.T) {
	tempDir := t.TempDir()

	file1, _ := os.CreateTemp(tempDir, "candidate_")
	file2, _ := os.CreateTemp(tempDir, "candidate_")

	filename1, filename2 := file1.Name(), file2.Name()

	tt := []struct {
		Name             string
		Candidates       []string
		ExpectedFileName string
		ExpectedErr      string
	}{
		{"Gets first, skip second candidate", []string{filename1, filename2}, filename1, ""},
		{"Skips first, gets second candidate", []string{"does_not_exists", filename2}, filename2, ""},
		{
			"Skips first, gets second, ignores third candidate",
			[]string{"does_not_exists", filename2, filename1},
			filename2,
			"",
		},
		{"No candidates (empty slice)", []string{}, "", "no candidate file available: []"},
		{"No candidates (nil slice)", nil, "", "no candidate file available: []"},
		{
			"Single nonexisting candidate",
			[]string{"does_not_exists"},
			"",
			"no candidate file available: [does_not_exists]",
		},
		{
			"Multiple nonexisting candidates",
			[]string{"does_not_exists", "this_either"},
			"",
			"no candidate file available: [does_not_exists this_either]",
		},
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
