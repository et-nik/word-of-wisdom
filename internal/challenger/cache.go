package challenger

import (
	"time"

	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/jellydator/ttlcache/v3"
)

type ChallengeCache struct {
	// Thread-safe cache
	cache      *ttlcache.Cache[string, domain.Challenge]
	defaultTTL time.Duration
}

func NewChallengeCache(defaultTTL time.Duration) *ChallengeCache {
	return &ChallengeCache{
		defaultTTL: defaultTTL,
		cache:      ttlcache.New[string, domain.Challenge](ttlcache.WithTTL[string, domain.Challenge](defaultTTL)),
	}
}

func (c *ChallengeCache) Get(key string) (domain.Challenge, bool) {
	if !c.cache.Has(key) {
		return domain.Challenge{}, false
	}
	return c.cache.Get(key).Value(), true
}

func (c *ChallengeCache) Set(key string, ch domain.Challenge) {
	c.cache.Set(key, ch, c.defaultTTL)
}

func (c *ChallengeCache) Delete(key string) {
	c.cache.Delete(key)
}
