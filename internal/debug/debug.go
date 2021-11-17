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

package debug // import "go.opentelemetry.io/otel/internal/debug"

import (
	"log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

var globalLoggger logr.Logger = stdr.New(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

func SetLogger(l logr.Logger) {
	globalLoggger = l
}

func Info(msg string, keysAndValues ...interface{}) {
	globalLoggger.V(1).Info(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	globalLoggger.Error(err, msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	globalLoggger.V(5).Info(msg, keysAndValues...)
}
