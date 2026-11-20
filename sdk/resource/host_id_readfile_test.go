// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build linux || dragonfly || freebsd || netbsd || openbsd || solaris

package resource

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFileExistent(t *testing.T) {
	fileContents := "foo"

	f, err := os.CreateTemp(t.TempDir(), "readfile_")
	require.NoError(t, err)

	defer os.Remove(f.Name())

	_, err = f.WriteString(fileContents)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	result, err := readFile(f.Name())
	require.NoError(t, err)
	require.Equal(t, result, fileContents)
}

func TestReadFileNonExistent(t *testing.T) {
	// create unique filename
	f, err := os.CreateTemp(t.TempDir(), "readfile_")
	require.NoError(t, err)

	// make file non-existent
	require.NoError(t, os.Remove(f.Name()))

	_, err = readFile(f.Name())
	require.ErrorIs(t, err, os.ErrNotExist)
}
