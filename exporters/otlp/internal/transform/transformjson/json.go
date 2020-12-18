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

// TODO:
//
// - proper error reporting (better error messages)
//
// - tests

package transformjson

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func Marshal(message proto.Message) ([]byte, error) {
	return marshal(message)
}

func Unmarshal(data []byte, message proto.Message) error {
	return unmarshal(data, message)
}

type result int

const (
	resultOK result = iota
	resultFailed
	resultFinished
)

type fixupType int

const (
	fixupToHex fixupType = iota
	fixupToBase64
)

func marshal(request proto.Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	marshaler := jsonpb.Marshaler{}
	if err := marshaler.Marshal(buffer, request); err != nil {
		return nil, err
	}
	return fixupJSON(buffer, fixupToHex)
}

func unmarshal(data []byte, request proto.Message) error {
	buffer := bytes.NewBuffer(data)
	revertedJSON, err := fixupJSON(buffer, fixupToBase64)
	if err != nil {
		return err
	}
	unmarshaler := jsonpb.Unmarshaler{}
	return unmarshaler.Unmarshal(bytes.NewBuffer(revertedJSON), request)
}

func fixupJSON(dataReader io.Reader, fixup fixupType) ([]byte, error) {
	state := newState(fixup)
	defer state.clean()
	decoder := json.NewDecoder(dataReader)
	for {
		processor := state.topProcessor()
		switch processor.process(decoder.Token()) {
		case resultOK:
			// continue processing
		case resultFailed:
			return nil, state.err
		case resultFinished:
			return json.Marshal(state.fixedMap)
		default:
			return nil, errors.New("invalid processor status")
		}
	}
}

type processor interface {
	process(token json.Token, err error) result
}

type jsonMap map[string]interface{}

type state struct {
	// self reference
	s          *state
	fixedMap   jsonMap
	processors []processor
	extra      []interface{}
	err        error
	tempOPS    []*objectProcessorState
	tempAPS    []*arrayProcessorState
	fixup      fixupType
}

func newState(fixup fixupType) *state {
	s := &state{
		fixup: fixup,
	}
	s.push((*toplevelProcessor)(s), nil)
	s.s = s
	return s
}

func (s *state) clean() {
	s.s = nil
}

func (s *state) fail(err error) result {
	s.err = err
	return resultFailed
}

func (s *state) pushObject(currentMap jsonMap) {
	proc := (*objectProcessor)(s)
	extra := s.getOPS(currentMap)
	s.push(proc, extra)
}

func (s *state) getOPS(currentMap jsonMap) *objectProcessorState {
	if len(s.tempOPS) > 0 {
		ops := s.tempOPS[len(s.tempOPS)-1]
		s.tempOPS[len(s.tempOPS)-1] = nil
		s.tempOPS = s.tempOPS[:len(s.tempOPS)-1]
		if len(s.tempOPS) == 0 {
			s.tempOPS = nil
		}
		ops.init(currentMap)
		return ops
	}
	ops := &objectProcessorState{}
	ops.init(currentMap)
	return ops
}

func (s *state) pushArray(slicePtr *[]interface{}) {
	proc := (*arrayProcessor)(s)
	extra := s.getAPS(slicePtr)
	s.push(proc, extra)
}

func (s *state) getAPS(slicePtr *[]interface{}) *arrayProcessorState {
	if len(s.tempAPS) > 0 {
		aps := s.tempAPS[len(s.tempAPS)-1]
		s.tempAPS[len(s.tempAPS)-1] = nil
		s.tempAPS = s.tempAPS[:len(s.tempAPS)-1]
		if len(s.tempAPS) == 0 {
			s.tempAPS = nil
		}
		aps.init(slicePtr)
		return aps
	}
	aps := &arrayProcessorState{}
	aps.init(slicePtr)
	return aps
}

func (s *state) push(proc processor, extra interface{}) {
	s.processors = append(s.processors, proc)
	s.extra = append(s.extra, extra)
}

func (s *state) topProcessor() processor {
	return s.processors[len(s.processors)-1]
}

func (s *state) topExtra() interface{} {
	return s.extra[len(s.extra)-1]
}

func (s *state) pop() {
	e := s.extra[len(s.extra)-1]
	switch v := e.(type) {
	case *objectProcessorState:
		s.tempOPS = append(s.tempOPS, v)
	case *arrayProcessorState:
		s.tempAPS = append(s.tempAPS, v)
	}

	s.processors[len(s.processors)-1] = nil
	s.processors = s.processors[:len(s.processors)-1]

	s.extra[len(s.extra)-1] = nil
	s.extra = s.extra[:len(s.extra)-1]

}

type toplevelProcessor state

var _ processor = (*toplevelProcessor)(nil)

func (p *toplevelProcessor) process(token json.Token, err error) result {
	if err == io.EOF {
		return resultFinished
	}
	if err != nil {
		return p.s.fail(err)
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return p.s.fail(errors.New("malformed JSON"))
	}
	switch delim.String() {
	case "{":
		p.fixedMap = make(jsonMap)
		p.s.pushObject(p.fixedMap)
	default:
		return p.s.fail(errors.New("unexpected closing delimiter"))
	}
	return resultOK
}

type objectTokenState int

const (
	key objectTokenState = iota
	value
)

type objectProcessorState struct {
	tokenState objectTokenState
	currentMap jsonMap
	currentKey string
	// for array values
	tempSlice []interface{}
}

func (s *objectProcessorState) init(currentMap jsonMap) {
	s.tokenState = key
	s.currentMap = currentMap
	s.currentKey = ""
	s.tempSlice = nil
}

type objectProcessor state

var _ processor = (*objectProcessor)(nil)

func (p *objectProcessor) process(token json.Token, err error) result {
	if err == io.EOF {
		return p.s.fail(errors.New("premature EOF"))
	}
	if err != nil {
		return p.s.fail(err)
	}
	extra := p.s.topExtra().(*objectProcessorState)
	switch extra.tokenState {
	case key:
		return p.processKey(token, extra)
	case value:
		return p.processValue(token, extra)
	default:
		return p.s.fail(errors.New("invalid token state"))
	}
}

var idKeys = map[string]bool{
	"traceId":      true,
	"spanId":       true,
	"parentSpanId": true,
}

func (p *objectProcessor) processKey(token json.Token, extra *objectProcessorState) result {
	extra.tokenState = value
	// If tempSlice is not nil, then we have just finished
	// processing an array value. Put it into the map, before
	// processing the next token.
	if extra.tempSlice != nil {
		extra.currentMap[extra.currentKey], extra.tempSlice = extra.tempSlice, nil
	}
	switch v := token.(type) {
	case string:
		extra.currentKey = v
	case json.Delim:
		switch v.String() {
		case "}":
			p.s.pop()
		default:
			return p.s.fail(errors.New("unexpected closing delimiter"))
		}
	default:
		return p.s.fail(errors.New("unexpected token"))
	}
	return resultOK
}

func (p *objectProcessor) processValue(token json.Token, extra *objectProcessorState) result {
	extra.tokenState = key
	if idKeys[extra.currentKey] {
		str, ok := token.(string)
		if !ok {
			return p.s.fail(errors.New("expected a string for an ID value"))
		}
		fixedID, err := fixID(p.fixup, str)
		if err != nil {
			return p.s.fail(err)
		}
		extra.currentMap[extra.currentKey] = fixedID
		return resultOK
	}
	switch v := token.(type) {
	case json.Delim:
		switch v.String() {
		case "{":
			m := make(jsonMap)
			extra.currentMap[extra.currentKey] = m
			p.s.pushObject(m)
		case "[":
			extra.tempSlice = []interface{}{}
			p.s.pushArray(&extra.tempSlice)
		default:
			return p.s.fail(errors.New("unexpected closing delimiter"))
		}
	default:
		v, err := handleLeafToken(token)
		if err != nil {
			return p.s.fail(err)
		}
		extra.currentMap[extra.currentKey] = v
	}
	return resultOK
}

func fixID(fixup fixupType, str string) (string, error) {
	switch fixup {
	case fixupToHex:
		return fixIDToHex(str)
	case fixupToBase64:
		return fixIDToBase64(str)
	default:
		return "", errors.New("invalid fixup type")
	}
}

func fixIDToHex(str string) (string, error) {
	strBuf := bytes.NewBufferString(str)
	dec := base64.NewDecoder(base64.StdEncoding, strBuf)
	b, err := ioutil.ReadAll(dec)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func fixIDToBase64(str string) (string, error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	bLen := len(b)
	bLen64 := (int64)(bLen)
	buf := new(bytes.Buffer)
	stdEnc := base64.StdEncoding
	buf.Grow(stdEnc.EncodedLen(bLen))
	enc := base64.NewEncoder(stdEnc, buf)
	written, err := io.Copy(enc, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	if written != bLen64 {
		return "", errors.New("failed to copy all bytes to base64 encoder")
	}
	enc.Close()
	return buf.String(), nil
}

type arrayProcessorState struct {
	slicePtr *[]interface{}
	// for subarrays
	tempSlice []interface{}
}

func (s *arrayProcessorState) init(slicePtr *[]interface{}) {
	s.slicePtr = slicePtr
	s.tempSlice = nil
}

type arrayProcessor state

var _ processor = (*arrayProcessor)(nil)

func (p *arrayProcessor) process(token json.Token, err error) result {
	if err == io.EOF {
		return p.s.fail(errors.New("premature EOF"))
	}
	if err != nil {
		return p.s.fail(err)
	}
	extra := p.s.topExtra().(*arrayProcessorState)
	// If tempSlice is not nil, then we have just finished
	// processing an array item in our array. Put it into the
	// slice, before processing the next token.
	if extra.tempSlice != nil {
		*extra.slicePtr = append(*extra.slicePtr, extra.tempSlice)
		extra.tempSlice = nil
	}
	switch v := token.(type) {
	case json.Delim:
		switch v.String() {
		case "]":
			p.s.pop()
		case "[":
			extra.tempSlice = []interface{}{}
			p.s.pushArray(&extra.tempSlice)
		case "{":
			m := make(jsonMap)
			*extra.slicePtr = append(*extra.slicePtr, m)
			p.s.pushObject(m)
		default:
			return p.s.fail(errors.New("invalid delimiter"))
		}
		return resultOK
	default:
		v, err := handleLeafToken(token)
		if err != nil {
			return p.s.fail(err)
		}
		*extra.slicePtr = append(*extra.slicePtr, v)
		return resultOK
	}
}

func handleLeafToken(token json.Token) (interface{}, error) {
	switch v := token.(type) {
	case bool:
		return v, nil
	case float64:
		return v, nil
	case json.Number:
		i, err := v.Int64()
		if err == nil {
			return i, nil
		}
		f, err := v.Float64()
		if err == nil {
			return f, nil
		}
		return v.String(), nil
	case string:
		return v, nil
	case nil:
		return nil, nil
	}
	return nil, errors.New("not a leaf token")
}
