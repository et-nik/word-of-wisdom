package handlers

import (
	"context"

	"github.com/et-nik/word-of-wisdom/internal/domain"
)

type Challenger interface {
	Prepare(n int, k int) domain.Challenge
}

type Verifier interface {
	Verify(challenge domain.Challenge, solution domain.Solution) bool
}

type Quoter interface {
	Quote(ctx context.Context) (string, error)
}
