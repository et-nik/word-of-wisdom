package challenger_test

import (
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/stretchr/testify/assert"
)

func TestChallenger(t *testing.T) {
	ch := challenger.New()

	result := ch.Prepare(54, 2)

	assert.Equal(t, "equihash", result.Algorithm)
	assert.Equal(t, 54, result.N)
	assert.Equal(t, 2, result.K)
	assert.Len(t, result.Seed, 32)
}
