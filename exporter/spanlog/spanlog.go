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

package spanlog

import (
	"os"
	"strings"

	"github.com/open-telemetry/opentelemetry-go/exporter/buffer"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
	"github.com/open-telemetry/opentelemetry-go/exporter/spandata"
	"github.com/open-telemetry/opentelemetry-go/exporter/spandata/format"
)

type (
	spanLog struct{}
)

func New() observer.Observer {
	return buffer.NewBuffer(1000, spandata.NewReaderObserver(&spanLog{}))
}

func (s *spanLog) Read(data *spandata.Span) {
	var buf strings.Builder
	buf.WriteString("----------------------------------------------------------------------\n")
	format.AppendSpan(&buf, data)
	os.Stdout.WriteString(buf.String())
}
