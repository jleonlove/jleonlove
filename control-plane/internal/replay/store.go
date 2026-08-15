package replay

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrReplay  = errors.New("execution envelope already consumed")
	ErrExpired = errors.New("execution envelope expired")
)

type Key struct {
	OrganizationID string
	RequestID      string
	Nonce          string
}

type Store interface {
	Consume(context.Context, Key, time.Time) error
}

type MemoryStore struct {
	mu   sync.Mutex
	used map[Key]time.Time
	now  func() time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{used: make(map[Key]time.Time), now: time.Now}
}

func (s *MemoryStore) Consume(_ context.Context, key Key, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	if !now.Before(expiresAt) {
		return ErrExpired
	}
	for k, expiry := range s.used {
		if !now.Before(expiry) {
			delete(s.used, k)
		}
	}
	if _, ok := s.used[key]; ok {
		return ErrReplay
	}
	s.used[key] = expiresAt
	return nil
}
