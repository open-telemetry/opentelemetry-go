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

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	want, err := os.ReadFile("testdata/want.go")
	require.NoError(t, err)
	wantText := string(want)

	out := filepath.Join(t.TempDir(), "out.go")

	err = Copy(out, "pkg", "testdata/in.go")
	require.NoError(t, err)

	got, err := os.ReadFile(out)
	require.NoError(t, err)
	gotText := string(got)

	assert.Equal(t, wantText, gotText)
}
