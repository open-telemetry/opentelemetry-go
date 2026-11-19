// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pprof

import (
	"slices"
	"strconv"
	"strings"
)

// Delta computes the interval delta profile from prev to cur for cumulative
// Go runtime profiles (such as allocs, mutex, and block).
//
// Matching samples are identified by their resolved stack frames and labels
// rather than pprof location IDs, because runtime/pprof reassigns location IDs
// on every snapshot. Values whose delta is non-positive are clamped to zero,
// and samples whose values are all zero are dropped.
func Delta(prev, cur *Profile) *Profile {
	if cur == nil {
		return nil
	}
	if prev == nil || !sameSampleTypes(prev, cur) {
		return cur
	}

	prevIdx := buildProfileIndex(prev)
	prevTotals := make(map[string][]int64, len(prev.Sample))
	var keyBuf []byte
	for i := range prev.Sample {
		s := &prev.Sample[i]
		var key string
		keyBuf, key = canonicalSampleKey(keyBuf[:0], prev, prevIdx, s)
		if existing, ok := prevTotals[key]; ok {
			for j, v := range s.Value {
				existing[j] += v
			}
		} else {
			prevTotals[key] = slices.Clone(s.Value)
		}
	}

	curIdx := buildProfileIndex(cur)
	type aggSample struct {
		sample Sample
		values []int64
	}
	curAgg := make([]aggSample, 0, len(cur.Sample))
	curPos := make(map[string]int, len(cur.Sample))
	for i := range cur.Sample {
		s := &cur.Sample[i]
		var key string
		keyBuf, key = canonicalSampleKey(keyBuf[:0], cur, curIdx, s)
		if pos, ok := curPos[key]; ok {
			for j, v := range s.Value {
				curAgg[pos].values[j] += v
			}
		} else {
			curPos[key] = len(curAgg)
			curAgg = append(curAgg, aggSample{
				sample: *s,
				values: slices.Clone(s.Value),
			})
		}
	}

	out := *cur
	out.Sample = make([]Sample, 0, len(curAgg))
	for key, pos := range curPos {
		agg := &curAgg[pos]
		prevVals, hasPrev := prevTotals[key]

		nonZero := false
		for j, v := range agg.values {
			d := v
			if hasPrev {
				d = v - prevVals[j]
			}
			if d > 0 {
				agg.values[j] = d
				nonZero = true
			} else {
				agg.values[j] = 0
			}
		}
		if !nonZero {
			continue
		}

		ns := agg.sample
		ns.Value = agg.values
		out.Sample = append(out.Sample, ns)
	}

	if prev.TimeNanos > 0 && cur.TimeNanos >= prev.TimeNanos {
		out.TimeNanos = prev.TimeNanos
		out.DurationNanos = cur.TimeNanos - prev.TimeNanos
	}
	return &out
}

func sameSampleTypes(a, b *Profile) bool {
	if len(a.SampleType) != len(b.SampleType) {
		return false
	}
	for i := range a.SampleType {
		if a.str(a.SampleType[i].Type) != b.str(b.SampleType[i].Type) ||
			a.str(a.SampleType[i].Unit) != b.str(b.SampleType[i].Unit) {
			return false
		}
	}
	return true
}

type profileIndex struct {
	locs  map[uint64]*Location
	funcs map[uint64]*Function
}

func buildProfileIndex(p *Profile) profileIndex {
	idx := profileIndex{
		locs:  make(map[uint64]*Location, len(p.Location)),
		funcs: make(map[uint64]*Function, len(p.Function)),
	}
	for i := range p.Location {
		idx.locs[p.Location[i].ID] = &p.Location[i]
	}
	for i := range p.Function {
		idx.funcs[p.Function[i].ID] = &p.Function[i]
	}
	return idx
}

func canonicalSampleKey(b []byte, p *Profile, idx profileIndex, s *Sample) ([]byte, string) {
	for _, locID := range s.LocationID {
		loc := idx.locs[locID]
		if loc == nil {
			b = append(b, "?;"...)
			continue
		}
		// Use PC address plus resolved lines so identity stays stable across
		// runtime/pprof snapshots where Location.ID is reassigned.
		b = strconv.AppendUint(b, loc.Address, 16)
		for _, ln := range loc.Line {
			b = append(b, '|')
			if fn := idx.funcs[ln.FunctionID]; fn != nil {
				b = append(b, p.str(fn.Name)...)
				b = append(b, ':')
				b = append(b, p.str(fn.Filename)...)
			}
			b = append(b, ':')
			b = strconv.AppendInt(b, ln.Line, 10)
			b = append(b, ':')
			b = strconv.AppendInt(b, ln.Column, 10)
		}
		b = append(b, ';')
	}

	if len(s.Label) == 0 {
		return b, string(b)
	}

	type resolvedLabel struct {
		key     string
		str     string
		num     int64
		numUnit string
		isStr   bool
	}
	labels := make([]resolvedLabel, len(s.Label))
	for i, lbl := range s.Label {
		labels[i] = resolvedLabel{
			key:     p.str(lbl.Key),
			str:     p.str(lbl.Str),
			num:     lbl.Num,
			numUnit: p.str(lbl.NumUnit),
			isStr:   lbl.Str != 0,
		}
	}
	slices.SortFunc(labels, func(a, b resolvedLabel) int {
		if c := strings.Compare(a.key, b.key); c != 0 {
			return c
		}
		if a.isStr != b.isStr {
			if a.isStr {
				return -1
			}
			return 1
		}
		if c := strings.Compare(a.str, b.str); c != 0 {
			return c
		}
		if a.num != b.num {
			if a.num < b.num {
				return -1
			}
			return 1
		}
		return strings.Compare(a.numUnit, b.numUnit)
	})

	b = append(b, '#')
	for _, lbl := range labels {
		b = append(b, lbl.key...)
		b = append(b, '=')
		if lbl.isStr {
			b = append(b, "s:"...)
			b = append(b, lbl.str...)
		} else {
			b = append(b, "n:"...)
			b = strconv.AppendInt(b, lbl.num, 10)
			b = append(b, '/')
			b = append(b, lbl.numUnit...)
		}
		b = append(b, '\x00')
	}
	return b, string(b)
}
