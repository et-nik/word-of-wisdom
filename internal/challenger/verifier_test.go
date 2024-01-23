package challenger_test

import (
	"encoding/binary"
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/pkg/equihash"
	"github.com/stretchr/testify/assert"
)

func TestNewVerifier_Positive(t *testing.T) {
	v := challenger.NewVerifier()

	result := v.Verify(
		domain.Challenge{
			N:    52,
			K:    3,
			Seed: equihash.Uint32ArrayToBytes([]uint32{1, 1, 1, 1}, binary.LittleEndian),
		},
		domain.Solution{
			Nonce: 2,
			Hash:  equihash.Uint32ArrayToBytes([]uint32{0x611, 0x1626, 0x1c37, 0x20cb, 0x241d, 0x30d7, 0x3811, 0x395c}, binary.LittleEndian),
		},
	)

	assert.True(t, result)
}

func TestNewVerifier_Negative(t *testing.T) {
	v := challenger.NewVerifier()

	result := v.Verify(
		domain.Challenge{
			N:    54,
			K:    2,
			Seed: []byte{4, 1, 2, 3, 41, 1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		domain.Solution{
			Nonce: 2,
			Hash:  []byte{1, 2, 3, 4},
		},
	)

	assert.False(t, result)
}

func TestNewVerifier_NegativeEmpty(t *testing.T) {
	v := challenger.NewVerifier()

	result := v.Verify(
		domain.Challenge{
			N:    54,
			K:    2,
			Seed: []byte{4, 1, 2, 3, 41, 1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		domain.Solution{
			Nonce: 2,
			Hash:  []byte{},
		},
	)

	assert.False(t, result)
}

func BenchmarkVerifier_Verify(b *testing.B) {
	b.ReportAllocs()

	v := challenger.NewVerifier()

	for i := 0; i < b.N; i++ {
		v.Verify(
			domain.Challenge{
				N:    54,
				K:    2,
				Seed: []byte{4, 1, 2, 3, 41, 1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			},
			domain.Solution{
				Nonce: 2,
				Hash:  []byte{},
			},
		)
	}
}
