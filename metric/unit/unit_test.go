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

package unit // import "go.opentelemetry.io/otel/metric/unit"

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	code   = "custom code"
	symbol = "custom print symbol"
)

var hella = prefix{code: "hella", symbol: "🤘"}

func TestPrefix(t *testing.T) {
	p := prefix{code: code, symbol: symbol}
	assert.Equal(t, code, p.Code(), "prefix code")
	assert.Equal(t, symbol, p.String(), "prefix symbol")
}

func TestNew(t *testing.T) {
	u := New(code)
	u = u.withPrefix(hella)

	assert.Equal(t, hella.code+code, u.Code(), "unit code")
	assert.Equal(t, hella.code+code, u.String(), "unit code as symbol")
}

func TestNewWithPrintSymbol(t *testing.T) {
	u := New(code, WithPrintSymbol(symbol))
	u = u.withPrefix(hella)

	assert.Equal(t, hella.code+code, u.Code(), "unit code")
	assert.Equal(t, hella.symbol+symbol, u.String(), "unit symbol")
}

func TestUnitJSONMarshalling(t *testing.T) {
	orig := New(code, WithPrintSymbol(symbol))
	orig = orig.withPrefix(hella)
	got, err := json.Marshal(orig)
	require.NoError(t, err)
	require.Equal(t, `"`+hella.code+code+`"`, string(got))

	decoded := new(Unit)
	require.NoError(t, json.Unmarshal([]byte(`"`+code+`"`), decoded))
	assert.Equal(t, code, decoded.Code())
}
