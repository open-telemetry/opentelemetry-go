# Copyright The OpenTelemetry Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM golang:alpine AS base
COPY . /go/src/github.com/open-telemetry/opentelemetry-go/
WORKDIR /go/src/github.com/open-telemetry/opentelemetry-go/

FROM base AS example-http-server
RUN go install ./example/http/server/server.go
CMD ["/go/bin/server"]

FROM base AS example-http-client
RUN go install ./example/http/client/client.go
CMD ["/go/bin/client"]

FROM base AS example-zipkin-client
RUN go install ./example/zipkin/main.go
CMD ["/go/bin/main"]
