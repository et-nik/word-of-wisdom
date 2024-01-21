package challenger

import (
	"crypto/rand"

	"github.com/et-nik/word-of-wisdom/internal/domain"
)

type Challenger struct {
}

func New() *Challenger {
	return &Challenger{}
}

func (c *Challenger) Prepare(n int, k int) domain.Challenge {
	seed := make([]byte, 32)
	_, _ = rand.Read(seed)

	return domain.Challenge{
		Algorithm: "equihash",
		N:         n,
		K:         k,
		Seed:      seed,
	}
}
