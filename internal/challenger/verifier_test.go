package challenger_test

import (
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/domain"
)

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
