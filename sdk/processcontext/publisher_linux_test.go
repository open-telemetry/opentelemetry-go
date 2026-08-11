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
	if len(payload) <= MaxPayloadSize {
		t.Skip("encoded payload fits; increase attribute count to trigger the error path")
	}
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

	mem := pub.impl.mem
	require.GreaterOrEqual(t, len(mem), headerSize)
	assert.Equal(t, []byte(signatureStr), mem[0:8])
	assert.Equal(t, formatVersion, binary.LittleEndian.Uint32(mem[offVersion:]))
}

func TestPublisherTimestampNonZero(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	assert.NotZero(t, readTS(pub.impl.mem))
}

func TestPublisherPayloadPtrPointsIntoMapping(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "ptr-test"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	mem := pub.impl.mem
	payloadPtr := binary.LittleEndian.Uint64(mem[offPayloadPtr:])
	mappingStart := uint64(uintptr(unsafe.Pointer(&mem[0])))
	mappingEnd := mappingStart + uint64(len(mem))

	assert.GreaterOrEqual(t, payloadPtr, mappingStart+headerSize)
	assert.Less(t, payloadPtr, mappingEnd)
}

func TestPublisherPayloadSizeMatchesEncoded(t *testing.T) {
	r := makeResource(t, attribute.String("service.name", "size-test"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	expected, err := encodeProcessContext(r)
	require.NoError(t, err)
	got := binary.LittleEndian.Uint32(pub.impl.mem[offPayloadSz:])
	assert.Equal(t, uint32(len(expected)), got)
}

// ---- Update ------------------------------------------------------------

func TestPublisherUpdate(t *testing.T) {
	r1 := makeResource(t, attribute.String("service.name", "v1"))
	pub, err := NewPublisher(r1)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	ts1 := readTS(pub.impl.mem)

	r2 := makeResource(t, attribute.String("service.name", "v2"))
	require.NoError(t, pub.Update(r2))

	ts2 := readTS(pub.impl.mem)
	assert.Greater(t, ts2, ts1, "timestamp must increase on update")

	expected, err := encodeProcessContext(r2)
	require.NoError(t, err)
	sz := binary.LittleEndian.Uint32(pub.impl.mem[offPayloadSz:])
	assert.Equal(t, uint32(len(expected)), sz)
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
	if len(payload) <= MaxPayloadSize {
		t.Skip("payload fits; increase attribute count")
	}

	err = pub.Update(r2)
	assert.Error(t, err)
}

func TestPublisherTimestampStrictlyIncreasing(t *testing.T) {
	r := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	prev := readTS(pub.impl.mem)
	for i := range 20 {
		require.NoError(t, pub.Update(r))
		ts := readTS(pub.impl.mem)
		assert.Greater(t, ts, prev, "iteration %d: timestamp must be strictly increasing", i)
		prev = ts
	}
}

// TestPublisherRemap verifies that a publisher correctly reallocates its
// mapping when the updated payload no longer fits within the original region.
func TestPublisherRemap(t *testing.T) {
	// Start with a minimal resource to get the smallest possible mapping.
	r1 := makeResource(t, attribute.String("k", "v"))
	pub, err := NewPublisher(r1)
	require.NoError(t, err)
	defer pub.Shutdown(t.Context()) //nolint:errcheck

	initialAddr := uintptr(unsafe.Pointer(&pub.impl.mem[0]))
	initialLen := len(pub.impl.mem)

	// Build a resource whose encoded payload exceeds the initial mapping.
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
	if headerSize+len(payload2) <= initialLen {
		t.Skip("payload fits in current mapping; no remap will occur — try larger attribute count")
	}

	require.NoError(t, pub.Update(r2))

	newAddr := uintptr(unsafe.Pointer(&pub.impl.mem[0]))
	assert.NotEqual(t, initialAddr, newAddr, "remap should produce a new mapping address")
	assert.GreaterOrEqual(t, len(pub.impl.mem), headerSize+len(payload2))
	assert.NotZero(t, readTS(pub.impl.mem))

	expected := uint32(len(payload2))
	got := binary.LittleEndian.Uint32(pub.impl.mem[offPayloadSz:])
	assert.Equal(t, expected, got)
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
	assert.NotZero(t, readTS(pub.impl.mem))
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
	assert.NotZero(t, readTS(pub.impl.mem))
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
