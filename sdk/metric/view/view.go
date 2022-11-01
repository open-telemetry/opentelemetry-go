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

package view // import "go.opentelemetry.io/otel/sdk/metric/view"

import (
	"fmt"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

// View provides users with the flexibility to customize the metrics that are
// output by the SDK. A View can be used to ignore, change the name,
// description, and aggregation of, and customize which attribute(s) are to be
// reported by Instruments.
//
// An empty View will match all instruments, and do no transformations.
type View struct {
	instrumentName *regexp.Regexp
	hasWildcard    bool
	scope          instrumentation.Scope
	instrumentKind InstrumentKind

	filter      attribute.Filter
	name        string
	description string
	agg         aggregation.Aggregation
}

// New returns a new configured View. If there are any duplicate Options passed,
// the last one passed will take precedence. The unique, de-duplicated,
// Options are all applied to the View. An instrument needs to match all of
// the match Options passed for the View to be applied to it. Similarly, all
// transform operation Options are applied to matched Instruments.
func New(opts ...Option) (View, error) {
	v := View{}

	for _, opt := range opts {
		v = opt.apply(v)
	}

	emptyScope := instrumentation.Scope{}
	if v.instrumentName == nil &&
		v.scope == emptyScope &&
		v.instrumentKind == undefinedInstrument {
		return View{}, fmt.Errorf("must provide at least 1 match option")
	}

	if v.hasWildcard && v.name != "" {
		return View{}, fmt.Errorf("invalid view: view name specified for multiple instruments")
	}

	return v, nil
}

// TransformInstrument will check if an instrument matches this view
// and will convert it if it does.
func (v View) TransformInstrument(inst Instrument) (transformed Instrument, match bool) {
	if !v.match(inst) {
		return Instrument{}, false
	}
	if v.name != "" {
		inst.Name = v.name
	}
	if v.description != "" {
		inst.Description = v.description
	}
	if v.agg != nil {
		inst.Aggregation = v.agg
	}
	return inst, true
}

// AttributeFilter returns a function that returns only attributes specified by
// WithFilterAttributes. If no filter was provided nil is returned.
func (v View) AttributeFilter() func(attribute.Set) attribute.Set {
	if v.filter == nil {
		return nil
	}
	return func(input attribute.Set) attribute.Set {
		out, _ := input.Filter(v.filter)
		return out
	}
}

func (v View) matchName(name string) bool {
	return v.instrumentName == nil || v.instrumentName.MatchString(name)
}

func (v View) matchScopeName(name string) bool {
	return v.scope.Name == "" || name == v.scope.Name
}

func (v View) matchScopeVersion(version string) bool {
	return v.scope.Version == "" || version == v.scope.Version
}

func (v View) matchScopeSchemaURL(schemaURL string) bool {
	return v.scope.SchemaURL == "" || schemaURL == v.scope.SchemaURL
}

func (v View) matchInstrumentKind(kind InstrumentKind) bool {
	return v.instrumentKind == undefinedInstrument || kind == v.instrumentKind
}

func (v View) match(i Instrument) bool {
	return v.matchName(i.Name) &&
		v.matchScopeName(i.Scope.Name) &&
		v.matchScopeSchemaURL(i.Scope.SchemaURL) &&
		v.matchScopeVersion(i.Scope.Version) &&
		v.matchInstrumentKind(i.Kind)
}

// Option applies a configuration option value to a View.
type Option interface {
	apply(View) View
}

type optionFunc func(View) View

func (f optionFunc) apply(v View) View {
	return f(v)
}

// MatchInstrumentName will match an instrument based on the its name.
// This will accept wildcards of * for zero or more characters, and ? for
// exactly one character. A name of "*" (default) will match all instruments.
func MatchInstrumentName(name string) Option {
	return optionFunc(func(v View) View {
		if strings.ContainsAny(name, "*?") {
			v.hasWildcard = true
		}
		name = regexp.QuoteMeta(name)
		name = "^" + name + "$"
		name = strings.ReplaceAll(name, "\\?", ".")
		name = strings.ReplaceAll(name, "\\*", ".*")
		v.instrumentName = regexp.MustCompile(name)
		return v
	})
}

// MatchInstrumentKind with match an instrument based on the instrument's kind.
// The default is to match all instrument kinds.
func MatchInstrumentKind(kind InstrumentKind) Option {
	return optionFunc(func(v View) View {
		v.instrumentKind = kind
		return v
	})
}

// MatchInstrumentationScope will do an exact match on any
// instrumentation.Scope field that is non-empty (""). The default is to match all
// instrumentation scopes.
func MatchInstrumentationScope(scope instrumentation.Scope) Option {
	return optionFunc(func(v View) View {
		v.scope = scope
		return v
	})
}

// WithRename will rename the instrument the view matches. If not used or empty the
// instrument name will not be changed. Must be used with a non-wildcard
// instrument name match. The default does not change the instrument name.
func WithRename(name string) Option {
	return optionFunc(func(v View) View {
		v.name = name
		return v
	})
}

// WithSetDescription will change the description of the instruments the view
// matches to desc. If not used or empty the description will not be changed.
func WithSetDescription(desc string) Option {
	return optionFunc(func(v View) View {
		v.description = desc
		return v
	})
}

// WithFilterAttributes will select attributes that have a matching key.  If not used
// or empty no filter will be applied.
func WithFilterAttributes(keys ...attribute.Key) Option {
	return optionFunc(func(v View) View {
		if len(keys) == 0 {
			return v
		}
		filterKeys := map[attribute.Key]struct{}{}
		for _, key := range keys {
			filterKeys[key] = struct{}{}
		}

		v.filter = attribute.Filter(func(kv attribute.KeyValue) bool {
			_, ok := filterKeys[kv.Key]
			return ok
		})
		return v
	})
}

// WithSetAggregation will use the aggregation a for matching instruments. If
// this option is not provided, the reader defined aggregation for the
// instrument will be used.
//
// If a is misconfigured, it will not be used and an error will be logged.
func WithSetAggregation(a aggregation.Aggregation) Option {
	cpA := a.Copy()
	if err := cpA.Err(); err != nil {
		global.Error(err, "not using aggregation with view", "aggregation", a)
		return optionFunc(func(v View) View { return v })
	}

	return optionFunc(func(v View) View {
		v.agg = cpA
		return v
	})
}
