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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// View returns true and the InstrumentStream to use for matching
// InstrumentProperties. Otherwise, if the view does not match, false is
// returned.
type View func(InstrumentProperties) (InstrumentStream, bool)

// NewView returns a View composed for the stream and all passed MatchFuncs.
// The returned View will only return stream if all match functions match the
// instrument properties they are passed.
func NewView(stream InstrumentStream, match ...MatchFunc) View {
	return func(p InstrumentProperties) (InstrumentStream, bool) {
		// Logically "AND" all matches.
		for _, m := range match {
			if !m(p) {
				return InstrumentStream{}, false
			}
		}
		return stream, true
	}
}

// MatchFunc returns true for properties that match a condition, false
// otherwise.
type MatchFunc func(InstrumentProperties) bool

// MatchName returns a MatchFunc that matches the exact name of an instrument.
func MatchName(name string) MatchFunc {
	return func(p InstrumentProperties) bool { return p.Name == name }
}

// MatchNameWildcard returns a MatchFunc that matches instrument names using
// the wildcard pattern. The wildcard pattern recognizes "*" as matching zero
// or more characters, and "?" as matching exactly one character. A pattern of
// just "*" will match all instruments.
func MatchNameWildcard(pattern string) MatchFunc {
	pattern = regexp.QuoteMeta(pattern)
	pattern = "^" + pattern + "$"
	pattern = strings.ReplaceAll(pattern, "\\?", ".")
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")
	re := regexp.MustCompile(pattern)
	return func(p InstrumentProperties) bool { return re.MatchString(p.Name) }
}

// MatchKind returns a MatchFunc that matches the exact kind of instruments.
func MatchKind(kind InstrumentKind) MatchFunc {
	return func(p InstrumentProperties) bool { return p.Kind == kind }
}

// MatchScope returns a MatchFunc that matches the exact scope of instruments.
func MatchScope(scope instrumentation.Scope) MatchFunc {
	return func(p InstrumentProperties) bool { return p.Scope == scope }
}
