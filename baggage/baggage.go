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

package baggage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	maxMembers               = 180
	maxBytesPerMembers       = 4096
	maxBytesPerBaggageString = 8192

	listDelimiter     = ","
	propertyDelimiter = ";"

	keyDef      = `([\x21\x23-\x27\x2A\x2B\x2D\x2E\x30-\x39\x41-\x5a\x5e-\x7a\x7c\x7e]+)`
	valueDef    = `([\x21\x23-\x2b\x2d-\x3a\x3c-\x5B\x5D-\x7e]*)`
	keyValueDef = `\s*` + keyDef + `\s*=\s*` + valueDef + `\s*`
)

var (
	keyRe      = regexp.MustCompile(`^` + keyDef + `$`)
	valueRe    = regexp.MustCompile(`^` + valueDef + `$`)
	keyValueRe = regexp.MustCompile(`^` + keyValueDef + `$`)
	propertyRe = regexp.MustCompile(`^(?:\s*` + keyDef + `\s*|` + keyValueDef + `)$`)
)

var (
	errInvalidProperty = errors.New("invalid baggage property")
	errInvalidMember   = errors.New("invalid baggage list-member")
	errMemberNumber    = errors.New("too many list-members in baggage-string")
	errMemberBytes     = errors.New("list-member too large")
	errBaggageBytes    = errors.New("baggage-string too large")
)

// Property is an additional metadata entry for a baggage list-member.
type Property struct {
	key, value string

	// hasValue indicates if a zero-value value means the property does not
	// have a value or if it was the zero-value.
	hasValue bool
}

func NewKeyProperty(key string) (Property, error) {
	p := Property{}
	if !keyRe.MatchString(key) {
		return p, fmt.Errorf("invalid key: %q", key)
	}
	p.key = key
	return p, nil
}

func NewKeyValueProperty(key, value string) (Property, error) {
	p := Property{}
	if !keyRe.MatchString(key) {
		return p, fmt.Errorf("invalid key: %q", key)
	}
	if !valueRe.MatchString(value) {
		return p, fmt.Errorf("invalid value: %q", value)
	}
	p.key = key
	p.value = value
	p.hasValue = true
	return p, nil
}

// parseProperty attempts to decode a Property from the passed string. It
// returns an error if the input is invalid according to the W3C Baggage
// specification.
func parseProperty(property string) (Property, error) {
	p := Property{}
	if property == "" {
		return p, nil
	}

	match := propertyRe.FindStringSubmatch(property)
	if len(match) != 4 {
		return p, fmt.Errorf("%w: %q", errInvalidProperty, property)
	}

	if match[1] != "" {
		p.key = match[1]
	} else {
		p.key = match[2]
		p.value = match[3]
		p.hasValue = true
	}
	return p, nil
}

// validate ensures p conforms to the W3C Baggage specification, returning an
// error otherwise.
func (p Property) validate() error {
	if !keyRe.MatchString(p.key) {
		return fmt.Errorf("%w: invalid key: %q", errInvalidProperty, p.key)
	}
	if p.hasValue && !valueRe.MatchString(p.value) {
		return fmt.Errorf("%w: invalid value: %q", errInvalidProperty, p.value)
	}
	if !p.hasValue && p.value != "" {
		return fmt.Errorf("%w: inconsistent value", errInvalidProperty)
	}
	return nil
}

// Key returns the Property key.
func (p Property) Key() string {
	return p.key
}

// Value returns the Property value. Additionally a boolean value is returned
// indicating if the returned value is the empty if the Property has a value
// that is empty or if the value is not set.
func (p Property) Value() (string, bool) {
	return p.value, p.hasValue
}

// String encodes Property into a string compliant with the W3C Baggage
// specification.
func (p Property) String() string {
	if p.hasValue {
		return fmt.Sprintf("%s=%v", p.key, p.value)
	}
	return p.key
}

type properties []Property

func (p properties) Copy() properties {
	props := make(properties, len(p))
	copy(props, p)
	return props
}

// validate ensures each Property in p conforms to the W3C Baggage
// specification, returning an error otherwise.
func (p properties) validate() error {
	for _, prop := range p {
		if err := prop.validate(); err != nil {
			return err
		}
	}
	return nil
}

// String encodes properties into a string compliant with the W3C Baggage
// specification.
func (p properties) String() string {
	props := make([]string, len(p))
	for i, prop := range p {
		props[i] = prop.String()
	}
	return strings.Join(props, propertyDelimiter)
}

// Member is a list-member of a baggage-string as defined by the W3C Baggage
// specification.
type Member struct {
	key, value string
	properties properties
}

func NewMember(key, value string, props ...Property) (Member, error) {
	m := Member{}
	if !keyRe.MatchString(key) {
		return m, fmt.Errorf("invalid key: %q", key)
	}
	if !valueRe.MatchString(value) {
		return m, fmt.Errorf("invalid value: %q", value)
	}
	p := properties(props)
	if err := p.validate(); err != nil {
		return m, err
	}

	return Member{key: key, value: value, properties: p.Copy()}, nil
}

// parseMember attempts to decode a Member from the passed string. It returns
// an error if the input is invalid according to the W3C Baggage
// specification.
func parseMember(member string) (Member, error) {
	if n := len(member); n > maxBytesPerMembers {
		return Member{}, fmt.Errorf("%w: %d", errBaggageBytes, n)
	}

	var (
		key, value string
		props      properties
	)

	parts := strings.SplitN(member, propertyDelimiter, 2)
	switch len(parts) {
	case 2:
		for _, pStr := range strings.Split(parts[1], propertyDelimiter) {
			p, err := parseProperty(pStr)
			if err != nil {
				return Member{}, err
			}
			props = append(props, p)
		}
		fallthrough
	case 1:
		match := keyValueRe.FindStringSubmatch(parts[0])
		if len(match) != 3 {
			return Member{}, fmt.Errorf("%w: %q", errInvalidMember, member)
		}
		key, value = match[1], match[2]
	default:
		return Member{}, fmt.Errorf("%w: %q", errInvalidMember, member)
	}

	return Member{key: key, value: value, properties: props}, nil
}

// validate ensures m conforms to the W3C Baggage specification, returning an
// error otherwise.
func (m Member) validate() error {
	if !keyRe.MatchString(m.key) {
		return fmt.Errorf("%w: invalid key: %q", errInvalidMember, m.key)
	}
	if !valueRe.MatchString(m.value) {
		return fmt.Errorf("%w: invalid value: %q", errInvalidMember, m.value)
	}
	return m.properties.validate()
}

// Key returns the Member key.
func (m Member) Key() string { return m.key }

// Value returns the Member value.
func (m Member) Value() string { return m.value }

// Properties returns a copy of the Member properties.
func (m Member) Properties() []Property { return m.properties.Copy() }

// String encodes Member into a string compliant with the W3C Baggage
// specification.
func (m Member) String() string {
	s := fmt.Sprintf("%s=%s", m.key, m.value)
	if len(m.properties) > 0 {
		s = fmt.Sprintf("%s%s%s", s, propertyDelimiter, m.properties.String())
	}
	return s
}

type value struct {
	value      string
	properties properties
}

// Baggage is a list of baggage members representing the baggage-string as
// defined by the W3C Baggage specification.
type Baggage struct { //nolint:golint
	list map[string]value
}

// Parse attempts to decode a baggage-string from the passed string. It
// returns an error if the input is invalid according to the W3C Baggage
// specification.
//
// If there are duplicate list-members contained in baggage, the last one
// defined (reading left-to-right) will be the only one kept. This diverges
// from the W3C Baggage specification which allows duplicate list-members, but
// conforms to the OpenTelemetry Baggage specification.
func Parse(baggage string) (Baggage, error) {
	var b Baggage
	if baggage == "" {
		return b, nil
	}

	if n := len(baggage); n > maxBytesPerBaggageString {
		return b, fmt.Errorf("%w: %d", errBaggageBytes, n)
	}

	for _, memberStr := range strings.Split(baggage, listDelimiter) {
		// Trim empty baggage members.
		if len(memberStr) == 0 {
			continue
		}

		m, err := parseMember(memberStr)
		if err != nil {
			return Baggage{}, err
		}

		b.list[m.key] = value{m.value, m.properties}
		if n := len(b.list); n > maxMembers {
			return Baggage{}, errMemberNumber
		}
	}

	return b, nil
}

// Member returns the baggage list-member identified by key and a boolean
// value indicating if the list-member existed or not.
func (b Baggage) Member(key string) (Member, bool) {
	v, ok := b.list[key]
	return Member{
		key:        key,
		value:      v.value,
		properties: v.properties.Copy(),
	}, ok
}

// Members returns all the baggage list-member.
// The order of the returned list-members does not have significance.
func (b Baggage) Members() []Member {
	members := make([]Member, len(b.list))
	for key := range b.list {
		// No need to verify the key exists.
		m, _ := b.Member(key)
		members = append(members, m)
	}
	return members
}

// SetMember returns a copy the Baggage with the member included. If the
// baggage contains a Member with the same key the existing Member is
// replaced.
//
// If member is invalid according to the W3C Baggage specification, an error
// is returned with the original Baggage.
func (b Baggage) SetMember(member Member) (Baggage, error) {
	if err := member.validate(); err != nil {
		return b, err
	}

	n := len(b.list)
	if _, ok := b.list[member.key]; !ok {
		n++
	}
	list := make(map[string]value, n)

	for k, v := range b.list {
		// Update instead of just copy and overwrite.
		if k == member.key {
			list[member.key] = value{
				value:      member.value,
				properties: member.properties.Copy(),
			}
			continue
		}
		list[k] = v
	}

	return Baggage{list: list}, nil
}

// DeleteMember returns a copy of the Baggage with the list-member identified
// by key removed.
func (b Baggage) DeleteMember(key string) Baggage {
	n := len(b.list)
	if _, ok := b.list[key]; ok {
		n--
	}
	list := make(map[string]value, n)

	for k, v := range b.list {
		if k == key {
			continue
		}
		list[k] = v
	}

	return Baggage{list: list}
}

// Len returns the number of list-members in the Baggage.
func (b Baggage) Len() int {
	return len(b.list)
}

// String encodes Baggage into a string compliant with the W3C Baggage
// specification. The returned string will be invalid if the Baggage contains
// any invalid list-members.
func (b Baggage) String() string {
	members := make([]string, 0, len(b.list))
	for k, v := range b.list {
		members = append(members, Member{
			key:        k,
			value:      v.value,
			properties: v.properties,
		}.String())
	}
	return strings.Join(members, listDelimiter)
}
