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

package trace // import "go.opentelemetry.io/otel/trace"

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	maxListMembers = 32

	listDelimiters  = ","
	memberDelimiter = "="

	errInvalidKey    errorConst = "invalid tracestate key"
	errInvalidValue  errorConst = "invalid tracestate value"
	errInvalidMember errorConst = "invalid tracestate list-member"
	errMemberNumber  errorConst = "too many list-members in tracestate"
	errDuplicate     errorConst = "duplicate list-member in tracestate"
)

type member struct {
	Key   string
	Value string
}

// [\x20-\x2b\x2d-\x3c\x3e-\x7e]*
func checkValueChar(v byte) bool {
	return v >= '\x20' && v <= '\x7e' && v != '\x2c' && v != '\x3d'
}

// [\x21-\x2b\x2d-\x3c\x3e-\x7e]
func checkValueLast(v byte) bool {
	return v >= '\x21' && v <= '\x7e' && v != '\x2c' && v != '\x3d'
}

func checkValue(val string) bool {
	n := len(val)
	if n == 0 || n > 256 {
		return false
	}
	// valueFormat         = `[\x20-\x2b\x2d-\x3c\x3e-\x7e]{0,255}[\x21-\x2b\x2d-\x3c\x3e-\x7e]`
	for i := 0; i < n-1; i++ {
		if !checkValueChar(val[i]) {
			return false
		}
	}
	return checkValueLast(val[n-1])
}

// [_0-9a-z\-\*\/]*
func checkKeyRemain(key string) bool {
	for _, v := range key {
		if (v >= '0' && v <= '9') || (v >= 'a' && v <= 'z') {
			continue
		}
		switch v {
		case '_', '-', '*', '/':
			continue
		}
		return false
	}
	return true
}

func checkKeyPart(key string, n int, tenant bool) bool {
	if len(key) == 0 {
		return false
	}
	first := key[0] // key first char
	ret := len(key[1:]) <= n
	if tenant {
		// [a-z0-9]
		ret = ret && ((first >= 'a' && first <= 'z') || (first >= '0' && first <= '9'))
	} else {
		// [a-z]
		ret = ret && first >= 'a' && first <= 'z'
	}
	return ret && checkKeyRemain(key[1:])
}

func checkKey(key string) bool {
	// noTenantKeyFormat   = `[a-z][_0-9a-z\-\*\/]{0,255}`
	// withTenantKeyFormat = `[a-z0-9][_0-9a-z\-\*\/]{0,240}@[a-z][_0-9a-z\-\*\/]{0,13}`
	tenant, system, ok := strings.Cut(key, "@")
	if !ok {
		return checkKeyPart(key, 255, false)
	}
	return checkKeyPart(tenant, 240, true) && checkKeyPart(system, 13, false)
}

// based on the W3C Trace Context specification, see
// https://www.w3.org/TR/trace-context-1/#tracestate-header
func newMember(key, value string) (member, error) {
	if !checkKey(key) {
		return member{}, fmt.Errorf("%w: %s", errInvalidKey, key)
	}
	if !checkValue(value) {
		return member{}, fmt.Errorf("%w: %s", errInvalidValue, value)
	}
	return member{Key: key, Value: value}, nil
}

func parseMember(m string) (member, error) {
	key, val, ok := strings.Cut(m, memberDelimiter)
	if !ok {
		return member{}, fmt.Errorf("%w: %s", errInvalidMember, m)
	}
	key = strings.TrimLeft(key, " \t")
	val = strings.TrimRight(val, " \t")
	result, e := newMember(key, val)
	if e != nil {
		return member{}, fmt.Errorf("%w: %s", errInvalidMember, m)
	}
	return result, nil
}

// String encodes member into a string compliant with the W3C Trace Context
// specification.
func (m member) String() string {
	return m.Key + "=" + m.Value
}

// TraceState provides additional vendor-specific trace identification
// information across different distributed tracing systems. It represents an
// immutable list consisting of key/value pairs, each pair is referred to as a
// list-member.
//
// TraceState conforms to the W3C Trace Context specification
// (https://www.w3.org/TR/trace-context-1). All operations that create or copy
// a TraceState do so by validating all input and will only produce TraceState
// that conform to the specification. Specifically, this means that all
// list-member's key/value pairs are valid, no duplicate list-members exist,
// and the maximum number of list-members (32) is not exceeded.
type TraceState struct { //nolint:revive // revive complains about stutter of `trace.TraceState`
	// list is the members in order.
	list []member
}

var _ json.Marshaler = TraceState{}

// ParseTraceState attempts to decode a TraceState from the passed
// string. It returns an error if the input is invalid according to the W3C
// Trace Context specification.
func ParseTraceState(ts string) (TraceState, error) {
	if ts == "" {
		return TraceState{}, nil
	}

	wrapErr := func(err error) error {
		return fmt.Errorf("failed to parse tracestate: %w", err)
	}

	var members []member
	found := make(map[string]struct{})
	for ts != "" {
		var memberStr string
		memberStr, ts, _ = strings.Cut(ts, listDelimiters)
		if len(memberStr) == 0 {
			continue
		}

		m, err := parseMember(memberStr)
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

	return TraceState{list: members}, nil
}

// MarshalJSON marshals the TraceState into JSON.
func (ts TraceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.String())
}

// String encodes the TraceState into a string compliant with the W3C
// Trace Context specification. The returned string will be invalid if the
// TraceState contains any invalid members.
func (ts TraceState) String() string {
	if len(ts.list) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(ts.list[0].Key)
	sb.WriteByte('=')
	sb.WriteString(ts.list[0].Value)
	for i := 1; i < len(ts.list); i++ {
		sb.WriteByte(listDelimiters[0])
		sb.WriteString(ts.list[i].Key)
		sb.WriteByte('=')
		sb.WriteString(ts.list[i].Value)
	}
	return sb.String()
}

// Get returns the value paired with key from the corresponding TraceState
// list-member if it exists, otherwise an empty string is returned.
func (ts TraceState) Get(key string) string {
	for _, member := range ts.list {
		if member.Key == key {
			return member.Value
		}
	}

	return ""
}

// Insert adds a new list-member defined by the key/value pair to the
// TraceState. If a list-member already exists for the given key, that
// list-member's value is updated. The new or updated list-member is always
// moved to the beginning of the TraceState as specified by the W3C Trace
// Context specification.
//
// If key or value are invalid according to the W3C Trace Context
// specification an error is returned with the original TraceState.
//
// If adding a new list-member means the TraceState would have more members
// then is allowed, the new list-member will be inserted and the right-most
// list-member will be dropped in the returned TraceState.
func (ts TraceState) Insert(key, value string) (TraceState, error) {
	m, err := newMember(key, value)
	if err != nil {
		return ts, err
	}
	n := len(ts.list)
	found := n
	for i := range ts.list {
		if ts.list[i].Key == key {
			found = i
		}
	}
	cTS := TraceState{}
	if found == n && n < maxListMembers {
		cTS.list = make([]member, n+1)
	} else {
		cTS.list = make([]member, n)
	}
	cTS.list[0] = m
	// When the number of members exceeds capacity, drop the "right-most".
	copy(cTS.list[1:], ts.list[0:found])
	if found < n {
		copy(cTS.list[1+found:], ts.list[found+1:])
	}
	return cTS, nil
}

// Delete returns a copy of the TraceState with the list-member identified by
// key removed.
func (ts TraceState) Delete(key string) TraceState {
	members := make([]member, ts.Len())
	copy(members, ts.list)
	for i, member := range ts.list {
		if member.Key == key {
			members = append(members[:i], members[i+1:]...)
			// TraceState should contain no duplicate members.
			break
		}
	}
	return TraceState{list: members}
}

// Len returns the number of list-members in the TraceState.
func (ts TraceState) Len() int {
	return len(ts.list)
}
