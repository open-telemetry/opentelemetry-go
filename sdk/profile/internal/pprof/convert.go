// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pprof

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	semconv "go.opentelemetry.io/otel/semconv/v1.45.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	labelTraceID = "otel.trace_id"
	labelSpanID  = "otel.span_id"
)

// ConvertOptions configures conversion from pprof to profiledata.
type ConvertOptions struct {
	// SampleTypes optionally restricts which pprof sample type names are
	// emitted. When empty, all sample types are emitted.
	SampleTypes []string
	// DefaultSampleType optionally overrides which sample type is placed at
	// index 0 (the default displayed profile in OTLP).
	DefaultSampleType string
	// FrameType optionally attaches a profile.frame.type attribute to each
	// converted Location.
	FrameType attribute.KeyValue
	// BuildIDs optionally enriches Mapping entries with OpenTelemetry build ID
	// attributes (htlhash, go, gnu).
	BuildIDs *BuildIDCache
	// Instantaneous forces Duration to 0 for snapshot profiles (such as live
	// heap and goroutine profiles).
	Instantaneous bool
}

// Convert transforms a decoded pprof Profile into one or more
// [profiledata.Profile] values (one per non-empty sample type) along with
// pprof scope attributes (default_sample_type and sample_type_order).
func Convert(p *Profile, opts ConvertOptions) ([]profiledata.Profile, []attribute.KeyValue) {
	if p == nil || len(p.SampleType) == 0 || len(p.Sample) == 0 {
		return nil, nil
	}

	selected := selectSampleTypes(p, opts)
	if len(selected) == 0 {
		return nil, nil
	}

	mappings := convertMappings(p, opts.BuildIDs)
	locations := convertLocations(p, mappings, opts.FrameType)
	profileAttrs := convertProfileAttrs(p)

	periodType := profiledata.ValueType{
		Type: p.str(p.PeriodType.Type),
		Unit: p.str(p.PeriodType.Unit),
	}
	var profTime time.Time
	if p.TimeNanos > 0 {
		profTime = time.Unix(0, p.TimeNanos).UTC()
	}
	var profDuration time.Duration
	if !opts.Instantaneous && p.DurationNanos > 0 {
		profDuration = time.Duration(p.DurationNanos)
	}

	// Pre-resolve each sample's stack, attributes, and span link once across
	// all sample types emitted from the same pprof profile.
	type resolvedSample struct {
		stack []*profiledata.Location
		attrs []profiledata.Attribute
		link  profiledata.Link
		key   string
	}
	resolved := make([]resolvedSample, len(p.Sample))
	var keyBuf []byte
	for i := range p.Sample {
		s := &p.Sample[i]
		stack := make([]*profiledata.Location, 0, len(s.LocationID))
		for _, locID := range s.LocationID {
			if loc := locations[locID]; loc != nil {
				stack = append(stack, loc)
			}
		}
		attrs, link := convertSampleLabels(p, s.Label)
		var key string
		keyBuf, key = mergedSampleKey(keyBuf[:0], s.LocationID, attrs, link)
		resolved[i] = resolvedSample{
			stack: stack,
			attrs: attrs,
			link:  link,
			key:   key,
		}
	}

	profiles := make([]profiledata.Profile, 0, len(selected))
	order := make([]int, 0, len(selected))
	for _, stIdx := range selected {
		st := p.SampleType[stIdx]
		samples := make([]profiledata.Sample, 0, len(p.Sample))
		samplePos := make(map[string]int, len(p.Sample))

		for i := range p.Sample {
			val := p.Sample[i].Value[stIdx]
			if val == 0 {
				continue
			}
			rs := &resolved[i]
			if pos, exists := samplePos[rs.key]; exists {
				samples[pos].Values[0] += val
				continue
			}
			samplePos[rs.key] = len(samples)
			samples = append(samples, profiledata.Sample{
				Stack:      rs.stack,
				Values:     []int64{val},
				Attributes: rs.attrs,
				Link:       rs.link,
			})
		}

		if len(samples) == 0 {
			continue
		}

		profiles = append(profiles, profiledata.Profile{
			SampleType: profiledata.ValueType{
				Type: p.str(st.Type),
				Unit: p.str(st.Unit),
			},
			Samples:    samples,
			Time:       profTime,
			Duration:   profDuration,
			PeriodType: periodType,
			Period:     p.Period,
			Attributes: profileAttrs,
		})
		order = append(order, stIdx)
	}

	if len(profiles) == 0 {
		return nil, nil
	}

	var scopeAttrs []attribute.KeyValue
	if defaultType := p.str(p.DefaultSampleType); defaultType != "" {
		scopeAttrs = append(scopeAttrs, semconv.PprofScopeDefaultSampleType(defaultType))
	} else if len(profiles) > 0 {
		scopeAttrs = append(scopeAttrs, semconv.PprofScopeDefaultSampleType(profiles[0].SampleType.Type))
	}
	if len(order) > 1 {
		scopeAttrs = append(scopeAttrs, semconv.PprofScopeSampleTypeOrder(order...))
	}
	return profiles, scopeAttrs
}

// selectSampleTypes picks the requested sample type indices and moves the
// default sample type to index 0 because OTLP visualization tools default to
// profiles[0], whereas pprof defaults to DefaultSampleType or the last entry.
func selectSampleTypes(p *Profile, opts ConvertOptions) []int {
	indices := make([]int, 0, len(p.SampleType))
	if len(opts.SampleTypes) == 0 {
		for i := range p.SampleType {
			indices = append(indices, i)
		}
	} else {
		for i, st := range p.SampleType {
			if slices.Contains(opts.SampleTypes, p.str(st.Type)) {
				indices = append(indices, i)
			}
		}
	}
	if len(indices) <= 1 {
		return indices
	}

	defaultName := opts.DefaultSampleType
	if defaultName == "" {
		defaultName = p.str(p.DefaultSampleType)
	}

	defaultPos := len(indices) - 1
	if defaultName != "" {
		for pos, idx := range indices {
			if p.str(p.SampleType[idx].Type) == defaultName {
				defaultPos = pos
				break
			}
		}
	}
	if defaultPos != 0 {
		indices[0], indices[defaultPos] = indices[defaultPos], indices[0]
	}
	return indices
}

func convertMappings(p *Profile, cache *BuildIDCache) map[uint64]*profiledata.Mapping {
	out := make(map[uint64]*profiledata.Mapping, len(p.Mapping))
	for i := range p.Mapping {
		m := &p.Mapping[i]
		filename := p.str(m.Filename)

		var attrs []profiledata.Attribute
		gnuBuildID := p.str(m.BuildID)
		var goBuildID, htlHash string
		if cache != nil {
			ids := cache.Lookup(filename)
			htlHash = ids.HTLHash
			goBuildID = ids.GoBuildID
			if gnuBuildID == "" {
				gnuBuildID = ids.GNUBuildID
			}
		}

		if gnuBuildID != "" {
			attrs = append(attrs, kvToAttr(semconv.ProcessExecutableBuildIDGNU(gnuBuildID)))
		}
		if goBuildID != "" {
			attrs = append(attrs, kvToAttr(semconv.ProcessExecutableBuildIDGo(goBuildID)))
		}
		if htlHash != "" {
			attrs = append(attrs, kvToAttr(semconv.ProcessExecutableBuildIDHtlhash(htlHash)))
		}
		if m.HasFunctions {
			attrs = append(attrs, kvToAttr(semconv.PprofMappingHasFunctions(true)))
		}
		if m.HasFilenames {
			attrs = append(attrs, kvToAttr(semconv.PprofMappingHasFilenames(true)))
		}
		if m.HasLineNumbers {
			attrs = append(attrs, kvToAttr(semconv.PprofMappingHasLineNumbers(true)))
		}
		if m.HasInlineFrames {
			attrs = append(attrs, kvToAttr(semconv.PprofMappingHasInlineFrames(true)))
		}

		memStart, memLimit := m.MemoryStart, m.MemoryLimit
		if memStart >= memLimit {
			// OTLP Profiles requires memory_start == 0 && memory_limit == 0 or
			// memory_start < memory_limit.
			memStart, memLimit = 0, 0
		}

		out[m.ID] = &profiledata.Mapping{
			MemoryStart: memStart,
			MemoryLimit: memLimit,
			FileOffset:  m.FileOffset,
			Filename:    filename,
			Attributes:  attrs,
		}
	}
	return out
}

func convertLocations(
	p *Profile,
	mappings map[uint64]*profiledata.Mapping,
	frameType attribute.KeyValue,
) map[uint64]*profiledata.Location {
	funcs := make(map[uint64]profiledata.Function, len(p.Function))
	for i := range p.Function {
		fn := &p.Function[i]
		startLine := max(fn.StartLine, 0)
		funcs[fn.ID] = profiledata.Function{
			Name:       p.str(fn.Name),
			SystemName: p.str(fn.SystemName),
			Filename:   p.str(fn.Filename),
			StartLine:  startLine,
		}
	}

	hasFrameType := frameType.Valid()
	out := make(map[uint64]*profiledata.Location, len(p.Location))
	for i := range p.Location {
		loc := &p.Location[i]
		lines := make([]profiledata.Line, 0, len(loc.Line))
		seenLines := make(map[profiledata.Line]struct{}, len(loc.Line))
		for _, ln := range loc.Line {
			pl := profiledata.Line{
				Function: funcs[ln.FunctionID],
				Line:     max(ln.Line, 0),
				Column:   max(ln.Column, 0),
			}
			// OTLP Profiles disallows duplicate lines within a Location.
			if _, dup := seenLines[pl]; dup {
				continue
			}
			seenLines[pl] = struct{}{}
			lines = append(lines, pl)
		}

		var attrs []profiledata.Attribute
		if loc.IsFolded {
			attrs = append(attrs, kvToAttr(semconv.PprofLocationIsFolded(true)))
		}
		if hasFrameType {
			attrs = append(attrs, kvToAttr(frameType))
		}

		out[loc.ID] = &profiledata.Location{
			Mapping:    mappings[loc.MappingID],
			Address:    loc.Address,
			Lines:      lines,
			Attributes: attrs,
		}
	}
	return out
}

func convertProfileAttrs(p *Profile) []profiledata.Attribute {
	var attrs []profiledata.Attribute
	if v := p.str(p.DropFrames); v != "" {
		attrs = append(attrs, kvToAttr(semconv.PprofProfileDropFrames(v)))
	}
	if v := p.str(p.KeepFrames); v != "" {
		attrs = append(attrs, kvToAttr(semconv.PprofProfileKeepFrames(v)))
	}
	if v := p.str(p.DocURL); v != "" {
		attrs = append(attrs, kvToAttr(semconv.PprofProfileDocURL(v)))
	}
	if len(p.Comment) > 0 {
		comments := make([]string, 0, len(p.Comment))
		for _, idx := range p.Comment {
			if c := p.str(idx); c != "" {
				comments = append(comments, c)
			}
		}
		if len(comments) > 0 {
			attrs = append(attrs, kvToAttr(semconv.PprofProfileComment(comments...)))
		}
	}
	return attrs
}

func convertSampleLabels(p *Profile, labels []Label) ([]profiledata.Attribute, profiledata.Link) {
	if len(labels) == 0 {
		return nil, profiledata.Link{}
	}

	var traceID trace.TraceID
	var spanID trace.SpanID
	traceIDIdx, spanIDIdx := -1, -1
	for i, lbl := range labels {
		if lbl.Str == 0 {
			continue
		}
		switch p.str(lbl.Key) {
		case labelTraceID:
			if tid, err := trace.TraceIDFromHex(p.str(lbl.Str)); err == nil {
				traceID = tid
				traceIDIdx = i
			}
		case labelSpanID:
			if sid, err := trace.SpanIDFromHex(p.str(lbl.Str)); err == nil {
				spanID = sid
				spanIDIdx = i
			}
		}
	}

	var link profiledata.Link
	if traceIDIdx >= 0 && spanIDIdx >= 0 {
		link = profiledata.Link{TraceID: traceID, SpanID: spanID}
	}

	// Keep unique attribute keys (last wins) to satisfy OTLP Profiles
	// attribute_indices key-uniqueness rules.
	attrs := make([]profiledata.Attribute, 0, len(labels))
	keyPos := make(map[string]int, len(labels))
	for i, lbl := range labels {
		if link.IsValid() && (i == traceIDIdx || i == spanIDIdx) {
			continue
		}
		key := p.str(lbl.Key)
		if key == "" {
			continue
		}

		var a profiledata.Attribute
		if lbl.Str != 0 {
			a = profiledata.Attribute{
				Key:   attribute.Key(key),
				Value: attribute.StringValue(p.str(lbl.Str)),
			}
		} else {
			a = profiledata.Attribute{
				Key:   attribute.Key(key),
				Value: attribute.Int64Value(lbl.Num),
				Unit:  p.str(lbl.NumUnit),
			}
		}

		if pos, exists := keyPos[key]; exists {
			attrs[pos] = a
		} else {
			keyPos[key] = len(attrs)
			attrs = append(attrs, a)
		}
	}

	// Sort attributes by key so sample identity keys and OTLP output are deterministic.
	slices.SortFunc(attrs, func(a, b profiledata.Attribute) int {
		return strings.Compare(string(a.Key), string(b.Key))
	})
	return attrs, link
}

func mergedSampleKey(
	b []byte,
	locIDs []uint64,
	attrs []profiledata.Attribute,
	link profiledata.Link,
) ([]byte, string) {
	for _, id := range locIDs {
		b = strconv.AppendUint(b, id, 16)
		b = append(b, ',')
	}
	b = append(b, '|')
	if link.IsValid() {
		b = append(b, link.TraceID.String()...)
		b = append(b, '/')
		b = append(b, link.SpanID.String()...)
	}
	b = append(b, '|')
	for _, a := range attrs {
		b = append(b, a.Key...)
		b = append(b, '=')
		b = append(b, a.Value.Type().String()...)
		b = append(b, ':')
		b = append(b, a.Value.String()...)
		if a.Unit != "" {
			b = append(b, '/')
			b = append(b, a.Unit...)
		}
		b = append(b, '\x00')
	}
	return b, string(b)
}

func kvToAttr(kv attribute.KeyValue) profiledata.Attribute {
	return profiledata.Attribute{Key: kv.Key, Value: kv.Value}
}
