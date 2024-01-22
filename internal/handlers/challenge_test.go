package handlers_test

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/et-nik/word-of-wisdom/internal/challenger"
	"github.com/et-nik/word-of-wisdom/internal/domain"
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
					challenger.NewChallengeCache(10*time.Second),
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
					challenger.NewChallengeCache(10*time.Second),
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

func TestChallengeHandlerRequest(t *testing.T) {
	// ARRANGE
	cache := &cacheMock{cache: make(map[string]domain.Challenge)}
	h := handlers.NewChallengeHandler(
		challenger.New(),
		cache,
		4,
		2,
	)
	msg := server.Message{
		Type: "challenge",
	}
	writer := bytes.NewBuffer(nil)

	// ACT
	h.Handle(context.Background(), msg, writer)

	// ASSERT
	assert.Contains(t, writer.String(), "4 2")
	assert.Len(t, cache.cache, 1)
	for _, v := range cache.cache {
		assert.Equal(t, 4, v.N)
		assert.Equal(t, 2, v.K)
		break
	}
}

type cacheMock struct {
	mu    sync.Mutex
	cache map[string]domain.Challenge
}

func (c *cacheMock) Get(key string) (domain.Challenge, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.cache[key]
	return v, ok
}

func (c *cacheMock) Set(key string, ch domain.Challenge) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = ch
}

func (c *cacheMock) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}
