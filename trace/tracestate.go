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

package trace

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	maxListMembers = 32

	listDelimiter = ","

	// based on the W3C Trace Context specification, see
	// https://www.w3.org/TR/trace-context-1/#tracestate-header
	noTenantKeyFormat   = `[a-z][_0-9a-z\-\*\/]{0,255}`
	withTenantKeyFormat = `[a-z0-9][_0-9a-z\-\*\/]{0,240}@[a-z][_0-9a-z\-\*\/]{0,13}`
	valueFormat         = `[\x20-\x2b\x2d-\x3c\x3e-\x7e]{0,255}[\x21-\x2b\x2d-\x3c\x3e-\x7e]`

	keyRe    = regexp.MustCompile(`^((` + noTenantKeyFormat + `)|(` + withTenantKeyFormat + `))$`)
	valueRe  = regexp.MustCompile(`^(` + valueFormat + `)$`)
	memberRe = regexp.MustCompile(`^\s*((` + noTenantKeyFormat + `)|(` + withTenantKeyFormat + `))=(` + valueFormat + `)\s*$`)

	errInvalidKey    errorConst = "invalid tracestate key"
	errInvalidValue  errorConst = "invalid tracestate value"
	errInvalidMember errorConst = "invalid tracestate list member"
	errMemberNumber  errorConst = "too many list members in tracestate"
	errDuplicate     errorConst = "duplicate list member in tracestate"
)

type member struct {
	Key   string
	Value string
}

func newMember(key, value string) (member, error) {
	if !keyRe.MatchString(key) {
		return member{}, fmt.Errorf("%w: %s", errInvalidKey, key)
	}
	if !valueRe.MatchString(value) {
		return member{}, fmt.Errorf("%w: %s", errInvalidValue, value)
	}
	return member{Key: key, Value: value}, nil
}

func parseMemeber(m string) (member, error) {
	matches := memberRe.FindStringSubmatch(m)
	if len(matches) != 5 {
		return member{}, fmt.Errorf("%w: %s", errInvalidMember, m)
	}

	return member{
		Key:   matches[1],
		Value: matches[4],
	}, nil

}

// String encodes member into a string compliant with the W3C tracecontext
// specification.
func (m member) String() string {
	return fmt.Sprintf("%s=%s", m.Key, m.Value)
}

// TraceState provides additional vendor-specific trace identification
// information across different distributed tracing systems. It represents an
// immutable list consisting of key/value pairs, each pair is referred to as a
// member. There can be a maximum of 32 list members.
//
// The key and value of each list member must be valid according to the W3C
// Trace Context specification (see https://www.w3.org/TR/trace-context-1/#key
// and https://www.w3.org/TR/trace-context-1/#value respectively).
//
// Trace state must be valid according to the W3C Trace Context specification
// at all times. All mutating operations validate their input and, in case of
// valid parameters, return a new TraceState.
type TraceState struct { //nolint:golint
	members []member
}

var _ json.Marshaler = TraceState{}

// ParseTraceState attempts to decode a TraceState from the passed
// string. It returns an error if the input is invalid according to the W3C
// tracecontext specification.
func ParseTraceState(tracestate string) (TraceState, error) {
	if tracestate == "" {
		return TraceState{}, nil
	}

	wrapErr := func(err error) error {
		return fmt.Errorf("failed to parse tracestate: %w", err)
	}

	var members []member
	found := make(map[string]struct{})
	for _, memberStr := range strings.Split(tracestate, listDelimiter) {
		if len(memberStr) == 0 {
			continue
		}

		m, err := parseMemeber(memberStr)
		if err != nil {
			return TraceState{}, wrapErr(err)
		}

		if _, ok := found[m.Key]; ok {
			return TraceState{}, wrapErr(errDuplicate)
		}
		found[m.Key] = struct{}{}

		members = append(members, m)
		if n := len(members); n > maxListMembers {
			return TraceState{}, wrapErr(errMemberNumber)
		}
	}

	return TraceState{members: members}, nil
}

// MarshalJSON marshals the TraceState into JSON.
func (ts TraceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.String())
}

// String encodes the TraceState into a string compliant with the W3C
// tracecontext specification. The returned string will be invalid if the
// TraceState contains any invalid members.
func (ts TraceState) String() string {
	members := make([]string, len(ts.members))
	for i, m := range ts.members {
		members[i] = m.String()
	}
	return strings.Join(members, listDelimiter)
}

// Get returns a value for given key from the trace state.
// If no key is found or provided key is invalid, returns an empty string.
func (ts TraceState) Get(key string) string {
	if !keyRe.MatchString(key) {
		return ""
	}

	for _, member := range ts.members {
		if member.Key == key {
			return member.Value
		}
	}

	return ""
}

// Insert adds a new key/value pair, if one doesn't exists; otherwise updates the existing entry.
// The new or updated entry is always inserted at the beginning of the TraceState, i.e.
// on the left side, as per the W3C Trace Context specification requirement.
func (ts TraceState) Insert(key, value string) (TraceState, error) {
	m, err := newMember(key, value)
	if err != nil {
		return ts, err
	}

	cTS := ts.remove(key)
	if cTS.Len()+1 > maxListMembers {
		return ts, fmt.Errorf("failed to insert: %w", errMemberNumber)
	}

	cTS.members = append(cTS.members, member{})
	copy(cTS.members[1:], cTS.members)
	cTS.members[0] = m

	return cTS, nil
}

// Delete removes specified entry from the trace state.
func (ts TraceState) Delete(key string) (TraceState, error) {
	if !keyRe.MatchString(key) {
		return ts, fmt.Errorf("%w: %s", errInvalidKey, key)
	}

	return ts.remove(key), nil
}

// Len returns the number of list-members in the TraceState.
func (ts TraceState) Len() int {
	return len(ts.members)
}

// remove returns a copy of ts with key removed, if it exists.
func (ts TraceState) remove(key string) TraceState {
	// allocate the same size in case key is not contained in ts.
	members := make([]member, ts.Len())
	copy(members, ts.members)

	for i, member := range ts.members {
		if member.Key == key {
			members = append(members[:i], members[i+1:]...)
			break
		}
	}

	return TraceState{members: members}
}
