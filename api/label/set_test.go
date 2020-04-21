// Copyright 2019, OpenTelemetry Authors
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

package label_test

import (
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/label"

	"github.com/stretchr/testify/require"
)

func TestSetDedup(t *testing.T) {
	tmp := &label.Sortable{}
	enc := label.NewDefaultEncoder()

	sl1 := []core.KeyValue{key.String("A", "1"), key.String("C", "D"), key.String("A", "B")}
	sl2 := []core.KeyValue{key.String("A", "2"), key.String("A", "B"), key.String("C", "D")}

	s1 := label.NewSet(sl1, tmp)
	s2 := label.NewSet(sl2, tmp)

	require.Equal(t, s1.Equivalent(), s2.Equivalent())

	require.Equal(t, sl1[0], key.String("A", "1"))
	require.Equal(t, sl2[0], key.String("A", "2"))

	require.Equal(t, "A=B,C=D", s1.Encoded(enc))
	require.Equal(t, "A=B,C=D", s2.Encoded(enc))
}
