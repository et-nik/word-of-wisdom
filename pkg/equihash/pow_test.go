package equihash

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: Improve test.
func Test(t *testing.T) {
	// ARRANGE
	start := time.Now()
	e := New(52, 3, []byte{4, 1, 2, 3, 41, 1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	//e := New(50, 2, [16]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})

	//p.seed = [SeedLength]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16}

	// ACT
	p := e.FindProof()
	result := p.Test()

	// ASSERT
	assert.True(t, result)
	elapsed := time.Since(start)
	log.Printf("Time elapsed: %s", elapsed)
}
