package memory

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/log"
)

const (
	// KiB refers to kibibyte.
	KiB = 1 << (10 * (iota + 1))
	// MiB refers to mebibyte.
	MiB = 1 << (10 * (iota + 1))
	// GiB refers to gibibyte.
	GiB = 1 << (10 * (iota + 1))
	// TiB refers to tebibyte.
	TiB = 1 << (10 * (iota + 1))
	// PiB refers to pebibyte.
	PiB = 1 << (10 * (iota + 1))
	// EiB refers to exbibyte.
	EiB = 1 << (10 * (iota + 1))
)

const (
	// AllocID is the identifier of the attack
	AllocID = "memory_allocation"

	// Options
	sizeKey = "size"
)

// Register the creator of the attack
func init() {
	attack.Register(AllocID, attack.CreatorFunc(func(o attack.Opts) (attack.Attacker, error) {
		return NewMemAllocationOpts(o)
	}))
}

// MemAllocation failer will apply a failure allocating memory.
type MemAllocation struct {
	Size uint64 // The size of the buffer in bytes.

	b    bytes.Buffer //  Structure that will allocate memory.
	log  log.Logger   // Logger.
	done bool         // Flag that marks the allocation was already made.
}

// NewMemAllocationOpts returns a new default memory allocation failer using options.
func NewMemAllocationOpts(opts attack.Opts) (*MemAllocation, error) {
	size, ok := opts[sizeKey].(int)
	if !ok {
		return nil, fmt.Errorf("invalid '%s' option with '%v' value", sizeKey, size)
	}

	return NewMemAllocation(uint64(size))
}

// NewMemAllocation returns a new default memory allocation failer.
func NewMemAllocation(size uint64) (*MemAllocation, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size can't be 0 or less")
	}

	return &MemAllocation{
		Size: size,
		log:  log.Base(),
	}, nil
}

// Apply will allocate and use the memory.
func (m *MemAllocation) Apply(ctx context.Context) error {
	// Check if its done before continuing.
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	if m.done {
		return errors.New("memory allocation already applied")
	}
	log.With("MiB", m.Size/MiB).Infof("Allocated memory")
	m.b.Write(make([]byte, m.Size))
	m.done = true
	return nil
}

// Revert will clean the memory usage letting the GC do its work.
func (m *MemAllocation) Revert() error {
	// Reset our buffer to be able to free the memory.
	m.b.Reset()
	m.b = bytes.Buffer{}

	// Return memory to the OS (Way better to do?)
	debug.FreeOSMemory()
	log.With("MiB", m.Size/MiB).Infof("Reverted allocated memory")
	m.done = false
	return nil
}
