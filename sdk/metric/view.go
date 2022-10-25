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

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

// View is an override to the default behavior of the SDK. It defines how data
// should be collected for certain instruments. It returns true and the exact
// DataStream to use for matching InstrumentProperties. Otherwise, if the view
// does not match, false is returned.
type View func(InstrumentProperties) (DataStream, bool)

// NewView returns a View that applies the DataStream mask for all instruments
// that match criteria. The returned View will only apply mask if all
// non-zero-value fields of criteria match the corresponding
// InstrumentProperties passed to the view. If no criteria are provided, all
// field of criteria are their zero-values, a view that matches no instruments
// is returned.
//
// The Name field of criteria supports wildcard pattern matching. The wildcard
// "*" is recognised as matching zero or more characters, and "?" is recognised
// as matching exactly one character. For example, a pattern of "*" will match
// all instrument names.
//
// The DataStream mask only applies updates for non-zero-value fields. By
// default, the InstrumentProperties the View matches against will be use for
// the returned DataStream and no Aggregation or AttributeFilter are set. If
// mask has a non-zero-value value for any of the Aggregation or
// AttributeFilter fields, or any of the InstrumentProperties fields, that
// value is used instead of the default. If you need to zero out an DataStream
// field returned from a View, create a View directly.
func NewView(criteria InstrumentProperties, mask DataStream) View {
	if criteria == zeroInstrumentProperties {
		return func(InstrumentProperties) (DataStream, bool) {
			return DataStream{}, false
		}
	}

	var matchFunc func(InstrumentProperties) bool
	if strings.ContainsAny(criteria.Name, "*?") {
		pattern := regexp.QuoteMeta(criteria.Name)
		pattern = "^" + pattern + "$"
		pattern = strings.ReplaceAll(pattern, "\\?", ".")
		pattern = strings.ReplaceAll(pattern, "\\*", ".*")
		re := regexp.MustCompile(pattern)
		matchFunc = func(p InstrumentProperties) bool {
			return re.MatchString(p.Name) &&
				criteria.matchesDescription(p) &&
				criteria.matchesKind(p) &&
				criteria.matchesUnit(p) &&
				criteria.matchesScope(p)
		}
	} else {
		matchFunc = criteria.matches
	}

	var agg aggregation.Aggregation
	if mask.Aggregation != nil {
		agg = mask.Aggregation.Copy()
		if err := agg.Err(); err != nil {
			global.Error(
				err, "not using aggregation with view",
				"aggregation", agg,
				"view", criteria,
			)
			agg = nil
		}
	}

	return func(p InstrumentProperties) (DataStream, bool) {
		if matchFunc(p) {
			stream := DataStream{
				InstrumentProperties: p.mask(mask.InstrumentProperties),
				Aggregation:          agg,
				AttributeFilter:      mask.AttributeFilter,
			}
			return stream, true
		}
		return DataStream{}, false
	}
}
