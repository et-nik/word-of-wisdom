package challenger

import (
	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/pkg/equihash"
)

type Verifier struct{}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(challenge domain.Challenge, solution domain.Solution) bool {
	proof := equihash.NewProof(
		challenge.N,
		challenge.K,
		challenge.Seed,
		solution.Nonce,
		solution.Hash,
	)

	return proof.Test()
}
