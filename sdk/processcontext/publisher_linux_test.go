// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build linux

package processcontext

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func makeResource(t *testing.T, attrs ...attribute.KeyValue) *resource.Resource {
	t.Helper()
	r, err := resource.New(t.Context(), resource.WithAttributes(attrs...))
	require.NoError(t, err)
	return r
}

// readTS reads MonotonicPublishedAtNs from a live mapping without going
// through Shutdown (so the memory is still mapped).
func readTS(mem []byte) uint64 {
	return atomic.LoadUint64((*uint64)(unsafe.Pointer(&mem[offTimestamp])))
}

// readPayload returns a copy of the payload bytes from a live publisher.
func readPayload(p *publisher) []byte {
	sz := int(binary.LittleEndian.Uint32(p.headerMem[offPayloadSz:]))
	var src []byte
	if p.payloadMem != nil {
		src = p.payloadMem[:sz]
	} else {
		src = p.headerMem[headerSize : headerSize+sz]
	}
	out := make([]byte, sz)
	copy(out, src)
	return out
}

// ---- Basic creation ----------------------------------------------------

func TestNewPublisher(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "test"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	require.NotNil(t, pub)
	assert.NoError(t, pub.Shutdown(t.Context()))
}

func TestNewPublisherPayloadTooLarge(t *testing.T) {
	attrs := make([]attribute.KeyValue, 0, 2000)
	for i := range 2000 {
		attrs = append(attrs, attribute.String(
			fmt.Sprintf("key.%04d.padding.padding.padding", i),
			strings.Repeat("value", 15),
		))
	}
	r := makeResource(t, attrs...)
	payload, encErr := encodeProcessContext(r)
	require.NoError(t, encErr)
	require.Greater(t, len(payload), MaxPayloadSize,
		"encoded payload must exceed MaxPayloadSize; increase attribute count if this fails")
	_, err := NewPublisher(r)
	assert.Error(t, err)
}

// ---- Mapping discoverability -------------------------------------------

func TestPublisherMappingDiscoverable(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "discoverable"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)

	found := mappingInMaps()
	assert.NoError(t, pub.Shutdown(t.Context()))
	assert.True(t, found, "OTEL_CTX mapping should appear in /proc/self/maps")
}

func TestPublisherShutdownRemovesMapping(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "removable"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	require.True(t, mappingInMaps(), "mapping must be present before Shutdown")

	require.NoError(t, pub.Shutdown(t.Context()))
	assert.False(t, mappingInMaps(), "mapping must be absent after Shutdown")
}

// ---- Header correctness ------------------------------------------------

func TestPublisherHeaderSignatureAndVersion(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	hdr := pub.impl.headerMem
	require.GreaterOrEqual(t, len(hdr), headerSize)
	assert.Equal(t, []byte(signatureStr), hdr[0:8])
	assert.Equal(t, formatVersion, binary.LittleEndian.Uint32(hdr[offVersion:]))
}

func TestPublisherTimestampNonZero(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	before := monotonicNs()
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	ts := readTS(pub.impl.headerMem)
	after := monotonicNs()
	assert.GreaterOrEqual(t, ts, before, "timestamp must not predate publisher creation")
	assert.LessOrEqual(t, ts, after, "timestamp must not exceed current time")
}

func TestPublisherPayloadPtrPointsIntoMapping(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "ptr-test"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	hdr := pub.impl.headerMem
	payloadPtr := binary.LittleEndian.Uint64(hdr[offPayloadPtr:])
	hdrStart := uint64(uintptr(unsafe.Pointer(&hdr[0])))
	hdrEnd := hdrStart + uint64(len(hdr))

	// Small resource: payload fits inline, so PayloadPtr is within the header page.
	assert.GreaterOrEqual(t, payloadPtr, hdrStart+headerSize)
	assert.Less(t, payloadPtr, hdrEnd)
}

func TestPublisherPayloadSizeMatchesEncoded(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "size-test"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	expected, err := encodeProcessContext(r)
	require.NoError(t, err)
	sz := binary.LittleEndian.Uint32(pub.impl.headerMem[offPayloadSz:])
	assert.Equal(t, uint32(len(expected)), sz)
	assert.Equal(t, expected, readPayload(pub.impl))
}

// ---- Update ------------------------------------------------------------

func TestPublisherUpdate(t *testing.T) {
	r1 := makeResource(t, attribute.String("service.name", "v1"))
	pub, err := NewPublisher(r1)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	ts1 := readTS(pub.impl.headerMem)

	r2 := makeResource(t, attribute.String("service.name", "v2"))
	require.NoError(t, pub.Update(r2))

	ts2 := readTS(pub.impl.headerMem)
	assert.Greater(t, ts2, ts1, "timestamp must increase on update")

	expected, err := encodeProcessContext(r2)
	require.NoError(t, err)
	sz := binary.LittleEndian.Uint32(pub.impl.headerMem[offPayloadSz:])
	assert.Equal(t, uint32(len(expected)), sz)
	assert.Equal(t, expected, readPayload(pub.impl))
}

func TestPublisherUpdatePayloadTooLarge(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	attrs := make([]attribute.KeyValue, 0, 2000)
	for i := range 2000 {
		attrs = append(attrs, attribute.String(
			fmt.Sprintf("key.%04d.padding.padding.padding", i),
			strings.Repeat("value", 15),
		))
	}
	r2 := makeResource(t, attrs...)
	payload, encErr := encodeProcessContext(r2)
	require.NoError(t, encErr)
	require.Greater(t, len(payload), MaxPayloadSize,
		"encoded payload must exceed MaxPayloadSize; increase attribute count if this fails")

	err = pub.Update(r2)
	assert.Error(t, err)
}

func TestPublisherTimestampStrictlyIncreasing(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	prev := readTS(pub.impl.headerMem)
	for i := range 20 {
		require.NoError(t, pub.Update(r))
		ts := readTS(pub.impl.headerMem)
		assert.Greater(t, ts, prev, "iteration %d: timestamp must be strictly increasing", i)
		prev = ts
	}
}

// TestPublisherPayloadSpill verifies that when the updated payload no longer
// fits inline in the header page the publisher spills it to a separate mapping
// while keeping the header mapping address stable.
func TestPublisherPayloadSpill(t *testing.T) {
	r1 := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r1)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	headerAddr := uintptr(unsafe.Pointer(&pub.impl.headerMem[0]))
	assert.Nil(t, pub.impl.payloadMem, "initially no separate payload mapping")

	// Build a resource whose encoded payload exceeds the inline area.
	attrs := make([]attribute.KeyValue, 0, 500)
	for i := range 500 {
		attrs = append(attrs, attribute.String(
			fmt.Sprintf("key.%04d.padding.padding", i),
			strings.Repeat("v", 10),
		))
	}
	r2 := makeResource(t, attrs...)
	payload2, encErr := encodeProcessContext(r2)
	require.NoError(t, encErr)
	require.Greater(t, headerSize+len(payload2), len(pub.impl.headerMem),
		"payload must exceed the inline area; increase attribute count if this fails")

	require.NoError(t, pub.Update(r2))

	// Header mapping address must be unchanged.
	assert.Equal(t, headerAddr, uintptr(unsafe.Pointer(&pub.impl.headerMem[0])),
		"header mapping address must remain stable after payload spill")

	// A separate payload mapping must now exist.
	require.NotNil(t, pub.impl.payloadMem, "payload must spill to a separate mapping")
	assert.GreaterOrEqual(t, len(pub.impl.payloadMem), len(payload2))

	// PayloadPtr must point into the separate payload mapping.
	payloadPtr := binary.LittleEndian.Uint64(pub.impl.headerMem[offPayloadPtr:])
	pmStart := uint64(uintptr(unsafe.Pointer(&pub.impl.payloadMem[0])))
	pmEnd := pmStart + uint64(len(pub.impl.payloadMem))
	assert.GreaterOrEqual(t, payloadPtr, pmStart)
	assert.Less(t, payloadPtr, pmEnd)

	assert.NotZero(t, readTS(pub.impl.headerMem))
	got := binary.LittleEndian.Uint32(pub.impl.headerMem[offPayloadSz:])
	assert.Equal(t, uint32(len(payload2)), got)
}

// ---- Shutdown ----------------------------------------------------------

func TestPublisherShutdownIdempotent(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)

	assert.NoError(t, pub.Shutdown(t.Context()))
	// Second Shutdown must not panic or error.
	assert.NoError(t, pub.Shutdown(t.Context()))
}

func TestPublisherUpdateAfterShutdown(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	require.NoError(t, pub.Shutdown(t.Context()))

	assert.Error(t, pub.Update(r))
}

// ---- Concurrent safety -------------------------------------------------

func TestPublisherConcurrentUpdate(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "concurrent"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	var wg sync.WaitGroup
	errs := make([]error, 20)
	for i := range 20 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r := makeResource(t, attribute.String("i", fmt.Sprint(idx)))
			errs[idx] = pub.Update(r)
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		assert.NoError(t, err, "goroutine %d: unexpected error", i)
	}
	assert.NotZero(t, readTS(pub.impl.headerMem))
}

// ---- Internal helpers --------------------------------------------------

func TestPageAlign(t *testing.T) {
	page := pageAlign(1)
	assert.Positive(t, page)

	// Result must be a multiple of the system page size.
	sysPage := unix.Getpagesize()
	assert.Equal(t, 0, page%sysPage, "result must be page-aligned")

	// Already-aligned value should not change.
	assert.Equal(t, page, pageAlign(page))
}

func TestMonotonicNs(t *testing.T) {
	ts := monotonicNs()
	assert.NotZero(t, ts, "monotonicNs must return a non-zero value")

	ts2 := monotonicNs()
	assert.GreaterOrEqual(t, ts2, ts, "subsequent calls must not go backward")
}

func TestAtomicStore64(t *testing.T) {
	mem := make([]byte, 64)
	atomicStore64(mem, 16, 0xDEADBEEFCAFEBABE)
	got := binary.LittleEndian.Uint64(mem[16:])
	assert.Equal(t, uint64(0xDEADBEEFCAFEBABE), got)
}

// ---- Anonymous mapping fallback (when memfd_create is unavailable) --------

type noopOption struct{ applied bool }

func (o *noopOption) apply(_ *config) { o.applied = true }

func TestNewPublisherWithOption(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	opt := &noopOption{}
	pub, err := NewPublisher(r, opt)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck
	assert.True(t, opt.applied, "option must be applied via newConfig")
}

func TestAllocMappingAnonymousFallback(t *testing.T) {
	// Simulate a kernel that lacks memfd_create by injecting a failing stub.
	orig := memfdCreateFunc
	memfdCreateFunc = func(_ string, _ int) (int, error) { return -1, unix.ENOSYS }
	defer func() { memfdCreateFunc = orig }()

	r := makeResource(t, attribute.String("service.name", "anon-fallback"))
	pub, err := NewPublisher(r)
	if err != nil {
		// On kernels without CONFIG_ANON_VMA_NAME, prctl also fails and the
		// function correctly returns "not discoverable".
		t.Logf("anonymous fallback not available on this kernel: %v", err)
		return
	}
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	assert.True(t, mappingInMaps(), "anonymous+prctl mapping must appear in /proc/self/maps")
	assert.NotZero(t, readTS(pub.impl.headerMem))
}

// TestNextTSForcedIncrement covers the branch where the clock value does not
// advance between calls, forcing a counter increment.
func TestNextTSForcedIncrement(t *testing.T) {
	// Set lastTS to a value far in the future so monotonicNs() < lastTS.
	future := monotonicNs() + 1_000_000_000_000 // 1000 seconds from now
	p := &publisher{lastTS: future}

	ts1 := p.nextTS()
	assert.Equal(t, future+1, ts1, "must increment past lastTS when clock is behind")

	ts2 := p.nextTS()
	assert.Equal(t, future+2, ts2, "must keep incrementing")
}

// ---- Injectable-var error paths ----------------------------------------

// TestAllocMappingFtruncateError exercises the ftruncate-failure branch in
// allocMapping by injecting a ftruncate error.
func TestAllocMappingFtruncateError(t *testing.T) {
	orig := ftruncateFunc
	ftruncateFunc = func(_ int, _ int64) error { return unix.ENOSPC }
	defer func() { ftruncateFunc = orig }()

	_, _, err := allocMapping(pageAlign(1))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ftruncate")
}

// TestAllocMappingMmapMemfdError exercises the mmap-failure branch in
// allocMapping when memfd is available.
func TestAllocMappingMmapMemfdError(t *testing.T) {
	orig := mmapFunc
	mmapFunc = func(_ int, _ int64, _, _, _ int) ([]byte, error) {
		return nil, unix.ENOMEM
	}
	defer func() { mmapFunc = orig }()

	_, _, err := allocMapping(pageAlign(1))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mmap (memfd)")
}

// TestNewPublisherAllocMappingFails exercises the anonymous-mmap-failure branch
// in allocMapping and the resulting error return in newPublisher by making both
// memfd_create and mmap fail.
func TestNewPublisherAllocMappingFails(t *testing.T) {
	origMemfd := memfdCreateFunc
	memfdCreateFunc = func(_ string, _ int) (int, error) { return -1, unix.ENOSYS }
	origMmap := mmapFunc
	mmapFunc = func(_ int, _ int64, _, _, _ int) ([]byte, error) {
		return nil, unix.ENOMEM
	}
	defer func() {
		memfdCreateFunc = origMemfd
		mmapFunc = origMmap
	}()

	r := makeResource(t, attribute.String("k", "v"))
	_, err := NewPublisher(r)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mmap (anonymous)")
}

// TestUpdateSpillAllocFails exercises the allocMapping-failure branch in update
// when the payload overflows the inline area and the spill allocation fails.
func TestUpdateSpillAllocFails(t *testing.T) {
	r1 := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r1)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	// Build a resource large enough to exceed the inline area of the header page.
	attrs := make([]attribute.KeyValue, 0, 200)
	for i := range 200 {
		attrs = append(attrs, attribute.String(
			fmt.Sprintf("key.%04d", i),
			strings.Repeat("v", 20),
		))
	}
	r2 := makeResource(t, attrs...)
	payload, encErr := encodeProcessContext(r2)
	require.NoError(t, encErr)
	require.Greater(t, headerSize+len(payload), len(pub.impl.headerMem),
		"payload must exceed the inline area; increase attribute count if this fails")

	origMemfd := memfdCreateFunc
	memfdCreateFunc = func(_ string, _ int) (int, error) { return -1, unix.ENOSYS }
	origMmap := mmapFunc
	mmapFunc = func(_ int, _ int64, _, _, _ int) ([]byte, error) {
		return nil, unix.ENOMEM
	}
	defer func() {
		memfdCreateFunc = origMemfd
		mmapFunc = origMmap
	}()

	updateErr := pub.Update(r2)
	require.Error(t, updateErr)
	assert.Contains(t, updateErr.Error(), "spill")
}

// TestNameVMASuccess exercises the success return path in nameVMA (line 238)
// by injecting a prctl implementation that always reports success.
func TestNameVMASuccess(t *testing.T) {
	orig := syscall6Func
	syscall6Func = func(_, _, _, _, _, _, _ uintptr) (uintptr, uintptr, unix.Errno) {
		return 0, 0, 0 // success
	}
	defer func() { syscall6Func = orig }()

	mem := make([]byte, pageAlign(1))
	assert.NoError(t, nameVMA(mem))
}

// TestMonotonicNsFallbackToMonotonic exercises the CLOCK_MONOTONIC fallback
// branch in monotonicNs (lines 256-258) by injecting a failure for
// CLOCK_BOOTTIME.
func TestMonotonicNsFallbackToMonotonic(t *testing.T) {
	orig := clockGettimeFunc
	clockGettimeFunc = func(clockid int32, ts *unix.Timespec) error {
		if clockid == unix.CLOCK_BOOTTIME {
			return unix.EINVAL
		}
		return unix.ClockGettime(clockid, ts)
	}
	defer func() { clockGettimeFunc = orig }()

	ts := monotonicNs()
	assert.NotZero(t, ts, "CLOCK_MONOTONIC fallback must return a non-zero timestamp")
}

// TestMonotonicNsBothClocksFail exercises the final fallback in monotonicNs
// (line 260) by injecting failures for both CLOCK_BOOTTIME and CLOCK_MONOTONIC.
func TestMonotonicNsBothClocksFail(t *testing.T) {
	orig := clockGettimeFunc
	clockGettimeFunc = func(_ int32, _ *unix.Timespec) error { return unix.EINVAL }
	defer func() { clockGettimeFunc = orig }()

	assert.Equal(t, uint64(1), monotonicNs(), "must return 1 when both clocks fail")
}

// mappingInMaps reports whether an OTEL_CTX named mapping is present in
// /proc/self/maps right now.
func mappingInMaps() bool {
	f, err := os.Open("/proc/self/maps")
	if err != nil {
		return false
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.Contains(sc.Text(), "OTEL_CTX") {
			return true
		}
	}
	return false
}
