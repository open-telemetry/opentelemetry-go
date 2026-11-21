// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pprof

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
)

const (
	// maxDecompressedSize bounds gzip expansion to protect the host application
	// against malformed or adversarial payloads.
	maxDecompressedSize = 64 * 1024 * 1024
	// maxPooledDecompressCap avoids pinning oversized buffers in the pool after spikes.
	maxPooledDecompressCap = 1 << 20
)

var (
	gzipReaderPool    sync.Pool
	decompressBufPool = sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}
)

// Protobuf wire types.
const (
	wireVarint  = 0
	wireFixed64 = 1
	wireBytes   = 2
	wireFixed32 = 5
)

var (
	errTruncated       = errors.New("pprof: truncated protobuf payload")
	errInvalidWireType = errors.New("pprof: invalid protobuf wire type")
	errInvalidFieldNum = errors.New("pprof: invalid protobuf field number")
	errPayloadTooLarge = errors.New("pprof: decompressed profile exceeds maximum size")
)

// Profile represents a decoded perftools.profiles.Profile protobuf message.
type Profile struct {
	SampleType        []ValueType
	Sample            []Sample
	Mapping           []Mapping
	Location          []Location
	Function          []Function
	StringTable       []string
	DropFrames        int64
	KeepFrames        int64
	TimeNanos         int64
	DurationNanos     int64
	PeriodType        ValueType
	Period            int64
	Comment           []int64
	DefaultSampleType int64
	DocURL            int64
}

// ValueType describes the type and unit of a pprof sample value or period.
type ValueType struct {
	Type int64
	Unit int64
}

// Sample records measurements for a call stack and optional labels.
type Sample struct {
	LocationID []uint64
	Value      []int64
	Label      []Label
}

// Label attaches a string or numeric label to a Sample.
type Label struct {
	Key     int64
	Str     int64
	Num     int64
	NumUnit int64
}

// Mapping describes an executable or library mapping in memory.
type Mapping struct {
	ID              uint64
	MemoryStart     uint64
	MemoryLimit     uint64
	FileOffset      uint64
	Filename        int64
	BuildID         int64
	HasFunctions    bool
	HasFilenames    bool
	HasLineNumbers  bool
	HasInlineFrames bool
}

// Location describes a single stack frame.
type Location struct {
	ID        uint64
	MappingID uint64
	Address   uint64
	Line      []Line
	IsFolded  bool
}

// Line describes a single source line within a Location.
type Line struct {
	FunctionID uint64
	Line       int64
	Column     int64
}

// Function describes a source function referenced by a Line.
type Function struct {
	ID         uint64
	Name       int64
	SystemName int64
	Filename   int64
	StartLine  int64
}

// Parse decodes a gzip-compressed or uncompressed pprof protobuf payload.
func Parse(data []byte) (*Profile, error) {
	raw := data
	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		buf, err := decompressGzip(data)
		if err != nil {
			return nil, err
		}
		// NOTE: decodeProfile copies strings and integers and retains no slice of
		// raw, so returning buf to the pool after decodeProfile finishes is safe.
		defer putDecompressBuf(buf)
		raw = buf.Bytes()
	}

	p := &Profile{}
	if err := decodeProfile(raw, p); err != nil {
		return nil, err
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func putDecompressBuf(buf *bytes.Buffer) {
	if buf.Cap() <= maxPooledDecompressCap {
		buf.Reset()
		decompressBufPool.Put(buf)
	}
}

func decompressGzip(data []byte) (*bytes.Buffer, error) {
	var gr *gzip.Reader
	var err error
	if v := gzipReaderPool.Get(); v != nil {
		gr = v.(*gzip.Reader)
		err = gr.Reset(bytes.NewReader(data))
	} else {
		gr, err = gzip.NewReader(bytes.NewReader(data))
	}
	if err != nil {
		return nil, fmt.Errorf("pprof: gzip reader: %w", err)
	}
	defer func() {
		_ = gr.Close()
		gzipReaderPool.Put(gr)
	}()

	buf := decompressBufPool.Get().(*bytes.Buffer)
	buf.Reset()
	// Read one byte past the cap to detect oversized streams cleanly.
	if _, err := buf.ReadFrom(io.LimitReader(gr, maxDecompressedSize+1)); err != nil {
		putDecompressBuf(buf)
		return nil, fmt.Errorf("pprof: gzip read: %w", err)
	}
	if buf.Len() > maxDecompressedSize {
		putDecompressBuf(buf)
		return nil, errPayloadTooLarge
	}
	return buf, nil
}

func (p *Profile) str(idx int64) string {
	if idx < 0 || idx >= int64(len(p.StringTable)) {
		return ""
	}
	return p.StringTable[idx]
}

func (p *Profile) validate() error {
	strLen := int64(len(p.StringTable))
	checkStr := func(indices ...int64) error {
		for _, idx := range indices {
			if idx < 0 || idx >= strLen {
				return fmt.Errorf("pprof: string index %d out of range [0, %d)", idx, strLen)
			}
		}
		return nil
	}

	if strLen == 0 {
		if len(p.SampleType) > 0 || len(p.Sample) > 0 || len(p.Mapping) > 0 ||
			len(p.Location) > 0 || len(p.Function) > 0 {
			return errors.New("pprof: empty string table with non-empty profile")
		}
		return nil
	}
	if p.StringTable[0] != "" {
		return fmt.Errorf("pprof: string_table[0] must be empty, got %q", p.StringTable[0])
	}

	for _, st := range p.SampleType {
		if err := checkStr(st.Type, st.Unit); err != nil {
			return err
		}
	}
	if err := checkStr(
		p.PeriodType.Type,
		p.PeriodType.Unit,
		p.DropFrames,
		p.KeepFrames,
		p.DefaultSampleType,
		p.DocURL,
	); err != nil {
		return err
	}
	if err := checkStr(p.Comment...); err != nil {
		return err
	}
	for _, m := range p.Mapping {
		if err := checkStr(m.Filename, m.BuildID); err != nil {
			return err
		}
	}
	for _, fn := range p.Function {
		if err := checkStr(fn.Name, fn.SystemName, fn.Filename); err != nil {
			return err
		}
	}
	for i := range p.Sample {
		s := &p.Sample[i]
		if len(s.Value) != len(p.SampleType) {
			return fmt.Errorf(
				"pprof: sample[%d] has %d values, want %d",
				i, len(s.Value), len(p.SampleType),
			)
		}
		for _, lbl := range s.Label {
			if err := checkStr(lbl.Key, lbl.Str, lbl.NumUnit); err != nil {
				return err
			}
		}
	}
	return nil
}

// wireReader decodes protobuf wire fields sequentially while latching the first
// error so message decoders stay concise without repeated error plumbing.
type wireReader struct {
	b   []byte
	err error
}

func (r *wireReader) next() (num, wt uint64, ok bool) {
	if r.err != nil || len(r.b) == 0 {
		return 0, 0, false
	}
	v, n := binary.Uvarint(r.b)
	if n <= 0 {
		r.err = errTruncated
		return 0, 0, false
	}
	r.b = r.b[n:]
	num = v >> 3
	if num == 0 {
		r.err = errInvalidFieldNum
		return 0, 0, false
	}
	return num, v & 7, true
}

func (r *wireReader) uint64(wt uint64) uint64 {
	if wt != wireVarint {
		r.err = errInvalidWireType
		return 0
	}
	v, n := binary.Uvarint(r.b)
	if n <= 0 {
		r.err = errTruncated
		return 0
	}
	r.b = r.b[n:]
	return v
}

func (r *wireReader) int64(wt uint64) int64 {
	return int64(r.uint64(wt)) //nolint:gosec // Protobuf int64 wire representation.
}

func (r *wireReader) bool(wt uint64) bool {
	return r.uint64(wt) != 0
}

func (r *wireReader) bytes(wt uint64) []byte {
	if wt != wireBytes {
		r.err = errInvalidWireType
		return nil
	}
	length, n := binary.Uvarint(r.b)
	if n <= 0 {
		r.err = errTruncated
		return nil
	}
	rest := r.b[n:]
	if length > uint64(len(rest)) {
		r.err = errTruncated
		return nil
	}
	r.b = rest[length:]
	return rest[:length]
}

func (r *wireReader) repeatedUint64(wt uint64, dst []uint64) []uint64 {
	switch wt {
	case wireVarint:
		return append(dst, r.uint64(wt))
	case wireBytes:
		sub := wireReader{b: r.bytes(wt)}
		for len(sub.b) > 0 && sub.err == nil {
			dst = append(dst, sub.uint64(wireVarint))
		}
		if sub.err != nil {
			r.err = sub.err
		}
		return dst
	default:
		r.err = errInvalidWireType
		return dst
	}
}

func (r *wireReader) repeatedInt64(wt uint64, dst []int64) []int64 {
	switch wt {
	case wireVarint:
		return append(dst, r.int64(wt))
	case wireBytes:
		sub := wireReader{b: r.bytes(wt)}
		for len(sub.b) > 0 && sub.err == nil {
			dst = append(dst, sub.int64(wireVarint))
		}
		if sub.err != nil {
			r.err = sub.err
		}
		return dst
	default:
		r.err = errInvalidWireType
		return dst
	}
}

func (r *wireReader) skip(wt uint64) {
	switch wt {
	case wireVarint:
		_ = r.uint64(wt)
	case wireFixed64:
		if len(r.b) < 8 {
			r.err = errTruncated
			return
		}
		r.b = r.b[8:]
	case wireBytes:
		_ = r.bytes(wt)
	case wireFixed32:
		if len(r.b) < 4 {
			r.err = errTruncated
			return
		}
		r.b = r.b[4:]
	default:
		r.err = errInvalidWireType
	}
}

func decodeProfile(b []byte, p *Profile) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			var vt ValueType
			if err := decodeValueType(r.bytes(wt), &vt); err != nil {
				return err
			}
			p.SampleType = append(p.SampleType, vt)
		case 2:
			var s Sample
			if err := decodeSample(r.bytes(wt), &s); err != nil {
				return err
			}
			p.Sample = append(p.Sample, s)
		case 3:
			var m Mapping
			if err := decodeMapping(r.bytes(wt), &m); err != nil {
				return err
			}
			p.Mapping = append(p.Mapping, m)
		case 4:
			var loc Location
			if err := decodeLocation(r.bytes(wt), &loc); err != nil {
				return err
			}
			p.Location = append(p.Location, loc)
		case 5:
			var fn Function
			if err := decodeFunction(r.bytes(wt), &fn); err != nil {
				return err
			}
			p.Function = append(p.Function, fn)
		case 6:
			p.StringTable = append(p.StringTable, string(r.bytes(wt)))
		case 7:
			p.DropFrames = r.int64(wt)
		case 8:
			p.KeepFrames = r.int64(wt)
		case 9:
			p.TimeNanos = r.int64(wt)
		case 10:
			p.DurationNanos = r.int64(wt)
		case 11:
			if err := decodeValueType(r.bytes(wt), &p.PeriodType); err != nil {
				return err
			}
		case 12:
			p.Period = r.int64(wt)
		case 13:
			p.Comment = r.repeatedInt64(wt, p.Comment)
		case 14:
			p.DefaultSampleType = r.int64(wt)
		case 15:
			p.DocURL = r.int64(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeValueType(b []byte, vt *ValueType) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			vt.Type = r.int64(wt)
		case 2:
			vt.Unit = r.int64(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeSample(b []byte, s *Sample) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			s.LocationID = r.repeatedUint64(wt, s.LocationID)
		case 2:
			s.Value = r.repeatedInt64(wt, s.Value)
		case 3:
			var lbl Label
			if err := decodeLabel(r.bytes(wt), &lbl); err != nil {
				return err
			}
			s.Label = append(s.Label, lbl)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeLabel(b []byte, lbl *Label) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			lbl.Key = r.int64(wt)
		case 2:
			lbl.Str = r.int64(wt)
		case 3:
			lbl.Num = r.int64(wt)
		case 4:
			lbl.NumUnit = r.int64(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeMapping(b []byte, m *Mapping) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			m.ID = r.uint64(wt)
		case 2:
			m.MemoryStart = r.uint64(wt)
		case 3:
			m.MemoryLimit = r.uint64(wt)
		case 4:
			m.FileOffset = r.uint64(wt)
		case 5:
			m.Filename = r.int64(wt)
		case 6:
			m.BuildID = r.int64(wt)
		case 7:
			m.HasFunctions = r.bool(wt)
		case 8:
			m.HasFilenames = r.bool(wt)
		case 9:
			m.HasLineNumbers = r.bool(wt)
		case 10:
			m.HasInlineFrames = r.bool(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeLocation(b []byte, loc *Location) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			loc.ID = r.uint64(wt)
		case 2:
			loc.MappingID = r.uint64(wt)
		case 3:
			loc.Address = r.uint64(wt)
		case 4:
			var ln Line
			if err := decodeLine(r.bytes(wt), &ln); err != nil {
				return err
			}
			loc.Line = append(loc.Line, ln)
		case 5:
			loc.IsFolded = r.bool(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeLine(b []byte, ln *Line) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			ln.FunctionID = r.uint64(wt)
		case 2:
			ln.Line = r.int64(wt)
		case 3:
			ln.Column = r.int64(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}

func decodeFunction(b []byte, fn *Function) error {
	r := wireReader{b: b}
	for num, wt, ok := r.next(); ok; num, wt, ok = r.next() {
		switch num {
		case 1:
			fn.ID = r.uint64(wt)
		case 2:
			fn.Name = r.int64(wt)
		case 3:
			fn.SystemName = r.int64(wt)
		case 4:
			fn.Filename = r.int64(wt)
		case 5:
			fn.StartLine = r.int64(wt)
		default:
			r.skip(wt)
		}
	}
	return r.err
}
