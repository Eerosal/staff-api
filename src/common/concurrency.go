package common

import (
	"go.uber.org/ratelimit"
	"sync"
)

func RunRateLimited[K any](keys []K, perSecond int, fn func(key K)) {
	rl := ratelimit.New(perSecond)
	var wg sync.WaitGroup

	wg.Add(len(keys))

	for _, key := range keys {
		key := key
		go func() {
			defer wg.Done()
			rl.Take()
			fn(key)
		}()
	}

	wg.Wait()
}
