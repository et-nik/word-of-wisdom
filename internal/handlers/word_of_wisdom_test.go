package handlers

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestWordsOfWisdomHandler(t *testing.T) {
	tests := []struct {
		name           string
		handlerFunc    func() *WordsOfWisdomHandler
		messagePayload string
		expected       string
	}{
		{
			name: "success",
			handlerFunc: func() *WordsOfWisdomHandler {
				q := mockQuoter{
					quote: "test",
				}

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
				)
			},
			messagePayload: "10 2 981a2b4f987a84a913bedf8ba067d2915454680e75351f5b4295c7db0c4211bd " +
				"2 0100000002000000050000000c000000",
			expected: "test",
		},
		{
			name: "quoter error",
			handlerFunc: func() *WordsOfWisdomHandler {
				q := mockQuoter{
					err: errors.New("quoter error"),
				}

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
				)
			},
			messagePayload: "10 2 981a2b4f987a84a913bedf8ba067d2915454680e75351f5b4295c7db0c4211bd " +
				"2 0100000002000000050000000c000000",
			expected: "I don't know what to say",
		},
		{
			name: "invalid message",
			handlerFunc: func() *WordsOfWisdomHandler {
				q := mockQuoter{
					quote: "test",
				}

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
				)
			},
			messagePayload: "abrakadabra",
			expected:       "0 invalid message",
		},
		{
			name: "success",
			handlerFunc: func() *WordsOfWisdomHandler {
				q := mockQuoter{
					quote: "test",
				}

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
				)
			},
			messagePayload: "10 2 981a2b4f987a84a913bedf8ba067d2915454680e75351f5b4295c7db0c4211bd " +
				"2 01000000000000000000000000000000",
			expected: "0 invalid solution",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// ARRANGE
			h := test.handlerFunc()
			msg := server.Message{
				Type:    "words_of_wisdom",
				Payload: test.messagePayload,
			}
			writer := bytes.NewBuffer(nil)

			// ACT
			h.Handle(context.Background(), msg, writer)

			// ASSERT
			assert.Equal(t, test.expected, writer.String())
		})
	}
}

type mockQuoter struct {
	quote string
	err   error
}

func (m mockQuoter) Quote(_ context.Context) (string, error) {
	return m.quote, m.err
}
