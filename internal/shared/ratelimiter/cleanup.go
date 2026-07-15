package ratelimiter

import (
	"time"
)

// start a background goroutine that removes stale limiters
func (s *Service) StartCleanup() {
	ticker := time.NewTicker(s.config.CleanupInterval)

	go func() {
		for range ticker.C {
			s.cleanup()
		}
	}()
}

// remove inactive limiters from the store
func (s *Service) cleanup() {
	now := time.Now()
	keys := s.store.Keys()

	for _, key := range keys {
		entry, ok := s.store.Get(key)
		if !ok {
			continue
		}

		// If idle for too long → delete
		if now.Sub(entry.LastSeen) > s.config.MaxIdleTime {
			s.store.Delete(key)
		}
	}
}
