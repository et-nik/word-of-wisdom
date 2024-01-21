package handlers

import (
	"context"
	"fmt"
	"io"

	"github.com/et-nik/word-of-wisdom/internal/server"
)

type ChallengeHandler struct {
	challenger Challenger
	n          int
	k          int
}

func NewChallengeHandler(challenger Challenger, n, k int) *ChallengeHandler {
	return &ChallengeHandler{
		challenger: challenger,
		n:          n,
		k:          k,
	}
}

func (h *ChallengeHandler) Handle(_ context.Context, _ server.Message, w io.Writer) {
	ch := h.challenger.Prepare(h.n, h.k)
	_, _ = w.Write([]byte(fmt.Sprintf("%d %d %x", ch.N, ch.K, ch.Seed)))
}
