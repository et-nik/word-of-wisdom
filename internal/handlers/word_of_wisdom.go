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
	verifier Verifier
	quoter   Quoter
}

func NewWordsOfWisdomHandler(verifier Verifier, quoter Quoter) *WordsOfWisdomHandler {
	return &WordsOfWisdomHandler{
		verifier: verifier,
		quoter:   quoter,
	}
}

func (h *WordsOfWisdomHandler) Handle(ctx context.Context, msg server.Message, w io.Writer) {
	log.Println("WordsOfWisdomHandler")

	ch, sl, err := parseMessage(msg)
	if err != nil {
		_, _ = w.Write([]byte("0 invalid message"))
		return
	}

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

func parseMessage(msg server.Message) (domain.Challenge, domain.Solution, error) {
	var ch domain.Challenge
	var sl domain.Solution

	_, err := fmt.Sscanf(msg.Payload, "%d %d %x %d %x", &ch.N, &ch.K, &ch.Seed, &sl.Nonce, &sl.Hash)
	if err != nil {
		return domain.Challenge{}, domain.Solution{}, err
	}

	return ch, sl, nil
}
