package handlers

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/internal/server"
	"github.com/pkg/errors"
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

				c := challenger.NewChallengeCache(10 * time.Second)
				c.Set("rEquEstId", domain.Challenge{
					N: 10,
					K: 2,
					Seed: []byte{
						39, 17, 226, 171, 209, 215, 180, 83, 223, 148, 34, 174, 237, 73, 199, 107, 254, 76,
						241, 92, 162, 129, 237, 58, 181, 219, 82, 238, 29, 234, 153, 141,
					},
				})

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
					c,
				)
			},
			messagePayload: "rEquEstId 2 01000000050000000c0000000f000000",
			expected:       "test",
		},
		{
			name: "quoter error",
			handlerFunc: func() *WordsOfWisdomHandler {
				q := mockQuoter{
					err: errors.New("some error"),
				}

				c := challenger.NewChallengeCache(10 * time.Second)
				c.Set("rEquEstId", domain.Challenge{
					N: 10,
					K: 2,
					Seed: []byte{
						39, 17, 226, 171, 209, 215, 180, 83, 223, 148, 34, 174, 237, 73, 199, 107, 254, 76,
						241, 92, 162, 129, 237, 58, 181, 219, 82, 238, 29, 234, 153, 141,
					},
				})

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
					c,
				)
			},
			messagePayload: "rEquEstId 2 01000000050000000c0000000f000000",
			expected:       "I don't know what to say",
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
					challenger.NewChallengeCache(10*time.Second),
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

				c := challenger.NewChallengeCache(10 * time.Second)
				c.Set("rEquEstId", domain.Challenge{
					N: 10,
					K: 2,
					Seed: []byte{
						20, 17, 226, 171, 209, 215, 180, 83, 223, 148, 34, 174, 237, 73, 199, 107, 254, 76,
						241, 92, 162, 129, 237, 58, 181, 219, 82, 238, 29, 234, 153, 141,
					},
				})

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
					c,
				)
			},
			messagePayload: "rEquEstId 2 01000000050000000c0000000f000000",
			expected:       "0 invalid solution",
		},
		{
			name: "no_saved_challenges",
			handlerFunc: func() *WordsOfWisdomHandler {
				q := mockQuoter{
					quote: "test",
				}

				return NewWordsOfWisdomHandler(
					challenger.NewVerifier(),
					q,
					challenger.NewChallengeCache(10*time.Second),
				)
			},
			messagePayload: "rEquEstId 2 01000000050000000c0000000f000000",
			expected:       "0 get a solution first",
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
