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

package debug // import "go.opentelemetry.io/otel/debug"

import (
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"

	"go.opentelemetry.io/otel/internal/debug"
)

// SetLogger configures the logger used internally to opentelemetry.
func SetLogger(logger logr.Logger) {
	debug.Log = logger
}

// SetDefaultLogger configures the internal logger to use stderr and show verbose logging messages.
func SetDefaultLogger() {
	SetLogger(stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)))
	stdr.SetVerbosity(5)
}
