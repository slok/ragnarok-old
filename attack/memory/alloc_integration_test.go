// +build integration

package memory

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryAllocationAttack(t *testing.T) {
	assert := assert.New(t)
	var size uint64 = 200 * MiB

	ma, err := NewMemAllocation(size)
	assert.NoError(err, "Creation of memory allocator shouldn't error")

	// Get current memory
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	startMem := mem.Alloc

	// Allocate memory and test if increased.
	ma.Apply(context.TODO())
	time.Sleep(1 * time.Millisecond)
	runtime.ReadMemStats(&mem)
	endMem := mem.Alloc

	// Let 10% margin delta from the wanted size
	assert.InDelta((endMem - startMem), size, float64(size)*0.15, "current memory allocation should be wanted allocation (15% deviation)")
	// Free memory and test if released.
	ma.Revert()
	time.Sleep(1 * time.Millisecond)
	runtime.ReadMemStats(&mem)

	// Let 10% margin delta from the wanted size
	assert.InDelta(startMem, mem.Alloc, float64(size)*0.15, "current memory and initial memory should be equal (15% deviation)")
}
