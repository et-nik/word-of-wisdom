package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/et-nik/word-of-wisdom/internal/server"
)

type ChallengeHandler struct {
	challenger Challenger
	challenges ChallengeRepository

	n int
	k int
}

func NewChallengeHandler(challenger Challenger, repository ChallengeRepository, n, k int) *ChallengeHandler {
	return &ChallengeHandler{
		challenger: challenger,
		challenges: repository,
		n:          n,
		k:          k,
	}
}

func (h *ChallengeHandler) Handle(_ context.Context, _ server.Message, w io.Writer) {
	ch := h.challenger.Prepare(h.n, h.k)

	requestID, err := generateRandomString(24)
	if err != nil {
		log.Println("failed to generate random string, err:", err)
		_, _ = w.Write([]byte("error"))
	}

	h.challenges.Set(requestID, ch)

	_, _ = w.Write([]byte(fmt.Sprintf("%s %d %d %x", requestID, ch.N, ch.K, ch.Seed)))
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
