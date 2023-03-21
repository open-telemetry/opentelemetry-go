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

package global

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/go-logr/logr"

	"github.com/stretchr/testify/assert"

	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
)

func TestRace(t *testing.T) {
	go SetLogger(stdr.New(log.New(os.Stderr, "", 0)))
	go Info("")
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		verbosity int
		logF      func()
		want      string
	}{
		{
			name:      "Verbosity 0 should log errors.",
			verbosity: 0,
			want:      `"msg"="foobar" "error"="foobar"`,
			logF: func() {
				Error(errors.New("foobar"), "foobar")
			},
		},
		{
			name:      "Verbosity 1 should log warnings",
			verbosity: 1,
			want:      `"level"=1 "msg"="foo"`,
			logF: func() {
				Warn("foo")
			},
		},
		{
			name:      "Verbosity 4 should log info",
			verbosity: 4,
			want:      `"level"=4 "msg"="bar"`,
			logF: func() {
				Info("bar")
			},
		},
		{
			name:      "Verbosity 8 should log debug",
			verbosity: 8,
			want:      `"level"=8 "msg"="baz"`,
			logF: func() {
				Debug("baz")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			SetLogger(newBuffLogger(&buf, test.verbosity))

			test.logF()

			assert.Equal(t, test.want, buf.String())
		})
	}
}

func newBuffLogger(buf *bytes.Buffer, verbosity int) logr.Logger {
	return funcr.New(func(prefix, args string) {
		_, _ = buf.Write([]byte(args))
	}, funcr.Options{
		Verbosity: verbosity,
	})
}
