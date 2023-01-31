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

package internal // import "go.opentelemetry.io/otel/exporters/otlp/internal"

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrappedError(t *testing.T) {
	e := WrapTracesError(context.Canceled)

	require.Equal(t, context.Canceled, errors.Unwrap(e))
	require.Equal(t, TracesExport, e.(wrappedExportError).kind)
	require.Equal(t, "traces export: context canceled", e.Error())
	require.True(t, errors.Is(e, context.Canceled))
}
