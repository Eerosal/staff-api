package avatar

import (
	"sync"
	"time"
)

var avatarCacheMutex sync.RWMutex
var avatarCache = make(map[string]*avatarCacheEntry)

type avatarCacheEntry struct {
	LastFetchTime time.Time
	LastAvatar    *string
}

func (c *avatarCacheEntry) clone() *avatarCacheEntry {
	avatar := *c.LastAvatar

	return &avatarCacheEntry{
		LastFetchTime: c.LastFetchTime,
		LastAvatar:    &avatar,
	}
}
