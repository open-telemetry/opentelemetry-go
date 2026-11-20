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
		// Batch single-element value slices into one backing array per profile
		// to avoid one []int64 heap allocation per sample.
		valSlab := make([]int64, len(p.Sample))

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
			pos := len(samples)
			samplePos[rs.key] = pos
			valSlab[pos] = val
			samples = append(samples, profiledata.Sample{
				Stack:      rs.stack,
				Values:     valSlab[pos : pos+1 : pos+1],
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
		// Per semconv, only record default_sample_type when present in the
		// source pprof profile and retained in the emitted profiles.
		for i := range profiles {
			if profiles[i].SampleType.Type == defaultType {
				scopeAttrs = append(scopeAttrs, semconv.PprofScopeDefaultSampleType(defaultType))
				break
			}
		}
	}
	if len(order) > 1 {
		scopeAttrs = append(scopeAttrs, semconv.PprofScopeSampleTypeOrder(denseSampleTypeOrder(order)...))
	}
	return profiles, scopeAttrs
}

// denseSampleTypeOrder maps emitted original pprof sample type indices into
// dense 0-based positions preserving their original relative order so every
// index stays within [0, len(order)).
func denseSampleTypeOrder(order []int) []int {
	sorted := slices.Clone(order)
	slices.Sort(sorted)
	dense := make([]int, len(order))
	for i, origIdx := range order {
		dense[i], _ = slices.BinarySearch(sorted, origIdx)
	}
	return dense
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
		def := indices[defaultPos]
		copy(indices[1:defaultPos+1], indices[:defaultPos])
		indices[0] = def
	}
	return indices
}

func convertMappings(p *Profile, cache *BuildIDCache) map[uint64]*profiledata.Mapping {
	out := make(map[uint64]*profiledata.Mapping, len(p.Mapping))
	for i := range p.Mapping {
		m := &p.Mapping[i]
		if m.ID == 0 {
			continue
		}
		filename := p.str(m.Filename)

		var attrs []profiledata.Attribute
		var gnuBuildID, goBuildID, htlHash string
		if cache != nil {
			ids := cache.Lookup(filename)
			htlHash = ids.HTLHash
			goBuildID = ids.GoBuildID
			gnuBuildID = ids.GNUBuildID
		}
		// Only fall back to pprof Mapping.BuildID for process.executable.build_id.gnu
		// when it is an even-length hex string (Windows PE mappings put path+modTime
		// in BuildID rather than a GNU ELF build ID).
		if gnuBuildID == "" {
			if rawID := p.str(m.BuildID); isHexBuildID(rawID) {
				gnuBuildID = rawID
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
		if fn.ID == 0 {
			continue
		}
		startLine := max(fn.StartLine, 0)
		funcs[fn.ID] = profiledata.Function{
			Name:       p.str(fn.Name),
			SystemName: p.str(fn.SystemName),
			Filename:   p.str(fn.Filename),
			StartLine:  startLine,
		}
	}

	hasFrameType := frameType.Valid()
	// Preallocate shared immutable attribute slices so per-location conversion
	// does not allocate a separate slice for every frame.
	var frameOnlyAttrs, foldedOnlyAttrs, foldedFrameAttrs []profiledata.Attribute
	if hasFrameType {
		frameOnlyAttrs = []profiledata.Attribute{kvToAttr(frameType)}
		foldedFrameAttrs = []profiledata.Attribute{
			kvToAttr(semconv.PprofLocationIsFolded(true)),
			kvToAttr(frameType),
		}
	}
	foldedOnlyAttrs = []profiledata.Attribute{kvToAttr(semconv.PprofLocationIsFolded(true))}

	out := make(map[uint64]*profiledata.Location, len(p.Location))
	locSlab := make([]profiledata.Location, len(p.Location))
	singleLines := 0
	for i := range p.Location {
		if len(p.Location[i].Line) == 1 {
			singleLines++
		}
	}
	lineSlab := make([]profiledata.Line, singleLines)
	linePos := 0
	for i := range p.Location {
		loc := &p.Location[i]
		if loc.ID == 0 {
			continue
		}

		var lines []profiledata.Line
		switch len(loc.Line) {
		case 0:
		case 1:
			ln := loc.Line[0]
			var fn profiledata.Function
			if ln.FunctionID != 0 {
				fn = funcs[ln.FunctionID]
			}
			lineSlab[linePos] = profiledata.Line{
				Function: fn,
				Line:     max(ln.Line, 0),
				Column:   max(ln.Column, 0),
			}
			lines = lineSlab[linePos : linePos+1 : linePos+1]
			linePos++
		default:
			lines = make([]profiledata.Line, 0, len(loc.Line))
			var seenLines map[profiledata.Line]struct{}
			if len(loc.Line) > 8 {
				seenLines = make(map[profiledata.Line]struct{}, len(loc.Line))
			}
			for _, ln := range loc.Line {
				var fn profiledata.Function
				if ln.FunctionID != 0 {
					fn = funcs[ln.FunctionID]
				}
				pl := profiledata.Line{
					Function: fn,
					Line:     max(ln.Line, 0),
					Column:   max(ln.Column, 0),
				}
				// OTLP Profiles disallows duplicate lines within a Location.
				if seenLines != nil {
					if _, dup := seenLines[pl]; dup {
						continue
					}
					seenLines[pl] = struct{}{}
				} else if slices.Contains(lines, pl) {
					continue
				}
				lines = append(lines, pl)
			}
		}

		var attrs []profiledata.Attribute
		switch {
		case loc.IsFolded && hasFrameType:
			attrs = foldedFrameAttrs
		case loc.IsFolded:
			attrs = foldedOnlyAttrs
		case hasFrameType:
			attrs = frameOnlyAttrs
		}

		var mapping *profiledata.Mapping
		if loc.MappingID != 0 {
			mapping = mappings[loc.MappingID]
		}
		locSlab[i] = profiledata.Location{
			Mapping:    mapping,
			Address:    loc.Address,
			Lines:      lines,
			Attributes: attrs,
		}
		out[loc.ID] = &locSlab[i]
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
	// Fast path for single-label samples (every Go heap/allocs sample carries
	// a single "bytes" label), avoiding map and sort overhead.
	if len(labels) == 1 {
		lbl := labels[0]
		key := p.str(lbl.Key)
		if key == "" {
			return nil, profiledata.Link{}
		}
		return []profiledata.Attribute{labelToAttr(p, key, lbl)}, profiledata.Link{}
	}

	// Keep unique attribute keys (last wins) to satisfy OTLP Profiles
	// attribute_indices key-uniqueness rules.
	attrs := make([]profiledata.Attribute, 0, len(labels))
	keyPos := make(map[string]int, len(labels))
	for _, lbl := range labels {
		key := p.str(lbl.Key)
		if key == "" {
			continue
		}
		a := labelToAttr(p, key, lbl)

		if pos, exists := keyPos[key]; exists {
			attrs[pos] = a
		} else {
			keyPos[key] = len(attrs)
			attrs = append(attrs, a)
		}
	}

	var link profiledata.Link
	if tPos, hasT := keyPos[labelTraceID]; hasT {
		if sPos, hasS := keyPos[labelSpanID]; hasS &&
			attrs[tPos].Value.Type() == attribute.STRING &&
			attrs[sPos].Value.Type() == attribute.STRING {
			tid, tErr := trace.TraceIDFromHex(attrs[tPos].Value.AsString())
			sid, sErr := trace.SpanIDFromHex(attrs[sPos].Value.AsString())
			if tErr == nil && sErr == nil {
				link = profiledata.Link{TraceID: tid, SpanID: sid}
				// Delete the higher index first so the lower index stays valid.
				hi, lo := max(tPos, sPos), min(tPos, sPos)
				attrs = slices.Delete(attrs, hi, hi+1)
				attrs = slices.Delete(attrs, lo, lo+1)
			}
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
	b = strconv.AppendInt(b, int64(len(locIDs)), 10)
	b = append(b, ';')
	for _, id := range locIDs {
		b = strconv.AppendUint(b, id, 16)
		b = append(b, ',')
	}
	if link.IsValid() {
		b = append(b, '1')
		b = append(b, link.TraceID[:]...)
		b = append(b, link.SpanID[:]...)
	} else {
		b = append(b, '0')
	}
	b = strconv.AppendInt(b, int64(len(attrs)), 10)
	b = append(b, ';')
	for _, a := range attrs {
		b = appendLenString(b, string(a.Key))
		b = strconv.AppendInt(b, int64(a.Value.Type()), 10)
		b = append(b, ',')
		switch a.Value.Type() {
		case attribute.STRING:
			b = appendLenString(b, a.Value.AsString())
		case attribute.INT64:
			b = strconv.AppendInt(b, a.Value.AsInt64(), 10)
			b = append(b, ',')
		default:
			b = appendLenString(b, a.Value.String())
		}
		b = appendLenString(b, a.Unit)
	}
	return b, string(b)
}

func labelToAttr(p *Profile, key string, lbl Label) profiledata.Attribute {
	if lbl.Str != 0 {
		return profiledata.Attribute{
			Key:   attribute.Key(key),
			Value: attribute.StringValue(p.str(lbl.Str)),
		}
	}
	return profiledata.Attribute{
		Key:   attribute.Key(key),
		Value: attribute.Int64Value(lbl.Num),
		Unit:  p.str(lbl.NumUnit),
	}
}

func kvToAttr(kv attribute.KeyValue) profiledata.Attribute {
	return profiledata.Attribute{Key: kv.Key, Value: kv.Value}
}

func isHexBuildID(s string) bool {
	if s == "" || len(s)%2 != 0 {
		return false
	}
	for i := range len(s) {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}
