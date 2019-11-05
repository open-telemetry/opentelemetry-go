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

package grpctrace // import "go.opentelemetry.io/otel/plugin/grpctrace"

import (
	"path"
	"strings"
)

const (
	componentKey    = "component"
	componentValue  = "grpc"
	peerServiceKey  = "peer.service"
	peerHostnameKey = "peer.hostname"
	peerPortKey     = "peer.port"
)

func serviceFromMethod(method string) string {
	return strings.TrimLeft(path.Ext(path.Dir(method)), ".")
}

func nameFromMethod(method string) string {
	return strings.TrimLeft(method, "/")
}
