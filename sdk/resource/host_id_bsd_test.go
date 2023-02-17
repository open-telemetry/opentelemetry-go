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

//go:build dragonfly || freebsd || netbsd || openbsd || solaris
// +build dragonfly freebsd netbsd openbsd solaris

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderValidPrimary(t *testing.T) {
	expectedHostID := "f2c668b579780554f70f72a063dc0864"
	reader := &hostIDReaderBSD{
		readFile: func(filename string) (string, error) {
			return expectedHostID + "\n", nil
		},
	}

	result, err := reader.read()
	require.NoError(t, err)
	require.Equal(t, expectedHostID, result)
}

func TestReaderInvalidPrimary(t *testing.T) {
	expectedHostID := "f2c668b579780554f70f72a063dc0864"
	reader := &hostIDReaderBSD{
		readFile: func(filename string) (string, error) {
			return "", errors.New("not found")
		},
		execCommand: func(string, ...string) (string, error) {
			return expectedHostID + "\n", nil
		},
	}

	result, err := reader.read()
	require.NoError(t, err)
	require.Equal(t, expectedHostID, result)
}

func TestReaderError(t *testing.T) {
	reader := &hostIDReaderBSD{
		readFile: func(string) (string, error) {
			return "", errors.New("could not parse host id")
		},
		execCommand: func(string, ...string) (string, error) {
			return "", errors.New("not found")
		},
	}

	result, err := reader.read()
	require.Error(t, err)
	require.Empty(t, result)
}
