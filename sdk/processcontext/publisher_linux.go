// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build linux

package processcontext

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/unix"

	"go.opentelemetry.io/otel/sdk/resource"
)

// Header layout — 32 bytes
//
//	[0:8]   Signature              [8]byte  "OTEL_CTX"
//	[8:12]  Version                uint32   2
//	[12:16] PayloadSize            uint32
//	[16:24] MonotonicPublishedAtNs uint64   written last (seqlock)
//	[24:32] PayloadPtr             uint64   virtual address of payload
const (
	headerSize    = 32
	signatureStr  = "OTEL_CTX"
	formatVersion = uint32(2)

	offVersion    = 8
	offPayloadSz  = 12
	offTimestamp  = 16
	offPayloadPtr = 24

	// PR_SET_VMA names an anonymous VMA region (Linux 5.17+ CONFIG_ANON_VMA_NAME).
	prSetVMA     = 0x53564d41
	prSetVMAAnon = 0

	// MFD_NOEXEC_SEAL seals the memfd against making it executable (Linux 6.3+).
	mfdNoExecSeal = 0x8
)

type publisher struct {
	mu     sync.Mutex
	mem    []byte
	lastTS uint64
	closed bool
}

func newPublisher(r *resource.Resource, _ config) (*publisher, error) {
	payload, err := encodeProcessContext(r)
	if err != nil {
		return nil, err
	}
	if len(payload) > MaxPayloadSize {
		return nil, fmt.Errorf(
			"processcontext: payload %d bytes exceeds MaxPayloadSize (%d)",
			len(payload), MaxPayloadSize,
		)
	}

	size := pageAlign(headerSize + len(payload))
	mem, hasMemfd, err := allocMapping(size)
	if err != nil {
		return nil, err
	}

	p := &publisher{mem: mem}
	p.writeInitial(payload)

	// nameVMA is called last: a profiler that hooks prctl to detect publication
	// must see a fully-initialised mapping.
	prctlErr := nameVMA(mem)
	if !hasMemfd && prctlErr != nil {
		_ = unix.Munmap(mem)
		return nil, errors.New(
			"processcontext: mapping not discoverable: memfd_create and prctl PR_SET_VMA both failed",
		)
	}

	return p, nil
}

func (p *publisher) update(r *resource.Resource) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("processcontext: publisher is shut down")
	}

	payload, err := encodeProcessContext(r)
	if err != nil {
		return err
	}
	if len(payload) > MaxPayloadSize {
		return fmt.Errorf("processcontext: payload %d bytes exceeds MaxPayloadSize (%d)", len(payload), MaxPayloadSize)
	}

	if headerSize+len(payload) > len(p.mem) {
		if err := p.remap(len(payload)); err != nil {
			return err
		}
	}

	// Seqlock update protocol: zero timestamp → write fields → new timestamp.
	atomicStore64(p.mem, offTimestamp, 0)
	p.writePayload(payload)
	atomicStore64(p.mem, offTimestamp, p.nextTS())

	_ = nameVMA(p.mem) // re-name per spec update protocol; failure is ignored
	return nil
}

func (p *publisher) shutdown() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}
	p.closed = true
	atomicStore64(p.mem, offTimestamp, 0)
	_ = unix.Munmap(p.mem)
	p.mem = nil
}

// writeInitial populates the header and payload on a freshly-zeroed mapping.
// Timestamp is written last to atomically signal readiness to readers.
func (p *publisher) writeInitial(payload []byte) {
	copy(p.mem[0:8], signatureStr)
	binary.LittleEndian.PutUint32(p.mem[offVersion:], formatVersion)
	p.writePayload(payload)
	atomicStore64(p.mem, offTimestamp, p.nextTS())
}

// writePayload updates PayloadSize, PayloadPtr, and the payload bytes.
// Caller must hold p.mu and ensure the mapping is large enough.
func (p *publisher) writePayload(payload []byte) {
	// len(payload) is always <= MaxPayloadSize (65536), checked by callers.
	sz := uint32(len(payload)) //nolint:gosec // G115: bounded by MaxPayloadSize check
	binary.LittleEndian.PutUint32(p.mem[offPayloadSz:], sz)
	// PayloadPtr is a virtual address in this process that external readers use.
	addr := uint64(uintptr(unsafe.Pointer(&p.mem[headerSize])))
	atomicStore64(p.mem, offPayloadPtr, addr)
	copy(p.mem[headerSize:], payload)
}

// remap replaces the current mapping with a larger one capable of holding
// payloadLen bytes. Allocates the new mapping before freeing the old one so
// that the old mapping remains readable until the swap is committed. Must be
// called with p.mu held.
func (p *publisher) remap(payloadLen int) error {
	size := pageAlign(headerSize + payloadLen)
	newMem, hasMemfd, err := allocMapping(size)
	if err != nil {
		return fmt.Errorf("processcontext: remap: %w", err)
	}
	prctlErr := nameVMA(newMem)
	if !hasMemfd && prctlErr != nil {
		_ = unix.Munmap(newMem)
		return errors.New("processcontext: remap: mapping not discoverable")
	}

	// Zero old timestamp before discarding; new mapping already has ts=0.
	atomicStore64(p.mem, offTimestamp, 0)
	_ = unix.Munmap(p.mem)

	// Initialize static header fields on the new mapping.
	copy(newMem[0:8], signatureStr)
	binary.LittleEndian.PutUint32(newMem[offVersion:], formatVersion)
	p.mem = newMem
	return nil
}

// nextTS returns a monotonically strictly increasing timestamp. Must be called
// with p.mu held (or before p is exposed to other goroutines).
func (p *publisher) nextTS() uint64 {
	ts := monotonicNs()
	if ts <= p.lastTS {
		ts = p.lastTS + 1
	}
	p.lastTS = ts
	return ts
}

// memfdCreateFunc is the syscall used to create a sealed file descriptor.
// Overridable in tests to exercise the anonymous-mapping fallback path.
var memfdCreateFunc = unix.MemfdCreate

// ftruncateFunc sets the size of a file descriptor. Overridable in tests.
var ftruncateFunc = unix.Ftruncate

// mmapFunc creates a memory mapping. Overridable in tests.
var mmapFunc = unix.Mmap

// clockGettimeFunc reads a clock value. Overridable in tests to exercise
// clock-fallback paths in monotonicNs.
var clockGettimeFunc = unix.ClockGettime

// syscall6Func issues a raw 6-argument syscall. Overridable in tests to
// exercise the nameVMA success and error paths without kernel dependency.
var syscall6Func = unix.Syscall6

// allocMapping creates the mmap'd region. It tries memfd first; on failure
// falls back to a MAP_PRIVATE|MAP_ANONYMOUS mapping. Returns whether memfd
// was used so the caller can decide whether prctl is required for
// discoverability.
func allocMapping(size int) (mem []byte, hasMemfd bool, err error) {
	fd, fdErr := memfdCreateFunc("OTEL_CTX", unix.MFD_CLOEXEC|unix.MFD_ALLOW_SEALING|mfdNoExecSeal)
	if fdErr != nil {
		fd, fdErr = memfdCreateFunc("OTEL_CTX", unix.MFD_CLOEXEC|unix.MFD_ALLOW_SEALING)
	}
	if fdErr == nil {
		hasMemfd = true
		if err = ftruncateFunc(fd, int64(size)); err != nil {
			_ = unix.Close(fd)
			return nil, false, fmt.Errorf("processcontext: ftruncate: %w", err)
		}
		mem, err = mmapFunc(fd, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE)
		_ = unix.Close(fd)
		if err != nil {
			return nil, false, fmt.Errorf("processcontext: mmap (memfd): %w", err)
		}
	} else {
		mem, err = mmapFunc(-1, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANONYMOUS)
		if err != nil {
			return nil, false, fmt.Errorf("processcontext: mmap (anonymous): %w", err)
		}
	}
	// MADV_DONTFORK prevents child processes from inheriting a stale copy of
	// this mapping after fork. The Go runtime does not survive a bare fork
	// (without exec), so a failure here is not treated as a publish failure
	// for Go processes. If fork support is ever added to Go, or if this code
	// is ported to a language that does support fork, this should become a
	// hard error so that forked processes do not expose an inconsistent view
	// of the context.
	_ = unix.Madvise(mem, unix.MADV_DONTFORK)
	return mem, hasMemfd, nil
}

// nameVMA calls prctl(PR_SET_VMA, PR_SET_VMA_ANON_NAME) to name the mapping
// so external readers can discover it via /proc/<pid>/maps. Requires Linux
// 5.17+ with CONFIG_ANON_VMA_NAME.
func nameVMA(mem []byte) error {
	name := append([]byte("OTEL_CTX"), 0) // NUL-terminated
	_, _, errno := syscall6Func(
		unix.SYS_PRCTL,
		prSetVMA,
		prSetVMAAnon,
		uintptr(unsafe.Pointer(&mem[0])),
		uintptr(len(mem)),
		uintptr(unsafe.Pointer(&name[0])),
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}

// atomicStore64 writes v to mem[off:off+8] with sequential-consistency
// semantics. mem[off] must be 8-byte aligned (guaranteed for offTimestamp=16
// and offPayloadPtr=24 when mem is page-aligned from mmap).
func atomicStore64(mem []byte, off int, v uint64) {
	atomic.StoreUint64((*uint64)(unsafe.Pointer(&mem[off])), v)
}

// monotonicNs returns the current CLOCK_BOOTTIME value in nanoseconds,
// falling back to CLOCK_MONOTONIC if unavailable.
func monotonicNs() uint64 {
	var ts unix.Timespec
	if clockGettimeFunc(unix.CLOCK_BOOTTIME, &ts) == nil {
		// Sec and Nsec are always non-negative for elapsed-time clocks.
		return uint64(ts.Sec)*1_000_000_000 + uint64(ts.Nsec) //nolint:gosec // G115: non-negative clock values
	}
	if clockGettimeFunc(unix.CLOCK_MONOTONIC, &ts) == nil {
		return uint64(ts.Sec)*1_000_000_000 + uint64(ts.Nsec) //nolint:gosec // G115: non-negative clock values
	}
	// This path should never be reached on any modern Linux kernel.
	return 1
}

func pageAlign(n int) int {
	page := unix.Getpagesize()
	return (n + page - 1) &^ (page - 1)
}
