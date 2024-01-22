package handlers_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/handlers"
	"github.com/et-nik/word-of-wisdom/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestChallengeHandler(t *testing.T) {
	tests := []struct {
		name        string
		handlerFunc func() *handlers.ChallengeHandler
		expected    string
	}{
		{
			name: "challenge handler",
			handlerFunc: func() *handlers.ChallengeHandler {
				return handlers.NewChallengeHandler(
					challenger.New(),
					4,
					2,
				)
			},
			expected: "4 2",
		},
		{
			name: "challenge handler",
			handlerFunc: func() *handlers.ChallengeHandler {
				return handlers.NewChallengeHandler(
					challenger.New(),
					10,
					3,
				)
			},
			expected: "10 3",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// ARRANGE
			h := test.handlerFunc()
			msg := server.Message{
				Type: "challenge",
			}
			writer := bytes.NewBuffer(nil)

			// ACT
			h.Handle(context.Background(), msg, writer)

			// ASSERT
			assert.Contains(t, writer.String(), test.expected)
		})
	}
}
