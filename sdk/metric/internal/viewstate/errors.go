package viewstate

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// ViewConflicts represents a per-reader summary of conflicts caused
// by creating an instrument after applying the view configuration.
// ViewConflicts may contain either or both of a list of duplicates
// and a semantic error.  Typically these conflicts will be the same
// for all readers, however since readers influence the defaults for
// aggregation kind and aggregator configuration, it is possible for
// different conflicts to arise.
//
// Full information about every conflict or error is always returned
// to the caller that registered the instrument along with a valid
// (potentially in-conflict) instrument.
type ViewConflictsBuilder map[string][]Conflict
type ViewConflictsError map[string][]Conflict

var _ error = ViewConflictsError{}

const noConflictsString = "no conflicts"

// Error shows one example Conflict with a summary of how many
// conflicts and readers experienced conflicts, in case of multiple
// readers and/or conflicts.
func (vc ViewConflictsError) Error() string {
	total := 0
	for _, l := range vc {
		total += len(l)
	}
	// These are almost always duplicative, so we print only examples for one Config.
	for rd, conflictsByReader := range vc {
		if len(conflictsByReader) == 0 {
			break
		}
		if len(vc) == 1 {
			if len(conflictsByReader) == 1 {
				return fmt.Sprintf("%v: %v", rd, conflictsByReader[0].Error())
			}
			return fmt.Sprintf("%d conflicts, e.g. %v: %v", total, rd, conflictsByReader[0])
		}
		return fmt.Sprintf("%d conflicts in %d readers, e.g. %v: %v", total, len(vc), rd, conflictsByReader[0])
	}
	return noConflictsString
}

func (ViewConflictsError) Is(err error) bool {
	_, ok := err.(ViewConflictsError)
	return ok
}

func (vc *ViewConflictsBuilder) Add(name string, c Conflict) {
	if *vc == nil {
		*vc = ViewConflictsBuilder{}
	}

	(*vc)[name] = append((*vc)[name], c)
}

func (vc *ViewConflictsBuilder) Combine(other ViewConflictsBuilder) {
	if *vc == nil {
		if len(other) == 0 {
			return
		}
		*vc = ViewConflictsBuilder{}
	}
	for k, v := range other {
		(*vc)[k] = v
	}
}

func (vc *ViewConflictsBuilder) AsError() error {
	if vc == nil || *vc == nil {
		return nil
	}
	return ViewConflictsError(*vc)
}

// Conflict describes both the duplicates instruments found and
// semantic errors caused when registering a new instrument.
type Conflict struct {
	// Duplicates
	Duplicates []Duplicate
	// Semantic will be a SemanticError if there was an instrument
	// vs. aggregation conflict or nil otherwise.
	Semantic error
}

var _ error = Conflict{}

// Duplicate is one of the other matching instruments when duplicate
// instrument registration conflicts arise.
type Duplicate interface {
	// Aggregation is the requested aggregation kind.  (If the
	// original aggregation caused a semantic error, this will
	// have been corrected to the default aggregation.)
	Aggregation() aggregation.Kind
	// Descriptor describes the output of the View (Name,
	// Description, Unit, Number Kind, InstrumentKind).
	Descriptor() sdkinstrument.Descriptor
	// Keys is non-nil with the distinct set of keys.  This uses
	// an attribute.Set where the Key is significant and the Value
	// is a meaningless Int(0), for simplicity.
	Keys() *attribute.Set
	// Config is the aggregator configuration.
	Config() aggregator.Config
	// OriginalName is the original name of the Duplicate
	// instrument before renaming.
	OriginalName() string
}

// Error shows the semantic error if non-nil and a summary of the
// duplicates if any were present.
func (c Conflict) Error() string {
	se := c.semanticError()
	de := c.duplicateError()
	if se == "" {
		return de
	}
	if de == "" {
		return se
	}
	return fmt.Sprintf("%s; %s", se, de)
}

func (c Conflict) semanticError() string {
	if c.Semantic == nil {
		return ""
	}
	return c.Semantic.Error()
}

func (c Conflict) duplicateError() string {
	if len(c.Duplicates) < 2 {
		return ""
	}
	// Note: choose the first and last element of the current conflicts
	// list because they are ordered, and if the conflict in question is
	// new it will be the last item.
	inst1 := c.Duplicates[0]
	inst2 := c.Duplicates[len(c.Duplicates)-1]
	name1 := fullNameString(inst1)
	name2 := renameString(inst2)
	conf1 := shortString(inst1)
	conf2 := shortString(inst2)

	var s strings.Builder
	s.WriteString(name1)

	if conf1 != conf2 {
		s.WriteString(fmt.Sprintf(" conflicts %v, %v%v", conf1, conf2, name2))
	} else if !equalConfigs(inst1.Config(), inst2.Config()) {
		s.WriteString(" has conflicts: different aggregator configuration")
	} else {
		s.WriteString(" has conflicts: different attribute filters")
	}

	if len(c.Duplicates) > 2 {
		s.WriteString(fmt.Sprintf(" and %d more", len(c.Duplicates)-2))
	}
	return s.String()
}

// SemanticError occurs when an instrument is paired with an
// incompatible aggregation.
type SemanticError struct {
	Instrument  sdkinstrument.Kind
	Aggregation aggregation.Kind
}

var _ error = SemanticError{}

func (s SemanticError) Error() string {
	return fmt.Sprintf("%v instrument incompatible with %v aggregation",
		strings.TrimSuffix(s.Instrument.String(), "Kind"),
		strings.TrimSuffix(s.Aggregation.String(), "Kind"),
	)
}

func (SemanticError) Is(err error) bool {
	_, ok := err.(SemanticError)
	return ok
}

// fullNameString helps rendering concise error descriptions by
// showing the original name only when it is different.
func fullNameString(d Duplicate) string {
	return fmt.Sprintf("name %q%v", d.Descriptor().Name, renameString(d))
}

// renameString is the fragment used by fullNameString when the
// original name is different than the output name.
func renameString(d Duplicate) string {
	if d.OriginalName() == d.Descriptor().Name {
		return ""
	}
	return fmt.Sprintf(" (original %q)", d.OriginalName())
}

// shortString concatenates the instrument kind, number kind,
// aggregation kind, and unit to summarize most of the potentially
// conflicting characteristics of an instrument.
func shortString(d Duplicate) string {
	s := fmt.Sprintf("%v-%v-%v",
		strings.TrimSuffix(d.Descriptor().Kind.String(), "Kind"),
		strings.TrimSuffix(d.Descriptor().NumberKind.String(), "Kind"),
		strings.TrimSuffix(d.Aggregation().String(), "Kind"),
	)
	if d.Descriptor().Unit != "" {
		s = fmt.Sprintf("%v-%v", s, d.Descriptor().Unit)
	}
	return s
}
