package solver

import (
	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/pkg/equihash"
)

type Solver struct {
}

func New() *Solver {
	return &Solver{}
}

func (s *Solver) Solve(challenge domain.Challenge) (domain.Solution, error) {
	eq := equihash.New(
		challenge.N,
		challenge.K,
		challenge.Seed,
	)

	proof := eq.FindProof()

	return domain.Solution{
		Nonce: int(proof.Nonce()),
		Hash:  proof.InputsBytes(),
	}, nil
}
