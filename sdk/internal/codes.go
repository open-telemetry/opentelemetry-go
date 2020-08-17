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

package internal

import (
	grpccodes "google.golang.org/grpc/codes"

	otelcodes "go.opentelemetry.io/otel/codes"
)

// conversions are the equivalence mapping from OpenTelemetry to gRPC codes.
// Even though the underlying value should be the same all mappings are
// explicit here to avoid any error.
var conversions = map[otelcodes.Code]grpccodes.Code{
	otelcodes.OK:                 grpccodes.OK,
	otelcodes.Canceled:           grpccodes.Canceled,
	otelcodes.Unknown:            grpccodes.Unknown,
	otelcodes.InvalidArgument:    grpccodes.InvalidArgument,
	otelcodes.DeadlineExceeded:   grpccodes.DeadlineExceeded,
	otelcodes.NotFound:           grpccodes.NotFound,
	otelcodes.AlreadyExists:      grpccodes.AlreadyExists,
	otelcodes.PermissionDenied:   grpccodes.PermissionDenied,
	otelcodes.ResourceExhausted:  grpccodes.ResourceExhausted,
	otelcodes.FailedPrecondition: grpccodes.FailedPrecondition,
	otelcodes.Aborted:            grpccodes.Aborted,
	otelcodes.OutOfRange:         grpccodes.OutOfRange,
	otelcodes.Unimplemented:      grpccodes.Unimplemented,
	otelcodes.Internal:           grpccodes.Internal,
	otelcodes.Unavailable:        grpccodes.Unavailable,
	otelcodes.DataLoss:           grpccodes.DataLoss,
	otelcodes.Unauthenticated:    grpccodes.Unauthenticated,
}

// ConvertCode converts an OpenTelemetry Code into the equivalent gRPC code.
func ConvertCode(code otelcodes.Code) grpccodes.Code {
	return conversions[code]
}
