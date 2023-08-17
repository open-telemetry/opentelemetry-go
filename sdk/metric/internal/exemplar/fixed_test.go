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

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/internal/exemplar"

import (
	"log"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TestDroppedHandlesInvalidGracefully(t *testing.T) {
	var msg string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) {
			msg = args
		}, funcr.Options{Verbosity: 20}))
		return func() { otel.SetLogger(orig) }
	}(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))))

	dest := []attribute.KeyValue{{}}
	// measured < recorded is invalid.
	dropped(&dest, fltrAlice, alice)

	assert.Contains(t, msg, "invalid measured attributes for exemplar, dropping")
	assert.Len(t, dest, 0, "dest not reset")
}
