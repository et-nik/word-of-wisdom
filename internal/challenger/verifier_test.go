package challenger_test

import (
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewVerifier_Positive(t *testing.T) {
	v := challenger.NewVerifier()

	result := v.Verify(
		domain.Challenge{
			N:    54,
			K:    2,
			Seed: []byte{4, 1, 2, 3, 41, 1, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		domain.Solution{
			Nonce: 2,
			Hash:  []byte{0x31, 0x3a, 0x1, 0x0, 0x67, 0xa9, 0x2, 0x0, 0xe7, 0x42, 0x4, 0x0, 0x5e, 0xd9, 0x6, 0x0},
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
