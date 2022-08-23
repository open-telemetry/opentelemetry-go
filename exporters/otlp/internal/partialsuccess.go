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

package internal // import "go.opentelemetry.io/otel/exporters/otlp/internal"

import "fmt"

// ErrPartialSuccess is the underlying error for all handling OTLP
// partial success messages.  Use `errors.Is(err, ErrPartialSuccess)`
// to test whether `err` belongs to this category.
var ErrPartialSuccess = fmt.Errorf("OTLP partial success")

// PartialSuccessToError produces an error suitable for passing to
// `otel.Handle()` out of the fields in a partial success response,
// independent of which signal produced the outcome.
func PartialSuccessToError(itemPlural string, itemsRejected int64, errorMessage string) error {
	if errorMessage == "" {
		errorMessage = "empty message"
	}
	return fmt.Errorf("%w: %s (%d %s rejected)", ErrPartialSuccess, errorMessage, itemsRejected, itemPlural)
}
