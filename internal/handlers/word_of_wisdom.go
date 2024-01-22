package handlers

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/internal/server"
)

type WordsOfWisdomHandler struct {
	verifier   Verifier
	quoter     Quoter
	challenges ChallengeRepository
}

func NewWordsOfWisdomHandler(verifier Verifier, quoter Quoter, repository ChallengeRepository) *WordsOfWisdomHandler {
	return &WordsOfWisdomHandler{
		verifier:   verifier,
		quoter:     quoter,
		challenges: repository,
	}
}

func (h *WordsOfWisdomHandler) Handle(ctx context.Context, msg server.Message, w io.Writer) {
	log.Println("WordsOfWisdomHandler")

	requestID, sl, err := parseMessage(msg)
	if err != nil {
		_, _ = w.Write([]byte("0 invalid message"))
		return
	}

	ch, ok := h.challenges.Get(requestID)
	if !ok {
		_, _ = w.Write([]byte("0 get a solution first"))
		return
	}

	h.challenges.Delete(requestID)

	if ok := h.verifier.Verify(ch, sl); !ok {
		_, _ = w.Write([]byte("0 invalid solution"))
		return
	}

	quote, err := h.quoter.Quote(ctx)
	if err != nil {
		// Graceful degradation
		quote = "I don't know what to say"
	}

	_, _ = w.Write([]byte(quote))
}

func parseMessage(msg server.Message) (string, domain.Solution, error) {
	var requestID string
	var sl domain.Solution

	_, err := fmt.Sscanf(msg.Payload, "%s %d %x", &requestID, &sl.Nonce, &sl.Hash)
	if err != nil {
		return "", domain.Solution{}, err
	}

	return requestID, sl, nil
}
