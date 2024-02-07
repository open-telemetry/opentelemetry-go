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

//go:build linux || dragonfly || freebsd || netbsd || openbsd || solaris

package resource

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFileExistent(t *testing.T) {
	fileContents := "foo"

	f, err := os.CreateTemp("", "readfile_")
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
	f, err := os.CreateTemp("", "readfile_")
	require.NoError(t, err)

	// make file non-existent
	require.NoError(t, os.Remove(f.Name()))

	_, err = readFile(f.Name())
	require.ErrorIs(t, err, os.ErrNotExist)
}
