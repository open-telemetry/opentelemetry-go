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

package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sUtil "go.opentelemetry.io/otel/schema/v1.1"
)

func TestRender(t *testing.T) {
	schemaF, err := os.Open("./testdata/schema.yaml")
	require.NoError(t, err)
	defer schemaF.Close()

	s, err := sUtil.Parse(schemaF)
	require.NoError(t, err)

	wantF, err := os.Open("./testdata/schema.go")
	require.NoError(t, err)

	var want bytes.Buffer
	_, err = io.Copy(&want, wantF)
	require.NoError(t, err)

	e, err := newEntry(s)
	require.NoError(t, err)

	var got bytes.Buffer
	render(&got, []entry{e})

	assert.Equal(t, want.String(), got.String())
}
