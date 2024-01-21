package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/config"
	"github.com/et-nik/word-of-wisdom/internal/handlers"
	"github.com/et-nik/word-of-wisdom/internal/server"
	"github.com/et-nik/word-of-wisdom/pkg/quotes"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		//nolint:gocritic
		log.Fatal(err)
	}

	s := server.NewServer(cfg)

	ch := challenger.New()
	verifier := challenger.NewVerifier()
	quoter := quotes.New()

	s.RegisterHandler(
		"challenge",
		handlers.NewChallengeHandler(ch, cfg.DifficultyWidth, cfg.DifficultyLength).Handle,
	)
	s.RegisterHandler("words-of-wisdom", handlers.NewWordsOfWisdomHandler(verifier, quoter).Handle)

	go func() {
		err := s.Run(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()
	err = s.Stop()
	if err != nil {
		log.Println(err)
	}
}
